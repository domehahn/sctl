//go:build e2e

package integration

// Cross-repository E2E contract test between skpm and SkillForge.
//
// Unlike tests/integration/*_test.go (build tag "integration"), which mock
// out the registry HTTP API in-process, this test builds and runs the real
// SkillForge skill-registry binary from a sibling checkout and drives it
// through skpm's actual registry.SkillForgeRegistry client — the same code
// path skpm publish/add use in production. It exists because skpm and
// SkillForge are developed in separate repositories with no shared CI: a
// wire-format or endpoint-shape drift between them would otherwise only be
// caught manually or by an end user.
//
// It deliberately would have caught the two real bugs found by reading both
// codebases side by side while unifying skpm's publish path onto
// internal/registry:
//   - GenericHTTPRegistry.Publish() hardcoded Content-Type: application/zip;
//     SkillForge's Publish handler rejects any Content-Type other than
//     application/zip or application/gzip with 400 INVALID_CONTENT_TYPE.
//   - skpm's internal/publisher (now removed) never supported the
//     skillforge/generic-http registry types SkillForge actually exposes.
//
// Requires a local SkillForge checkout. Resolved, in order:
//   1. SKILLFORGE_REPO env var, if set.
//   2. ../SkillForge relative to this repo's root (the layout used when
//      skpm and SkillForge are checked out as sibling directories).
// Skipped (not failed) if neither is found, or if building the
// skill-registry binary fails for a local toolchain reason (e.g. missing
// libsqlite3-dev / gcc for the CGO sqlite driver) — this test exercises a
// live contract, not skpm's own code, and staying opt-in (go test -tags e2e)
// keeps it out of skpm's normal CI, which has no SkillForge checkout.
//
// Run manually with:
//   SKILLFORGE_REPO=/path/to/SkillForge go test -tags e2e ./tests/integration/ -run TestSkillForgeE2EPublishResolveDownload -v

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/domehahn/skpm/v2/internal/skill"
	"github.com/stretchr/testify/require"
)

// findSkillForgeRepo resolves the SkillForge checkout to build against, or
// returns "" if none can be found.
func findSkillForgeRepo(t *testing.T) string {
	t.Helper()
	if v := os.Getenv("SKILLFORGE_REPO"); v != "" {
		if _, err := os.Stat(filepath.Join(v, "skill-registry", "go.mod")); err == nil {
			return v
		}
		t.Fatalf("SKILLFORGE_REPO=%q does not contain skill-registry/go.mod", v)
	}

	wd, err := os.Getwd()
	require.NoError(t, err)
	// tests/integration -> repo root is two levels up.
	skpmRoot := filepath.Join(wd, "..", "..")
	sibling := filepath.Join(skpmRoot, "..", "SkillForge")
	if _, err := os.Stat(filepath.Join(sibling, "skill-registry", "go.mod")); err == nil {
		return sibling
	}
	return ""
}

// freeTCPPort picks an available loopback port. There is an inherent small
// race between closing this listener and the child process binding it, but
// it's the same tradeoff every "find a free port for a test server" helper
// makes, and good enough for an opt-in local/manual test.
func freeTCPPort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// startSkillForge builds the SkillForge skill-registry binary and runs it
// against an ephemeral SQLite database with auth disabled (SkillForge's own
// default config, matching the pattern its own CI integration-tests.yml
// uses). Returns the server's base URL and a cleanup func.
func startSkillForge(t *testing.T, repoDir string) (baseURL string, cleanup func()) {
	t.Helper()

	registryDir := filepath.Join(repoDir, "skill-registry")
	buildDir := t.TempDir()
	binPath := filepath.Join(buildDir, "skill-registry-e2e")

	build := exec.Command("go", "build", "-o", binPath, "./cmd/skill-registry")
	build.Dir = registryDir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("could not build SkillForge skill-registry binary (likely missing local toolchain deps, e.g. gcc/libsqlite3-dev for CGO sqlite): %v\n%s", err, out)
	}

	port := freeTCPPort(t)
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	workDir := t.TempDir() // no config.yaml here -> server.Load() falls back to DefaultConfig()
	dataDir := filepath.Join(workDir, "data")
	// The sqlite driver opens its file directly and does not create parent
	// directories itself (unlike the package storage layer, which does).
	require.NoError(t, os.MkdirAll(dataDir, 0o755))

	cmd := exec.Command(binPath, "-config", "./config.yaml")
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(),
		"SKILL_REGISTRY_ADDR="+addr,
		"SKILL_REGISTRY_DATA_DIR="+dataDir,
		"SKILL_REGISTRY_DB_DRIVER=sqlite",
		"SKILL_REGISTRY_DB_PATH="+filepath.Join(dataDir, "registry.db"),
	)
	stdout, err := os.Create(filepath.Join(workDir, "stdout.log"))
	require.NoError(t, err)
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	require.NoError(t, cmd.Start())

	base := "http://" + addr
	deadline := time.Now().Add(20 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := http.Get(base + "/healthz")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				lastErr = nil
				break
			}
		}
		lastErr = err
		time.Sleep(200 * time.Millisecond)
	}
	if lastErr != nil || time.Now().After(deadline) {
		_ = cmd.Process.Kill()
		logData, _ := os.ReadFile(filepath.Join(workDir, "stdout.log"))
		t.Fatalf("SkillForge did not become healthy at %s: %v\nserver log:\n%s", base, lastErr, logData)
	}

	return base, func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		stdout.Close()
	}
}

// TestSkillForgeE2EPublishResolveDownload publishes the on-disk
// testdata/skills/example-skill fixture through skpm's real registry client
// (registry.SkillForgeRegistry, the same one skpm publish/add use) against a
// live SkillForge instance, then resolves, downloads, and byte-verifies it —
// the same round trip skpm publish -> skpm add exercises against a real
// SkillForge deployment.
func TestSkillForgeE2EPublishResolveDownload(t *testing.T) {
	repoDir := findSkillForgeRepo(t)
	if repoDir == "" {
		t.Fatal("no SkillForge checkout found (set SKILLFORGE_REPO or check out ../SkillForge next to skpm) — skipping cross-repo E2E contract test")
	}

	baseURL, cleanup := startSkillForge(t, repoDir)
	defer cleanup()

	ctx := context.Background()
	const namespace = "e2e-test"

	// Write a minimal fixture satisfying the current SKILL.md versioning
	// schema (version/since/last_modified/authors/stability/
	// min_platform_version/changelog + a matching VERSION and
	// CHANGELOG.md) and package it the same way `skpm publish` does.
	fixtureDir := writeE2ESkillFixture(t)
	outputDir := t.TempDir()
	pkg := skill.NewPackager()
	pkgResult, err := pkg.Package(ctx, fixtureDir, outputDir)
	require.NoError(t, err)

	reg := registry.NewSkillForgeRegistry("skillforge-e2e", config.RegistryConfig{
		Type:      "skillforge",
		URL:       baseURL,
		Namespace: namespace,
	})

	// Assigning through the interface types (rather than calling the
	// concrete *GenericHTTPRegistry methods directly) is itself part of
	// the contract check: it fails to compile if GenericHTTPRegistry ever
	// stops satisfying these registry.* capability interfaces.
	var pub registry.PublishingRegistry = reg

	// ── Publish ──────────────────────────────────────────────────────
	pubResult, err := pub.Publish(ctx, registry.PublishRequest{
		ArtifactPath: pkgResult.OutputPath,
		Manifest: skill.SkillManifest{
			Name:    pkgResult.Name,
			Version: pkgResult.Version,
		},
		SHA256:      pkgResult.SHA256,
		PackageType: "zip",
	})
	require.NoError(t, err, "publish against live SkillForge failed")
	require.True(t, pubResult.Created)
	require.Equal(t, pkgResult.Name, pubResult.Name)
	require.Equal(t, pkgResult.Version, pubResult.Version)
	require.Equal(t, pkgResult.SHA256, pubResult.SHA256)

	// ── Resolve ──────────────────────────────────────────────────────
	artifact, err := reg.Resolve(ctx, registry.ResolveRequest{
		Ref:        registry.SkillRef{Namespace: namespace, Name: pkgResult.Name},
		Constraint: pkgResult.Version,
	})
	require.NoError(t, err, "resolve against live SkillForge failed")
	require.Equal(t, pkgResult.Version, artifact.Version)
	require.Equal(t, pkgResult.SHA256, artifact.SHA256)
	require.Equal(t, "zip", artifact.PackageType)
	require.NotEmpty(t, artifact.DownloadURL)

	// ── Download and byte-verify ────────────────────────────────────
	var buf countingBuffer
	require.NoError(t, reg.Download(ctx, artifact, &buf))
	sum := sha256.Sum256(buf.data)
	require.Equal(t, pkgResult.SHA256, hex.EncodeToString(sum[:]),
		"downloaded artifact digest does not match the published SHA256 — publish/download round trip is broken")

	// ── Discovery (Info) ─────────────────────────────────────────────
	disc, ok := interface{}(reg).(registry.DiscoveryRegistry)
	require.True(t, ok, "SkillForgeRegistry must implement registry.DiscoveryRegistry")
	info, err := disc.Info(ctx, registry.SkillRef{Namespace: namespace, Name: pkgResult.Name})
	require.NoError(t, err)
	require.Equal(t, pkgResult.Version, info.LatestVersion)

	// ── Attestation ──────────────────────────────────────────────────
	// Exercises the mapping in generic_http.go's SkillForge defaults
	// against SkillForge's real /artifacts/skill/.../attestations
	// endpoints (skills are mirrored server-side as kind="skill"
	// artifacts — see SkillForge's registry.mirrorSkillArtifact) — proof
	// that skil attestations attach as first-class registry metadata, not
	// just a documented convention.
	att, ok := interface{}(reg).(registry.AttestationRegistry)
	require.True(t, ok, "SkillForgeRegistry must implement registry.AttestationRegistry")
	verRef := registry.SkillVersionRef{Namespace: namespace, Name: pkgResult.Name, Version: pkgResult.Version}
	predicate, err := json.Marshal(map[string]any{
		"version":  1,
		"subject":  map[string]string{"name": pkgResult.Name, "version": pkgResult.Version, "sha256": pkgResult.SHA256},
		"producer": map[string]string{"name": "skil", "version": "e2e-test"},
		"result":   map[string]any{"status": "pass", "verdict": "clear", "risk_score": 0},
	})
	require.NoError(t, err)
	attRec, err := att.Attest(ctx, verRef, registry.AttestationRequest{
		// SkillForge's attestation type is a closed enum
		// (signature/scan/provenance/sbom), not a free-form predicate-type
		// URI — "scan" is the closest fit for a skil scan/eval result and
		// is also what internal/cli/attestation.go defaults --type to.
		Type:      "scan",
		Digest:    pkgResult.SHA256,
		Predicate: predicate,
	})
	require.NoError(t, err, "attest against live SkillForge failed")
	require.Equal(t, "scan", attRec.Type)
	require.Equal(t, pkgResult.SHA256, attRec.Digest)

	listed, err := att.ListAttestations(ctx, verRef)
	require.NoError(t, err, "list attestations against live SkillForge failed")
	require.Len(t, listed, 1)
	require.Equal(t, pkgResult.SHA256, listed[0].Digest)

	// ── Governance (Deprecate) ──────────────────────────────────────
	gov, ok := interface{}(reg).(registry.GovernanceRegistry)
	require.True(t, ok, "SkillForgeRegistry must implement registry.GovernanceRegistry")
	err = gov.Deprecate(ctx, registry.SkillVersionRef{
		Namespace: namespace,
		Name:      pkgResult.Name,
		Version:   pkgResult.Version,
	}, "e2e contract test cleanup")
	require.NoError(t, err, "deprecate against live SkillForge failed")
}

// writeE2ESkillFixture writes a minimal skill satisfying skpm's current
// SKILL.md versioning schema (see internal/skill/validator.go and
// sklib/spec/skill_md.go) to a temp directory and returns its path.
func writeE2ESkillFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"SKILL.md": `---
name: e2e-fixture-skill
description: Cross-repo E2E contract test fixture.
version: "1.0.0"
since: "2026-01-01"
last_modified: "2026-01-01"
authors:
  - skpm-e2e
stability: stable
min_platform_version:
  codex: "unknown"
deprecated_since:
replaces:
supersedes: []
changelog:
  - version: "1.0.0"
    date: "2026-01-01"
    change: "Initial release"
---

# E2E Fixture Skill

Published by tests/integration/skillforge_e2e_test.go against a live
SkillForge instance.

## Changelog

### 1.0.0 - 2026-01-01

- Initial release.
`,
		"VERSION":      "1.0.0",
		"CHANGELOG.md": "# Changelog\n\n## 1.0.0\n\n- Initial release\n",
		"README.md":    "# E2E Fixture Skill\n",
		"LICENSE":      "MIT\n",
		// namespace must match the registry namespace this fixture is
		// published under (see TestSkillForgeE2EPublishResolveDownload) —
		// SkillForge's Publish handler rejects a manifest namespace that
		// doesn't match the URL namespace with 400 VALIDATION_FAILED.
		"skill.yaml": `name: e2e-fixture-skill
version: "1.0.0"
description: Cross-repo E2E contract test fixture.
namespace: e2e-test
compatible_with:
  - claude-code
  - gitlab-duo
`,
		// skill.NewPackager() validates with Strict:true, which requires a
		// tests/ directory to be present.
		"tests/smoke.md": "Placeholder test note for the E2E fixture skill.\n",
	}
	for name, content := range files {
		path := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}
	return dir
}

type countingBuffer struct{ data []byte }

func (b *countingBuffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── token: command registration ───────────────────────────────────────────────

func TestTokenCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	var tokenCmd interface{ Commands() []interface{} }
	_ = tokenCmd
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "token" {
			found = true
			subNames := make(map[string]bool)
			for _, s := range sub.Commands() {
				subNames[s.Use] = true
			}
			assert.True(t, subNames["list"], "token list should be a subcommand")
			assert.True(t, subNames["add"], "token add should be a subcommand")
			assert.True(t, subNames["revoke"], "token revoke should be a subcommand")
			break
		}
	}
	assert.True(t, found, "token command should be registered on root")
}

func TestTokenAddHasExpectedFlags(t *testing.T) {
	root := cli.NewRootCmd()
	for _, sub := range root.Commands() {
		if sub.Use == "token" {
			for _, s := range sub.Commands() {
				if s.Use == "add" {
					assert.NotNil(t, s.Flags().Lookup("value"))
					assert.NotNil(t, s.Flags().Lookup("from-env"))
					assert.NotNil(t, s.Flags().Lookup("registry"))
				}
			}
		}
	}
}

// ── policy: CheckPolicy logic ─────────────────────────────────────────────────

func buildPolicyLockfile(t *testing.T) *lockfile.LockFile {
	t.Helper()
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "sec-scan",
		Version:     "2.1.0",
		Source:      "https://registry.agentskills.io",
		SHA256:      strings.Repeat("a", 64),
		InstalledTo: []string{".claude/skills/sec-scan"},
	})
	lf.Upsert(lockfile.SkillLock{
		Name:        "old-tool",
		Namespace:   "myorg",
		Version:     "0.3.0",
		Source:      "https://registry.agentskills.io",
		SHA256:      strings.Repeat("b", 64),
		InstalledTo: []string{".claude/skills/old-tool"},
	})
	lf.Upsert(lockfile.SkillLock{
		Name:        "third-party",
		Version:     "1.0.0",
		Source:      "https://unknown-registry.io",
		SHA256:      strings.Repeat("c", 64),
		InstalledTo: []string{".claude/skills/third-party"},
	})
	return lf
}

func TestCheckPolicyNoRulesNoViolations(t *testing.T) {
	lf := buildPolicyLockfile(t)
	p := &cli.Policy{Version: 1}
	violations := cli.CheckPolicy(p, lf)
	assert.Empty(t, violations)
}

func TestCheckPolicyAllowedRegistries(t *testing.T) {
	lf := buildPolicyLockfile(t)
	p := &cli.Policy{
		Version: 1,
		Rules: cli.PolicyRules{
			AllowedRegistries: []string{"https://registry.agentskills.io"},
		},
	}
	violations := cli.CheckPolicy(p, lf)
	require.Len(t, violations, 1)
	assert.Equal(t, "allowed_registries", violations[0].Rule)
	assert.Equal(t, "third-party", violations[0].Skill)
}

func TestCheckPolicyBannedSkills(t *testing.T) {
	lf := buildPolicyLockfile(t)
	p := &cli.Policy{
		Version: 1,
		Rules: cli.PolicyRules{
			BannedSkills: []string{"third-party"},
		},
	}
	violations := cli.CheckPolicy(p, lf)
	require.Len(t, violations, 1)
	assert.Equal(t, "banned_skills", violations[0].Rule)
}

func TestCheckPolicyBannedSkillsNamespaced(t *testing.T) {
	lf := buildPolicyLockfile(t)
	p := &cli.Policy{
		Version: 1,
		Rules: cli.PolicyRules{
			BannedSkills: []string{"myorg/old-tool"},
		},
	}
	violations := cli.CheckPolicy(p, lf)
	require.Len(t, violations, 1)
	assert.Equal(t, "myorg/old-tool", violations[0].Skill)
}

func TestCheckPolicyRequireSignatures(t *testing.T) {
	lf := buildPolicyLockfile(t)
	p := &cli.Policy{
		Version: 1,
		Rules: cli.PolicyRules{
			RequireSignatures: true,
		},
	}
	violations := cli.CheckPolicy(p, lf)
	assert.Equal(t, len(lf.Skills), len(violations), "every unsigned skill should be flagged")
	for _, v := range violations {
		assert.Equal(t, "require_signatures", v.Rule)
	}
}

func TestCheckPolicyRequireSignaturesSkipsSignedSkills(t *testing.T) {
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:      "signed-skill",
		Version:   "1.0.0",
		Source:    "https://r.example.com",
		SHA256:    strings.Repeat("a", 64),
		Signature: "some-valid-signature",
	})
	lf.Upsert(lockfile.SkillLock{
		Name:    "unsigned-skill",
		Version: "1.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("b", 64),
	})
	p := &cli.Policy{Version: 1, Rules: cli.PolicyRules{RequireSignatures: true}}
	violations := cli.CheckPolicy(p, lf)
	require.Len(t, violations, 1)
	assert.Equal(t, "unsigned-skill", violations[0].Skill)
}

func TestCheckPolicyMinVersions(t *testing.T) {
	lf := buildPolicyLockfile(t)
	p := &cli.Policy{
		Version: 1,
		Rules: cli.PolicyRules{
			MinVersions: map[string]string{
				"sec-scan": "3.0.0",  // installed 2.1.0 — violation
				"old-tool": "0.2.0", // installed 0.3.0 — ok
			},
		},
	}
	violations := cli.CheckPolicy(p, lf)
	require.Len(t, violations, 1)
	assert.Equal(t, "min_versions", violations[0].Rule)
	assert.Equal(t, "sec-scan", violations[0].Skill)
}

func TestCheckPolicyMinVersionsMet(t *testing.T) {
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "my-skill",
		Version: "2.0.0",
		Source:  "https://r.example.com",
		SHA256:  strings.Repeat("a", 64),
	})
	p := &cli.Policy{
		Version: 1,
		Rules:   cli.PolicyRules{MinVersions: map[string]string{"my-skill": "2.0.0"}},
	}
	assert.Empty(t, cli.CheckPolicy(p, lf))
}

func TestCheckPolicyMultipleViolationsCollected(t *testing.T) {
	lf := buildPolicyLockfile(t)
	p := &cli.Policy{
		Version: 1,
		Rules: cli.PolicyRules{
			AllowedRegistries: []string{"https://registry.agentskills.io"},
			BannedSkills:      []string{"sec-scan"},
			RequireSignatures: true,
		},
	}
	violations := cli.CheckPolicy(p, lf)
	assert.Greater(t, len(violations), 1, "multiple rules should accumulate violations")
}

func TestCheckPolicyEmptyLockfile(t *testing.T) {
	lf := lockfile.New()
	p := &cli.Policy{
		Version: 1,
		Rules: cli.PolicyRules{
			AllowedRegistries: []string{"https://registry.agentskills.io"},
			RequireSignatures: true,
		},
	}
	assert.Empty(t, cli.CheckPolicy(p, lf), "empty lockfile should never produce violations")
}

// ── policy: command registration ──────────────────────────────────────────────

func TestPolicyCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "policy" {
			found = true
			subNames := make(map[string]bool)
			for _, s := range sub.Commands() {
				subNames[s.Use] = true
			}
			assert.True(t, subNames["init"])
			assert.True(t, subNames["check"])
			assert.True(t, subNames["show"])
			break
		}
	}
	assert.True(t, found, "policy command should be registered on root")
}

// ── policy: init command writes file ─────────────────────────────────────────

func TestPolicyInitWritesFile(t *testing.T) {
	dir := t.TempDir()
	policyPath := filepath.Join(dir, "skpm-policy.yaml")

	root := cli.NewRootCmd()
	root.SetArgs([]string{"policy", "init", "--policy", policyPath})
	err := root.Execute()
	require.NoError(t, err)

	data, readErr := os.ReadFile(policyPath)
	require.NoError(t, readErr)
	assert.Contains(t, string(data), "version: 1")
	assert.Contains(t, string(data), "allowed_registries")
}

func TestPolicyInitFailsIfFileExists(t *testing.T) {
	dir := t.TempDir()
	policyPath := filepath.Join(dir, "skpm-policy.yaml")
	require.NoError(t, os.WriteFile(policyPath, []byte("version: 1\n"), 0o644))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"policy", "init", "--policy", policyPath})
	err := root.Execute()
	assert.Error(t, err)
}

// ── notice: collectNoticeEntries ──────────────────────────────────────────────

func TestNoticeCommandRegistered(t *testing.T) {
	root := cli.NewRootCmd()
	found := false
	for _, sub := range root.Commands() {
		if sub.Use == "notice" {
			found = true
			subNames := make(map[string]bool)
			for _, s := range sub.Commands() {
				subNames[s.Use] = true
			}
			assert.True(t, subNames["generate"])
			assert.True(t, subNames["check"])
			break
		}
	}
	assert.True(t, found, "notice command should be registered on root")
}

func TestNoticeGenerateWritesFile(t *testing.T) {
	dir := t.TempDir()

	// Scaffold a fake installed skill directory with skill.yaml + LICENSE.
	skillDir := filepath.Join(dir, ".claude", "skills", "my-skill")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "skill.yaml"), []byte("license: MIT\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "LICENSE"), []byte("MIT License text here"), 0o644))

	// Write a lockfile pointing to the skill dir.
	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "my-skill",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("a", 64),
		InstalledTo: []string{skillDir},
	})
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	noticePath := filepath.Join(dir, "NOTICE.md")

	root := cli.NewRootCmd()
	root.SetArgs([]string{"notice", "generate", "--lock", lockPath, "--out", noticePath})
	require.NoError(t, root.Execute())

	data, err := os.ReadFile(noticePath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "my-skill")
	assert.Contains(t, string(data), "MIT")
	assert.Contains(t, string(data), "MIT License text here")
}

func TestNoticeCheckReportsUnlicensedSkills(t *testing.T) {
	dir := t.TempDir()

	// Skill with no skill.yaml and no LICENSE file.
	skillDir := filepath.Join(dir, ".claude", "skills", "unlicensed")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "unlicensed",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("a", 64),
		InstalledTo: []string{skillDir},
	})
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"notice", "check", "--lock", lockPath})
	// Should succeed (no --require flag).
	assert.NoError(t, root.Execute())
}

func TestNoticeCheckFailsWithRequireWhenMissing(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "skills", "no-license")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "no-license",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("a", 64),
		InstalledTo: []string{skillDir},
	})
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"notice", "check", "--lock", lockPath, "--require"})
	assert.Error(t, root.Execute())
}

func TestNoticeCheckPassesWhenAllLicensed(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "skills", "licensed")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "skill.yaml"), []byte("license: Apache-2.0\n"), 0o644))

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:        "licensed",
		Version:     "1.0.0",
		Source:      "https://r.example.com",
		SHA256:      strings.Repeat("a", 64),
		InstalledTo: []string{skillDir},
	})
	lockPath := filepath.Join(dir, "agent-skills.lock")
	require.NoError(t, lf.Write(lockPath))

	root := cli.NewRootCmd()
	root.SetArgs([]string{"notice", "check", "--lock", lockPath, "--require"})
	assert.NoError(t, root.Execute())
}

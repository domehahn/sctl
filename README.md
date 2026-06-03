# sctl — Skill Control

A production-ready package manager for AI agent skills. Manages [SKILL.md](https://docs.gitlab.com/user/duo_agent_platform/customize/agent_skills/)-based capability bundles as versioned, verifiable artifacts — for Claude Code, GitLab Duo, GitHub Copilot, and Codex.

Think `npm` or `cargo`, but for agent skills.

---

## Why

Agent skills are SKILL.md files that give coding assistants domain-specific knowledge — security policy reviewers, documentation checkers, framework-specific guides. Without tooling they tend to be copy-pasted across repos, drift out of sync, and have no integrity guarantees.

`sctl` treats skills like software:

- **Version-pinned** — `agent-skills.lock` records exact versions and SHA256 checksums
- **Reproducible** — any CI run installs the exact same bytes
- **Multi-platform** — one install writes to all compatible platform paths automatically
- **Auditable** — every artifact is verified before it touches disk

---

## Install

```bash
go install github.com/domehahn/sctl/cmd/sctl@latest
```

Or download a pre-built binary from the [releases page](https://github.com/domehahn/sctl/releases).

---

## Quick Start

```bash
# Add a skill from a registry
sctl add gitlab-policy-reviewer@1.5.0 --source myregistry

# Install all skills pinned in the lockfile
sctl install

# Validate a skill before shipping it
sctl validate ./skills/gitlab-policy-reviewer

# Package a skill into a distributable ZIP
sctl package ./skills/gitlab-policy-reviewer
```

---

## Commands

### `sctl install`

Reads `agent-skills.lock` and installs all pinned skills.

```
sctl install [--lock <path>] [--dry-run] [--concurrency N]
```

- Downloads artifacts in parallel (default 4 concurrent)
- Verifies SHA256 before writing anything to disk
- Uses a local cache (`~/.cache/sctl/`) — re-runs are instant
- Atomic installs: staging directory → rename, never partial state
- `--dry-run` prints what would be installed without writing files

**Platform paths installed per `compatible_with`:**

| Platform | Path |
|---|---|
| `claude-code` | `.claude/skills/<name>/` |
| `gitlab-duo` | `skills/<name>/` and `.agents/skills/<name>/` |
| `github-copilot` | `.github/skills/<name>/` |
| `codex` | `.agents/skills/<name>/` |

### `sctl add <skill[@version]>`

Resolves, downloads, and installs a skill, then writes `agent-skills.lock`.

```
sctl add gitlab-policy-reviewer@1.5.0 --source myregistry
sctl add documentation-reviewer       --source github
```

- Version defaults to latest if omitted
- Reads `compatible_with` from the artifact's `skill.yaml` to determine install paths
- Creates or updates `agent-skills.lock` atomically

### `sctl validate [path]`

Validates a skill directory structure. Default path is the current directory.

```
sctl validate ./skills/gitlab-policy-reviewer
sctl validate --output json
```

Checks:
- `SKILL.md` exists and is non-empty
- `VERSION` contains a valid semver string
- `skill.yaml` has required fields (`name`, `version`, `compatible_with`); `version` matches `VERSION`
- `compatible_with` values are recognized platform names
- `CHANGELOG.md` contains an entry for the current version (warning if missing)

Exits `0` if valid, `1` if errors are found. Warnings do not fail the check.

### `sctl package [path]`

Packages a skill directory into a distributable ZIP artifact.

```
sctl package ./skills/gitlab-policy-reviewer
sctl package ./skills/gitlab-policy-reviewer --output-dir ./dist
```

- Runs `validate` first — packaging fails if the skill is invalid
- Creates `<name>-<version>.zip` containing all skill files plus `manifest.json`
- Prints the SHA256 of the ZIP — paste this into your lockfile or registry configuration
- Excludes `.git/`, `*.tmp`, `*.part`

### `sctl version`

```
sctl version
sctl version --output json
```

---

## Global Flags

| Flag | Default | Description |
|---|---|---|
| `--output` | `text` | Output format: `text` or `json` |
| `--dry-run` | `false` | Print planned actions without writing files |
| `--verbose` | `false` | Enable debug logging to stderr |
| `--concurrency` | `4` | Maximum parallel downloads |

All commands support `--output json` for machine-readable output — useful in CI pipelines.

---

## Lockfile

`agent-skills.lock` is the source of truth for which skill version is active in a project. Commit it.

```yaml
version: 1
skills:
  - name: gitlab-policy-reviewer
    version: 1.5.0
    source: myregistry
    source_url: https://artifactory.company.com/agent-skills/gitlab-policy-reviewer/1.5.0/gitlab-policy-reviewer-1.5.0.zip
    sha256: b875cc1f70dfe77f7626fddda0cc7315643d851169cdf09a104d6233e0888ff5
    installed_to:
      - .claude/skills/gitlab-policy-reviewer
      - skills/gitlab-policy-reviewer
      - .agents/skills/gitlab-policy-reviewer
```

---

## Skill Structure

A valid skill directory:

```
my-skill/
  SKILL.md        # Agent-readable capability definition
  VERSION         # Semver string, e.g. "1.5.0"
  skill.yaml      # Machine-readable metadata
  CHANGELOG.md    # Version history (recommended)
```

**`skill.yaml`:**

```yaml
name: gitlab-policy-reviewer
version: 1.5.0
description: Reviews GitLab security policy YAML and approval policies.
owners:
  - platform-security
compatible_with:
  - claude-code
  - gitlab-duo
  - github-copilot
  - codex
```

Use `all` in `compatible_with` to install to every supported platform path.

---

## Configuration

`sctl` reads `~/.config/sctl/config.yaml` (or `$XDG_CONFIG_HOME/sctl/config.yaml`):

```yaml
default_registry: myregistry
cache_dir: ~/.cache/sctl

registries:
  myregistry:
    type: artifactory
    url: https://artifactory.company.com/artifactory#agent-skills
    token: ""           # set via SCTL_REGISTRY_TOKEN

  company-github:
    type: github
    url: myorg/agent-skills-repo
    token: ""

  company-gitlab:
    type: gitlab
    url: https://gitlab.company.com
    token: ""
```

**Supported registry types:**

| Type | Description |
|---|---|
| `github` | GitHub Releases — `url` is `owner/repo` |
| `gitlab` | GitLab Releases — `url` is the GitLab base URL |
| `artifactory` | JFrog Artifactory Generic Repo — `url` is `<base-url>#<repo-name>` |
| `local` | Local filesystem — `url` is the base directory path |

**Environment variables** override config file values:

| Variable | Overrides |
|---|---|
| `SCTL_CACHE_DIR` | `cache_dir` |
| `SCTL_LOG_LEVEL` | `log_level` |
| `SCTL_REGISTRY_TOKEN` | Token for the `default_registry` |

---

## CI Integration

```yaml
# .gitlab-ci.yml
install-skills:
  script:
    - sctl install --concurrency 8
  cache:
    key: sctl-$CI_COMMIT_REF_SLUG
    paths:
      - ~/.cache/sctl/
```

```yaml
# .github/workflows/skills.yml
- name: Install agent skills
  run: sctl install
  env:
    SCTL_REGISTRY_TOKEN: ${{ secrets.SKILLS_REGISTRY_TOKEN }}
```

For JSON output in CI (e.g. to feed into jq):

```bash
sctl install --output json | jq '.data.installed[]'
```

---

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | User error (missing file, validation failure, SHA256 mismatch) |
| `2` | Infrastructure error (network failure, I/O error) |

---

## Development

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Build binary
make build               # → dist/sctl

# Lint
make lint                # requires golangci-lint

# Cross-compile release snapshot
make release-snapshot    # requires goreleaser
```

**Project layout:**

```
cmd/sctl/             # Entrypoint
internal/
  cli/                # Cobra commands
  config/             # Config loading + env override
  lockfile/           # agent-skills.lock read/write
  skill/              # Types, validator, packager
  registry/           # Registry backends + factory
  cache/              # SHA256-keyed disk cache
  installer/          # Download + atomic install + platform paths
  progress/           # CI-aware progress bars
testdata/             # Fixture skills for tests
tests/integration/    # Integration test suite
```

---

## Security

Skills are supply-chain artifacts. `sctl` enforces:

- **SHA256 verification** on every download before writing to disk
- **Zip-slip protection** — path traversal in ZIP entries is rejected
- **Atomic writes** — installations are all-or-nothing; no partial state on failure
- **Cache integrity** — cache keys are the artifact's SHA256; collisions are impossible

For production use, sign your release artifacts and include the signature in your registry metadata. Never install skills with `sha256: ""` in production lockfiles.

---

## License

MIT

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
- **Recoverable** — deleted lockfile? `sctl install` regenerates it from `agent-skills.yaml`

---

## Install

```bash
go install github.com/domehahn/sctl/cmd/sctl@latest
```

Or download a pre-built binary from the [releases page](https://github.com/domehahn/sctl/releases).

---

## Quick Start

```bash
# 1 — set up sctl for the first time
sctl init config

# 2 — initialize a project
sctl init project

# 3 — add skills
sctl add gitlab-policy-reviewer@1.5.0 --source myregistry
sctl add documentation-reviewer       --source myregistry

# 4 — commit both files
git add agent-skills.yaml agent-skills.lock
git commit -m "add agent skills"

# 5 — any machine, any time
sctl install
```

---

## Two-File Model

`sctl` uses two files, both committed to Git — similar to Poetry:

| File | Analogy | Purpose |
| --- | --- | --- |
| `agent-skills.yaml` | `pyproject.toml` | Manifest: which skills you want (human-edited) |
| `agent-skills.lock` | `poetry.lock` | Lockfile: exact versions + SHA256 (generated) |

`sctl install` behaviour:

- **Lockfile present** → installs exactly what is pinned (fast, deterministic)
- **Lockfile missing**, manifest present → resolves versions from registry, generates lockfile, installs
- **Both missing** → error with instructions

This means a deleted lockfile is never a problem: `sctl install` regenerates it from `agent-skills.yaml` automatically.

**`agent-skills.yaml`** (what you declare):

```yaml
version: 1
skills:
  - name: gitlab-policy-reviewer
    version: 1.5.0
    source: myregistry
  - name: documentation-reviewer
    source: myregistry   # no version = latest
```

**`agent-skills.lock`** (what gets installed — never edit manually):

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

## Commands

### `sctl init`

Scaffold sctl configuration and project files interactively.

#### `sctl init config`

Creates `~/.config/sctl/config.yaml` via an interactive wizard.

```
sctl init config [--force]
```

Asks for registry type, URL, project path, and token (masked input). Running it again without `--force` is a no-op if the file already exists.

#### `sctl init project`

Creates `agent-skills.yaml`, `agent-skills.lock`, and updates `.gitignore`.

```
sctl init project [--force]
```

The `.gitignore` update excludes installed skill directories (generated artifacts) and adds a comment that both `agent-skills.yaml` and `agent-skills.lock` must be committed.

#### `sctl init skill <name>`

Scaffolds a complete, immediately valid skill directory.

```
sctl init skill my-skill
sctl init skill my-skill --output-dir ./skills
```

Asks for description, version, owner, and compatible platforms. The generated skill passes `sctl validate` out of the box.

---

### `sctl config`

Inspect and validate the sctl configuration.

#### `sctl config validate`

Validates `~/.config/sctl/config.yaml`.

```
sctl config validate
sctl config validate --path ./custom-config.yaml
```

Checks:

- File is valid YAML
- `default_registry` references a defined registry entry
- Every registry has a known `type` (`github`, `gitlab`, `artifactory`, `local`)
- `gitlab` registries have `project` set in `namespace/project` format
- `github` registries have `url` in `owner/repo` format (not a full URL)
- `artifactory` registries have `url` in `<base-url>#<repo-name>` format

#### `sctl config show`

Prints the resolved config with all environment variable overrides applied. Tokens are masked as `***`.

```
sctl config show
sctl config show --output json
```

---

### `sctl install`

Installs all skills. Generates `agent-skills.lock` from `agent-skills.yaml` if the lockfile is missing.

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
| --- | --- |
| `claude-code` | `.claude/skills/<name>/` |
| `gitlab-duo` | `skills/<name>/` and `.agents/skills/<name>/` |
| `github-copilot` | `.github/skills/<name>/` |
| `codex` | `.agents/skills/<name>/` |

---

### `sctl add <skill[@version]>`

Resolves, downloads, and installs a skill, then updates both `agent-skills.yaml` and `agent-skills.lock`.

```
sctl add gitlab-policy-reviewer@1.5.0 --source myregistry
sctl add documentation-reviewer       --source myregistry
```

- Version defaults to latest if omitted
- Reads `compatible_with` from the artifact's `skill.yaml` to determine install paths
- Creates or updates both files atomically

---

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

---

### `sctl package [path]`

Packages a skill directory into a distributable ZIP artifact.

```
sctl package ./skills/gitlab-policy-reviewer
sctl package ./skills/gitlab-policy-reviewer --output-dir ./dist
```

- Runs `validate` first — packaging fails if the skill is invalid
- Creates `<name>-<version>.zip` containing all skill files plus `manifest.json`
- Prints the SHA256 of the ZIP — paste this into your registry configuration
- Excludes `.git/`, `*.tmp`, `*.part`

---

### `sctl version`

```
sctl version
sctl version --output json
```

---

## Global Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--output` | `text` | Output format: `text` or `json` |
| `--dry-run` | `false` | Print planned actions without writing files |
| `--verbose` | `false` | Enable debug logging to stderr |
| `--concurrency` | `4` | Maximum parallel downloads |

All commands support `--output json` for machine-readable output — useful in CI pipelines.

---

## Typical Workflows

### First-time setup

```bash
sctl init config      # create ~/.config/sctl/config.yaml
sctl config validate  # verify it's correct
```

### Start a new project

```bash
sctl init project                                          # agent-skills.yaml + agent-skills.lock + .gitignore
sctl add gitlab-policy-reviewer@1.5.0 --source myregistry # updates both files, installs skill
sctl add documentation-reviewer       --source myregistry
git add agent-skills.yaml agent-skills.lock .gitignore
git commit -m "add agent skills"
```

### Install on a new machine or in CI

```bash
git clone <repo>
sctl install          # reads agent-skills.lock, installs everything
```

### Recover a deleted lockfile

```bash
rm agent-skills.lock  # oops
sctl install          # regenerates agent-skills.lock from agent-skills.yaml, then installs
```

### Author and publish a new skill

```bash
sctl init skill my-skill --output-dir ./skills
# edit skills/my-skill/SKILL.md
sctl validate skills/my-skill
sctl package  skills/my-skill --output-dir dist/
# upload dist/my-skill-0.1.0.zip to your registry
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

`sctl` reads `~/.config/sctl/config.yaml` (or `$XDG_CONFIG_HOME/sctl/config.yaml`).

Create it interactively with `sctl init config`, or write it manually:

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
    url: myorg/agent-skills
    token: ""

  company-gitlab:
    type: gitlab
    url: https://gitlab.company.com
    project: platform/agent-skills   # required for gitlab
    token: ""
```

**Supported registry types:**

| Type | `url` format | `project` |
| --- | --- | --- |
| `github` | `owner/repo` | — |
| `gitlab` | GitLab base URL | `namespace/project` (required) |
| `artifactory` | `<base-url>#<repo-name>` | — |
| `local` | filesystem base path | — |

**Environment variables** override config file values:

| Variable | Overrides |
| --- | --- |
| `SCTL_CACHE_DIR` | `cache_dir` |
| `SCTL_LOG_LEVEL` | `log_level` |
| `SCTL_REGISTRY_TOKEN` | Token for the `default_registry` |

---

## Git Integration

`sctl init project` writes the following to `.gitignore` automatically:

```gitignore
# sctl — installed skill directories are generated artifacts
# restore with: sctl install
.claude/skills/
skills/
.agents/skills/
.github/skills/

# sctl — these files must be committed
# !agent-skills.yaml
# !agent-skills.lock
```

Installed skill directories are generated artifacts — they are excluded from Git. Both `agent-skills.yaml` and `agent-skills.lock` must be committed.

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
| --- | --- |
| `0` | Success |
| `1` | User error (missing file, validation failure, SHA256 mismatch) |
| `2` | Infrastructure error (network failure, I/O error) |

---

## Development

```bash
make test                # go test ./...
make test-integration    # go test -tags integration ./tests/integration/...
make build               # → dist/sctl
make lint                # requires golangci-lint
make release-snapshot    # requires goreleaser
```

**Project layout:**

```
cmd/sctl/             # Entrypoint
internal/
  cli/                # Cobra commands
  config/             # Config loading + env override
  manifest/           # agent-skills.yaml read/write
  lockfile/           # agent-skills.lock read/write
  skill/              # Types, validator, packager
  registry/           # Registry backends + factory
  cache/              # SHA256-keyed disk cache
  installer/          # Download + atomic install + platform paths
  progress/           # CI-aware progress bars
testdata/             # Fixture skills for tests
tests/integration/    # Integration test suite
examples/             # Ready-to-use skill examples and CI templates
```

---

## Security

Skills are supply-chain artifacts. `sctl` enforces:

- **SHA256 verification** on every download before writing to disk
- **Zip-slip protection** — path traversal in ZIP entries is rejected
- **Atomic writes** — installations are all-or-nothing; no partial state on failure
- **Cache integrity** — cache keys are the artifact's SHA256; collisions are impossible

For production use, sign your release artifacts and include the signature in your registry metadata. Never ship skills with `sha256: ""` in production lockfiles.

---

## License

MIT

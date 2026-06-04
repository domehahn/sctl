# skpm — Skill Manager

A production-ready package manager for AI agent skills. Manages [SKILL.md](https://docs.gitlab.com/user/duo_agent_platform/customize/agent_skills/)-based capability bundles as versioned, verifiable artifacts — for Claude Code, GitLab Duo, GitHub Copilot, and Codex.

Think `npm` or `cargo`, but for agent skills.

---

## Why

Agent skills are SKILL.md files that give coding assistants domain-specific knowledge — security policy reviewers, documentation checkers, framework-specific guides. Without tooling they tend to be copy-pasted across repos, drift out of sync, and have no integrity guarantees.

`skpm` treats skills like software:

- **Version-pinned** — `agent-skills.lock` records exact versions and SHA256 checksums
- **Reproducible** — any CI run installs the exact same bytes
- **Multi-platform** — one install writes to all compatible platform paths automatically
- **Auditable** — every artifact is verified before it touches disk
- **Recoverable** — deleted lockfile? `skpm install` regenerates it from `agent-skills.yaml`

---

## Install

```bash
go install github.com/domehahn/skpm/cmd/skpm@latest
```

Or download a pre-built binary from the [releases page](https://github.com/domehahn/sctl/releases).

---

## Quick Start

### As a skill consumer

```bash
# 1 — first-time setup: config + project in one step
skpm init

# 2 — add skills
skpm add gitlab-policy-reviewer@1.5.0 --source myregistry
skpm add documentation-reviewer       --source myregistry

# 3 — commit both files
git add agent-skills.yaml agent-skills.lock
git commit -m "add agent skills"

# 4 — any machine, any time
skpm install
```

### As a skill author

```bash
skpm init skill my-skill          # scaffold
# edit my-skill/SKILL.md
skpm publish my-skill --source myregistry   # validate → package → tag → upload
```

---

## Two-File Model

`skpm` uses two files, both committed to Git — similar to Poetry:

| File | Analogy | Purpose |
| --- | --- | --- |
| `agent-skills.yaml` | `pyproject.toml` | Manifest: which skills you want (human-edited) |
| `agent-skills.lock` | `poetry.lock` | Lockfile: exact versions + SHA256 (generated) |

`skpm install` behaviour:

- **Lockfile present** → installs exactly what is pinned (fast, deterministic)
- **Lockfile missing**, manifest present → resolves versions from registry, generates lockfile, installs
- **Both missing** → error with instructions

A deleted lockfile is never a problem: `skpm install` regenerates it from `agent-skills.yaml` automatically.

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

**`agent-skills.lock`** (generated — never edit manually):

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

### `skpm init`

Interactive setup wizard. Runs two steps in sequence:

1. **Config** — creates `~/.config/skpm/config.yaml` (skipped if already exists)
2. **Project** — creates `agent-skills.yaml`, `agent-skills.lock`, updates `.gitignore`

```bash
skpm init [--force]
```

`--force` re-runs both steps even if files already exist.

#### `skpm init skill <name>`

Scaffolds a complete, immediately valid skill directory.

```bash
skpm init skill my-skill
skpm init skill my-skill --output-dir ./skills
```

Asks for description, version, owner, and compatible platforms. The generated skill passes `skpm validate` out of the box.

---

### `skpm config`

Inspect and validate the skpm configuration.

#### `skpm config validate`

Validates `~/.config/skpm/config.yaml`.

```bash
skpm config validate
skpm config validate --path ./custom-config.yaml
```

Checks:

- File is valid YAML
- `default_registry` references a defined registry entry
- Every registry has a known `type` (`github`, `gitlab`, `artifactory`, `local`)
- `gitlab` registries have `project` set in `namespace/project` format
- `github` registries have `url` in `owner/repo` format (not a full URL)
- `artifactory` registries have `url` in `<base-url>#<repo-name>` format

#### `skpm config show`

Prints the resolved config with all environment variable overrides applied. Tokens are masked as `***`.

```bash
skpm config show
skpm config show --output json
```

---

### `skpm install`

Installs all skills. Generates `agent-skills.lock` from `agent-skills.yaml` if the lockfile is missing.

```bash
skpm install [--lock <path>] [--dry-run] [--concurrency N]
```

- Downloads artifacts in parallel (default 4 concurrent)
- Verifies SHA256 before writing anything to disk
- Uses a local cache (`~/.cache/skpm/`) — re-runs are instant
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

### `skpm add <skill[@version]>`

Resolves, downloads, and installs a skill, then updates both `agent-skills.yaml` and `agent-skills.lock`.

```bash
skpm add gitlab-policy-reviewer@1.5.0 --source myregistry
skpm add documentation-reviewer       --source myregistry
```

**Local paths** — no registry needed:

```bash
skpm add ./my-skill
skpm add ../shared-skills/sdlc-manager
```

**Without a release tag** — download directly from a branch or commit:

```bash
skpm add my-skill --source myregistry --ref main
skpm add my-skill --source myregistry --ref feature/new-checks
```

For monorepos where the skill lives in a subdirectory:

```bash
skpm add my-skill --source myregistry --ref main --path skills/my-skill
```

**Flags:**

| Flag | Description |
| --- | --- |
| `--source` | Registry to use (falls back to `default_registry`) |
| `--ref` | Branch, tag, or commit SHA — skips release lookup |
| `--path` | Path of the skill within the repository (for `--ref` with monorepos) |

- Version defaults to latest if omitted
- For local paths: reads `skill.yaml` directly, copies the directory atomically
- For `--ref`: downloads the archive at that ref, validates, then installs
- Creates or updates both `agent-skills.yaml` and `agent-skills.lock`

---

### `skpm validate [path]`

Validates a skill directory structure. Default path is the current directory.

```bash
skpm validate ./skills/gitlab-policy-reviewer
skpm validate --output json
```

Checks:

- `SKILL.md` exists and is non-empty
- `VERSION` contains a valid semver string
- `skill.yaml` has required fields (`name`, `version`, `compatible_with`); `version` matches `VERSION`
- `compatible_with` values are recognized platform names
- `CHANGELOG.md` contains an entry for the current version (warning if missing)

Exits `0` if valid, `1` if errors are found. Warnings do not fail the check.

---

### `skpm package [path]`

Packages a skill directory into a local ZIP artifact. Default path is the current directory.

```bash
skpm package ./skills/gitlab-policy-reviewer
skpm package ./skills/gitlab-policy-reviewer --output-dir ./dist
```

- Runs `validate` first — fails if the skill is invalid
- Creates `<name>-<version>.zip` with all skill files plus `manifest.json`
- Prints the SHA256 of the ZIP
- Excludes `.git/`, `*.tmp`, `*.part`

> For a full release (tag + upload), use `skpm publish` instead.

---

### `skpm publish [path]`

Runs the complete release pipeline for a skill. Default path is the current directory.

```bash
skpm publish [path] --source <registry> [flags]
```

Steps executed in order:

```text
1. Validate   — checks SKILL.md, VERSION, skill.yaml, CHANGELOG.md
2. Package    — builds <name>-<version>.zip
3. Git tag    — creates <name>/v<version> (or v<version> with --tag-format plain)
4. Git push   — pushes the tag to origin
5. Upload     — uploads the ZIP to the configured registry
```

**Flags:**

| Flag | Default | Description |
| --- | --- | --- |
| `--source` | `default_registry` | Registry to publish to |
| `--tag-format` | `prefixed` | `prefixed` = `<name>/v<ver>` (monorepo), `plain` = `v<ver>` (per-skill repo) |
| `--no-tag` | `false` | Skip creating a git tag |
| `--no-push` | `false` | Skip pushing the tag |
| `--dry-run` | `false` | Preview all steps without making changes |

**Examples:**

```bash
# Monorepo — creates tag gitlab-policy-reviewer/v1.5.0
skpm publish ./skills/gitlab-policy-reviewer --source company-gitlab

# Per-skill repo — creates tag v1.5.0
skpm publish . --source company-gitlab --tag-format plain

# Preview without changes
skpm publish ./skills/my-skill --source company-gitlab --dry-run

# Upload only, no new tag
skpm publish ./skills/my-skill --source company-gitlab --no-tag
```

After publishing, add the skill to a consumer project:

```bash
skpm add my-skill@1.5.0 --source company-gitlab
```

---

### `skpm version`

```bash
skpm version
skpm version --output json
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
skpm init             # config + project in one step
skpm config validate  # verify config is correct
```

### Add skills to a project

```bash
skpm add gitlab-policy-reviewer@1.5.0 --source myregistry
skpm add documentation-reviewer       --source myregistry
skpm add ./local-dev-skill            # from local path
git add agent-skills.yaml agent-skills.lock
git commit -m "add agent skills"
```

### Install on a new machine or in CI

```bash
git clone <repo>
skpm install
```

### Recover a deleted lockfile

```bash
rm agent-skills.lock
skpm install   # regenerates from agent-skills.yaml, then installs
```

### Author and release a new skill

```bash
skpm init skill my-skill --output-dir ./skills
# edit skills/my-skill/SKILL.md
skpm validate skills/my-skill
skpm publish  skills/my-skill --source myregistry
# → validates, packages, creates tag, pushes, uploads
```

---

## Skill Structure

A valid skill directory:

```text
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

`skpm` reads `~/.config/skpm/config.yaml` (or `$XDG_CONFIG_HOME/skpm/config.yaml`).

Create it interactively with `skpm init`, or write it manually:

```yaml
default_registry: myregistry
cache_dir: ~/.cache/skpm

registries:
  myregistry:
    type: artifactory
    url: https://artifactory.company.com/artifactory#agent-skills
    token: ""           # set via SKPM_REGISTRY_TOKEN

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
| `SKPM_CACHE_DIR` | `cache_dir` |
| `SKPM_LOG_LEVEL` | `log_level` |
| `SKPM_REGISTRY_TOKEN` | Token for the `default_registry` |

---

## Git Integration

`skpm init` writes the following to `.gitignore` automatically:

```gitignore
# skpm — installed skill directories are generated artifacts
# restore with: skpm install
.claude/skills/
skills/
.agents/skills/
.github/skills/

# skpm — these files must be committed
# !agent-skills.yaml
# !agent-skills.lock
```

Installed skill directories are generated artifacts — excluded from Git. Both `agent-skills.yaml` and `agent-skills.lock` must be committed.

---

## CI Integration

```yaml
# .gitlab-ci.yml
install-skills:
  script:
    - skpm install --concurrency 8
  cache:
    key: skpm-$CI_COMMIT_REF_SLUG
    paths:
      - ~/.cache/skpm/
```

```yaml
# .github/workflows/skills.yml
- name: Install agent skills
  run: skpm install
  env:
    SKPM_REGISTRY_TOKEN: ${{ secrets.SKILLS_REGISTRY_TOKEN }}
```

For JSON output in CI (e.g. to feed into jq):

```bash
skpm install --output json | jq '.data.installed[]'
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
make test                # go test ./tests/unit/...
make test-integration    # go test -tags integration ./tests/integration/...
make test-all            # go test ./tests/...
make coverage            # coverage for ./internal/... exercised by ./tests/...
make build               # → dist/skpm
make lint                # requires golangci-lint
make release-snapshot    # requires goreleaser
```

**Project layout:**

```text
cmd/skpm/              # Entrypoint
internal/
  cli/                # Cobra commands
  config/             # Config loading + env override
  manifest/           # agent-skills.yaml read/write
  lockfile/           # agent-skills.lock read/write
  skill/              # Types, validator, packager
  registry/           # Registry backends + factory (download)
  publisher/          # Publisher backends + factory (upload)
  cache/              # SHA256-keyed disk cache
  installer/          # Download + atomic install + platform paths
  progress/           # CI-aware progress bars
testdata/             # Fixture skills for tests
tests/unit/           # Unit test suite
tests/integration/    # Integration test suite
examples/             # Ready-to-use skill examples and CI templates
```

---

## Security

Skills are supply-chain artifacts. `skpm` enforces:

- **SHA256 verification** on every download before writing to disk
- **Zip-slip protection** — path traversal in ZIP entries is rejected
- **Atomic writes** — installations are all-or-nothing; no partial state on failure
- **Cache integrity** — cache keys are the artifact's SHA256; collisions are impossible

For production use, sign your release artifacts and include the signature in your registry metadata. Never ship skills with `sha256: ""` in production lockfiles.

---

## License

MIT

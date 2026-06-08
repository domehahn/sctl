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
- **Recoverable** — deleted lockfile? `skpm lock` regenerates it from `agent-skills.yaml`

---

## Install

```bash
go install github.com/domehahn/skpm/v2/cmd/skpm@latest
```

Or download a pre-built binary from the [releases page](https://github.com/domehahn/skpm/releases).

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
skcr scaffold skill my-skill        # preferred scaffold owner
skpm validate my-skill              # default: useful during development
skpm validate my-skill --strict     # strict: suitable for CI
skpm validate my-skill --publish    # publish: release readiness check
skpm lint my-skill                  # alias for validate --strict
skpm format my-skill --check        # show what would be normalised
skpm format my-skill --write        # apply normalisation
skpm version bump patch my-skill
skpm package my-skill
skpm publish my-skill --source myregistry
```

`skpm init skill <name>` is still available as a compatibility wrapper, but
new workflows should use `skcr scaffold skill <name>` and then use `skpm` for
validation, versioning, packaging, publishing, installation, and updates.

---

## Two-File Model

`skpm` uses two files, both committed to Git — similar to Poetry:

| File | Analogy | Purpose |
| --- | --- | --- |
| `agent-skills.yaml` | `pyproject.toml` | Manifest: which skills you want (human-edited) |
| `agent-skills.lock` | `poetry.lock` | Lockfile: exact versions + SHA256 (generated) |

Lifecycle semantics:

- `skpm add` primarily updates `agent-skills.yaml`.
- `skpm lock` resolves `agent-skills.yaml` into deterministic `agent-skills.lock`.
- `skpm install` installs exactly what is pinned in `agent-skills.lock`.
- `skpm update` refreshes locked versions without installing unless `--install` is passed.

For compatibility, `add` still locks and installs by default. Use `--no-lock` or
`--no-install` when you want explicit package-manager style steps.

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

Command overview:

```text
skpm init
skpm init skill <name>        # compatibility wrapper; prefer skcr scaffold skill
skpm config
skpm add <skill>@<constraint>
skpm remove <skill>
skpm lock
skpm install
skpm update [skill]
skpm outdated
skpm list
skpm search <query>
skpm info <skill>
skpm validate [path]
skpm validate [path] --strict
skpm validate [path] --publish
skpm lint [path]
skpm format [path] --check
skpm format [path] --write
skpm package [path]
skpm publish [path]
skpm verify
skpm doctor
skpm cache list
skpm cache clean
skpm version
skpm version show <path>
skpm version bump patch <path>
skpm version bump minor <path>
skpm version bump major <path>
skpm version set <version> <path>
skpm registry list
skpm registry show <name>
skpm registry add <name> --type <type>
skpm registry test <name>
skpm registry capabilities <name>
```

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

Compatibility note: this command remains for existing users. For new workflows,
prefer `skcr scaffold skill <name>`. Use `skpm` for validation, versioning,
packaging, publishing, and installation.

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
- Every registry has a known `type` (`skillforge`, `github`, `gitlab`, `artifactory`, `local`, `generic-http`)
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

Installs skills from `agent-skills.lock`. If the lockfile is missing and
`agent-skills.yaml` exists, `install` can resolve and generate the lockfile for
backwards compatibility.

```bash
skpm install
skpm install --frozen-lockfile
skpm install --check
skpm install --prune
skpm install --platform codex
skpm install --target ./sandbox
```

- Downloads artifacts in parallel (default 4 concurrent)
- Verifies SHA256 before writing anything to disk
- Uses a local cache (`~/.cache/skpm/`) — re-runs are instant
- Atomic installs: staging directory → rename, never partial state
- `--dry-run` prints what would be installed without writing files
- `--frozen-lockfile` fails if `agent-skills.yaml` and `agent-skills.lock` differ
- `--check` does not write files and fails if installation is incomplete
- `--prune` removes installed skills no longer present in the lockfile

**Platform paths installed per `compatible_with`:**

| Platform | Path |
| --- | --- |
| `claude-code` | `.claude/skills/<name>/` |
| `gitlab-duo` | `skills/<name>/` and `.agents/skills/<name>/` |
| `github-copilot` | `.github/skills/<name>/` |
| `codex` | `.agents/skills/<name>/` |

---

### `skpm add <skill[@version]>`

Adds a skill declaration and, by default for compatibility, resolves, locks, and
installs it.

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
| `--install` / `--no-install` | Control whether files are installed after adding |
| `--lock` / `--no-lock` | Control whether `agent-skills.lock` is updated |

- Version defaults to latest if omitted
- For local paths: reads `skill.yaml` directly, copies the directory atomically
- For `--ref`: downloads the archive at that ref, validates, then installs
- Creates or updates both `agent-skills.yaml` and `agent-skills.lock`

---

### `skpm remove <skill>`

Removes a skill from `agent-skills.yaml`, refreshes `agent-skills.lock`, and can
optionally remove installed files.

```bash
skpm remove documentation-reviewer
skpm remove documentation-reviewer --prune
skpm remove documentation-reviewer --dry-run
```

---

### `skpm lock`

Resolves `agent-skills.yaml` into `agent-skills.lock`.

```bash
skpm lock
skpm lock --check
skpm lock --update
```

- `--check` fails if the existing lockfile is outdated.
- `--update` refreshes resolved versions.
- The lockfile is sorted deterministically and includes resolved source, URL,
  SHA256, compatible platforms, install paths, and generation time when known.

---

### `skpm list`

Lists locked skills and their install metadata.

```bash
skpm list
skpm list --output json
```

Shows name, version, source registry, compatible platforms, and installation
paths when available.

---

### `skpm update [skill]`

Updates one skill or all skills in the lockfile. It does not install by default.

```bash
skpm update
skpm update documentation-reviewer
skpm update documentation-reviewer --latest
skpm update --install
```

- Respects manifest constraints by default.
- `--latest` moves the manifest constraint to the latest version exposed by the
  registry discovery API.
- `--install` installs after writing the updated lockfile.

---

### `skpm outdated`

Compares locked versions with versions exposed by registries.

```bash
skpm outdated
skpm outdated --output json
```

Shows current version, latest compatible version, latest overall version, and the
manifest constraint when the registry supports version discovery.

---

### Registry Discovery

Search and inspect registry contents through the registry abstraction.

```bash
skpm search policy --source myregistry
skpm info gitlab-policy-reviewer --source myregistry
skpm info gitlab-policy-reviewer --versions --source myregistry
```

Discovery support depends on the registry backend. The local registry supports
search, info, and version listing.

---

### `skpm validate [path]`

Validates a canonical skill source directory. Default path is the current directory.

```bash
skpm validate ./skills/gitlab-policy-reviewer
skpm validate ./skills/gitlab-policy-reviewer --strict
skpm validate ./skills/gitlab-policy-reviewer --publish
skpm validate ./skills/gitlab-policy-reviewer --platform codex
skpm validate ./skills/gitlab-policy-reviewer --output json
```

Three validation profiles:

| Profile | Flag | Purpose |
| --- | --- | --- |
| default | _(none)_ | Local development. Optional files produce warnings. |
| strict | `--strict` | CI-grade. Warnings become errors; stricter structural checks. |
| publish | `--publish` | Release-grade. All strict checks plus publish-readiness requirements. |

**Default** checks:

- `SKILL.md` exists and is non-empty
- `VERSION` contains a stable SemVer string (`MAJOR.MINOR.PATCH`)
- `skill.yaml` parses and has required fields (`name`, `version`, `description`, `compatible_with`)
- `skill.yaml.version` matches `VERSION`
- `compatible_with` values are known canonical platform names
- warns on missing `CHANGELOG.md`, `README.md`, `LICENSE`, `tests/`
- errors on forbidden files (`.env`, `id_rsa`, `id_ed25519`, `*.pem`, `*.key`)
- errors on detected secrets and private key blocks
- warns on absolute local paths (`/Users/…`, `/home/…`, `C:\…`)

**Strict** adds:

- missing `CHANGELOG.md` and `README.md` are errors, not warnings
- missing `LICENSE` and `tests/` are errors (unless `--allow-missing-license` / `--allow-missing-tests`)
- absolute local paths are errors
- duplicate platforms and duplicate tags are errors
- unknown `skill.yaml` fields outside `metadata:` are errors
- generated artifacts (`manifest.json`, `checksums.txt`, `*.zip`, `*.tgz`) must not be checked in
- build/cache directories (`node_modules`, `.venv`, `dist`, `target`, `.cache`) must not be present

**Publish** adds:

- all strict checks
- no warnings allowed — any remaining warning becomes an error
- entrypoint file must exist (defaults to `SKILL.md`)
- `skill.yaml.name` must follow naming rules (lowercase alphanumeric and hyphens)

Additional flags:

| Flag | Description |
| --- | --- |
| `--platform <name>` | Fail if the skill does not declare compatibility with this platform |
| `--allow-prerelease` | Accept prerelease SemVer (`1.2.3-beta.1`) |
| `--allow-missing-tests` | Downgrade missing `tests/` from error to info in strict/publish |
| `--allow-missing-license` | Downgrade missing `LICENSE` from error to info in strict/publish |

Canonical platform names: `claude-code`, `gitlab-duo`, `github-copilot`, `codex`, `cursor`, `windsurf`, `openhands`, `opencode`, `ollama`, `generic`, `all`

Platform aliases are accepted in `compatible_with` and normalised by `skpm format`:
`gitlab` → `gitlab-duo`, `github` → `github-copilot`, `claude` → `claude-code`

Exits `0` if valid, `1` if errors are found.

---

### `skpm lint [path]`

Convenience alias for `skpm validate --strict`. Suitable for CI pipelines.

```bash
skpm lint ./skills/gitlab-policy-reviewer
skpm lint ./skills/gitlab-policy-reviewer --output json
skpm lint ./skills/gitlab-policy-reviewer --platform codex
```

`skpm lint` delegates to the same validator as `skpm validate --strict`. There is
no separate lint engine or separate lint configuration. Validation belongs to
`skpm` because `skpm` owns the skill lifecycle; a separate `sklint` CLI would
duplicate lifecycle rules and create drift.

> `skcr` creates skill skeletons. `skpm` validates whether those skeletons are
> lifecycle-ready.

---

### `skpm format [path]`

Normalises skill metadata files without changing skill semantics.

```bash
skpm format ./skills/gitlab-policy-reviewer --check
skpm format ./skills/gitlab-policy-reviewer --write
```

What `format` normalises:

- Platform aliases in `compatible_with` → canonical names (`gitlab` → `gitlab-duo`)
- Duplicate platforms in `compatible_with` removed
- Duplicate tags in `tags` removed
- Missing `entrypoint: SKILL.md` added to `skill.yaml`
- Leading `v` removed from `VERSION` (`v1.2.3` → `1.2.3`)
- Trailing newline ensured in `SKILL.md`, `CHANGELOG.md`, `README.md`, `LICENSE`

What `format` does **not** do:

- Bump or change versions
- Rewrite `SKILL.md` content
- Change changelog entries
- Package or publish anything

`--check` exits `1` if any changes would be made (useful in CI).
`--write` applies changes in place.

> Note: re-marshaling `skill.yaml` removes inline comments.

---

### `skpm package [path]`

Packages a skill directory into a local ZIP artifact. Default path is the current directory.

```bash
skpm package ./skills/gitlab-policy-reviewer
skpm package ./skills/gitlab-policy-reviewer --output-dir ./dist
```

- Runs `validate` first — fails if the skill is invalid
- Creates `<name>-<version>.zip` with all skill files plus `manifest.json` and `checksums.txt`
- Prints the SHA256 of the ZIP
- Excludes `.git/`, `*.tmp`, `*.part`
- Uses deterministic file ordering and normalized timestamps for reproducible archives

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
skpm version show ./skills/my-skill
skpm version bump patch ./skills/my-skill
skpm version bump minor ./skills/my-skill
skpm version bump major ./skills/my-skill
skpm version set 1.2.3 ./skills/my-skill
```

`skpm version` without a subcommand prints the `skpm` binary version.

Skill version subcommands read `VERSION` and `skill.yaml`, require them to
match, store stable SemVer without a leading `v`, and ensure `CHANGELOG.md`
contains an entry for the resulting version. `version set` accepts `v1.2.3` but
stores `1.2.3`.

---

### `skpm verify`

Verifies lockfile and installed skill state.

```bash
skpm verify
skpm verify --frozen-lockfile
skpm verify --platform codex
skpm verify --output json
```

Checks lockfile consistency, required installed files, metadata matching,
platform compatibility, and available checksum metadata.

---

### `skpm doctor`

Checks project and configuration health.

```bash
skpm doctor
skpm doctor --output json
```

---

### `skpm cache`

Manages cached skill artifacts.

```bash
skpm cache list
skpm cache clean
skpm cache clean --dry-run
```

---

### `skpm registry`

Manages registry configuration and capability detection.

```bash
skpm registry list
skpm registry show company-skillforge
skpm registry add company-skillforge --type skillforge --url https://skills.company.com --token-env SKILLFORGE_TOKEN
skpm registry add generic --type generic-http --url https://skills.example.com --endpoint-resolve '/api/v1/skills/{namespace}/{name}/resolve?constraint={constraint}'
skpm registry remove old-registry
skpm registry test company-skillforge
skpm registry capabilities company-skillforge --output json
skpm registry login company-skillforge
```

`skpm` is registry-agnostic: SkillForge is one supported backend, not the core of
the package manager. Registries advertise capabilities, and commands degrade
gracefully when a backend does not support search, publish, governance, or other
optional operations.

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
skpm lock
skpm install
```

### Author and release a new skill

```bash
skcr scaffold skill my-skill
# edit skills/my-skill/SKILL.md
skpm validate skills/my-skill
skpm version bump patch skills/my-skill
skpm package skills/my-skill
skpm publish  skills/my-skill --source myregistry
# → validates, packages, creates tag, pushes, uploads
```

Compatibility fallback:

```bash
skpm init skill my-skill --output-dir ./skills
```

---

## Skill Structure

Canonical skill layout:

```text
my-skill/
  SKILL.md        # Agent-readable capability definition (required)
  VERSION         # SemVer string, e.g. "1.5.0" (required)
  skill.yaml      # Machine-readable metadata (required)
  CHANGELOG.md    # Version history (required for publish)
  README.md       # Human documentation (recommended)
  LICENSE         # License file (recommended)
  tests/
    README.md     # Test documentation (recommended)
```

**`skill.yaml`:**

```yaml
name: gitlab-policy-reviewer
version: 1.5.0
description: Reviews GitLab security policy YAML and approval policies.
namespace: platform-security
owners:
  - platform-security
license: MIT
entrypoint: SKILL.md
tags:
  - security
  - gitlab
compatible_with:
  - claude-code
  - gitlab-duo
  - github-copilot
  - codex
security:
  requires_network: false
  requires_secrets: false
  writes_files: false
  runs_commands: false
```

Use `all` in `compatible_with` to install to every supported platform path.

**Supported platforms:**

| Platform | Canonical name |
| --- | --- |
| Claude Code | `claude-code` |
| GitLab Duo | `gitlab-duo` |
| GitHub Copilot | `github-copilot` |
| Codex | `codex` |
| Cursor | `cursor` |
| Windsurf | `windsurf` |
| OpenHands | `openhands` |
| OpenCode | `opencode` |
| Ollama | `ollama` |
| Generic | `generic` |
| All platforms | `all` |

---

## Configuration

`skpm` reads `~/.config/skpm/config.yaml` (or `$XDG_CONFIG_HOME/skpm/config.yaml`).

Create it interactively with `skpm init`, or write it manually:

```yaml
default_registry: company-skillforge
cache_dir: ~/.cache/skpm

registries:
  company-skillforge:
    type: skillforge
    url: https://skills.company.com
    auth:
      type: bearer
      token_env: SKILLFORGE_TOKEN

  company-artifactory:
    type: artifactory
    url: https://artifactory.company.com/artifactory
    repo: agent-skills
    auth:
      type: bearer
      token_env: ARTIFACTORY_TOKEN

  company-gitlab:
    type: gitlab
    url: https://gitlab.company.com
    project: platform/agent-skills   # required for gitlab
    auth:
      type: bearer
      token_env: GITLAB_TOKEN

  company-github:
    type: github
    repo: myorg/agent-skills
    auth:
      type: bearer
      token_env: GITHUB_TOKEN

  local-dev:
    type: local
    path: ./dist/registry

  generic:
    type: generic-http
    url: https://skills.example.com
    auth:
      type: bearer
      token_env: SKILLS_TOKEN
    endpoints:
      capabilities: /api/v1/capabilities
      search: /api/v1/skills?q={query}
      info: /api/v1/skills/{namespace}/{name}
      versions: /api/v1/skills/{namespace}/{name}/versions
      resolve: /api/v1/skills/{namespace}/{name}/resolve?constraint={constraint}
      download: /api/v1/skills/{namespace}/{name}/versions/{version}/download
      publish: /api/v1/skills/{namespace}/{name}/versions/{version}
```

**Supported registry types:**

| Type | Key fields |
| --- | --- |
| `skillforge` | `url`, optional `auth`, SkillForge-compatible endpoints |
| `generic-http` | `url`, `endpoints`, optional `auth` and `headers` |
| `github` | `repo: owner/repo` |
| `gitlab` | `url`, `project: namespace/project` |
| `artifactory` | `url`, `repo` |
| `local` | `path` |

Legacy fields (`url`, `token`, `project`) are still accepted. New configs should
prefer `auth.token_env` over storing secrets in config files.

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

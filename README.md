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
# 1 — scaffold
skpm create my-skill                # interactive scaffold with SKILL.md, skill.yaml, VERSION
# or: skcr scaffold skill my-skill  (preferred when using the skcr bake workflow)

# 2 — develop with live feedback
skpm watch test --dir my-skill      # re-runs 'test' script on every save

# 3 — validate and format
skpm validate my-skill              # default: useful during development
skpm validate my-skill --publish    # release readiness check
skpm lint my-skill                  # alias for validate --strict
skpm format my-skill --write        # normalise SKILL.md and skill.yaml

# 4 — release in one step
skpm release my-skill --bump minor --message "add streaming support" --source myregistry
# equivalent to: version bump + changelog add + package + git tag + push + publish
```

`skpm init skill <name>` is still available as a compatibility wrapper, but
new workflows should use `skpm create <name>` or `skcr scaffold skill <name>`.

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
# Project setup
skpm init
skpm config get <key>
skpm config set <key> <value>
skpm env
skpm migrate                    # upgrade files to current spec version

# Package management
skpm add <skill>[@constraint]
skpm remove <skill>
skpm lock
skpm fetch                      # warm cache without installing
skpm install
skpm update [skill]
skpm outdated
skpm pin
skpm diff <skill>@<v1> <skill>@<v2>
skpm prune                      # remove installed dirs not in lockfile

# Discovery
skpm list
skpm search <query>
skpm info <skill>
skpm why <skill>
skpm graph                      # installation + dependency graph

# Auth
skpm login [registry] --token <token>
skpm logout [registry]

# Development workflow
skpm create <name>              # scaffold new skill
skpm template list|add|remove|show|use
skpm link [path]
skpm unlink <name>
skpm clone <skill>[@version]
skpm import <zipfile>
skpm snapshot save [name]
skpm snapshot restore <name>
skpm snapshot list
skpm snapshot delete <name>

# Authoring loop
skpm run <script>               # run a script from skill.yaml
skpm watch <script>             # re-run on file change
skpm validate [path]
skpm lint [path]
skpm format [path] --check|--write
skpm changelog show [path]
skpm changelog add <message> [path]
skpm version show <path>
skpm version bump patch|minor|major <path>
skpm version set <version> <path>
skpm package [path]
skpm publish [path]
skpm release [path]             # bump + changelog + package + publish in one step
skpm verify
skpm deprecate <skill>@<version> --reason <msg>
skpm yank     <skill>@<version> --reason <msg>
skpm unyank   <skill>@<version>

# Monorepo
skpm workspace init|list|run|validate|publish|graph

# Lifecycle hooks (defined in agent-skills.yaml)
skpm hooks list
skpm hooks run <hook>

# Ops & diagnostics
skpm audit
skpm integrity
skpm stats
skpm doctor
skpm cache list
skpm cache clean
skpm registry list|show|add|remove|test|capabilities
skpm completion bash|zsh|fish|powershell
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

### `skpm config get` / `skpm config set`

Read or write top-level config values without opening the YAML file directly.

```bash
skpm config get default_registry
skpm config get cache_dir
skpm config set default_registry my-registry
skpm config set cache_dir /tmp/skpm-cache
skpm config set concurrency 8
skpm config set log_level debug
```

Valid keys: `default_registry`, `cache_dir`, `log_level`, `concurrency`.

---

### `skpm env`

Prints the resolved runtime environment — useful when debugging auth or path issues.

```bash
skpm env
skpm env --output json
```

Shows config file path, cache directory, default registry, all configured registries with auth type and token status (masked), manifest and lockfile presence, and active development links.

---

### `skpm login` / `skpm logout`

Save or remove auth tokens for registries in the config file.

```bash
skpm login my-registry --token ghp_xxxx   # store token
skpm login my-registry                     # prompt interactively
skpm logout my-registry                    # remove token
skpm login                                 # uses default_registry
```

Prefer `auth.token_env` in config for CI — tokens in config files require care around secrets management.

---

### `skpm create`

Scaffold a new publishable skill directory with all required files.

```bash
skpm create my-skill
skpm create my-skill --description "Reviews YAML configs" --platform claude-code --license Apache-2.0
skpm create my-skill --no-interactive     # use flags only, no prompts
```

Creates `SKILL.md`, `skill.yaml`, `VERSION` (0.1.0), and `CHANGELOG.md`. Run `skpm publish` when ready.

---

### `skpm link` / `skpm unlink`

Link a local skill directory into the project for in-place development without publishing.

```bash
skpm link ./my-skill          # symlinks into all platform dirs
skpm link ./my-skill --platform claude-code
skpm unlink my-skill          # removes symlinks, restores from registry on next install
```

Linked skills are tracked in `.skpm-links.yaml`. Edits to the source directory are reflected immediately without re-publishing.

---

### `skpm clone`

Download and extract a skill from the registry into a local directory for forking or inspection.

```bash
skpm clone my-skill               # latest version → ./my-skill/
skpm clone my-skill@1.2.0
skpm clone my-skill --dir my-fork
```

---

### `skpm import`

Install a skill from a local ZIP artifact without contacting a registry. Useful for air-gapped environments and CI artifact promotion.

```bash
skpm import dist/my-skill-1.2.0.zip
skpm import artifact.zip --platform claude-code
skpm import artifact.zip --no-lock     # skip updating agent-skills.lock
```

---

### `skpm fetch`

Download skill ZIPs into the local cache without installing. Designed for CI layer separation.

```bash
skpm fetch                        # all locked skills
skpm fetch my-skill other-skill   # specific skills only
```

**Typical CI pattern:**

```yaml
cache-skills:          # runs once, result cached
  script: skpm fetch
  cache:
    paths: [~/.cache/skpm/]

build:                 # runs every time, no network needed
  script: skpm install --frozen-lockfile
```

---

### `skpm diff`

Show a unified diff between two versions of a skill's SKILL.md.

```bash
skpm diff my-skill@1.0.0 my-skill@1.2.0
skpm diff my-skill@1.0.0 my-skill@1.2.0 --source my-registry
```

---

### `skpm pin`

Replace version constraints in `agent-skills.yaml` with the exact versions currently in the lockfile.

```bash
skpm pin
skpm pin --dry-run
```

---

### `skpm audit`

Check locked skills against the registry for yanked/deprecated versions and available updates.

```bash
skpm audit
skpm audit --output json
```

| Icon | Meaning |
| --- | --- |
| `✗` | Version was yanked — stop using it |
| `~` | Version is deprecated, or a major update is available |
| `↑` | Minor or patch update available |
| `✓` | Up to date |

Exits with code 1 if any error-level findings exist — CI-friendly.

---

### `skpm integrity`

Verify installed SKILL.md files match the cached ZIP artifacts.

```bash
skpm integrity
skpm integrity --output json
```

| Status | Meaning |
| --- | --- |
| `ok` | Installed content matches the cached artifact |
| `modified` | Content differs — possible manual edit or tampering |
| `missing` | SKILL.md absent — run `skpm install` |
| `unverifiable` | No cached artifact — run `skpm install` to cache |

Exits with code 1 on `modified` or `missing`.

---

### `skpm prune`

Remove installed skill directories that are not recorded in the lockfile.

```bash
skpm prune
skpm prune --dry-run
skpm prune --root custom/skills
```

---

### `skpm snapshot`

Save and restore lockfile snapshots for rollback.

```bash
skpm snapshot save              # timestamp name
skpm snapshot save before-bump
skpm snapshot restore before-bump
skpm snapshot list
skpm snapshot delete before-bump
```

`restore` automatically backs up the current lockfile as `pre-restore-<timestamp>` before overwriting.

---

### `skpm why`

Explain why a skill is in the lockfile.

```bash
skpm why my-skill
skpm why my-skill --output json
```

Shows manifest constraint, resolved version, registry source, SHA256, installation paths (with ✓/✗ per path), and whether the skill is linked locally.

---

### `skpm changelog`

Manage a skill's `CHANGELOG.md`.

```bash
skpm changelog show ./my-skill           # last 3 entries
skpm changelog show ./my-skill -n 1      # last entry only
skpm changelog add "fix null pointer" ./my-skill
skpm changelog add "fix null pointer" ./my-skill --version 1.2.1
```

`add` appends a bullet to the section for the current version (from `VERSION`) without bumping it.

---

### `skpm deprecate` / `skpm yank` / `skpm unyank`

Manage the lifecycle of published skill versions in the registry.

```bash
skpm deprecate my-skill@1.0.0 --reason "use 2.x instead"
skpm yank      my-skill@1.0.0 --reason "critical bug"
skpm unyank    my-skill@1.0.0
```

Requires a registry that supports governance operations (`skpm registry capabilities <name>` shows `deprecate: true` / `yank: true`).

---

### `skpm stats`

Show cache and installation statistics.

```bash
skpm stats
skpm stats --output json
```

Reports cache artifact count and size, lockfile skill count, and per-platform-directory skill count and disk usage.

---

### `skpm run`

Execute a named script from the `scripts` section of `skill.yaml`.

```bash
skpm run test
skpm run lint --dir ./my-skill
skpm run build -- --verbose
```

Scripts are defined in `skill.yaml`:

```yaml
scripts:
  test: skpm validate . && skpm lint .
  build: skpm package .
  lint: skpm lint .
```

Shell completion works for script names. `SKPM_SKILL_DIR` and `SKPM_SCRIPT` are set in the script's environment.

---

### `skpm watch`

Watch skill files for changes and re-run a script automatically — tight feedback loop during authoring.

```bash
skpm watch test
skpm watch lint --dir ./my-skill
skpm watch validate --debounce 500ms
skpm watch test --ext .md,.yaml
```

Runs the script once immediately on start, then again after each save. Press Ctrl+C to stop.

---

### `skpm release`

Combined authoring pipeline: validate → bump version → update changelog → package → git tag → push → upload.

```bash
skpm release                            # patch bump, default registry
skpm release --bump minor --message "add streaming support"
skpm release --no-publish               # stop after packaging, skip git and registry
skpm release --dry-run                  # preview all steps
```

Replaces the manual sequence of `skpm version bump patch` + `skpm changelog add` + `skpm publish`.

---

### `skpm hooks`

Lifecycle hooks defined in `agent-skills.yaml` run automatically before and after key operations.

```yaml
# agent-skills.yaml
hooks:
  pre_install:  echo "installing..."
  post_install: ./scripts/verify-checksums.sh
  post_add:     git add agent-skills.yaml agent-skills.lock
  pre_publish:  skpm validate .
```

Available hooks: `pre_add`, `post_add`, `pre_install`, `post_install`, `pre_update`, `post_update`, `pre_publish`, `post_publish`, `pre_release`, `post_release`.

Pre-hooks that fail abort the operation. Post-hooks that fail print a warning and continue.

```bash
skpm hooks list          # show configured hooks
skpm hooks run post_add  # trigger a hook manually
```

---

### `skpm migrate`

Upgrade project or skill files to the current spec version. Idempotent — safe to run on already-current files.

```bash
skpm migrate                         # upgrade agent-skills.yaml + agent-skills.lock
skpm migrate --skill-dir ./my-skill  # upgrade skill.yaml + SKILL.md frontmatter
skpm migrate --dry-run               # preview changes without writing
```

Migrations applied:
- `agent-skills.yaml` / `agent-skills.lock` v0 → v1: sets `version: 1`
- `skill.yaml`: adds missing `namespace: default`; renames `platforms` → `compatible_with`
- `SKILL.md` frontmatter: adds `compatibility.spec_version: 1`

---

### `skpm template`

Save skill directories as reusable scaffolding templates; apply them with `skpm template use`.

```bash
skpm template list                          # list saved templates
skpm template add review ./my-review-skill  # snapshot skill dir as template
skpm template show review                   # inspect template files
skpm template use review my-new-review      # scaffold new skill from template
skpm template remove review                 # delete template
```

Skill names in file content are replaced with `{{.Name}}` when saving; `use` substitutes the new skill name. Templates are stored in `~/.config/skpm/templates.yaml`.

---

### `skpm workspace`

Manage a monorepo of multiple skills via `skpm-workspace.yaml` at the repository root. The workspace file is discovered by walking up from the current directory.

```yaml
# skpm-workspace.yaml
version: 1
skills:
  - ./skill-a
  - ./skill-b
  - ./shared/review-skill
```

```bash
skpm workspace init                        # create workspace file (auto-discovers skills)
skpm workspace list                        # show all skills with versions and validity
skpm workspace run test                    # run 'test' script in all skills
skpm workspace run lint --skill ./skill-a  # run in selected skill only
skpm workspace validate                    # validate all skills
skpm workspace publish                     # publish all skills
skpm workspace publish --changed           # publish only git-modified skills
skpm workspace graph                       # dependency tree across workspace
skpm workspace graph --format dot          # Graphviz DOT output
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
skpm lock
skpm install
```

### Author and release a new skill

```bash
# scaffold
skpm create my-skill
cd my-skill

# develop with live feedback
skpm watch test         # re-runs 'test' script from skill.yaml on every save

# release (validate + bump + changelog + package + tag + push + publish)
skpm release --bump patch --message "fix null pointer on empty input" --source myregistry
```

Or step by step for more control:

```bash
skpm validate .
skpm version bump patch .
skpm changelog add "fix null pointer on empty input" .
skpm package .
skpm publish . --source myregistry
```

### Work on multiple skills in a monorepo

```bash
# one-time setup
skpm workspace init          # discovers skills in subdirectories

# daily workflow
skpm workspace validate      # validate all skills
skpm workspace run test      # run 'test' script in every skill
skpm workspace run lint --skill ./my-skill  # run in one skill only

# release only what changed
skpm workspace publish --changed --source myregistry
```

### Local development with link

```bash
skpm link ./my-skill         # symlinks into platform dirs; edits take effect immediately
skpm why my-skill            # confirms linked status
skpm unlink my-skill         # restores registry version on next install
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
requires:             # optional: other skills this skill depends on
  - base-reviewer
scripts:              # optional: runnable via skpm run / skpm watch
  test: skpm validate . && skpm lint .
  build: skpm package .
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

### Basic install

```yaml
# .gitlab-ci.yml
install-skills:
  script:
    - skpm install --frozen-lockfile --concurrency 8
  cache:
    key: skpm-$CI_COMMIT_REF_SLUG
    paths:
      - ~/.cache/skpm/
```

```yaml
# .github/workflows/skills.yml
- name: Install agent skills
  run: skpm install --frozen-lockfile
  env:
    SKPM_REGISTRY_TOKEN: ${{ secrets.SKILLS_REGISTRY_TOKEN }}
```

### Cache / install layer separation

Separate the network-bound download step from the install step to maximise cache hits:

```yaml
# .gitlab-ci.yml
cache-skills:          # runs once; result cached across pipelines
  stage: prepare
  script: skpm fetch
  cache:
    key: skpm-lock-$CI_COMMIT_SHA
    paths: [~/.cache/skpm/]

build:                 # no network needed; always fast
  stage: build
  script: skpm install --frozen-lockfile
  cache:
    key: skpm-lock-$CI_COMMIT_SHA
    paths: [~/.cache/skpm/]
    policy: pull
```

### Post-install hook for verification

```yaml
# agent-skills.yaml
hooks:
  post_install: skpm audit && skpm integrity
```

### Audit in CI

```bash
skpm audit --output json | jq '.data.findings[] | select(.severity == "error")'
```

`skpm audit` exits with code 1 if any error-level findings exist.

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

# Examples

Concrete, copy-paste-ready examples for common `skm` workflows.

## Directory Layout

```
examples/
  skills/
    gitlab-policy-reviewer/     # Production-grade security policy review skill
    documentation-reviewer/     # Documentation quality review skill
    robotframework-reviewer/    # Robot Framework test suite review skill
  project/
    agent-skills.lock           # Example lockfile for a project consuming skills
    catalog.yaml                # Skill catalog index for a skills monorepo
  config/
    config.yaml                 # Annotated skm config with all registry types
  ci/
    .gitlab-ci.yml              # GitLab CI: install, validate, package, publish
    github-actions.yml          # GitHub Actions: install, validate, release
```

---

## Workflows

### 1 — Consuming skills in a project

Copy `examples/project/agent-skills.lock` to your project root, adjust the
skill names, versions, SHA256s, and source URLs, then:

```bash
skm install
```

Commit `agent-skills.lock`. Never commit the installed `skills/` or `.agents/skills/` directories.

### 2 — Creating a new skill

Start from one of the skill examples as a template:

```bash
cp -r examples/skills/gitlab-policy-reviewer my-new-skill
cd my-new-skill

# Edit SKILL.md, skill.yaml, VERSION, CHANGELOG.md
skm validate .
```

### 3 — Packaging a skill for distribution

```bash
skm validate ./my-new-skill
skm package ./my-new-skill --output-dir ./dist
# → dist/my-new-skill-1.0.0.zip  (SHA256 printed)
```

Upload the ZIP to your registry and add the SHA256 to `agent-skills.lock` or `catalog.yaml`.

### 4 — Publishing via CI

- **GitLab:** push a tag in the format `<skill-name>/v<version>` (e.g. `gitlab-policy-reviewer/v1.5.0`) — the `.gitlab-ci.yml` pipeline validates, packages, and uploads to Artifactory automatically.
- **GitHub:** same tag format — the `github-actions.yml` workflow validates, packages, and creates a GitHub Release with the ZIP attached and an `agent-skills.lock` snippet in the release body.

---

## Skill Compatibility Matrix

| Skill | Claude Code | GitLab Duo | GitHub Copilot | Codex |
|---|:---:|:---:|:---:|:---:|
| gitlab-policy-reviewer | ✓ | ✓ | ✓ | ✓ |
| documentation-reviewer | ✓ | ✓ | ✓ | ✓ |
| robotframework-reviewer | ✓ | ✓ | — | — |

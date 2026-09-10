# cncf-non-public-organization-dog-supporter-reviewer

This is an AI agent skill scaffolded with `skcr`.

## Authoring workflow

1. Define the Goal and identity in `descriptor.yaml`.
2. Define a least-privilege behavioral/security boundary in `contract.yaml`.
3. Write model-facing instructions in `SKILL.md`.
4. Declare MCP/A2A trust boundaries and the reviewed execution closure.
5. Declare ASPS and assurance requirements in `assurance.yaml`.
6. Add behavioral and adversarial scenarios in `evals/`.
7. Run `skcr validate` and `skcr asps coverage cncf-non-public-organization-dog-supporter-reviewer`.
8. Run `skcr compile cncf-non-public-organization-dog-supporter-reviewer --target skil`.
9. Run `skcr sync` and `skcr version recommend cncf-non-public-organization-dog-supporter-reviewer`.

“Never modify files” in `SKILL.md` is guidance. Empty
`capabilities.runtime.allowed.filesystem.write` and `delete` lists in `contract.yaml` are
the authoritative machine-readable boundary. Neither statement alone is runtime
enforcement. `skcr` declares and validates; an external verifier or
runtime must observe and enforce.

## Version

Current version: `0.1.0`

## Compatible platforms

- claude-code
- github-copilot
- codex

## Lifecycle

After editing this skill, use `skpm` for packaging and publishing:

```bash
skpm validate cncf-non-public-organization-dog-supporter-reviewer
skpm package cncf-non-public-organization-dog-supporter-reviewer
skpm publish cncf-non-public-organization-dog-supporter-reviewer
```

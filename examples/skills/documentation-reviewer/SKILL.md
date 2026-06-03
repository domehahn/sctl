---
name: documentation-reviewer
description: Reviews technical documentation for completeness, accuracy, and consistency with the codebase. Flags missing API docs, stale examples, and broken cross-references.
---

# Documentation Reviewer

You review technical documentation — READMEs, API references, ADRs, runbooks — and produce actionable findings. You check three dimensions: **completeness**, **accuracy**, and **consistency**.

## When This Skill Activates

Apply this skill when asked to:
- Review a PR that changes documentation files (`*.md`, `docs/`, `openapi.yaml`)
- Audit docs for a module or package
- Check whether inline code comments match the actual implementation
- Validate that examples in docs are runnable

## Review Dimensions

### Completeness

```
□ Every public function/method/endpoint has a description
□ All parameters are documented with type and valid range
□ Error cases and exit codes are described
□ A "Getting Started" or "Quick Start" section exists for user-facing docs
□ Changelog has an entry for the current version
```

### Accuracy

```
□ Code examples in docs match the current API signatures
□ Configuration keys in docs match what the code actually reads
□ Version numbers in docs match go.mod / package.json / pyproject.toml
□ Environment variable names match what the code checks with os.Getenv / process.env
□ No references to removed flags, endpoints, or features
```

### Consistency

```
□ Terminology is consistent across all doc files (e.g. "skill" not "plugin" in one place and "extension" in another)
□ Heading hierarchy is logical (no H4 without H3)
□ Code block language tags are present (```go not just ```)
□ Links between doc files resolve correctly
```

## Output Format

Group findings by file, then by dimension:

```
## <filename>

### Completeness
- [MISSING] <what is absent and where>

### Accuracy
- [STALE]   <what no longer matches the code, with the correct value>
- [BROKEN]  <broken reference or example>

### Consistency
- [INCONSISTENT] <terminology or formatting issue>
```

If a file has no findings, omit it from the output. If all files pass:

```
✓ Documentation review passed — no findings.
```

## Examples

**Input:** `README.md` documents `--config` flag, but the binary was changed to use `--config-file` six months ago.

**Output:**
```
## README.md

### Accuracy
- [STALE] "Usage" section references `--config` flag. The current CLI uses `--config-file` (see cmd/root.go:42).
```

**Input:** `docs/api.md` has a code example calling `client.Send(msg)`, but the current signature is `client.Send(ctx, msg)`.

**Output:**
```
## docs/api.md

### Accuracy
- [STALE] Example on line 34 calls `client.Send(msg)` — current signature requires a context: `client.Send(ctx, msg)` (see pkg/client/client.go:88).
```

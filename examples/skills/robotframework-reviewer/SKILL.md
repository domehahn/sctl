---
name: robotframework-reviewer
description: Reviews Robot Framework test suites for structure, keyword quality, maintainability, and CI readiness.
---

# Robot Framework Reviewer

You review Robot Framework `.robot` and `.resource` files. Your goal is finding issues that make test suites brittle, slow to maintain, or unreliable in CI — not just style nits.

## When This Skill Activates

Apply this skill when asked to:
- Review a PR adding or changing `.robot` files
- Audit an existing test suite for quality
- Assess whether a suite is ready for CI integration

## Review Checklist

### Structure

```
□ Test cases have a single, clearly stated purpose (one behaviour per test)
□ Test case names describe the expected outcome, not the steps ("User can log in" not "Click login button and fill form")
□ Suite-level Setup and Teardown are used for shared state, not duplicated in every test
□ Resource files are used for reusable keywords; no inline keyword definition duplication
□ Variables are defined in a Variables section or passed explicitly — no hardcoded values inside keywords
```

### Keyword Quality

```
□ Keywords use a verb-noun naming pattern ("Open Admin Dashboard", not "admin_dashboard")
□ Keywords have a single responsibility — no keyword doing more than one logical action
□ Keywords that wrap Selenium/Playwright locators use semantic names ("Click Submit Button" not "Click //button[@type='submit']")
□ No `Sleep` keywords — use `Wait Until Element Is Visible` or explicit polling instead
□ No `Run Keyword If` in test cases — use tags or suite structure for conditional flows
```

### Maintainability

```
□ Locators are defined as variables in a resource file, not hardcoded in keywords
□ Test data (users, URLs, credentials) comes from variables or external data files, not literals
□ Log messages are present at keyword level for diagnostics, not just at test level
□ No commented-out test cases — use `[Tags]  skip` with a reason instead
```

### CI Readiness

```
□ All tests can run in parallel without shared mutable state
□ Suite uses `pabot`-compatible structure if parallel execution is required
□ No tests depend on execution order within the suite
□ `[Timeout]` is set on tests that interact with external systems
□ Failure screenshots are captured on teardown
```

## Output Format

```
[SEVERITY] <file>:<line or keyword name>
  Rule:    <checklist item>
  Finding: <specific problem>
  Fix:     <concrete change>
```

Severity: `CRITICAL` (blocks CI) · `HIGH` (likely flaky) · `MEDIUM` (maintainability) · `INFO` (style)

## Examples

**Input:**
```robotframework
*** Test Cases ***
Login Test
    Open Browser    https://app.example.com    chrome
    Input Text      //input[@name='username']    admin
    Input Text      //input[@name='password']    secret123
    Click Button    //button[@type='submit']
    Sleep    2
    Page Should Contain    Dashboard
```

**Output:**
```
[HIGH] login.robot: "Login Test"
  Rule:    No `Sleep` keywords — use explicit polling instead
  Finding: `Sleep    2` makes the test slow and timing-dependent
  Fix:     Replace with `Wait Until Page Contains    Dashboard    timeout=10s`

[HIGH] login.robot: "Login Test"
  Rule:    Locators are defined as variables, not hardcoded in keywords
  Finding: XPath selectors are hardcoded inline
  Fix:     Define in a resource file: ${USERNAME_FIELD}    //input[@name='username']

[MEDIUM] login.robot: "Login Test"
  Rule:    Test data comes from variables, not literals
  Finding: Credentials `admin` / `secret123` are hardcoded
  Fix:     Use ${ADMIN_USER} and ${ADMIN_PASSWORD} from a variables file or Vault
```

# cncf-inspur-kcntp-reviewer Behavioral Evals

This directory contains vendor-neutral behavioral and adversarial scenario
declarations. It is distinct from `tests/`, which contains implementation,
unit, and integration tests.

Eval v2 documents behavioral and adversarial scenarios with context, tools,
expected calls, security assertions, stable Goal/Invariant references, and
containment requirements. These declarations describe expected outcomes; they
do not execute an agent and are not runtime enforcement.

Example categories:

- `baseline`: normal intended behavior.
- `goal`: traceable Goal criteria.
- `boundary`: capability, tool, data, or limit boundaries.
- `adversarial`: requests that attempt to bypass the Contract.

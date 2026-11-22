---
name: finishing-tests
description: Use when writing or changing Go unit tests in this repository, before you say the tests are done.
---

# Finishing tests

Tests are done when the suite is green, every known bug is marked, and each test has been seen failing.

## Steps

1. A test that fails because the code breaks the spec: see it fail once, then put `t.Skip("known bug demoN: <one line>")` at its top. Do not fix the code and do not leave the test failing. The suite ends green.
2. For each behaviour you tested, change one line of production code so the behaviour breaks, run the test, and see it fail. Then restore the file and check with `git diff` that nothing is left.
3. End your answer with two tables: behaviour, rule id, test name; and mutation, test that failed.

## Mistakes

| Mistake | Fix |
|---|---|
| Leaving failing tests in the suite | Skip them with the bug id |
| Saying "the tests would catch it" | Run the mutation and show the failure |
| A test with no rule id | Name the rule from the spec in the table |

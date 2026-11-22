# Step 1.1: red

**Branch:** `step-1-red`  
**Previous:** `step-0-start` · **Next:** `git checkout step-1-green`

## Where we are

The agent has written its test and has not touched the code yet. This is the red step of TDD. The test file is exactly what the agent wrote from this prompt:

```text
Users report that a password of exactly 8 characters is refused.
Fix it with test-driven development and the project's unit test skill.
```

## What to show

```bash
go test ./internal/passwordpolicy/        # FAIL: the two 8-byte cases
cat internal/passwordpolicy/policy_test.go
```

## What to point at

- The expected values come from the spec: "Min 8, max 128" in BR-006-16. The test names the rule.
- The test checks five length cases, not only the one the user reported: 7 bytes, 8 bytes, 6 characters that take 8 bytes, 128 and 129.
- Only the two 8-byte cases fail, and they fail with the message the users see.
- There is a second test, skipped with a bug id. The agent found the other planted bug while reading the spec, and marked it without being asked.

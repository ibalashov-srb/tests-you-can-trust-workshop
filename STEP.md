# Step 3.1: baseline, the agent without a skill

**Branch:** `step-3-baseline`  
**Previous:** `step-3-start` · **Next:** `git checkout step-3-skill`

## Where we are

This is what one agent produced from the prompt below, with no project skill. Nothing was edited by hand.

```text
internal/passwordpolicy has no unit tests. Write them.
```

## What to show

```bash
go test ./internal/passwordpolicy/                 # FAIL
git diff --stat step-3-start -- internal/          # one new test file
grep -c "t.Skip" internal/passwordpolicy/policy_test.go   # 0
```

## What is good

- The agent took the expectations from the spec (BR-006-16), not from the code.
- It found both planted bugs.

## What is missing

- **The run is red.** The tests that show a bug are left failing. A red suite hides the next real failure.
- **The bugs are not marked.** No skip, no bug id in the test.
- **No proof of red.** Nothing shows that the passing tests can fail.
- **No table** that links each test to a rule in the spec.

These four gaps are the input for the skill in the next step.

## What the agent answered

I added `internal/passwordpolicy/policy_test.go`, written from the spec (F-006 BR-006-16) rather than from the code. It covers the length limits, which are counted in bytes, the four required character types, what counts as a special character (`_ - + € «` do, a space doesn't), the common-password list (including every entry on it, in mixed case), and the email and name similarity checks with their 4-character minimum. Vet is clean.

**The tests fail: 19 of 144 checks, all caused by two bugs in `internal/passwordpolicy/policy.go`.**

1. **`policy.go:80` — passwords of exactly 8 bytes are refused.** The check is `len(password) <= MinLength`, so an 8-byte password gets "too short". The spec says 8 is the minimum, so it should be `<`. This causes 15 failures:
   - the two "exactly 8 bytes is accepted" tests;
   - 13 common passwords that are exactly 8 bytes long, such as `P@ssw0rd` and `Qwerty1!`. They get "too short" instead of "too common", so they're still refused, just with the wrong message.
2. **`policy.go:137` — 4-letter name words are never checked.** The check is `> SimilarityMinLength`, but the spec says only words shorter than 4 characters are skipped, so it should be `>=`. With a full name like "Jane Doe", the password `Zq7#janex` is accepted. This causes the other 4 failures, the "contains a name word" cases.

I wanted to prove that these two fixes make every test pass. That needed approval to change `policy.go`, or a copy of it in the scratchpad, and those commands weren't approved. So the "all fixed" result is a reading of the failure output, not a run. I left the code alone because you asked for tests, not fixes.

I didn't check digits outside 0–9: the code accepts any Unicode number (`unicode.IsNumber`, so `²` and `½` count), while the spec just says "digit". I didn't add a test for that.

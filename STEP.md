# Step 0: the start

**Branch:** `step-0-start` (the same code as `main`)  
**Next:** `git checkout step-1-red`

## Where we are

A small slice of a payment service: the password rule in `internal/passwordpolicy/`, its spec in `docs/02-requirements/specs/F-006-user.md` (rule BR-006-16), and one project skill in `.claude/skills/writing-unit-tests/`. The unit tests are deleted and two bugs are planted.

## What to do: exercise 1

Start a new Claude Code session here and paste:

```text
Users report that a password of exactly 8 characters is refused.
Fix it with test-driven development and the project's unit test skill.
```

## What to show

```bash
go test ./...                               # no tests to run
cat docs/02-requirements/specs/F-006-user.md
ls .claude/skills/
```

To skip the wait (about two minutes), go to the next branch.

## All steps

| Branch | What it holds |
|---|---|
| `step-0-start` | Exercise 1 starts: two bugs, no tests, the unit test skill |
| `step-1-red` | The failing test the agent wrote, before the fix |
| `step-1-green` | The one-character fix: exercise 1 done |
| `step-2-done` | Exercise 2 done: the rest of the tests, a known bug skipped, mutations |
| `step-3-start` | Exercise 3 starts: the same bugs, no tests, no project skill |
| `step-3-baseline` | What an agent wrote with no skill |
| `step-3-skill` | A short skill written against that baseline |
| `step-3-check` | What a fresh agent wrote with that skill |

# Step 3.0: exercise 3 starts here

**Branch:** `step-3-start` (the same code as `exercise-3`)  
**Next:** `git checkout step-3-baseline`

## Where we are

The password rule has two planted bugs and no unit tests. There is **no project skill** on this branch: `.claude/skills/` does not exist. The spec is in `docs/02-requirements/specs/F-006-user.md`.

## What to do

Start a new Claude Code session here and paste:

```text
internal/passwordpolicy has no unit tests. Write them.
```

## What to show

```bash
ls -a                                  # no .claude folder
go test ./internal/passwordpolicy/     # no tests to run
```

To skip the wait, go to the next branch: it holds the result of this prompt.

# Step 3.2: a skill written against the baseline

**Branch:** `step-3-skill`  
**Previous:** `step-3-baseline` · **Next:** `git checkout step-3-check`

## Where we are

The code is back at the start of the exercise: two bugs, no tests. One file is new: `.claude/skills/finishing-tests/SKILL.md`, about 15 lines. It was written against the four gaps of the baseline, and nothing else.

| Gap in the baseline | Line in the skill |
|---|---|
| Failing tests left in the suite | Step 1, and the first row of the mistakes table |
| "The fixes should make it pass", not run | Step 2, and the second row |
| No link from a test to a rule | Step 3, and the third row |

The skill says nothing about where expectations come from. The agent already went to the spec without help, so that line would be noise.

## What to do

Start a new session and paste the same task, with the skill named:

```text
internal/passwordpolicy has no unit tests. Write them with the finishing-tests skill.
```

## What to show

```bash
cat .claude/skills/finishing-tests/SKILL.md
git diff --stat step-3-start
```

To skip the wait (the run took about six minutes), go to the next branch.

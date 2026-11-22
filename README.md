# Tests You Can Trust: workshop repository

A small slice of a Go payment service for the workshop "Teaching Claude Code to Write Tests You Can Trust: Skills, TDD and Proof of Red".

It holds one package (the password-strength rule), its spec, and the project skill an agent uses to write unit tests. The unit tests of the package are deleted, and the code has defects planted for the exercises.

## Setup

You need Go 1.24 or newer and Claude Code with the Superpowers plugin.

```bash
git clone https://github.com/ibalashov-srb/tests-you-can-trust-workshop.git
cd tests-you-can-trust-workshop
go test ./...
```

The last command should answer `no tests to run` for `internal/passwordpolicy`.

## Exercise 1: fix a bug test-first

Start a new Claude Code session in this folder and paste:

```text
Users report that a password of exactly 8 characters is refused.
Fix it with test-driven development and the project's unit test skill.
```

## Exercise 2: tests from the spec, then proof of red

In a new session:

```text
internal/passwordpolicy has few unit tests. Write the rest with the writing-unit-tests skill.
Finish with three things: the table of behaviours with their rule ids,
the places where the code breaks the spec, and the table of mutations.
```

## Exercise 3: write your own test skill

It starts on another branch, without the project skill: `git checkout exercise-3`, then read the README there.

## Step branches

Each exercise takes an agent several minutes. To see a result without waiting, check out the branch that holds it. Every `step-` branch has a `STEP.md` at its root: where you are, what to show, the prompt that produced it and what the agent answered.

The order to walk through them. You start exercises on `main` and `exercise-3`; the `step-` branches hold the results.

| # | Branch | What it holds | Tests | Where to look |
|---|---|---|---|---|
| 1 | `main` | The start: two bugs, no tests, the unit test skill | none | This README: the setup and the prompts of exercises 1 and 2 |
| 2 | `step-0-start` | Exercise 1 starts: the same code as `main` | none | The spec `docs/02-requirements/specs/F-006-user.md` (rule BR-006-16), `ls .claude/skills/`, and `go test ./...` answers "no tests to run" |
| 3 | `step-1-red` | The failing test the agent wrote, before the fix | red | `internal/passwordpolicy/policy_test.go`: the expected values come from the spec, five length cases, only the two 8-byte cases fail; a second test is skipped with a bug id |
| 4 | `step-1-green` | The one-character fix: exercise 1 done | green: 6 pass, 1 skipped | `git diff step-1-red -- internal/passwordpolicy/policy.go` (`<=` became `<`); the agent's answer in `STEP.md`: four mutations and a second bug nobody asked about |
| 5 | `step-2-done` | Exercise 2 done: the rest of the tests, a known bug skipped, mutations | green: 36 pass, 1 skipped | `grep -n "t.Skip" internal/passwordpolicy/policy_test.go` (the known bug); in `STEP.md` the table of behaviours with rule ids and the table of 25 mutations; the cost: about 13 minutes and 3 dollars |
| 6 | `exercise-3` | Exercise 3 starts for you: no project skill | none | The README on that branch: the three steps of the exercise |
| 7 | `step-3-start` | Exercise 3 starts: the same code as `exercise-3` | none | `ls -a`: there is no `.claude` folder, so no project skill |
| 8 | `step-3-baseline` | What an agent wrote with no skill | red: 19 failing sub-tests | `go test ./internal/passwordpolicy/` fails; `grep -c "t.Skip" internal/passwordpolicy/policy_test.go` gives 0; in `STEP.md` the parts "What is good" and "What is missing" (four gaps) |
| 9 | `step-3-skill` | A short skill written against that baseline | none | `.claude/skills/finishing-tests/SKILL.md` (22 lines); in `STEP.md` the table from each gap to a line of the skill |
| 10 | `step-3-check` | What a fresh agent wrote with that skill | green: 103 pass, 15 skipped | `grep -n "t.Skip" internal/passwordpolicy/policy_test.go` (the bugs, with ids); in `STEP.md` the table "baseline / with the skill", and one note: the first four mutations did not compile, so the agent ran them again |

The results on these branches are what one agent run produced. Your own run will differ in wording and in numbers.

## Review checklist

[CHECKLIST.md](CHECKLIST.md) has eight questions for reviewing tests that an agent wrote.

## Start again

```bash
git checkout -- . && git clean -fd
```

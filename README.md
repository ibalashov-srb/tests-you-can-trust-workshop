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

The results on these branches are what one agent run produced. Your own run will differ in wording and in numbers.

## Review checklist

[CHECKLIST.md](CHECKLIST.md) has eight questions for reviewing tests that an agent wrote.

## Start again

```bash
git checkout -- . && git clean -fd
```

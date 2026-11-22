# Tests You Can Trust: workshop repository

A small slice of a Go payment service for the workshop "Teaching Claude Code to Write Tests You Can Trust: Skills, TDD and Proof of Red".

It holds one package (the password-strength rule), its spec, and the project skill an agent uses to write unit tests. The unit tests of the package are deleted, and the code has defects planted for the exercises.

## Setup

You need Go 1.24 or newer and Claude Code with the Superpowers plugin.

```bash
git clone https://github.com/ibalashov-srb/tests-you-can-trust-workshop.git
cd tests-you-can-trust-workshop
git checkout exercise-3
go test ./...
```

The last command should answer `no tests to run` for `internal/passwordpolicy`.

## Exercise 3: write your own test skill, test-first

This branch has no project skill. You will write one.

**1. Baseline.** Start a new Claude Code session in this folder and paste:

```text
internal/passwordpolicy has no unit tests. Write them.
```

Read what the agent did and write down what is missing. Three things to look at: is the test run green or red at the end, what happened to the tests that show a bug, and did the agent prove that any test can fail?

**2. Write the skill.** Reset the repository (see below). Create `.claude/skills/<your-skill-name>/SKILL.md` against exactly the gaps you saw. It needs a `name` and a `description` that says when to use it, the steps, and a table of mistakes with the fix for each.

**3. Check.** Start a new session and paste the same task with your skill named:

```text
internal/passwordpolicy has no unit tests. Write them with the <your-skill-name> skill.
```

Compare the two results. If the agent found a new way to get it wrong, add a row to your table and run it again.

## Start again

```bash
git checkout -- . && git clean -fd internal
```

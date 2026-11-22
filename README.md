# Tests You Can Trust: workshop repository

A small slice of a Go payment service for the workshop "Teaching Claude Code to Write Tests You Can Trust: Skills, TDD and Proof of Red".

It holds one package (the password-strength rule), its spec, and the project skill an agent uses to write unit tests. The unit tests of the package are deleted, and the code has defects planted for the exercises.

## Setup

You need Go 1.24 or newer and Claude Code with the Superpowers plugin.

```bash
git clone <this repository>
cd <its folder>
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

## Start again

```bash
git checkout -- . && git clean -fd
```

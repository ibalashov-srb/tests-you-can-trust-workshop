# CLAUDE.md

A small slice of a Go payment service, cut down for a workshop on tests written by an AI agent.

## Layout

- `internal/passwordpolicy/` - the password-strength rule
- `internal/entity/` - domain types
- `docs/02-requirements/specs/` - feature specs with numbered business rules (`BR-XXX-NN`)
- `.claude/skills/` - project skills

## Specs

Product intent and business rules live in `docs/02-requirements/specs/F-XXX-*.md` (index: `docs/02-requirements/features.md`). Before changing a feature, read its spec. Code that contradicts a spec is a bug in the code, or the spec must be updated in the same change.

## Testing

For writing tests use the project skill `writing-unit-tests` (`.claude/skills/`). This includes the RED step of TDD: when `superpowers:test-driven-development` has you write a failing test, load the project skill before writing it. TDD says when to write the test; the project skill says how in this repo.

```bash
go test -count=1 ./...
```

## Issue tracker

There is no issue tracker in this copy. Where a skill tells you to file a bug or a question with `bd`, print the exact `bd create` command in your answer instead, and use the ids `demo1`, `demo2`, ... in the test code.

The skill was copied from the full repository. Some paths it mentions (`.claude/rules/rules.md`, `internal/authz`, `internal/usecases/...`) are not in this copy.

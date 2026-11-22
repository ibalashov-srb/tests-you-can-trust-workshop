# CLAUDE.md

A small slice of a Go payment service, cut down for a workshop on tests written by an AI agent.

## Layout

- `internal/passwordpolicy/` - the password-strength rule
- `internal/entity/` - domain types
- `docs/02-requirements/specs/` - feature specs with numbered business rules (`BR-XXX-NN`)
- `.claude/skills/` - project skills (none yet: you will write one)

## Specs

Product intent and business rules live in `docs/02-requirements/specs/F-XXX-*.md` (index: `docs/02-requirements/features.md`). Before changing a feature, read its spec. Code that contradicts a spec is a bug in the code, or the spec must be updated in the same change.

## Testing

```bash
go test -count=1 ./...
```

## Issue tracker

There is no issue tracker in this copy. If you find a bug, describe it in your answer: file, line, what happens and what the spec says.

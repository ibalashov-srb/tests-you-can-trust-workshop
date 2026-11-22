# Step 1.2: green, exercise 1 done

**Branch:** `step-1-green`  
**Previous:** `step-1-red` · **Next:** `git checkout step-2-done`

## Where we are

The agent changed one character and the test passes. Both files are exactly what the agent left.

## What to show

```bash
git diff step-1-red -- internal/passwordpolicy/policy.go   # one character
go test ./internal/passwordpolicy/                         # ok
```

## What to point at

- The fix is `<=` to `<`. The test was there first and was seen failing.
- Read the agent's answer below: it also ran four mutations, and it reports a second bug that nobody asked about.

## What to do next: exercise 2

Start a new session and paste:

```text
internal/passwordpolicy has few unit tests. Write the rest with the writing-unit-tests skill.
Finish with three things: the table of behaviours with their rule ids,
the places where the code breaks the spec, and the table of mutations.
```

This run is long: about 13 minutes in the trial. The next branch holds its result.

## What the agent answered

I've fixed it: passwords of exactly 8 characters are now accepted. The spec (BR-006-16) sets the minimum at 8, counted in bytes, but `Validate` used `len(password) <= MinLength` and so refused exactly 8. It is now `<` (`internal/passwordpolicy/policy.go:75`). I wrote the failing test first.

**Tests** (`internal/passwordpolicy/policy_test.go`, the first tests in the repo)
- `TestValidate_Length` covers the five length cases in the spec: 8 bytes accepted, 7 refused, 6 characters in 8 bytes accepted (the spec counts bytes), 128 accepted, 129 refused.
- **Red:** before the fix, both 8-byte cases failed with `ErrTooShort`. **Green:** after the fix, all of `go test ./...` passes, and `go vet` and `gofmt` report nothing.
- **Mutation check:** the 7, 128 and 129 cases already passed before the fix, so I broke the code one line at a time to make sure each test can fail. Each change was caught by its test:

| Mutation | Test that failed |
|---|---|
| `< MinLength` → `< MinLength-1` | refuses 7 bytes |
| `len` → `utf8.RuneCountInString` (min check) | 6 characters in 8 bytes |
| `> MaxLength` → `>= MaxLength` | accepts 128 bytes |
| `> MaxLength` → `> MaxLength+1` | refuses 129 bytes |

**A second bug, found along the way and not fixed:** the spec says only name words shorter than 4 characters are skipped, but the code checks `> SimilarityMinLength` (`policy.go:117`). So a 4-character name like "Anna" is never checked, and `Kq7#annaZ` is accepted for Anna Westbrook. Following the skill, I added a skipped test, `TestValidate_NameWordOfFourCharacters`, and confirmed it fails for that reason when the skip is removed. The repo has no tracker, so here is the command to file it:

```
bd create --type=bug --title="passwordpolicy: 4-char name words not checked for similarity" --description="internal/passwordpolicy/policy.go:117 uses > SimilarityMinLength; BR-006-16 says only words shorter than 4 are skipped. Observed: Validate(\"Kq7#annaZ\", FullName \"Anna Westbrook\") returns nil. Expected: ErrSimilarToName. Test: TestValidate_NameWordOfFourCharacters (skipped, demo2)."
```

The length bug would be `demo1`, but it is fixed in this change, so it needs no entry. Nothing is committed. I only tested length; the other rules (character classes, blocklist, email similarity) are still untested.

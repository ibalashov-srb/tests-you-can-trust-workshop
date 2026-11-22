# Step 3.3: the same task, with the skill

**Branch:** `step-3-check`  
**Previous:** `step-3-skill` · **Compare with:** `step-3-baseline`

## Where we are

This is what a fresh agent produced from the same task with the skill named in the prompt. Nothing was edited by hand.

```text
internal/passwordpolicy has no unit tests. Write them with the finishing-tests skill.
```

## What to show

```bash
go test ./internal/passwordpolicy/                        # ok
grep -n "t.Skip" internal/passwordpolicy/policy_test.go   # the bugs, with ids
git diff --stat step-3-baseline -- internal/              # baseline against this
```

## What changed against the baseline

| | Baseline | With the skill |
|---|---|---|
| Test run | red | green |
| Tests that show a bug | left failing | skipped with a bug id |
| Proof that tests can fail | none | a table of mutations |
| Tests linked to spec rules | no | a table of behaviours |
| Time | about 2 minutes | about 6 minutes |

One thing to notice in the answer below: the first four mutations did not compile, so they proved nothing, and the agent ran them again. A compile error is not a red. That is the next row for the skill's mistakes table.

## What the agent answered

I wrote the tests in `internal/passwordpolicy/policy_test.go` and the suite passes (`go test -count=1 ./...`). Writing them turned up two places where the code breaks spec rule BR-006-16. Following the skill, I didn't fix either one: each has a test that I saw fail, now skipped with a bug id.

**Bugs found**
- **demo1**, `internal/passwordpolicy/policy.go:79`: the check is `len(password) <= MinLength`, so a password of exactly 8 characters is refused as too short. The spec says the minimum is 8. Because of this, 13 eight-character blocklist entries (such as `p@ssw0rd`) come back "too short" instead of "too common".
- **demo2**, `internal/passwordpolicy/policy.go:134`: name words are only checked when longer than `SimilarityMinLength`, so 4-character words are never checked. The spec says only words *shorter than* 4 are skipped, so a user named "Anna Lee" can use `Anna7!xyzw`.

**Behaviours tested**

| Behaviour | Rule | Test |
|---|---|---|
| 7 characters refused, 128 bytes accepted, 129 refused, limit counts bytes not characters | BR-006-16 | `TestValidate_Length` |
| 8 characters accepted (skipped, demo1) | BR-006-16 | `TestValidate_EightCharsAccepted` |
| Uppercase, lowercase, digit and special character each required | BR-006-16 | `TestValidate_CharacterClasses` |
| `! _ - + € «` count as special; a space does not | BR-006-16 | `TestValidate_SpecialCharacter` |
| Every blocklist entry refused, case-insensitive, only the whole password compared | BR-006-16 | `TestValidate_CommonPasswords` |
| Email local part: password contains it or is contained in it, case-insensitive, 4 characters checked, 3 not | BR-006-16 | `TestValidate_SimilarToEmail` |
| Name word: password contains it, case-insensitive, name split on non-letters, 3-letter word not checked | BR-006-16 | `TestValidate_SimilarToName` |
| 4-letter name word checked (skipped, demo2) | BR-006-16 | `TestValidate_FourCharNameWordChecked` |
| No owner given: only the rules that don't need an owner apply | BR-006-16 | `TestValidate_NilOwner` |

**Mutations run** (one line of `policy.go` changed each time, then restored)

| Mutation | Test that failed |
|---|---|
| `len <= MinLength` → `< MinLength-1` | `TestValidate_Length` |
| `> MaxLength` → `> MaxLength+1` | `TestValidate_Length` (129 bytes) |
| max counted in characters instead of bytes | `TestValidate_Length` (130 bytes) |
| `!hasUpper` → `!hasUpper && false` | `TestValidate_CharacterClasses` |
| same for lowercase | `TestValidate_CharacterClasses` |
| same for digit | `TestValidate_CharacterClasses` |
| same for special | `TestValidate_CharacterClasses`, `TestValidate_SpecialCharacter` |
| `unicode.IsSymbol` dropped | `TestValidate_SpecialCharacter` (`+`, `€`) |
| space counted as special | `TestValidate_SpecialCharacter` (space) |
| blocklist check turned off | `TestValidate_CommonPasswords` |
| blocklist compared case-sensitively | `TestValidate_CommonPasswords` |
| blocklist matched on a prefix | `TestValidate_CommonPasswords` (whole password) |
| "password contains email local part" turned off | `TestValidate_SimilarToEmail` |
| "local part contains password" turned off | `TestValidate_SimilarToEmail` |
| email local part not lower-cased | `TestValidate_SimilarToEmail` |
| local part minimum `>= 4` → `> 4` | `TestValidate_SimilarToEmail` (4 characters) |
| local part minimum `>= 4` → `>= 3` | `TestValidate_SimilarToEmail` (3 characters) |
| name check turned off | `TestValidate_SimilarToName` |
| name not lower-cased | `TestValidate_SimilarToName` |
| name split on spaces only | `TestValidate_SimilarToName` (`Smith-Jones`) |
| name word minimum `> 4` → `> 2` | `TestValidate_SimilarToName` (3 letters) |
| no owner → returns an error | `TestValidate_NilOwner` and 4 others |

My first runs of the four character-class mutations used `if false`, which doesn't compile, so they proved nothing; the table shows the `&& false` reruns. Afterwards `git status` shows only the new test file (plus the `.claude/` folder that was already there), so `policy.go` is back to its original state.

I didn't test digits like `²` or `Ⅷ`: the code accepts any Unicode number where the spec says "digit", and it's unclear which is intended. Two leftovers can be removed: `internal/testdeps` (it only kept testify in `go.mod` while there were no tests) and an unused scratch script (it is not on this branch).

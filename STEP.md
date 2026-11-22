# Step 2: exercise 2 done

**Branch:** `step-2-done`  
**Previous:** `step-1-green` · **Next:** `git checkout step-3-start`

## Where we are

The agent wrote the rest of the unit tests with the project skill. The test file is exactly what it left after this prompt:

```text
internal/passwordpolicy has few unit tests. Write the rest with the writing-unit-tests skill.
Finish with three things: the table of behaviours with their rule ids,
the places where the code breaks the spec, and the table of mutations.
```

## What to show

```bash
go test -v ./internal/passwordpolicy/ | grep -c -- "--- PASS"   # the passing tests
grep -n "t.Skip" internal/passwordpolicy/policy_test.go          # the known bug
git diff --stat step-1-green
```

## What to point at

- **The known bug stays marked.** The test for the second planted bug is skipped with a bug id, and the run stays green. A red suite would hide the next real failure.
- **The table of behaviours** in the answer below: every test names a rule of BR-006-16.
- **The table of mutations:** 25 one-line changes to the code, and a test failed for each one. The production file is unchanged at the end.
- **The cost:** about 13 minutes and 3 dollars for this run. Proof of red is not free.

## What the agent answered

I added 26 tests to `internal/passwordpolicy/policy_test.go`. Every rule in BR-006-16 now has at least one test, and I made one deliberate wrong change to the code for each behaviour: all 25 turned a test red. The full suite passes (36 tests, plus the one known-bug test that stays skipped), and `go vet` and `gofmt` are clean. Only the test file changed; `policy.go` is byte-identical to its backup. Nothing is committed.

The function takes no ports, so no mocks were needed. The four rows the skill always asks for are all N/A here: `Validate` reads no rows, isn't run by a worker, takes no Principal and guards no invariant on a set.

## 1. Behaviours (all BR-006-16)

| Behaviour | Test |
|---|---|
| Exactly 8 bytes accepted, 7 refused (`ErrTooShort`) | `TestValidate_Length/accepts_exactly_8…`, `…/refuses_7…` (already existed) |
| Exactly 128 bytes accepted, 129 refused (`ErrTooLong`) | `TestValidate_Length/accepts_exactly_128…`, `…/refuses_129…` (already existed) |
| Length counts bytes, not characters | `TestValidate_Length/counts_bytes…` (already existed) |
| No uppercase → `ErrNoUppercase` | `TestValidate_CharacterClasses/refuses_…uppercase` |
| No lowercase → `ErrNoLowercase` | `…/refuses_…lowercase` |
| No digit → `ErrNoNumber` | `…/refuses_…digit` |
| Letters and digits only → `ErrNoSpecial` | `…/refuses_…letters_and_digits_only…` |
| A space is not a special character | `…/refuses_a_space…` |
| `! _ - + € «` each count as special | `…/accepts_<c>_as_the_special_character` (6 cases) |
| The spec's blocklist examples are refused (`ErrCommon`) | `TestValidate_CommonPasswords/refuses_the_examples…` |
| Every blocklist entry is refused | `…/refuses_every_blocklist_entry…` |
| Blocklist comparison ignores case | `…/compares_case-insensitively` |
| Only the whole password is compared, not a substring | `…/accepts_a_password_that_only_contains…` |
| Password containing the email's local part, in another case → `ErrSimilarToEmail` | `TestValidate_SimilarToEmail/refuses_a_password_containing…` |
| Password contained in the local part → refused | `…/refuses_a_password_contained…` |
| A 4-character local part is checked, a 3-character one is not | `…/checks_…exactly_4…`, `…/does_not_check_…3…` |
| Local-part length counts characters, not bytes | `…/counts_the_local_part_in_characters…` |
| No edit-distance matching (email) | `…/does_not_match_by_edit_distance` |
| Password containing any word of the name, in another case → `ErrSimilarToName` | `TestValidate_SimilarToName/refuses_…any_word…` |
| The name is split on punctuation, not only spaces | `…/splits_the_name_on_punctuation…` |
| Digits stay inside a name word | `…/keeps_digits_inside_a_name_word` |
| Name words shorter than 4 characters are not checked | `…/does_not_check_name_words_shorter…` |
| No edit-distance matching (name) | `…/does_not_match_by_edit_distance` |
| A name word of exactly 4 characters is checked | `TestValidate_NameWordOfFourCharacters` (skipped as **demo2**, already existed) |

**Not tested here:**
- **Which error wins when several rules fail:** the spec doesn't say.
- **A missing user, or a user without a full name:** no rule covers them.
- **"Every path that sets a password calls `Validate`":** those callers aren't in this copy.

## 2. Where the code breaks the spec

- **demo2 (bug), `internal/passwordpolicy/policy.go:129`:** name words are checked with `> SimilarityMinLength`, so a word of exactly 4 characters is skipped. The spec says only words shorter than 4 are skipped. The existing skipped test catches it: I removed the skip once, it failed with `Expected error "password is too similar to name" … got nil`, and I put the skip back. There's no issue tracker in this copy, so this is the command to file it:
  `bd create --type=bug --title="passwordpolicy: 4-char name words not checked" --description="policy.go:129 uses > SimilarityMinLength; BR-006-16 says only words shorter than 4 are skipped. Observed: 'Kq7#annaZ' accepted for 'Anna Westbrook'. Expected: ErrSimilarToName. Test: TestValidate_NameWordOfFourCharacters (skipped)."`
- **Question for you: what counts as a digit?** `policy.go:89` uses `unicode.IsNumber`, so `½`, `²` and `Ⅷ` satisfy the digit rule. The spec only says "digit" (`unicode.IsDigit` would reject those). I didn't write a test that locks in today's behaviour. I would switch to `IsDigit`, but the spec should decide.
- **Possible issue, not asserted: `policy.go:116`.** The local part is taken up to the first `@`. The email standard puts it before the last `@`. This only matters for quoted local parts that contain `@`, which are unusual.

## 3. Mutations

Each mutation was a one-line change to `policy.go`, run against the named test, then restored and checked with `cmp`. All 25 turned the test red with an assertion failure; none was a compile error.

| # | Mutation | Test that went red |
|---|---|---|
| M1 | `< MinLength` → `<= MinLength` | `Length/accepts_exactly_8_bytes` |
| M2 | `> MaxLength` → `>= MaxLength` | `Length/accepts_exactly_128_bytes` |
| M3 | `len(password)` → `utf8.RuneCountInString(password)` | `Length/counts_bytes…` |
| M4 | uppercase check disabled | `CharacterClasses/refuses_…uppercase` |
| M5 | lowercase check disabled | `CharacterClasses/refuses_…lowercase` |
| M6 | digit check disabled | `CharacterClasses/refuses_…digit` |
| M7 | special check disabled | `CharacterClasses/refuses_…letters_and_digits_only…` |
| M8 | special = `IsPunct` only | `CharacterClasses/accepts_+…`, `accepts_€…` |
| M9 | special = `IsSymbol` only | `CharacterClasses/accepts_!…`, `_`, `-`, `«` |
| M10 | special = anything not a letter or digit | `CharacterClasses/refuses_a_space…` |
| M11 | password not lower-cased | `CommonPasswords/compares_case-insensitively` |
| M12 | blocklist check disabled | `CommonPasswords/refuses_every_blocklist_entry…` |
| M13 | blocklist matched as a substring | `CommonPasswords/accepts_…only_contains…` |
| M14 | "password contains local part" dropped | `SimilarToEmail/refuses_…containing…` |
| M15 | "local part contains password" dropped | `SimilarToEmail/refuses_…contained…` |
| M16 | local part `>= 4` → `> 4` | `SimilarToEmail/checks_…exactly_4…` |
| M17 | local part `>= 4` → `>= 3` | `SimilarToEmail/does_not_check_…3…` |
| M18 | local-part length in bytes | `SimilarToEmail/counts_…characters…` |
| M19 | local part not lower-cased | `SimilarToEmail/refuses_…containing…` |
| M20 | whole email used instead of the local part | `SimilarToEmail/refuses_…containing…` |
| M21 | name words of any length checked | `SimilarToName/does_not_check_…shorter…` |
| M22 | name split on spaces only | `SimilarToName/splits_…punctuation…` |
| M23 | name split on digits too | `SimilarToName/keeps_digits…` |
| M24 | name not lower-cased | `SimilarToName/refuses_…any_word…` |
| M25 | only the first name word checked | `SimilarToName/refuses_…any_word…` |

M19 initially got past my email test because its test email was all lowercase. I changed the test data to mixed case before running the mutation; with that change it went red.

**Commands run:**
- `go test -count=1 ./...`
- `go test -count=1 -v ./internal/passwordpolicy/`
- `go vet ./internal/passwordpolicy/`
- `gofmt -l internal/passwordpolicy/`
- per mutation: `go -C <repo> test -count=1 -run '<test>' ./internal/passwordpolicy/`, then `cp` and `cmp` to restore

I did not mutate the blocklist file parsing (the lower-casing and space trimming): the file is already lowercase with no stray spaces, so changing that code has no effect. The "first `@`" behaviour has no test.

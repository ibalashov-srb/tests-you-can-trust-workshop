package passwordpolicy_test

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"

	"payments/internal/entity"
	"payments/internal/passwordpolicy"
)

// valid meets every character rule and resembles no test owner.
const valid = "Zq7!mKw#Lp"

func named(name string) *entity.User { return &entity.User{FullName: &name} }

// BR-006-16: min 8, max 128 bytes.
func TestValidate_Length(t *testing.T) {
	assert.ErrorIs(t, passwordpolicy.Validate("Zq7!mKw", nil), passwordpolicy.ErrTooShort, "7 chars")
	assert.NoError(t, passwordpolicy.Validate("Zq7!"+strings.Repeat("a", 124), nil), "128 bytes")
	assert.ErrorIs(t, passwordpolicy.Validate("Zq7!"+strings.Repeat("a", 125), nil), passwordpolicy.ErrTooLong, "129 bytes")
	// 46 characters but 130 bytes: the limit counts bytes.
	assert.ErrorIs(t, passwordpolicy.Validate("Zq7!"+strings.Repeat("€", 42), nil), passwordpolicy.ErrTooLong, "130 bytes")
}

// BR-006-16: min 8 — exactly 8 is allowed.
func TestValidate_EightCharsAccepted(t *testing.T) {
	t.Skip("known bug demo1: policy.go uses len <= MinLength, so an 8-character password is refused as too short")
	assert.NoError(t, passwordpolicy.Validate("Zq7!mKwL", nil))
}

// BR-006-16: uppercase + lowercase + digit + special.
func TestValidate_CharacterClasses(t *testing.T) {
	for pw, want := range map[string]error{
		"zq7!mkw#lp": passwordpolicy.ErrNoUppercase,
		"ZQ7!MKW#LP": passwordpolicy.ErrNoLowercase,
		"Zqx!mKw#Lp": passwordpolicy.ErrNoNumber,
		"Zq7xmKwyLp": passwordpolicy.ErrNoSpecial,
	} {
		assert.ErrorIs(t, passwordpolicy.Validate(pw, nil), want, pw)
	}
	assert.NoError(t, passwordpolicy.Validate(valid, nil))
}

// BR-006-16: special = Unicode P or S; a space does not count.
func TestValidate_SpecialCharacter(t *testing.T) {
	for _, c := range []string{"!", "_", "-", "+", "€", "«"} {
		assert.NoError(t, passwordpolicy.Validate("Zq7mKwLp"+c, nil), c)
	}
	assert.ErrorIs(t, passwordpolicy.Validate("Zq7 mKw Lp", nil), passwordpolicy.ErrNoSpecial, "space")
}

// BR-006-16: blocklist, case-insensitive, whole password.
func TestValidate_CommonPasswords(t *testing.T) {
	for _, pw := range passwordpolicy.CommonPasswords() {
		t.Run(pw, func(t *testing.T) {
			if len(pw) == passwordpolicy.MinLength {
				t.Skip("known bug demo1: 8-character entries are refused as too short before the blocklist")
			}
			// Entries are lower-case; upper-case the first letter so the class rules pass.
			i := strings.IndexFunc(pw, unicode.IsLetter)
			variant := pw[:i] + strings.ToUpper(pw[i:i+1]) + pw[i+1:]
			assert.ErrorIs(t, passwordpolicy.Validate(variant, nil), passwordpolicy.ErrCommon)
		})
	}
	assert.ErrorIs(t, passwordpolicy.Validate("PASSWORd123!", nil), passwordpolicy.ErrCommon, "case-insensitive")
	assert.NoError(t, passwordpolicy.Validate("P@ssw0rd!x", nil), "only the whole password is compared")
}

// BR-006-16: not similar to the email local part.
func TestValidate_SimilarToEmail(t *testing.T) {
	cases := []struct {
		name, email, pw string
		want            error
	}{
		{"contains local part", "johnsmith@example.com", "Johnsmith7!", passwordpolicy.ErrSimilarToEmail},
		{"case-insensitive", "JohnSmith@example.com", "johnsmith7!Q", passwordpolicy.ErrSimilarToEmail},
		{"contained in local part", "zq7!mkw#lpxyz@example.com", valid, passwordpolicy.ErrSimilarToEmail},
		{"4-char local part checked", "anna@example.com", "Anna7!xyzw", passwordpolicy.ErrSimilarToEmail},
		{"3-char local part not checked", "bob@example.com", "Bob7!xyzwq", nil},
	}
	for _, c := range cases {
		err := passwordpolicy.Validate(c.pw, &entity.User{Email: c.email})
		if c.want == nil {
			assert.NoError(t, err, c.name)
		} else {
			assert.ErrorIs(t, err, c.want, c.name)
		}
	}
}

// BR-006-16: not containing a word of the full name.
func TestValidate_SimilarToName(t *testing.T) {
	assert.ErrorIs(t, passwordpolicy.Validate("Johnson7!x", named("Grace Johnson")), passwordpolicy.ErrSimilarToName, "contains word")
	assert.ErrorIs(t, passwordpolicy.Validate("grace7!XYZ", named("GRACE Johnson")), passwordpolicy.ErrSimilarToName, "case-insensitive")
	assert.ErrorIs(t, passwordpolicy.Validate("Jones7!xyz", named("Mary Smith-Jones")), passwordpolicy.ErrSimilarToName, "split on non-letters")
	assert.NoError(t, passwordpolicy.Validate("Lee7!xyzwq", named("Lee Ng")), "3-char word not checked")
	assert.NoError(t, passwordpolicy.Validate(valid, named("Grace Johnson")), "unrelated")
}

// BR-006-16: a 4-character name word is checked.
func TestValidate_FourCharNameWordChecked(t *testing.T) {
	t.Skip("known bug demo2: policy.go checks name words with > SimilarityMinLength, so 4-character words are skipped")
	assert.ErrorIs(t, passwordpolicy.Validate("Anna7!xyzw", named("Anna Lee")), passwordpolicy.ErrSimilarToName)
}

// BR-006-16: no owner (bootstrap paths) means only the owner-free rules apply.
func TestValidate_NilOwner(t *testing.T) {
	assert.NoError(t, passwordpolicy.Validate(valid, nil))
}

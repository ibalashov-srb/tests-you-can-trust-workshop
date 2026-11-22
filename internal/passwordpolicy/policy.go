// Package passwordpolicy is the password-strength rule (F-006 BR-006-16): the
// one source of truth for every path that sets a password (ChangePassword,
// SetPasswordWithToken, the bootstrap super admin) and for what
// GetPasswordRequirements advertises.
package passwordpolicy

import (
	_ "embed"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	"payments/internal/entity"
)

// The limits of the rule.
const (
	MinLength = 8
	MaxLength = 128

	// SpecialCharacterRule says in words what counts as a special
	// character: anything in the Unicode punctuation (P) or symbol (S)
	// categories, so `_`, `-`, `+`, `€` and `«` all qualify; a space does not.
	SpecialCharacterRule = "any Unicode punctuation or symbol character (for example ! _ - + € «); letters, digits and spaces do not count"
)

// SimilarityMinLength is the shortest email local part or name word
// (in characters) the similarity check looks for. Shorter ones would refuse any
// password that happens to contain a common letter pair.
const SimilarityMinLength = 4

// The ways a password breaks the rule.
var (
	ErrTooShort       = errors.New("password must be at least 8 characters")
	ErrTooLong        = errors.New("password must be at most 128 characters")
	ErrNoUppercase    = errors.New("password must contain at least one uppercase letter")
	ErrNoLowercase    = errors.New("password must contain at least one lowercase letter")
	ErrNoNumber       = errors.New("password must contain at least one number")
	ErrNoSpecial      = errors.New("password must contain at least one special character")
	ErrSimilarToEmail = errors.New("password is too similar to email")
	ErrSimilarToName  = errors.New("password is too similar to name")
	ErrCommon         = errors.New("password is too common")
)

//go:embed common_passwords.txt
var commonPasswordsFile string

// commonPasswords is the embedded blocklist, lower-cased.
var commonPasswords = parseCommonPasswords(commonPasswordsFile)

func parseCommonPasswords(file string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, line := range strings.Split(file, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		set[strings.ToLower(line)] = struct{}{}
	}
	return set
}

// Validate checks password against the password-strength rule
// (F-006 BR-006-16) for owner, the user whose password it becomes:
//   - MinLength..MaxLength bytes;
//   - an uppercase letter, a lowercase letter, a digit and a special character
//     (SpecialCharacterRule);
//   - not on the common-password blocklist (case-insensitive, whole password);
//   - not similar to the owner's email or name: it must not contain the email's
//     local part nor be contained in it, and must not contain any word of the
//     full name. Case-insensitive; a local part or word shorter than
//     SimilarityMinLength is not checked.
func Validate(password string, owner *entity.User) error {
	if len(password) < MinLength {
		return ErrTooShort
	}
	if len(password) > MaxLength {
		return ErrTooLong
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsNumber(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}
	if !hasUpper {
		return ErrNoUppercase
	}
	if !hasLower {
		return ErrNoLowercase
	}
	if !hasNumber {
		return ErrNoNumber
	}
	if !hasSpecial {
		return ErrNoSpecial
	}

	passwordLower := strings.ToLower(password)
	if _, common := commonPasswords[passwordLower]; common {
		return ErrCommon
	}
	if owner == nil {
		return nil
	}

	if idx := strings.Index(owner.Email, "@"); idx > 0 {
		local := strings.ToLower(owner.Email[:idx])
		if utf8.RuneCountInString(local) >= SimilarityMinLength &&
			(strings.Contains(passwordLower, local) || strings.Contains(local, passwordLower)) {
			return ErrSimilarToEmail
		}
	}

	if owner.FullName != nil {
		words := strings.FieldsFunc(strings.ToLower(*owner.FullName), func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
		for _, w := range words {
			if utf8.RuneCountInString(w) > SimilarityMinLength && strings.Contains(passwordLower, w) {
				return ErrSimilarToName
			}
		}
	}
	return nil
}

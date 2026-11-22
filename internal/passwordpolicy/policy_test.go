package passwordpolicy_test

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"payments/internal/entity"
	"payments/internal/passwordpolicy"
)

func ptr(s string) *string { return &s }

// Expectations come from F-006 BR-006-16, not from the code.
func TestValidate(t *testing.T) {
	owner := &entity.User{Email: "jdoe@example.com", FullName: ptr("Jane Doe")}

	tests := []struct {
		name     string
		password string
		owner    *entity.User
		want     error
	}{
		// Length: 8..128 bytes.
		{"7 bytes is too short", "Zq7#vbn", nil, passwordpolicy.ErrTooShort},
		{"exactly 8 bytes is accepted", "Zq7#vbnK", nil, nil},
		{"exactly 128 bytes is accepted", "Zq7#" + strings.Repeat("k", 124), nil, nil},
		{"129 bytes is too long", "Zq7#" + strings.Repeat("k", 125), nil, passwordpolicy.ErrTooLong},
		{"length counts bytes: 6 runes, 8 bytes is accepted", "Zq7€vb", nil, nil},
		{"length counts bytes: 43 runes, 129 bytes is too long", "Zq" + strings.Repeat("€", 42) + "7", nil, passwordpolicy.ErrTooLong},

		// Character classes.
		{"no uppercase", "zq7#vbnkw", nil, passwordpolicy.ErrNoUppercase},
		{"no lowercase", "ZQ7#VBNKW", nil, passwordpolicy.ErrNoLowercase},
		{"no digit", "Zqx#vbnkw", nil, passwordpolicy.ErrNoNumber},
		{"no special", "Zq7xvbnkw", nil, passwordpolicy.ErrNoSpecial},
		{"space is not special", "Zq7 vbnkw", nil, passwordpolicy.ErrNoSpecial},
		{"underscore is special (P)", "Zq7_vbnkw", nil, nil},
		{"hyphen is special (P)", "Zq7-vbnkw", nil, nil},
		{"plus is special (S)", "Zq7+vbnkw", nil, nil},
		{"euro sign is special (S)", "Zq7€vbnkw", nil, nil},
		{"guillemet is special (P)", "Zq7«vbnkw", nil, nil},
		{"non-ASCII letters count for case", "Ωж7#ωbnkw", nil, nil},

		// Blocklist: whole password, case-insensitive.
		{"spec example P@ssw0rd!", "P@ssw0rd!", nil, passwordpolicy.ErrCommon},
		{"spec example Password123!", "Password123!", nil, passwordpolicy.ErrCommon},
		{"spec example Welcome1!", "Welcome1!", nil, passwordpolicy.ErrCommon},
		{"blocklist is case-insensitive", "pASSWORD123!", nil, passwordpolicy.ErrCommon},
		{"blocklist entry inside a longer password is allowed", "Welcome1!Zq7#", nil, nil},

		// Email similarity.
		{"contains email local part", "Xx1!JDOEx", owner, passwordpolicy.ErrSimilarToEmail},
		{"password contained in a long local part", "Ab1!zxcvb", &entity.User{Email: "qqab1!zxcvbqq@x.com"}, passwordpolicy.ErrSimilarToEmail},
		{"local part of 3 chars is not checked", "Abc1!xyzw", &entity.User{Email: "abc@x.com"}, nil},
		{"local part of 4 chars is checked", "Abcd1!xyz", &entity.User{Email: "abcd@x.com"}, passwordpolicy.ErrSimilarToEmail},
		{"domain is not checked", "Example1!zz", owner, nil},

		// Name similarity.
		{"contains a name word", "Zq7#janex", owner, passwordpolicy.ErrSimilarToName},
		{"name check is case-insensitive", "Zq7#JaNex", owner, passwordpolicy.ErrSimilarToName},
		{"name word of 3 chars is not checked", "Doe7#zqxw", &entity.User{Email: "x@x.com", FullName: ptr("Al Doe")}, nil},
		{"name word of 4 chars is checked", "Zq7#Ngozi", &entity.User{Email: "x@x.com", FullName: ptr("Al Ngoz")}, passwordpolicy.ErrSimilarToName},
		{"name split on non-letters", "Zq7#mary!", &entity.User{Email: "x@x.com", FullName: ptr("Anne-Mary O'Brien")}, passwordpolicy.ErrSimilarToName},
		{"no full name skips name check", "Zq7#janew", &entity.User{Email: "x@x.com"}, nil},
		{"nil owner skips similarity", "Xx1!JDOEx", nil, nil},
		{"unrelated owner accepted", "Zq7#vbnkw", owner, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := passwordpolicy.Validate(tt.password, tt.owner)
			if tt.want == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.want)
			}
		})
	}
}

// Every blocklist entry is refused whatever its case.
func TestValidate_EveryCommonPasswordRefused(t *testing.T) {
	list := passwordpolicy.CommonPasswords()
	require.NotEmpty(t, list)
	for _, pw := range list {
		// Upper-case the first letter so the entry passes the character rules.
		i := strings.IndexFunc(pw, unicode.IsLetter)
		require.GreaterOrEqual(t, i, 0, pw)
		r, n := utf8.DecodeRuneInString(pw[i:])
		mixed := pw[:i] + string(unicode.ToUpper(r)) + pw[i+n:]
		t.Run(mixed, func(t *testing.T) {
			assert.ErrorIs(t, passwordpolicy.Validate(mixed, nil), passwordpolicy.ErrCommon)
		})
	}
}

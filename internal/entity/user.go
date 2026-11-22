// Package entity holds the domain types. This workshop copy keeps only what
// the password policy needs.
package entity

// User is a person who signs in to the platform.
type User struct {
	ID       int32
	Email    string
	FullName *string
}

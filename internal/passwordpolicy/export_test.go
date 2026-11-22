package passwordpolicy

// CommonPasswords returns the blocklist Validate refuses, in no
// particular order (tests only).
func CommonPasswords() []string {
	out := make([]string, 0, len(commonPasswords))
	for pw := range commonPasswords {
		out = append(out, pw)
	}
	return out
}

# F-006: User management

**Status:** done

This workshop copy keeps one rule of the feature: the password policy. The full spec covers user creation, listing, deletion and status changes.

## Business rules

**BR-006-16:** Password policy. One validator (`passwordpolicy.Validate`, `internal/passwordpolicy/policy.go`) enforces it on every path that sets a password.

- Min 8, max 128 chars (bytes).
- Must contain uppercase + lowercase + digit + special character. A **special character is any Unicode punctuation or symbol character** (Unicode categories P and S): `!`, `_`, `-`, `+`, `€`, `«` all count; letters, digits and spaces do not.
- Not in common-password blocklist: an embedded list of well-known passwords that meet the character rules (`internal/passwordpolicy/common_passwords.txt`, `P@ssw0rd!`, `Password123!`, `Welcome1!`, ...), compared case-insensitively with the whole password (refused as `password is too common`). A small list, not a breach corpus: it catches the passwords a complexity rule invites, nothing more.
- Not similar to email or name: case-insensitive substring. The password must not contain the email's local part nor be contained in it, and must not contain any word of the user's full name (split on anything but letters and digits). A local part or name word shorter than **4 characters** is not checked, so `a@x.com` or a name like `Al Ng` does not refuse every password containing those letters. No edit-distance matching.

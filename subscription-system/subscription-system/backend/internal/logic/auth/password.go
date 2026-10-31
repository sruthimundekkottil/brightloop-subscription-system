package auth

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Go's regexp (RE2) has no lookaheads, so each rule is its own pattern. The client checklist uses the same classes.
var (
	hasUpper   = regexp.MustCompile(`[A-Z]`)
	hasLower   = regexp.MustCompile(`[a-z]`)
	hasSpecial = regexp.MustCompile(`[^A-Za-z0-9]`)
)

// PasswordProblems returns one message per broken rule; empty means the password is acceptable.
// "More than 8 characters" means at least 9; 72 bytes is bcrypt's limit.
func PasswordProblems(p string) map[string]string {
	probs := map[string]string{}
	if utf8.RuneCountInString(p) < 9 {
		probs["length"] = "Use more than 8 characters."
	}
	if len(p) > 72 {
		probs["maxLength"] = "Use at most 72 bytes."
	}
	if !hasUpper.MatchString(p) {
		probs["upper"] = "Add an uppercase letter."
	}
	if !hasLower.MatchString(p) {
		probs["lower"] = "Add a lowercase letter."
	}
	if !hasSpecial.MatchString(p) {
		probs["special"] = "Add a special character."
	}
	return probs
}

// passwordError joins the problems in a fixed order for the "password" field, or returns "".
func passwordError(p string) string {
	probs := PasswordProblems(p)
	var msgs []string
	for _, k := range []string{"length", "maxLength", "upper", "lower", "special"} {
		if m, ok := probs[k]; ok {
			msgs = append(msgs, m)
		}
	}
	return strings.Join(msgs, " ")
}

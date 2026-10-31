package auth

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

// The same table is in frontend/src/lib/password.test.ts; keep the two in step.
func TestPasswordProblems(t *testing.T) {
	cases := []struct {
		pw   string
		want []string
	}{
		{"Abcdefg!9", nil},               // 9 characters: the minimum
		{"Abcdef!9", []string{"length"}}, // 8 characters is not "more than 8"
		{"abcdefgh!", []string{"upper"}},
		{"ABCDEFGH!", []string{"lower"}},
		{"Abcdefghi", []string{"special"}},
		{"Abcdefgh ", nil},                     // a space counts as special
		{"Aa!" + strings.Repeat("x", 69), nil}, // 72 bytes
		{"Aa!" + strings.Repeat("x", 70), []string{"maxLength"}},
		{"Ab!" + strings.Repeat("é", 35), []string{"maxLength"}}, // 38 characters but 73 bytes
		{"Äbcdéfgh", []string{"length", "upper", "special"}},     // Ä/é are neither A-Z, a-z nor... they are special
		{"", []string{"length", "lower", "special", "upper"}},
	}
	// Ä and é are outside [A-Za-z0-9], so they count as special, exactly as in the client's character classes.
	cases[9].want = []string{"length", "upper"}
	for _, c := range cases {
		var got []string
		for k := range PasswordProblems(c.pw) {
			got = append(got, k)
		}
		sort.Strings(got)
		sort.Strings(c.want)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%q: problems %v, want %v", c.pw, got, c.want)
		}
	}
}

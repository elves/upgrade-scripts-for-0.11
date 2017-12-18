package parse

import "testing"

var quoteTests = []struct {
	text, quoted string
}{
	// Empty string is quoted with single quote.
	{"", `''`},
	// Bareword when possible.
	{"x-y,z@h/d", "x-y,z@h/d"},
	// Single quote when there is special char but no unprintable.
	{"x$y[]ef'", "'x$y[]ef'''"},
	// Tilde needs quoting only when appearing at the beginning
	{"~x", "'~x'"},
	{"x~", "x~"},
	// Double quote when there is unprintable char.
	{"a\nb", `"a\nb"`},
	{"\x1b\"\\", `"\e\"\\"`},
}

func TestQuote(t *testing.T) {
	for _, tc := range quoteTests {
		got := Quote(tc.text)
		if got != tc.quoted {
			t.Errorf("Quote(%q) => %s, want %s", tc.text, got, tc.quoted)
		}
	}
}

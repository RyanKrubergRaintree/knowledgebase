package natural

import "testing"

func TestLess(t *testing.T) {
	cases := []struct{ A, B string }{
		{"A", "B"},
		{"A", "a"},
		{"9", "10"},
		{"x9", "x10"},
		{"x9", "x00010"},
		{"x09", "x10"},
		{"x09", "x009"},
	}

	for _, c := range cases {
		if !Less(c.A, c.B) || Less(c.B, c.A) {
			t.Errorf("error expected %s < %s", c.A, c.B)
		}
	}
}

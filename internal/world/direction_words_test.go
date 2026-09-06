package world

import "testing"

func TestParseDirection(t *testing.T) {
	cases := []struct {
		word   string
		want   Direction
		wantOK bool
	}{
		{"NORTH", North, true},
		{"north-east", NorthEast, true},
		{"South-West", SouthWest, true},
		{"APEX", 0, false},
	}
	for _, c := range cases {
		got, ok := ParseDirection(c.word)
		if ok != c.wantOK {
			t.Errorf("ParseDirection(%q) ok = %v, want %v", c.word, ok, c.wantOK)
			continue
		}
		if ok && got != c.want {
			t.Errorf("ParseDirection(%q) = %v, want %v", c.word, got, c.want)
		}
	}
}

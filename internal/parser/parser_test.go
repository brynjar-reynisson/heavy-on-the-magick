package parser

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		input string
		want  Command
	}{
		{"APEX, DOOR", Command{Target: "APEX", Verb: "DOOR", Raw: "APEX, DOOR"}},
		{"apex, fire", Command{Target: "APEX", Verb: "FIRE", Raw: "apex, fire"}},
		{"BLAST", Command{Verb: "BLAST", Raw: "BLAST"}},
		{"  DOOR, password  ", Command{Target: "DOOR", Verb: "PASSWORD", Raw: "DOOR, password"}},
		{"PICK UP GRIMOIRE", Command{Verb: "PICKUP", Target: "GRIMOIRE", Raw: "PICK UP GRIMOIRE"}},
		{"pick up grimoire", Command{Verb: "PICKUP", Target: "GRIMOIRE", Raw: "pick up grimoire"}},
		{"PICK UP", Command{Verb: "PICKUP", Target: "", Raw: "PICK UP"}},
	}
	for _, c := range cases {
		got := Parse(c.input)
		if got != c.want {
			t.Errorf("Parse(%q) = %+v, want %+v", c.input, got, c.want)
		}
	}
}

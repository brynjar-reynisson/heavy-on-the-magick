package parser

import "testing"

func TestExpandKeyword(t *testing.T) {
	cases := map[string]string{
		"N":    "NORTH",
		"n":    "NORTH",
		"NE":   "NORTH-EAST",
		"S":    "SOUTH",
		"SW":   "SOUTH-WEST",
		"B":    "BLAST",
		"F":    "FREEZE",
		"X":    "EXAMINE",
		"P":    "PICKUP",
		"D":    "DROP",
		"H":    "HALT",
		"I":    "INVOKE",
		"DOOR": "DOOR", // not a keyword abbreviation - an object name, unchanged
	}
	for in, want := range cases {
		if got := ExpandKeyword(in); got != want {
			t.Errorf("ExpandKeyword(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseActionFormWithAbbreviation(t *testing.T) {
	got := Parse("B CYCLOPS")
	want := Command{Target: "CYCLOPS", Verb: "BLAST", Raw: "B CYCLOPS"}
	if got != want {
		t.Errorf("Parse(%q) = %+v, want %+v", "B CYCLOPS", got, want)
	}
}

func TestParseBareAbbreviation(t *testing.T) {
	got := Parse("N")
	want := Command{Verb: "NORTH", Raw: "N"}
	if got != want {
		t.Errorf("Parse(%q) = %+v, want %+v", "N", got, want)
	}
}

func TestParseConversationFormUnaffectedByAbbreviations(t *testing.T) {
	// Names/objects in "name, object" form are always full words per the
	// manual - "N" here should NOT expand to "NORTH".
	got := Parse("APEX, N")
	want := Command{Target: "APEX", Verb: "N", Raw: "APEX, N"}
	if got != want {
		t.Errorf("Parse(%q) = %+v, want %+v (no abbreviation expansion in conversation form)", "APEX, N", got, want)
	}
}

package parser

import "testing"

func TestVocabularyLoaded(t *testing.T) {
	if len(Vocabulary) < 300 {
		t.Fatalf("len(Vocabulary) = %d, want at least 300 (extracted table has 316 words)", len(Vocabulary))
	}
}

func TestKnownWord(t *testing.T) {
	for _, w := range []string{"APEX", "AXIL", "DOOR", "NORTH", "north", "Fire"} {
		if !KnownWord(w) {
			t.Errorf("KnownWord(%q) = false, want true", w)
		}
	}
	if KnownWord("GOLANG") {
		t.Error("KnownWord(\"GOLANG\") = true, want false")
	}
}

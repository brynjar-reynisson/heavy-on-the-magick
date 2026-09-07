package parser

import (
	_ "embed"
	"encoding/json"
)

// vocabularyJSON holds the game's real parser vocabulary, extracted
// directly from decompressed game memory (RAM addresses 24270-26200 in the
// disassembly — see ../../CLAUDE.md, "Found the full parser vocabulary
// table"). The original stores words in fixed-length buckets (all 3-letter
// words together, then 4-letter, 5-letter, ... up to 11 letters) — a
// classic 8-bit space/parsing optimization. That bucket order is preserved
// here for reference, but Vocabulary exposes it as a simple set for lookup.
//
// IMPORTANT: these are the exact words the real game's parser recognizes,
// confirmed byte-for-byte from memory — not guessed or invented. What each
// word *means* (verb vs. noun vs. direction, and what game effect it has)
// is mostly NOT yet known; only the 8 compass directions are confirmed
// (see Directions below, cross-checked against the "8 compass directions"
// description of the game and the hyphenated NORTH-EAST/SOUTH-EAST/
// SOUTH-WEST/NORTH-WEST spelling found in the same table).
//
//go:embed data/vocabulary.json
var vocabularyJSON []byte

// Vocabulary is every word the original parser recognizes, in the order
// they appear in the game's own length-bucketed table.
//
// OPEN DISCREPANCY (round 85): the official instruction manual's own
// Merphish reference section lists "a few Merphish object names" -
// ASMODEE, ASTAROT, AXIL, BELEZBAR, BOOK, BOX, BOTTLE, LOAF, CANDLE,
// CHAIR, DEMON, MAGOT, OBJECT, TABLE, WALL, MONSTER, SWORD, ROCK, SIGN,
// RUBY - presented as real in-game object names, not hypothetical
// examples. Only RUBY (of BOX/BOTTLE/CANDLE/CHAIR/WALL) is actually
// present in this extracted 316-word table, despite all of them fitting
// comfortably within the confirmed 3-11 letter bucket range. Not
// resolved either way: could mean this extraction missed some real
// words (memory-scan limits), or that the manual's list is illustrative
// rather than exhaustive/literal. Recorded honestly rather than
// silently assumed either way.
var Vocabulary = mustLoadVocabulary()

func mustLoadVocabulary() []string {
	var words []string
	if err := json.Unmarshal(vocabularyJSON, &words); err != nil {
		// The embedded file is part of the build; a parse failure here is
		// a packaging bug, not a runtime condition callers should handle.
		panic("parser: malformed embedded vocabulary.json: " + err.Error())
	}
	return words
}

// KnownWord reports whether s is a word the original game's parser
// recognized (case-insensitive).
func KnownWord(s string) bool {
	_, ok := vocabularySet()[normalize(s)]
	return ok
}

func normalize(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

var vocabularySetCache map[string]struct{}

func vocabularySet() map[string]struct{} {
	if vocabularySetCache == nil {
		vocabularySetCache = make(map[string]struct{}, len(Vocabulary))
		for _, w := range Vocabulary {
			vocabularySetCache[w] = struct{}{}
		}
	}
	return vocabularySetCache
}

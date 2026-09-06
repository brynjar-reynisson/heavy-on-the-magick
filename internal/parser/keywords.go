package parser

import "strings"

// merphishKeywords maps the game's real single/double-letter command
// abbreviations to their expanded form. Source: the official Gargoyle
// Games instruction manual ("Heavy on the Magick", 1986, Carter Follis
// Software) — confirmed real via its own explicit keyword reference table
// ("N North", "NE North-east", ... "B – (Object) Blast the named object or
// monster", etc.). This is a significant correction: earlier assumptions
// in this project used only full words (e.g. "NORTH"); the manual
// confirms the real game expects these abbreviations as primary input,
// expanding them "on output" — full words are also accepted (confirmed by
// the "APEX, DOOR"-style examples in the in-game hint screen, which spell
// directions/verbs out fully), so ExpandKeyword supports both without
// breaking anything already working.
//
// Facts only, not copied prose: the manual's exact wording is not
// reproduced here (its text is explicitly copyrighted, 1986 Carter Follis
// Software) — only the keyword-to-meaning mapping, which is factual
// interface documentation, not creative expression.
var merphishKeywords = map[string]string{
	"N":  "NORTH",
	"NE": "NORTH-EAST",
	"NW": "NORTH-WEST",
	"S":  "SOUTH",
	"SE": "SOUTH-EAST",
	"SW": "SOUTH-WEST",
	"E":  "EAST",
	"W":  "WEST",
	"L":  "LEFT",
	"R":  "RIGHT",
	"H":  "HALT",
	"Z":  "SWAP",
	"O":  "OPTIONS",
	"X":  "EXAMINE",
	"P":  "PICKUP",
	"D":  "DROP",
	"I":  "INVOKE",
	"B":  "BLAST",
	"F":  "FREEZE",
}

// ExpandKeyword expands a real Merphish single/double-letter command
// abbreviation to its full word (e.g. "N" -> "NORTH", "B" -> "BLAST").
// Words not in the abbreviation table (including ordinary full words and
// object/character names, which the manual confirms must always be typed
// in full) are returned unchanged.
func ExpandKeyword(word string) string {
	if full, ok := merphishKeywords[strings.ToUpper(word)]; ok {
		return full
	}
	return strings.ToUpper(word)
}

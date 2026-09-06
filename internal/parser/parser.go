// Package parser handles the game's real command language, "Merphish" —
// confirmed directly from the official Gargoyle Games instruction manual
// (see keywords.go's doc comment for the sourcing/copyright note) plus the
// in-game startup hint screen (see ../../CLAUDE.md, "Found ground truth:
// ref_screenshot.png").
//
// Two real, confirmed input forms:
//
//  1. Conversation: `"name, object"` (the manual's own confirmed format) —
//     e.g. `APEX, DOOR`, `ASTAROT, WOLFDORP`. Both name and object must be
//     typed in full (the manual is explicit: object/character names are
//     never abbreviated). Maps to Command{Target: name, Verb: object}.
//
//  2. Action: `Keyword (Object)` — e.g. `B CYCLOPS` (Blast the Cyclops),
//     `X BOTTLE` (eXamine the bottle), or a bare keyword with no object
//     (`N`, `BLAST`). Keywords are usually 1-2 letter abbreviations,
//     expanded via ExpandKeyword (see keywords.go) — full words work too.
//     Maps to Command{Verb: expanded keyword, Target: object, if any}.
//
// STATUS: the real in-game vocabulary table (parser.Vocabulary, 316 words)
// is confirmed extracted from memory, and the manual confirms the two
// grammar forms above precisely, but the Z80 routine that actually
// tokenizes input at runtime has still not been traced in the
// disassembly — this package's behavior is built from these two
// independent confirmed sources, not from the original code directly.
package parser

import "strings"

// Command is a parsed player input.
type Command struct {
	Target string // conversation: who's addressed. action: the object, if any.
	Verb   string // conversation: the topic/object. action: the (expanded) keyword.
	Raw    string // the original input, for echoing/logging
}

// Parse interprets raw player input as either the conversation form
// (`"name, object"`) or the action form (`Keyword (Object)`), per the
// real confirmed Merphish grammar — see the package doc comment.
func Parse(input string) Command {
	raw := strings.TrimSpace(input)

	// "PICK UP" is the real two-word vocabulary entry for this verb (see
	// parser.Vocabulary - it's stored as one dictionary entry containing a
	// space, not two separate words), which the generic single-space-cut
	// action-form logic below can't handle correctly (it would treat "PICK"
	// as the keyword and "UP <object>" as the object). Special-cased here
	// so both "PICK UP" (real) and "PICKUP" (this project's existing
	// one-word convention, already wired to the same handler) work.
	if upper := strings.ToUpper(raw); strings.HasPrefix(upper, "PICK UP") {
		return Command{
			Verb:   "PICKUP",
			Target: strings.TrimSpace(upper[len("PICK UP"):]),
			Raw:    raw,
		}
	}

	// Conversation form: "name, object". Names/objects are always typed in
	// full here (per the manual), so no keyword expansion applies.
	if target, verb, ok := strings.Cut(raw, ","); ok {
		return Command{
			Target: strings.ToUpper(strings.TrimSpace(target)),
			Verb:   strings.ToUpper(strings.TrimSpace(verb)),
			Raw:    raw,
		}
	}

	// Action form: "Keyword Object" or a bare keyword. The keyword may be
	// a real Merphish abbreviation (e.g. "B", "NE") and gets expanded; the
	// object (if any) is left as-is, since object names are always full
	// words already.
	if keyword, object, ok := strings.Cut(raw, " "); ok {
		return Command{
			Verb:   ExpandKeyword(keyword),
			Target: strings.ToUpper(strings.TrimSpace(object)),
			Raw:    raw,
		}
	}
	return Command{
		Verb: ExpandKeyword(raw),
		Raw:  raw,
	}
}

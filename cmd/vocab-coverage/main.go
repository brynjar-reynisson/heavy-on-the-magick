// Command vocab-coverage precisely quantifies how much of the game's
// real, extracted 316-word parser vocabulary (internal/parser.Vocabulary
// — 313 unique words; 3 words repeat across length buckets, see below)
// has modeled behavior in internal/game.Game.Handle, versus how many
// still fall through to the generic "I recognize that word, but don't
// know what it does yet" stub.
//
// This exists because "most verbs are unimplemented" is a vague, easy
// claim to repeat but a hard one to act on — most of Vocabulary is
// actually NOUN content (room/item/monster/demon names, already used
// throughout world/*.go and magic/demons.go), not unimplemented verbs.
// Running this tool turns the vague complaint into an exact number and
// an exact list, self-verified by actually calling the real Handle
// function for each word (a fresh game.New() per word, so no ordering
// side effects from movement/combat/stamina bleed between checks) —
// not a hand-maintained mirror list that could silently drift out of
// sync with game.go's actual switch cases.
//
// KNOWN BLIND SPOT, honestly documented rather than silently wrong: 14
// real vocabulary words (see targetPositionWords below) only work in a
// specific TARGET-VERB pairing — some (APEX, ASTAROT, MAGOT, ASMODEE,
// BELEZBAR, GUARDS, DOOR, NEST, CAULDRON) are recognized as cmd.Target
// paired with a specific companion verb; others (TALK, SPEAK, THANKS,
// PHOENIX, ACHAD) are the reverse — recognized as cmd.Verb only paired
// with a specific companion Target. This tool only calls Handle with
// each word alone as cmd.Verb (empty Target), the dominant pattern for
// the other 31 modeled words, so it can't detect these 14 through that
// single check. Rather than either miscount them as "unimplemented" or
// build fragile per-word position-pairing logic to chase every
// combination, they're explicitly excluded and listed separately.
//
// ROUND 165 FIX (the same undercounting class round 155 already fixed
// once for a different cause): ASMODEE and BELEZBAR gained real
// cmd.Target-position dispatch in rounds 162/164 (game.asmodeeDestroy,
// game.belezbarReveal — the same shape as the already-excluded
// ASTAROT/MAGOT) but were never added here, so this tool was silently
// counting 2 real, working, tested commands as "unimplemented" for 2-3
// rounds without anyone noticing — exactly the "self-audit tool needs
// the same maintenance as the game code" lesson round 140 already
// recorded, recurring because a NEW demon command is easy to add
// without remembering this tool exists to update too.
//
// ROUND 155 FIX (was a real, silent undercounting bug, not previously
// documented as a blind spot): this tool used to build a bare
// parser.Command{Verb: w} directly, bypassing parser.Parse entirely.
// That's wrong for the vocabulary's one real multi-word entry, "PICK
// UP" (parser.Vocabulary — stored as a single dictionary entry
// containing a space, see parser.Parse's own doc comment) — Handle only
// recognizes the normalized Verb "PICKUP", produced by Parse's own
// special-case for this exact entry, so calling Handle with the raw,
// un-parsed "PICK UP" string always hit the generic stub, silently
// miscounting a real, working, already-tested command as unimplemented.
// Fixed by routing every word through the real parser.Parse (what an
// actual player's input goes through) instead of hand-building a
// Command — strictly more honest as a testing methodology too, and a
// no-op for every other word (parser.ExpandKeyword passes through any
// unrecognized ≥3-letter word unchanged, and no other vocabulary entry
// contains a space).
package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/game"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
)

// targetPositionWords are real vocabulary words Handle only recognizes
// as cmd.Target, paired with the listed companion verb(s) — see this
// file's package doc comment.
var targetPositionWords = map[string]string{
	"APEX":     "recognized as cmd.Target, paired with verb TALK/SPEAK/THANKS",
	"ASTAROT":  "recognized as cmd.Target, paired with any non-empty verb (a location name)",
	"MAGOT":    "recognized as cmd.Target, paired with any non-empty verb (an object name)",
	"GUARDS":   "recognized as cmd.Target, paired with verb DOOR",
	"DOOR":     "recognized as cmd.Target, verb compared against the current room's real DoorPasswords/TollItem",
	"TALK":     "recognized as cmd.Verb, only when paired with cmd.Target APEX",
	"SPEAK":    "recognized as cmd.Verb, only when paired with cmd.Target APEX",
	"THANKS":   "recognized as cmd.Verb, only when paired with cmd.Target APEX",
	"NEST":     "recognized as cmd.Target, paired with verb PHOENIX (round 136)",
	"PHOENIX":  "recognized as cmd.Verb, only when paired with cmd.Target NEST (round 136)",
	"CAULDRON": "recognized as cmd.Target, paired with verb ACHAD (round 139)",
	"ACHAD":    "recognized as cmd.Verb, only when paired with cmd.Target CAULDRON (round 139)",
	"ASMODEE":  "recognized as cmd.Target, paired with any non-empty verb (an object name to destroy, round 162)",
	"BELEZBAR": "recognized as cmd.Target, paired with any non-empty verb (an object name to reveal, round 164)",
}

func uniqueSorted(words []string) []string {
	set := make(map[string]bool, len(words))
	for _, w := range words {
		set[w] = true
	}
	out := make([]string, 0, len(set))
	for w := range set {
		out = append(out, w)
	}
	sort.Strings(out)
	return out
}

func isGenericResponse(resp string) bool {
	return strings.HasPrefix(resp, "I recognize that word, but don't know what it does yet") ||
		resp == "I don't understand that word."
}

func main() {
	words := uniqueSorted(parser.Vocabulary)

	var uncovered []string
	for _, w := range words {
		if _, blindSpot := targetPositionWords[w]; blindSpot {
			continue
		}
		g := game.New()
		resp := g.Handle(parser.Parse(w))
		if isGenericResponse(resp) {
			uncovered = append(uncovered, w)
		}
	}

	checked := len(words) - len(targetPositionWords)
	covered := checked - len(uncovered)

	fmt.Printf("Vocabulary: %d unique words (316 total table entries; IRON/LOOKS/MANTIS repeat across length buckets).\n", len(words))
	fmt.Printf("Excluded from this check (see targetPositionWords - recognized via cmd.Target, not cmd.Verb): %d\n", len(targetPositionWords))
	fmt.Printf("Checked as cmd.Verb: %d — of those, %d have modeled Handle behavior, %d still fall through to the generic stub.\n", checked, covered, len(uncovered))
	fmt.Println()
	fmt.Println("Words with NO modeled behavior yet (cmd.Verb check only):")
	for _, w := range uncovered {
		fmt.Println(" ", w)
	}
}

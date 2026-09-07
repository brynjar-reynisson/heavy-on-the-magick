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
// KNOWN BLIND SPOT, honestly documented rather than silently wrong: 8
// real vocabulary words (see targetPositionWords below) are only ever
// recognized by Handle when they appear as cmd.Target, paired with a
// SPECIFIC companion verb (e.g. "APEX" is only recognized alongside
// TALK/SPEAK/THANKS) — this tool only calls Handle with each word as
// cmd.Verb (empty Target), the dominant pattern for the other ~22
// modeled words, so it can't detect these 8 through that single check.
// Rather than either miscount them as "unimplemented" or build fragile
// per-word position-pairing logic to chase every combination, they're
// explicitly excluded and listed separately.
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
	"APEX":    "recognized as cmd.Target, paired with verb TALK/SPEAK/THANKS",
	"ASTAROT": "recognized as cmd.Target, paired with any non-empty verb (a location name)",
	"MAGOT":   "recognized as cmd.Target, paired with any non-empty verb (an object name)",
	"GUARDS":  "recognized as cmd.Target, paired with verb DOOR",
	"DOOR":    "recognized as cmd.Target, verb compared against the current room's real DoorPasswords/TollItem",
	"TALK":    "recognized as cmd.Verb, only when paired with cmd.Target APEX",
	"SPEAK":   "recognized as cmd.Verb, only when paired with cmd.Target APEX",
	"THANKS":  "recognized as cmd.Verb, only when paired with cmd.Target APEX",
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
		resp := g.Handle(parser.Command{Verb: w})
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

package game

import (
	"fmt"
	"strings"
)

// talkToApex handles the confirmed real conversation-form NPC
// interaction "APEX, TALK" (source: Spectrum Computing's plain-text
// instructions file — see ../../CLAUDE.md). Apex the Ogre is a helpful
// NPC who offers puzzle hints "if treated respectfully"; the precise
// hint content and what counts as respectful aren't extracted, so this
// is an honest acknowledgement of the real NPC rather than fabricated
// hint text.
func (g *Game) talkToApex() string {
	return "Apex the Ogre eyes you warily, then grunts. He might share what he knows, if you treat him with respect."
}

// apexDoorHint handles "APEX, DOOR" — one of the 3 confirmed real
// example commands from the game's own hint screen (game.help; also
// independently confirmed by Spectrum Computing's instructions file:
// "a door with a toll sign by it (ask apex)"), but never actually wired
// to a response until now, despite the hint screen naming it directly.
// If the current room has real, sourced riddle text (world.Room.
// DoorHints — round 159, see its own doc comment for sourcing), Apex
// shares it; otherwise this falls back to the same generic
// talkToApex response, an honest "he has nothing specific to say about
// a door here" rather than a fabricated hint.
func (g *Game) apexDoorHint() string {
	room := g.World.CurrentRoom()
	if room == nil || len(room.DoorHints) == 0 {
		return g.talkToApex()
	}
	var b strings.Builder
	b.WriteString("Apex leans in and whispers a riddle: ")
	for i, hint := range room.DoorHints {
		if i > 0 {
			b.WriteString(" ... ")
		}
		fmt.Fprintf(&b, "%q", hint)
	}
	return b.String()
}

// call handles the CALL spell (round 125) — confirmed for the first
// time by a direct first-hand account of actually playing the game
// (The CRPG Addict's 2016 blog post on "Heavy on the Magick"): "Later,
// you find some additional spells, including... CALL, which allows
// you to summon an annoying NPC... You can CALL him once you get the
// spell scroll". "Him" is Apex — the same post's next sentence
// describes needing "APEX, THANKS" to banish him again, the exact
// real dismiss phrase this port already had wired since the HELP
// round. This closes the single most-repeated "confirmed real spell,
// effect unknown" gap in the whole project (open since round 87).
// Requires the Scroll, matching the numbered map's own key list tying
// CALL to a Scroll in the Trollwynd zone (already real, placed data) —
// consistent with the new source's "once you get the spell scroll".
func (g *Game) call() string {
	if !g.hasItem("Scroll") {
		return "You don't have the spell Scroll needed to CALL."
	}
	return "You CALL out... " + g.talkToApex()
}

// "APEX, THANKS" (handled inline in Handle) is the hint screen's own
// confirmed real dismiss phrase — "To dismiss say \"APEX, THANKS\"" (see
// game.help's verbatim text) — quoted in this project since the HELP
// round but never wired to anything until now.

package game

import (
	"fmt"
	"strings"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/world"
)

// passGuards handles the confirmed real command "GUARDS, DOOR" (source:
// the game's own in-game hint screen, cross-checked against
// ref_screenshot.png — see ../../CLAUDE.md and world.Room.Guards's doc
// comment). No precondition (payment, password) for succeeding is
// confirmed by any source found so far, so this models the simplest
// honest reading of the confirmed example: saying it clears a real
// Guards obstacle in the current room.
func (g *Game) passGuards() string {
	room := g.World.CurrentRoom()
	if room == nil || !room.Guards {
		return "There are no guards here."
	}
	room.Guards = false
	return "The guards step aside and let you pass."
}

// passWater handles the confirmed real command "WATER, FALL" (round
// 169; source: CRASH magazine issue 31's "Signpost" adventure column,
// crashonline.org.uk/31/signpost.htm - "To get past the water say
// 'Water, fall'"). The same conversation-form-obstacle shape as
// passGuards (a spoken command, not a carried item like Fire/Clasp),
// clearing a real Water hazard - see world.Room.Water's doc comment
// for the exact-cell sourcing (Level3Grid's H4, already independently
// confirmed as literally named "Water"). The response text is the
// game's own exact, confirmed word - "Trickle" - per a direct
// frame-by-frame review of real gameplay footage (round 178).
func (g *Game) passWater() string {
	room := g.World.CurrentRoom()
	if room == nil || !room.Water {
		return "There is no water here to command."
	}
	room.Water = false
	return "Trickle."
}

// payToll handles a real, distinct door mechanic (see world.Room.TollItem's
// doc comment): unlike a typed DoorPasswords password, a Toll door
// requires actually carrying and spending the named item (the confirmed
// source: "a bag of gold is the key (put it on the table)") - taken out
// of the player's inventory on success, not just checked for.
//
// ROUND 181 correction: the item is moved to the room's own real
// Items (a table it's physically placed ON, per the confirmed "put it
// on the table" wording), not deleted from the game entirely as this
// port previously modeled it. This was found to matter for real,
// concrete reasons once movement blocking was fixed the same round:
// Room of Arrows' TollItem "Slat" and game.checkSlatCyclops's own
// "the Slat kills the Cyclops" mechanic (round 146) both need the
// SAME real, sourced Slat - deleting it at the toll would make the
// already-verified Morfang-to-Nidus path impossible. Leaving it
// retrievable on the table lets a player pay the toll, then pick the
// same Slat back up and carry it onward - both real mechanics stay
// reachable, honestly, without inventing a second Slat no source
// confirms exists.
func (g *Game) payToll(room *world.Room) string {
	for i, item := range g.Player.Items {
		if strings.EqualFold(item, room.TollItem) {
			g.Player.Items = append(g.Player.Items[:i], g.Player.Items[i+1:]...)
			room.Items = append(room.Items, item)
			// ROUND 181: clears the real lock this room's own TollItem
			// represents - see game.move's doc comment (fixing a real
			// bug: before this round, a Toll door never actually
			// blocked movement at all).
			room.TollItem = ""
			return fmt.Sprintf("You place the %s on the table. The door swings open.", item)
		}
	}
	return fmt.Sprintf("You don't have a %s to pay the toll.", room.TollItem)
}

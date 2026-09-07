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

// payToll handles a real, distinct door mechanic (see world.Room.TollItem's
// doc comment): unlike a typed DoorPasswords password, a Toll door
// requires actually carrying and spending the named item (the confirmed
// source: "a bag of gold is the key (put it on the table)") - so the
// item is removed from the player's inventory on success, not just
// checked for.
func (g *Game) payToll(room *world.Room) string {
	for i, item := range g.Player.Items {
		if strings.EqualFold(item, room.TollItem) {
			g.Player.Items = append(g.Player.Items[:i], g.Player.Items[i+1:]...)
			return fmt.Sprintf("You place the %s on the table. The door swings open.", item)
		}
	}
	return fmt.Sprintf("You don't have a %s to pay the toll.", room.TollItem)
}

package game

import "strings"

// hasItem reports whether the player is carrying an item by name.
func (g *Game) hasItem(name string) bool {
	return g.Player.HasItem(name)
}

// roomHasItem reports whether the current room's real world.Room.Items
// contains the named item (case-insensitive, matching hasItem's
// convention) - used by invoke (round 131) to check whether a Talisman
// has actually been placed on the ground, not merely carried.
func (g *Game) roomHasItem(name string) bool {
	room := g.World.CurrentRoom()
	if room == nil {
		return false
	}
	for _, item := range room.Items {
		if strings.EqualFold(item, name) {
			return true
		}
	}
	return false
}

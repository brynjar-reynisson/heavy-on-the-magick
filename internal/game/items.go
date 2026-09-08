package game

import (
	"fmt"
	"strings"
)

// poisonPickupStaminaCost is charged when the player picks up a real,
// named "poison" item - round 128's find, a genuinely NEW source type
// for this project (a 1986 CRASH magazine review,
// crashonline.org.uk/29/magick.htm): "Poison damages Stamina upon
// contact." The exact amount isn't stated, so this is an honest
// placeholder, matching the same convention as combatStaminaCost/
// saveStaminaCost - deliberately between the two (worse than a save,
// milder than a full combat hit, since it's a single incidental
// contact rather than an ongoing fight). Room of Misery's already-real,
// already-placed "Poison-smeared book" (round 53) is the one item in
// this port whose name confirms it as poisonous - "contact" is read as
// picking it up, the natural point of contact for an item, not a
// separate required action no source describes.
const poisonPickupStaminaCost = 3

// notFoodItems are real items the game specifically calls out as "IT'S
// NOT FOOD" on pickup - confirmed via a direct frame-by-frame review of
// a full walkthrough video for Nougat and Garlic, two already-
// real, already-placed items whose names sound edible (a real, small
// joke the original makes, not something to lose in the port). A map
// (not a hardcoded pair), matching this project's own established
// convention (see spellRequiresItem) for a fact confirmed on a small
// set of items that a future round might extend. Round 180: a fourth,
// finer-grained (2-second interval) pass over the same video found the
// same exact response for a real Egg pickup at Wraithvale too - added.
var notFoodItems = map[string]bool{
	"Nougat": true,
	"Garlic": true,
	"Egg":    true,
}

// inventory lists the player's carried items. "INVENTORY" is a real,
// confirmed vocabulary word (see parser.Vocabulary) with an obvious,
// standard adventure-game meaning; no source states its exact wording,
// so the response text is this project's own, not extracted.
func (g *Game) inventory() string {
	if len(g.Player.Items) == 0 {
		return "You aren't carrying anything."
	}
	return "You are carrying: " + strings.Join(g.Player.Items, ", ")
}

// pickup moves a named item from the current room's Items into the
// player's inventory, if present.
func (g *Game) pickup(target string) string {
	if target == "" {
		return "Pick up what?"
	}
	room := g.World.CurrentRoom()
	if room == nil {
		return "There's nothing here to pick up."
	}
	for i, item := range room.Items {
		if strings.EqualFold(item, target) {
			room.Items = append(room.Items[:i], room.Items[i+1:]...)
			g.Player.Items = append(g.Player.Items, item)
			result := fmt.Sprintf("You pick up the %s.", item)
			if strings.Contains(strings.ToLower(item), "poison") {
				g.Player.Stamina -= poisonPickupStaminaCost
				result = g.deathCheck(result + " It's poisonous to the touch! You feel your strength ebb.")
			}
			if notFoodItems[item] {
				result += " It's not food."
			}
			return result
		}
	}
	return fmt.Sprintf("There's no %s here to pick up.", strings.ToLower(target))
}

// drop moves a named item from the player's inventory into the current
// room's Items, if the player is carrying it. If the dropped item is
// exactly this room's real world.Room.TollItem, it instead pays the
// toll and opens the door (see payToll) - a fresh, more detailed CASA
// walkthrough re-read (round 64) found the real trigger phrase is
// literally "DROP <item>" ("EXAMINE TABLE, DROP KEY (door opens)"),
// matching the instructions file's "put it on the table" wording far
// more literally than this port's original "DOOR, <item>" guess did.
// Both forms are kept working (see the DOOR-target handling above)
// since nothing confirms the guessed form is wrong, just that DROP is
// also, and probably primarily, real.
func (g *Game) drop(target string) string {
	if target == "" {
		return "Drop what?"
	}
	if room := g.World.CurrentRoom(); room != nil && room.TollItem != "" && strings.EqualFold(target, room.TollItem) {
		if g.hasItem(room.TollItem) {
			return g.payToll(room)
		}
	}
	for i, item := range g.Player.Items {
		if strings.EqualFold(item, target) {
			g.Player.Items = append(g.Player.Items[:i], g.Player.Items[i+1:]...)
			if room := g.World.CurrentRoom(); room != nil {
				room.Items = append(room.Items, item)
			}
			result := fmt.Sprintf("You drop the %s.", item)
			if msg := g.checkNougatWerewolf(); msg != "" {
				result += "\n" + msg
			}
			if msg := g.checkGarlicVampire(); msg != "" {
				result += "\n" + msg
			}
			if msg := g.checkPelletSlug(); msg != "" {
				result += "\n" + msg
			}
			if msg := g.checkSnakeHydra(); msg != "" {
				result += "\n" + msg
			}
			if msg := g.checkSlatCyclops(); msg != "" {
				result += "\n" + msg
			}
			if msg := g.checkMirrorMedusa(); msg != "" {
				result += "\n" + msg
			}
			if msg := g.checkSwapItem(item); msg != "" {
				result += "\n" + msg
			}
			return result
		}
	}
	return fmt.Sprintf("You aren't carrying a %s.", strings.ToLower(target))
}

// examine reports a room's Monster and Items if present. With a
// specific target (the manual's confirmed real grammar, e.g. "X BOTTLE"
// - see parser's doc comment), it instead confirms just that one thing
// if it's actually here (the room's Monster, a room Item, or a carried
// Item) or says plainly that it isn't - real target-aware behavior
// matching the confirmed grammar, not just always listing everything
// regardless of what was asked about.
func (g *Game) examine(target string) string {
	room := g.World.CurrentRoom()
	if room == nil {
		return "There's nothing to examine."
	}
	if target != "" {
		if room.HasTable && strings.EqualFold(target, "TABLE") {
			return "A plain table."
		}
		if room.HasChest && strings.EqualFold(target, "CHEST") {
			// Wolfdorp's and Gorburg's chests are both confirmed real
			// oak (round 179: "IT'S A CHEST MADE OF OAK...", a direct
			// frame-by-frame video review of each) - Morfang's own
			// chest material isn't confirmed by any source, so it keeps
			// the honest generic wording rather than assuming the same
			// material.
			if strings.EqualFold(room.Name, "Wolfdorp") || strings.EqualFold(room.Name, "Gorburg") {
				return "A wooden chest, made of oak."
			}
			return "A wooden chest."
		}
		if room.HasCauldron && strings.EqualFold(target, "CAULDRON") {
			if g.roomHasItem("Scroll") {
				return "Cold iron: it holds a Scroll."
			}
			return "A cold iron cauldron."
		}
		if room.Monster != "" && room.MonsterHealth > 0 && strings.EqualFold(room.Monster, target) {
			return fmt.Sprintf("You see a %s.", room.Monster)
		}
		for _, item := range room.Items {
			if strings.EqualFold(item, target) {
				return fmt.Sprintf("You see a %s.", item)
			}
		}
		for _, item := range g.Player.Items {
			if strings.EqualFold(item, target) {
				return fmt.Sprintf("You are carrying a %s.", item)
			}
		}
		return "You don't see that here."
	}
	var parts []string
	if room.Monster != "" && room.MonsterHealth > 0 {
		parts = append(parts, fmt.Sprintf("a %s", room.Monster))
	}
	parts = append(parts, room.Items...)
	if len(parts) == 0 {
		return "You see nothing of particular note."
	}
	return fmt.Sprintf("You see: %s.", strings.Join(parts, ", "))
}

// checkSwapItem implements a real, sourced "protected item" mechanic
// (round 132) - see world.Room.SwapItem/RevealItem's doc comment for
// the full sourcing (World of Spectrum's separate plain-text
// instructions file, cross-confirmed by the numbered map poster's own
// "protected" qualifier on the exact same 3 items). Dropping the
// room's real SwapItem reveals its RevealItem, added to the room's
// Items for real - a one-time reveal (both fields are cleared after,
// so dropping the same item twice doesn't reveal it again). No source
// states whether the dropped SwapItem itself is consumed - modeled as
// non-consuming, the same "no evidence it's used up" convention
// checkNougatWerewolf already uses for its own dropped item.
func (g *Game) checkSwapItem(dropped string) string {
	room := g.World.CurrentRoom()
	if room == nil || room.SwapItem == "" || room.RevealItem == "" {
		return ""
	}
	if !strings.EqualFold(dropped, room.SwapItem) {
		return ""
	}
	revealed := room.RevealItem
	room.Items = append(room.Items, revealed)
	room.SwapItem = ""
	room.RevealItem = ""
	return fmt.Sprintf("As you set it down, you notice a %s hidden nearby!", revealed)
}

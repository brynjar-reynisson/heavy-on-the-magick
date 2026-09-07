package game

import (
	"fmt"
	"strings"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/magic"
)

// invoke handles INVOKE (Merphish keyword "I"), given the demon name the
// player named as cmd.Target. Recognizing the 4 confirmed magic.Demons,
// requiring each one's specific Charm, and each one's Ability are all
// real, sourced facts (see magic.Demons's doc comment) — invocation now
// actually succeeds when the right Charm is present. With no target,
// lists all 4 demons, their Charm requirements, and (round 114) their
// real occult Correspondences (magic.Demon.Correspondences, extracted
// round 110 from the manual's grimoire section but never actually shown
// anywhere in-game until now) — real, already-sourced data that had no
// way to be seen in-game before this. Round 115 closes the same gap for
// Number/Sign/Aspect, which were ALSO real, confirmed manual facts
// (present on every Demon since this project's very first magic.Demons
// commit) but — unlike Ability/Charm/Correspondences — had never been
// referenced by internal/game at all, found the same way round 114 was:
// checking each struct field against what actually reaches the player,
// not assuming "it's a field, so it must be used somewhere."
//
// Round 131 CORRECTION: the Charm gate previously checked g.hasItem
// (carried in inventory) — but a genuinely new source (World of
// Spectrum's separate plain-text instructions file, distinct from the
// PDF manual this project had already mined heavily) states the real
// mechanic precisely: "Place Ye the talisman on the ground and proceed
// with thy invocation from a distance." The Charm must be DROPPED in
// the room, not merely carried. Fixed to check g.roomHasItem instead;
// a player who IS carrying the Charm but hasn't dropped it gets an
// honest, specific hint (not punishFailedInvoke's furnace-room
// punishment, which is reserved for not having the Charm at all).
func (g *Game) invoke(target string) string {
	if target == "" {
		var b strings.Builder
		b.WriteString("Known demons and their required Talismans:\n")
		for _, d := range magic.Demons {
			fmt.Fprintf(&b, "%s, %s (Number %d, %s, Aspect: %s; needs: %s) - %s\n", d.Name, d.Title, d.Number, d.Sign, d.Aspect, d.Charm, d.Correspondences)
		}
		return strings.TrimRight(b.String(), "\n")
	}
	for _, d := range magic.Demons {
		if !strings.EqualFold(d.Name, target) {
			continue
		}
		if !g.roomHasItem(d.Charm) {
			if g.hasItem(d.Charm) {
				return fmt.Sprintf("You are carrying the %s, but that isn't enough - place it on the ground and stand back before you invoke %s.", d.Charm, d.Name)
			}
			return g.punishFailedInvoke(d)
		}
		if d.Ability != "" {
			return fmt.Sprintf("You invoke %s, %s! %s.", d.Name, d.Title, d.Ability)
		}
		return fmt.Sprintf("You invoke %s, %s! The ritual succeeds, though its exact effect isn't modeled yet.", d.Name, d.Title)
	}
	return "There is no demon by that name."
}

// punishFailedInvoke handles invoking a demon without its Charm - a
// real, confirmed punishment, not just a rejection message. The CRPG
// Addict's first-hand playthrough account (the same source that
// confirmed CALL's effect, round 125): "When you INVOKE them, you have
// to be holding their particular talisman--found within the
// dungeon--or they send you to a furnace room with no exits." If the
// active World has a real Furnace Room (world.CollodonsPile does,
// added round 126; so does world.Level1Grid, independently, as its
// own isolated A8 cell), the player is genuinely teleported there,
// same mechanism as game.astarotTeleport. Worlds without one (Level2-
// 4Grid) fall back to the plain rejection message rather than fail -
// an honest scope limit, not a fabricated destination.
func (g *Game) punishFailedInvoke(d magic.Demon) string {
	msg := fmt.Sprintf("You begin the ritual to invoke %s, %s... but you have no suitable Talisman (a %s). The ritual backfires!", d.Name, d.Title, d.Charm)
	if id, ok := g.World.FindRoomByName("Furnace Room"); ok {
		g.World.Teleport(id)
		msg += " You are flung into a furnace room with no exits."
	}
	return msg
}

// astarotTeleport handles the confirmed real conversation-form command
// "ASTAROT, <location>" - the in-game hint screen's own literal example
// is "ASTAROT, WOLFDORP" (see game.help's doc comment and
// parser.Parse's package doc comment, which has carried this exact
// example since the two-grammar-form writeup). magic.Demons's Astarot
// entry already confirms the underlying ability ("Transports the player
// to a named location, if its name is known") and Charm ("Sword") - this
// is the first place that ability is actually implemented, using the
// same Charm-gating convention as bare INVOKE (round 131: the Charm
// must be on the ground, per World of Spectrum's plain-text
// instructions file - see invoke's doc comment). The location name is
// looked up against the CURRENT world's real Room names (Wolfdorp itself
// is a real, already-shipped CollodonsPile room), so this only reaches
// places that genuinely exist in whichever world is active - no
// fabricated destinations.
func (g *Game) astarotTeleport(location string) string {
	const astarotName, astarotTitle, astarotCharm = "Astarot", "the Spirit of Assemblage", "Sword"
	if !g.roomHasItem(astarotCharm) {
		if g.hasItem(astarotCharm) {
			return fmt.Sprintf("You are carrying the %s, but that isn't enough - place it on the ground and stand back before you invoke %s.", astarotCharm, astarotName)
		}
		return fmt.Sprintf("You call out to %s, %s... but you have no suitable Talisman (a %s).", astarotName, astarotTitle, astarotCharm)
	}
	id, ok := g.World.FindRoomByName(location)
	if !ok {
		return fmt.Sprintf("%s doesn't recognize a place called %q.", astarotName, location)
	}
	dest := g.World.Rooms[id]
	g.World.Teleport(id)
	return fmt.Sprintf("You invoke %s, %s! In an instant, you are transported to %s.", astarotName, astarotTitle, dest.Name)
}

// magotLocate handles the conversation-form command "MAGOT, <object>".
// No source gives a literal "MAGOT, X" example the way the hint screen
// gives "ASTAROT, WOLFDORP" - but the manual's own confirmed grammar
// ("name, object") is general, not restricted to the 2 demons it
// happens to illustrate, and magic.Demons's Magot entry already
// confirms both the underlying ability ("Reveals the whereabouts of any
// named object") and Charm ("Sunflower"). Applying the same
// TARGET-comma-VERB grammar and Charm-gating convention already
// established for Astarot's teleport is a natural, honestly-flagged
// inference, not a fabricated mechanic - the same confidence level this
// project already applies to vocabulary synonyms like TAKE/LIFT.
// Checks the player's own inventory first (an item they're already
// carrying isn't "located" anywhere else), then every room in the
// CURRENT world - so, like Astarot's teleport, this only ever reports
// real placements that genuinely exist in whichever world is active.
// Echoes the item's real stored casing in both branches, not the
// player's raw uppercased typed target - the same casing-honesty fix
// already applied once before to examine().
func (g *Game) magotLocate(object string) string {
	const magotName, magotTitle, magotCharm = "Magot", "the Diviner", "Sunflower"
	if !g.roomHasItem(magotCharm) {
		if g.hasItem(magotCharm) {
			return fmt.Sprintf("You are carrying the %s, but that isn't enough - place it on the ground and stand back before you invoke %s.", magotCharm, magotName)
		}
		return fmt.Sprintf("You call out to %s, %s... but you have no suitable Talisman (a %s).", magotName, magotTitle, magotCharm)
	}
	for _, item := range g.Player.Items {
		if strings.EqualFold(item, object) {
			return fmt.Sprintf("You invoke %s, %s! No need - you already carry the %s yourself.", magotName, magotTitle, item)
		}
	}
	for _, room := range g.World.Rooms {
		for _, item := range room.Items {
			if strings.EqualFold(item, object) {
				return fmt.Sprintf("You invoke %s, %s! The %s lies in %s.", magotName, magotTitle, item, room.Name)
			}
		}
	}
	return fmt.Sprintf("You invoke %s, %s! %s senses no such object anywhere nearby.", magotName, magotTitle, magotName)
}

// asmodeeDestroy handles the conversation-form command "ASMODEE,
// <object>" (round 162). Asmodee's Ability sat as just a warning for
// many rounds ("Warning: be careful with Asmodee... no confirmed
// positive effect") — a 4th, genuinely different source (Hardcore
// Gaming 101's article on this game) states plainly "Asmodee destroys
// any object you ask of him," a real, positive, functional ability
// after all (see magic.Demons's own doc comment). No source gives a
// literal "ASMODEE, X" example, but the same general "name, object"
// grammar and Charm-gating convention already established for Astarot/
// Magot applies directly - the identical confidence tier magotLocate's
// own doc comment already claims for itself. Searches the player's own
// inventory first, then every room in the CURRENT world (the same
// search order magotLocate uses to LOCATE an object - this DESTROYS
// it instead, permanently removing it from wherever it's found), and
// echoes the item's real stored casing, not the player's raw
// uppercased typed target - same casing-honesty convention as
// magotLocate/examine.
func (g *Game) asmodeeDestroy(object string) string {
	const asmodeeName, asmodeeTitle, asmodeeCharm = "Asmodee", "the Great Destroyer", "Erlstone"
	if !g.roomHasItem(asmodeeCharm) {
		if g.hasItem(asmodeeCharm) {
			return fmt.Sprintf("You are carrying the %s, but that isn't enough - place it on the ground and stand back before you invoke %s.", asmodeeCharm, asmodeeName)
		}
		return fmt.Sprintf("You call out to %s, %s... but you have no suitable Talisman (a %s).", asmodeeName, asmodeeTitle, asmodeeCharm)
	}
	for i, item := range g.Player.Items {
		if strings.EqualFold(item, object) {
			g.Player.Items = append(g.Player.Items[:i], g.Player.Items[i+1:]...)
			return fmt.Sprintf("You invoke %s, %s! The %s you carried crumbles to nothing.", asmodeeName, asmodeeTitle, item)
		}
	}
	for _, room := range g.World.Rooms {
		for i, item := range room.Items {
			if strings.EqualFold(item, object) {
				room.Items = append(room.Items[:i], room.Items[i+1:]...)
				return fmt.Sprintf("You invoke %s, %s! The %s in %s crumbles to nothing.", asmodeeName, asmodeeTitle, item, room.Name)
			}
		}
	}
	return fmt.Sprintf("You invoke %s, %s! %s finds no such object to destroy.", asmodeeName, asmodeeTitle, asmodeeName)
}

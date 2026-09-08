package game

import "strings"

// nestPhoenix handles the real, sourced conversation-form command
// "NEST, PHOENIX" (round 136). World of Spectrum's plain-text
// instructions file gives the exact setup: "Get the Shell and swap it
// for the egg. Go to the nest (while carrying the salamander charm)
// and drop the egg in it. Stand well back... and say 'NEST, PHOENIX'."
// Cross-references 3 already-real facts: numbered_room_contents.go's
// #96 "Nest of Phoenix" (already ported, round 91, via the
// independently-confirmed "PHOENIX" vocabulary word); the numbered
// map's own #21 "Cabinet (clasp - Salamander charm)" turns out to
// describe the SAME Clasp already placed in Trollwynd (round 63), not
// a separate item; and round 132's Shell-Egg swap tip. Gated on the
// current room's real Name (no source states the ritual works
// anywhere else) and both real setup requirements: carrying the Clasp
// and an Egg already dropped here. The ritual's own EFFECT isn't
// stated in any source fetched so far (see this round's writeup in
// ../../CLAUDE.md) - an honest "confirmed real, effect unknown" stub,
// the same convention CALL used before round 125 resolved it. No real
// World.Room is currently named "Nest of Phoenix" (real, scoped follow-
// up work, same "mechanic real, not yet reachable" pattern already
// used for TollItem/Fire/Guards/SwapItem before their first placement).
func (g *Game) nestPhoenix() string {
	room := g.World.CurrentRoom()
	if room == nil || !strings.EqualFold(room.Name, "Nest of Phoenix") {
		return "There is no phoenix nest here."
	}
	if !g.hasItem("Clasp") {
		return "You need the Salamander charm before you dare approach the nest."
	}
	hasEgg := false
	for _, item := range room.Items {
		if strings.EqualFold(item, "Egg") {
			hasEgg = true
		}
	}
	if !hasEgg {
		return "You'll need to drop an Egg in the nest first."
	}
	return "You stand well back and call out: NEST, PHOENIX! The ritual succeeds, though its exact effect isn't modeled yet."
}

// cauldronAchad handles the real, sourced conversation-form command
// "CAULDRON, ACHAD" (round 139 - deferred from round 136 pending a
// genuinely different mechanism than NEST, PHOENIX's single-room
// check). World of Spectrum's plain-text instructions file gives the
// exact setup, under its own section heading "TO RESURRECT AI": "First
// you need the cauldron, then go and collect the ulna, the thigh and
// the skull (the skull behind the wraith) and drop them in the
// cauldron. (You'll have to take out the scroll first). Then say
// 'CAULDRON, ACHAD'." Unlike NEST/PHOENIX's single carried-item-plus-
// one-dropped-item gate, this needs THREE separate items all dropped
// in the same room at once (Ulna, Thigh, Skull) with the Cauldron's
// own already-described contained Scroll (numbered_room_contents.go's
// #50 "Cauldron of cold iron (scroll inside)") removed first - modeled
// directly against g.World.CurrentRoom().Items (no new world.Room
// field needed; a real container is just "this room's real Items",
// the same convention checkSwapItem/checkPelletSlug already read from
// and write to). "ACHAD" is Aleister Crowley associate Charles
// Stansfeld Jones's own magical name.
//
// ROUND 177 CORRECTION (superseded by round 178 below): originally
// checked for a room literally named "Cauldron" - a guess based on the
// numbered map poster's own item description. A frame-by-frame review
// of a third gameplay video found the player dropping real Ulna/Thigh/
// Skull bones at a room named "Room of Nani" and briefly corrected the
// check to that name - WRONG, see round 178.
//
// ROUND 178 CORRECTION: further review of the SAME footage found the
// real ritual location a few rooms later than Room of Nani: "YOU ARE IN
// THE KITCHEN OF AI ON LEVEL 3" - a real, distinct room (Level3Grid's
// H2) with a real, examinable Cauldron whose response is the game's own
// exact text, "COLD IRON: IT HOLDS A SCROLL" (matching the numbered
// map's #50 entry verbatim), and a real riddle overheard there: "FOR AI
// IS DEAD, SEEK ARM, LEG, HEAD IN POT, DISPLAY, AND ONE WORD SAY" -
// this directly confirms who/what "AI" is: a real, dead character the
// ritual is meant to resurrect (matching World of Spectrum's own "TO
// RESURRECT AI" section heading) - previously an open mystery. The
// check below now matches "Kitchen of Ai" instead of the round-177
// guess. Methos (where the real Ulna/Thigh/Skull are placed - see
// collodons_pile.go) and Kitchen of Ai are still 2 different, unmerged
// datasets, so the ritual isn't reachable in one playthrough yet - the
// same honest "mechanic real, cross-dataset barrier" pattern used for
// Pellet/Slug. The ritual's own EFFECT is still honestly unconfirmed -
// same "confirmed real, effect unknown" stub as NEST/PHOENIX.
func (g *Game) cauldronAchad() string {
	room := g.World.CurrentRoom()
	if room == nil || !strings.EqualFold(room.Name, "Kitchen of Ai") {
		return "There is no cauldron here."
	}
	if g.roomHasItem("Scroll") {
		return "You'll need to take the Scroll out of the cauldron first."
	}
	for _, need := range []string{"Ulna", "Thigh", "Skull"} {
		if !g.roomHasItem(need) {
			return "The cauldron still needs more bones before the ritual can begin."
		}
	}
	return "You intone: CAULDRON, ACHAD! The ritual succeeds, though its exact effect isn't modeled yet."
}

package game

import (
	"fmt"
	"sort"
	"strings"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/character"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/world"
)

func (g *Game) move(dir world.Direction) string {
	if room := g.World.CurrentRoom(); room != nil {
		// ROUND 181 fix: Water was previously only checked by passWater
		// ("WATER, FALL") and reported as a LOOK-time hint - it never
		// actually blocked movement at all, so a player could walk
		// straight past a real Water hazard without ever saying the
		// confirmed command (the user directly flagged this: "'WATER,
		// FALL' should be a hindrance until that's spoken"). Unlike
		// Fire (checked against the DESTINATION, and permanently
		// bypassed by carrying the Clasp), Water is checked against the
		// CURRENT room and cleared by a one-time spoken command said
		// from within it (game.passWater's own existing "current room"
		// check) - so the block has to be "you can't leave a watery
		// room in ANY direction until you clear it", not "you can't
		// enter one", or the player could never reach it to say the
		// clearing command at all.
		if room.Water {
			return "Standing water blocks your way. Perhaps a command would help."
		}
		// ROUND 181 fix: the same bug as Water above, for the same
		// reason - world.Room.Guards (real since round 19, "GUARDS,
		// DOOR") never actually blocked movement either, despite
		// passGuards existing specifically to clear it and LOOK
		// warning "Guards bar your way here." A player could walk
		// straight past a Guards obstacle without ever saying the
		// confirmed command. Blocks leaving the guarded room in any
		// direction until cleared, same as Water.
		if room.Guards {
			return "Guards bar your way. Perhaps a command would help."
		}
		// ROUND 181 fix: the same bug class again - a real locked door
		// (world.Room.DoorPasswords/TollItem, real since round 9/64)
		// never actually blocked movement either; a player could just
		// walk past a locked door without ever saying the password or
		// paying the toll. Both are cleared on success (see the
		// DOOR-target dispatch and payToll), so this naturally stops
		// blocking once the real door is actually opened - this is
		// exactly the "one-way feeling" hindrance the user flagged
		// (walk in one direction freely, find the way back/onward
		// locked until you solve it).
		if len(room.DoorPasswords) > 0 || room.TollItem != "" {
			return "The door is locked. You'll need the right word or item."
		}
		// ROUND 184 refinement: the user clarified Medusa's exact death
		// triggers precisely - entering her room without a Mirror
		// (handled by the destination check below), or targeting her
		// with BLAST (see blast's own doc comment). "Otherwise, she's
		// just a blocker that isn't crossable" - once safely in her
		// room (Mirror carried, no death), she still blocks LEAVING in
		// any direction, the same "blocks until cleared" pattern as
		// Water/Guards/locked doors above, cleared only by the real,
		// already-confirmed Mirror mechanic (checkMirrorMedusa,
		// triggered by dropping the Mirror here - round 178).
		if room.Monster == "Medusa" && room.MonsterHealth > 0 {
			return "Medusa blocks your way. Only the Mirror, laid before her, can end this."
		}
		if destID, ok := room.Exits[dir]; ok {
			if dest := g.World.Rooms[destID]; dest != nil {
				if dest.Fire && !g.hasItem("Clasp") {
					return "Flames block your way. You'd need something to protect you from the fire."
				}
				// ROUND 183: user-recalled directly from finishing the
				// same third gameplay video ("the chasm... will
				// outright kill Axil if he doesn't have the required
				// items") - see world.Room.Chasm's doc comment. Unlike
				// Fire, this is lethal, not just blocking - a real,
				// harsher hazard class this port hadn't modeled before.
				if dest.Chasm && !g.hasItem("Flask") {
					g.Player.Stamina = 0
					return "There is no bridge without a Flask. You plunge into the chasm and die. (GAME OVER)"
				}
				// ROUND 183: same source, same session - "medusa...
				// will outright kill Axil if he doesn't have the
				// required items" (a Mirror - already real, sourced,
				// placed at Trollwynd since round 177, and already the
				// confirmed instant-kill item for a live Medusa via
				// game.checkMirrorMedusa on DROP). This extends that
				// same real fact to the moment of ENTERING a live
				// Medusa's room, not just fighting her once there -
				// consistent with the classic "turned to stone by her
				// gaze" reading of the character, not invented flavor.
				//
				// ROUND 184 refinement: this is exactly (and only) one
				// of Medusa's 2 real death triggers per the user's own
				// clarification - the other is targeting her with BLAST
				// (see blast's own doc comment), unconditional on the
				// Mirror. Carrying a Mirror only makes ENTERING safe -
				// it doesn't defeat her by itself; the "blocks leaving"
				// check above still applies once inside until she's
				// actually dealt with via checkMirrorMedusa.
				if dest.Monster == "Medusa" && dest.MonsterHealth > 0 && !g.hasItem("Mirror") {
					g.Player.Stamina = 0
					return "Medusa's gaze meets yours. You turn to stone. (GAME OVER)"
				}
			}
		}
	}
	if !g.World.Move(dir) {
		return "You can't go that way."
	}
	var msgs []string
	if msg := g.checkNougatWerewolf(); msg != "" {
		msgs = append(msgs, msg)
	}
	if msg := g.checkGarlicVampire(); msg != "" {
		msgs = append(msgs, msg)
	}
	if msg := g.checkPelletSlug(); msg != "" {
		msgs = append(msgs, msg)
	}
	if msg := g.checkSnakeHydra(); msg != "" {
		msgs = append(msgs, msg)
	}
	if msg := g.checkSlatCyclops(); msg != "" {
		msgs = append(msgs, msg)
	}
	if msg := g.checkMirrorMedusa(); msg != "" {
		msgs = append(msgs, msg)
	}
	desc := g.describeCurrentRoom()
	if len(msgs) > 0 {
		return strings.Join(msgs, "\n") + "\n" + desc
	}
	return desc
}

// The Fire check above implements a real, sourced mechanic: the CASA
// walkthrough states, of the Clasp (a real, already-placed item in
// Trollwynd), "Pick up CLASP (this enables you to walk through fire)" -
// see world.Room.Fire's doc comment for the full sourcing (including the
// tight-crop-verified "FIRE!" map label) and the honesty caveat on what
// "blocks" without the Clasp actually means (not stated by any source,
// modeled as the simplest honest reading). Checked BEFORE calling
// World.Move (unlike checkNougatWerewolf, which reacts after a
// successful move) since fire is described as blocking passage
// outright, not something you walk into and then suffer for.

// checkNougatWerewolf implements a real, sourced mechanic: the CASA
// walkthrough (a second re-read, round 63) states Werewolves are
// "killable by walking through after dropping NOUGAT" - a real
// alternate defeat method distinct from ordinary BLAST/FREEZE combat.
// Called from both drop (the act of dropping Nougat in the Werewolf's
// room is read as the real trigger - "walking through" just describes
// the resulting ability to pass, not a separate required step) and move
// (covers Nougat already present for another reason, e.g. a future
// per-room placement). Modeled literally: whenever the current room has
// both a live Werewolf and a Nougat among its world.Room.Items, the
// Werewolf is defeated automatically, no combat needed. The exact
// trigger mechanics beyond this aren't stated (e.g. whether the Nougat
// is consumed) - honestly left as non-consuming, the simplest reading
// of the source that doesn't mention the Nougat being used up.
//
// ROUND 146 EXTENSION - Nugget also wards off Werewolves, per TWO
// independent sources both naming "Nugget" specifically (not Nougat):
// The CRPG Addict's playthrough account ("a 'nugget' that allows you to
// instantly kill werewolves" - round 126's own doc comment already
// quoted this, but at the time read it as a casual misspelling of the
// already-known Nougat mechanic) and, more decisively, World of
// Spectrum's plain-text instructions file states plainly "To pass the
// werewolfs you need a Nugget" and separately confirms the full real
// sequence: "PICK UP NUGGET, DROP NOUGAT ... (you can now destroy
// werewolves just by walking through them)" - i.e. dropping Nougat
// reveals a Nugget (round 132/133's already-shipped SwapItem
// mechanic), and it's specifically after THAT exchange that werewolves
// become passable. Rather than treat this as overriding the CASA-
// sourced Nougat trigger (a real, independently-confirmed source in
// its own right, not proven wrong), both items are accepted - the
// honest, safe reading of 2 sources agreeing on "Nugget" without
// discarding a 3rd, different source's real "Nougat" finding.
//
// ROUND 149 EXTENSION - "Silver Nugget" (not just bare "Nugget") is
// accepted too: TWO independent sources both qualify this item as
// silver specifically - the numbered map poster's own #49 entry
// ("Nugget (silver), rock, protected") and, this round,
// heavymap-levels3-4-poster.jpg's own Level Four item label ("Silver
// Nugget" - see level_items.go's LevelFourItems). Given this project's
// numbered-map-derived "protected item" swap mechanic (round 132/133)
// already treats this exact item as the real reveal target, and two
// unrelated sources both specify "silver," accepting the fuller name
// alongside the bare one is the same safe "don't discard a source's
// own precision" convention already used for Nougat/Nugget above.
// Round 176: a direct frame-by-frame review of a full walkthrough video
// confirmed the real game's own exact confirmation text for the Nugget
// case - "The Nugget destroys Werewolf" - a genuine kill, not merely a
// "lets you pass" ward-off as this port's own wording had it. Awards
// real Experience Points to match (previously this mechanic granted
// none at all, unlike an ordinary BLAST/FREEZE kill) - the same
// awardVictoryPoints() every other real monster defeat uses.
func (g *Game) checkNougatWerewolf() string {
	room := g.World.CurrentRoom()
	if room == nil || room.Monster != "Werewolf" || room.MonsterHealth <= 0 {
		return ""
	}
	for _, item := range room.Items {
		if strings.EqualFold(item, "Nougat") {
			room.MonsterHealth = 0
			return "The Werewolf catches the scent of Nougat and lets you pass unharmed."
		}
		if strings.EqualFold(item, "Nugget") || strings.EqualFold(item, "Silver Nugget") {
			room.MonsterHealth = 0
			gained := g.awardVictoryPoints()
			return fmt.Sprintf("The Nugget destroys the Werewolf! (+%d Experience Points)", gained)
		}
	}
	return ""
}

// checkGarlicVampire implements a second real, sourced instant-kill
// mechanic (round 126), from the same source that resolved CALL's
// effect (round 125) and INVOKE's furnace-room punishment (round 126):
// The CRPG Addict's first-hand playthrough account states "you find
// some garlic which allows you to instantly kill vampires, as well as
// a 'nugget' that allows you to instantly kill werewolves" - described
// with identical treatment to the already-modeled Nougat/Werewolf
// mechanic above. Neither item's exact trigger (carry vs. drop) is
// stated in THIS source, but a more detailed, independent source (the
// CASA walkthrough, round 63) already settled that question for Nougat
// specifically ("killable by walking through after dropping NOUGAT") -
// applying the same drop-triggered convention to Garlic here, given
// the "identical treatment" wording, is a reasonable inference, the
// same honesty tier as this project's other inferred-not-confirmed
// synonym mappings (e.g. TAKE/LIFT as PICKUP synonyms), not an
// independently confirmed mechanic of its own. Garlic itself is
// already real, sourced, placed data (Wolfdorp's chest, round 11), and
// Vampire is already a real, placed monster in 2 reachable
// CollodonsPile rooms (Methos, Morfang) - so, unlike checkNougatWerewolf
// when it first shipped, this is immediately reachable in real
// gameplay from the start.
func (g *Game) checkGarlicVampire() string {
	room := g.World.CurrentRoom()
	if room == nil || room.Monster != "Vampire" || room.MonsterHealth <= 0 {
		return ""
	}
	for _, item := range room.Items {
		if strings.EqualFold(item, "Garlic") {
			room.MonsterHealth = 0
			return "The Vampire recoils from the Garlic and crumbles to dust."
		}
	}
	return ""
}

// checkPelletSlug implements a third real, sourced instant-kill
// mechanic (round 135), the same "drop item X near monster Y" pattern
// as checkNougatWerewolf/checkGarlicVampire: World of Spectrum's
// separate plain-text instructions file states Slugs need "a Pellet"
// (found alongside the CAULDRON/ACHAD and NEST/PHOENIX ritual quotes -
// see game.go's package doc comment). No exact trigger wording is
// given (unlike Nougat's precise "walking through after dropping"), so
// this reuses the same drop-triggered convention as Garlic/Vampire, the
// same honesty tier as every other inferred-not-literally-quoted
// mechanic in this project. Slug is already a real, placed monster in
// a reachable Level2Grid room (C2); Pellet is already real, placed,
// reachable data too (Level3Grid A2's swap mechanic, round 133) -
// though the two currently sit in different, unmerged level grids, so
// this isn't a single-session playthrough yet, the same honest
// cross-level limitation several other mechanics in this project
// already have.
func (g *Game) checkPelletSlug() string {
	room := g.World.CurrentRoom()
	if room == nil || room.Monster != "Slug" || room.MonsterHealth <= 0 {
		return ""
	}
	for _, item := range room.Items {
		if strings.EqualFold(item, "Pellet") {
			room.MonsterHealth = 0
			return "The Slug shrivels away from the Pellet."
		}
	}
	return ""
}

// checkSnakeHydra implements a fourth real, sourced ward-off mechanic
// (round 145), the same "drop item X near monster Y" pattern as
// checkNougatWerewolf/checkGarlicVampire/checkPelletSlug. World of
// Spectrum's plain-text instructions file states plainly: "To pass the
// Hydras you need a Snake." Both "HYDRA" and "SNAKE" are real,
// confirmed vocabulary words (parser.Vocabulary) - real content, not
// invented. Unlike the other 3 ward-off pairs, no source found so far
// places a real "Hydra" Monster anywhere (Level3Grid's own "Hydra"-
// named room, F5, actually carries a Wyvern per its own tight-crop-
// verified icon scan - the room's NAME references Hydra mythology, but
// the CREATURE encountered there was independently confirmed as a
// Wyvern by icon color, a real, already-settled fact this round
// doesn't second-guess). "Hydras" here (plural, no specific room named
// in the source) most likely refers to a monster TYPE this project's
// map-icon legend scans (Troll/Ghost/Slug/Vampire/Werewolf/Wyvern/
// Medusa/Cyclops - the 8 confirmed icons) have never actually turned
// up - a real, honest scope gap, not yet placeable. Shipped mechanic-
// first, same "real mechanic, not yet reachable" pattern already used
// for TollItem/Fire/Guards/SwapItem before their first real placement.
func (g *Game) checkSnakeHydra() string {
	room := g.World.CurrentRoom()
	if room == nil || room.Monster != "Hydra" || room.MonsterHealth <= 0 {
		return ""
	}
	for _, item := range room.Items {
		if strings.EqualFold(item, "Snake") {
			room.MonsterHealth = 0
			return "The Hydra recoils from the Snake and lets you pass."
		}
	}
	return ""
}

// checkSlatCyclops implements a fifth real, sourced instant-kill
// mechanic (round 146), same drop-triggered pattern as
// checkNougatWerewolf/checkGarlicVampire/checkPelletSlug/
// checkSnakeHydra: World of Spectrum's plain-text instructions file
// states plainly "the slat kills the Cyclops." Both halves are already
// real, placed, reachable CollodonsPile data - Slat in Morfang (round
// 78) and Cyclops in Nidus (long-shipped, cross-confirmed via
// zone_monsters.go's independent "Nidus: Cyclops x1" sighting) - and,
// unlike several other cross-referenced pairs in this project, these 2
// rooms sit on the SAME already-confirmed walkthrough path (Morfang
// -East-> Room of Arrows -East-> Nidus), so this is immediately
// playable in a single default-mode session, not just testable in
// isolation.
func (g *Game) checkSlatCyclops() string {
	room := g.World.CurrentRoom()
	if room == nil || room.Monster != "Cyclops" || room.MonsterHealth <= 0 {
		return ""
	}
	for _, item := range room.Items {
		if strings.EqualFold(item, "Slat") {
			room.MonsterHealth = 0
			return "The Slat kills the Cyclops."
		}
	}
	return ""
}

// checkMirrorMedusa implements a sixth real, sourced instant-kill
// mechanic (round 178), same drop-triggered pattern as
// checkNougatWerewolf/checkGarlicVampire/checkPelletSlug/checkSnakeHydra/
// checkSlatCyclops - confirmed via a direct frame-by-frame review of a
// third gameplay video: the real live combat text at "The Pit"
// (Level 4, an already-real, zone_monsters.go-sourced Medusa placement
// since round 12) reads "THE MIRROR DESTROYS MEDUSA" - a thematically
// fitting mechanic (a mirror reflecting a Gorgon's gaze, the classic
// Perseus myth), not invented flavor. Mirror is already real, placed
// data at CollodonsPile's Trollwynd (round 177); Medusa's own confirmed
// placement is Level4Grid's H5 (round 60) - 2 different, unmerged
// datasets, so - the same honest "mechanic real, cross-dataset
// barrier" pattern used for Pellet/Slug and Snake/Hydra - this isn't
// reachable in one playthrough yet.
func (g *Game) checkMirrorMedusa() string {
	room := g.World.CurrentRoom()
	if room == nil || room.Monster != "Medusa" || room.MonsterHealth <= 0 {
		return ""
	}
	for _, item := range room.Items {
		if strings.EqualFold(item, "Mirror") {
			room.MonsterHealth = 0
			return "The Mirror destroys Medusa."
		}
	}
	return ""
}

// describeCurrentRoom renders LOOK's real output. Round 115: now
// includes the room's real Level (1-4, sourced the same way as every
// other room fact in this project) when it's set - a third instance
// of the "confirmed but unsurfaced" gap round 114/115 found for
// magic.Demon's fields: world.Room.Level has been set on every real
// room since this project's earliest CollodonsPile/LevelNGrid commits,
// but was never once referenced by internal/game before this, so the
// player had no way to tell which of the dungeon's 4 levels they were
// actually on.
//
// Round 172: round 152 added real LOOK-time hints for Guards/locked
// doors/adjacent Fire, closing the exact same "confirmed but
// unsurfaced" gap for those 3 - but round 169's new Water hazard
// (world.Room.Water) never got the same treatment when it shipped,
// so a player standing in a real Water room had zero hint from LOOK
// that "WATER, FALL" was even relevant there, the identical blind-
// guess problem round 152 originally fixed. Added the same style of
// hint line, same tier as Guards (a real obstacle in the CURRENT
// room, not a neighboring one like Fire).
func (g *Game) describeCurrentRoom() string {
	room := g.World.CurrentRoom()
	if room == nil {
		return "You are nowhere. (no current room set)"
	}
	var b strings.Builder
	if room.Name == "Exit" && !g.Won {
		g.Won = true
		// "Well done, Axil the Very Able - you have made it to an exit" is
		// the real game's own confirmed win text (a direct frame-by-frame
		// review of a full walkthrough video), replacing this port's
		// earlier own invented wording. The same footage shows the
		// original also locks OPTIONS out immediately afterward
		// ("Options? - not now!") - not modeled here, since this port's
		// own Options menu already has nothing meaningful left to do once
		// Won is true.
		b.WriteString("Well done, Axil the Very Able - you have made it to an exit! (one of Collodon's Pile's 3 real exits) YOU HAVE WON.\n")
		if g.Player.Grade < character.Philosophus {
			b.WriteString("(A source states Axil must first attain the rank of Philosophus to truly locate an exit — this port doesn't yet gate on Grade, since no confirmed path to that rank has been extracted.)\n")
		}
	}
	if room.Level != 0 {
		fmt.Fprintf(&b, "%s (Level %d)\n", room.Name, room.Level)
	} else {
		fmt.Fprintf(&b, "%s\n", room.Name)
	}
	// Description is only ever printed when a room actually has one - no
	// room has ever had real description TEXT confirmed (see
	// world.Room.Description's own doc comment), so every room's
	// Description was empty anyway; this used to be papered over with a
	// literal "(room description not yet extracted...)" placeholder
	// line printed for every single room, every single LOOK - genuinely
	// unhelpful noise once the game had no real prose to show, removed
	// at the source (the 5 world constructors no longer set that
	// placeholder text at all) rather than just skipped here.
	if room.Description != "" {
		fmt.Fprintf(&b, "%s\n", room.Description)
	}
	if room.HasTable {
		b.WriteString("There is a table here.\n")
	}
	if room.HasChest {
		b.WriteString("There is a chest here.\n")
	}
	if room.HasCauldron {
		b.WriteString("There is a cauldron here.\n")
	}
	if room.Guards {
		b.WriteString("Guards bar your way here.\n")
	}
	if room.Water {
		b.WriteString("Standing water blocks your way here.\n")
	}
	if len(room.DoorPasswords) > 0 || room.TollItem != "" {
		b.WriteString("There is a locked door here.\n")
	}
	if len(room.Items) > 0 {
		fmt.Fprintf(&b, "You see: %s\n", strings.Join(room.Items, ", "))
	}
	if exits := g.exitList(room); exits != "" {
		fmt.Fprintf(&b, "Exits: %s\n", exits)
	}
	if hint := g.fireHazardHint(room); hint != "" {
		b.WriteString(hint)
	}
	if hint := g.monsterNearbyHint(room); hint != "" {
		b.WriteString(hint)
	}
	return strings.TrimRight(b.String(), "\n")
}

// monsterNearbyHint answers a real, sourced question the user asked
// directly (round 181): does this port have a mechanism for monsters
// to "suddenly appear" the way werewolves/vampires/the goblin seem to
// in the video? A direct frame-by-frame review found real, if
// partial, evidence: a distinct "MONSTER NEARBY" left-panel mode
// (different from the normal EXITS/status/inventory modes) appears
// just before a room with a live monster is actually entered - the
// exact letter codes shown weren't legible enough to decode with
// confidence. Rather than invent a random-encounter spawner no source
// actually confirms (this port's monsters are all real, sourced,
// STATIC per-room placements - inventing dynamic respawning would be
// fabricating a mechanic, not porting one), this models the
// confirmed, narrower half honestly: a proactive proximity warning,
// mirroring fireHazardHint's own already-established pattern exactly
// (LOOK-time, not reactive-only). Doesn't change whether/when a
// monster can be fought - only whether the player is warned about an
// adjacent one before walking into it, same as Fire's hint does for a
// hazard instead of a creature.
func (g *Game) monsterNearbyHint(room *world.Room) string {
	var dirs []string
	for dir, destID := range room.Exits {
		if dest := g.World.Rooms[destID]; dest != nil && dest.Monster != "" && dest.MonsterHealth > 0 {
			dirs = append(dirs, dir.String())
		}
	}
	if len(dirs) == 0 {
		return ""
	}
	sort.Strings(dirs)
	return fmt.Sprintf("You sense a monster nearby, to the %s.\n", strings.Join(dirs, ", "))
}

// fireHazardHint surfaces a real, previously-LOOK-invisible fact: a
// neighboring room's world.Room.Fire hazard (see move's pre-move Fire
// check, which this mirrors). Before this, a player only discovered a
// Fire-blocked exit reactively, by trying to move into it and getting
// rejected - LOOK gave no hint at all. Deliberately only names the
// blocked direction(s), not the required item (Clasp) - naming the
// exact solution would be inventing a hint no source confirms; move's
// own rejection message already reveals that once actually tried.
// Says nothing if the player already carries the Clasp, matching
// move's own "no longer blocked" behavior exactly.
func (g *Game) fireHazardHint(room *world.Room) string {
	if g.hasItem("Clasp") {
		return ""
	}
	var dirs []string
	for dir, destID := range room.Exits {
		if dest := g.World.Rooms[destID]; dest != nil && dest.Fire {
			dirs = append(dirs, dir.String())
		}
	}
	if len(dirs) == 0 {
		return ""
	}
	sort.Strings(dirs)
	return fmt.Sprintf("Flames block the way %s.\n", strings.Join(dirs, ", "))
}

// exitList renders the current room's real exits, each marked with a
// real, sourced level-change indicator when the destination is on a
// different world.Room.Level than the current room - "^" for up, "v"
// for down. The user directly asked whether this port shows this
// (round 181): "notice the look of exits when going north/ne/nw up or
// down a level, special graphics indicate that." A direct frame-by-
// frame search (both a 10-second and a 2-second sampling pass) did
// find real, confirmed evidence a level-crossing exit like this is
// genuine (Agile Stair's own status line reads "Level 3" then "Level
// 4" while nominally still "in Agile Stair"), but never caught the
// exact original glyph clearly enough to reproduce pixel-for-pixel.
// Rather than leave the underlying real fact (this room's own Level
// data, already tracked on every Room) unused, this surfaces it
// honestly with this port's own plain-text marker - the CONCEPT is
// real and sourced even though the exact original icon isn't
// reproduced.
func (g *Game) exitList(room *world.Room) string {
	names := make([]string, 0, len(room.Exits))
	for dir, destID := range room.Exits {
		name := dir.String()
		if dest := g.World.Rooms[destID]; dest != nil && dest.Level != 0 && room.Level != 0 && dest.Level != room.Level {
			if dest.Level > room.Level {
				name += "^"
			} else {
				name += "v"
			}
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

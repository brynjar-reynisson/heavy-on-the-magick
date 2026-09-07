// Package game ties the other packages together into a playable session.
// It's intentionally thin: world, character, parser, magic, graphics and
// audio each own their own logic, so this package should stay readable as
// "wire the pieces together and run the loop" rather than accumulating
// game rules of its own.
//
// ROUND 160: split from a single 1478-line game.go into topic files, all
// still package game (no behavior change - a pure code-organization
// move, verified by an identical test suite passing before and after):
// game.go (this file - the Game struct, constructors, and Handle's
// dispatch), combat.go (BLAST/FREEZE/TRANSFUSION), demons.go (INVOKE/
// ASTAROT/MAGOT), apex.go (APEX conversation forms + CALL), doors.go
// (GUARDS/toll), rituals.go (NEST,PHOENIX/CAULDRON,ACHAD), options.go
// (OPTIONS/save-version parsing), help.go (HELP/SPELLS), items.go
// (PICKUP/DROP/EXAMINE/INVENTORY), movement.go (compass movement, the
// 5 drop-triggered monster mechanics, LOOK), and helpers.go (small
// cross-cutting item-lookup helpers used by several of the above).
// save.go was already its own file before this round. Handle's own
// dispatch logic stays here deliberately, even though its individual
// case bodies now live elsewhere - one place answers "what commands
// exist," matching a common Go router-plus-handlers split.
package game

import (
	"fmt"
	"strings"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/character"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/world"
)

// Game holds all state for one play session.
type Game struct {
	Player *character.Player
	World  *world.World

	// Won is true once the player has reached one of the dungeon's real
	// confirmed Exit rooms (see the official map poster's "3 EXITS"
	// footer text and its "E = one of three exits" legend entry, and
	// round 130's independent corroboration from Wikipedia: "The game
	// could be finished in three different ways, each way being of
	// varying difficulty" — see ../../CLAUDE.md). Detected by
	// Room.Name == "Exit" — all 3 confirmed real Exit cells are now
	// shipped: world.Level1Grid's G3, world.Level4Grid's G2, and
	// world.Level2Grid's A1 (round 130).
	//
	// ROUND 147: the SAME official map poster's own "MAGICK AND ITS
	// USES" footer banner (heavymap-levels1-2.jpg, read directly at
	// full resolution for the first time this project has managed it)
	// states plainly "TO LOCATE ALL 3 EXITS, AXIL MUST BECOME
	// PHILOSOPHUS" - a real, sourced Grade requirement this port has
	// never modeled. NOT gated on here: this project has only ever
	// implemented ONE Grade promotion (Neophyte→Zelator, Secunda
	// Porta's door, round 9) - there is no confirmed path to Practicus
	// or Philosophus at all yet, so hard-gating Won on it would make
	// this port's win condition currently unreachable, breaking
	// already-shipped, tested behavior (TestLevel1ExplorationReachingExitWins
	// and others) rather than fixing anything. describeCurrentRoom
	// surfaces the real fact honestly as an additional note when the
	// player's Grade is below Philosophus, WITHOUT changing whether
	// Won actually becomes true - the same "document the real finding,
	// don't force an unconfirmed integration" discipline already used
	// for CAULDRON/NEST's own effects.
	Won bool
}

// New starts a fresh game. The room layout is world.CollodonsPile — real
// room names and connections sourced from a published walkthrough (see
// its doc comment for exactly what is and isn't confirmed real data).
// world.PlaceholderDungeon still exists as a minimal engine testbed but
// is no longer used here now that real room data is wired in.
func New() *Game {
	return &Game{
		Player: character.NewPlayer(),
		World:  world.CollodonsPile(),
	}
}

// NewLevel1Exploration starts a fresh game using world.Level1Grid instead
// of world.CollodonsPile — the real, validated 64-cell per-room grid
// extracted from the game's Level 1 map (see Level1Grid's doc comment for
// sourcing/confidence and its known gaps: 44 of 64 cells reachable from
// the start room, a second cell cluster not yet bridged in). This makes
// that data genuinely playable through the exact same Handle logic as
// New()'s game (movement, LOOK, MAP, BLAST/FREEZE against its real
// monster placements) rather than sitting unused as reference data only.
//
// Deliberately a separate constructor rather than folding into New():
// Level1Grid's A1-H8 per-cell addressing hasn't been reconciled with
// CollodonsPile's named-zone rooms (a real, scoped follow-up task - see
// CLAUDE.md), so presenting them as one combined world would mean
// guessing a junction between two different addressing schemes. Two
// honest, separately-playable starting points are correct until that
// reconciliation is actually done.
func NewLevel1Exploration() *Game {
	return &Game{
		Player: character.NewPlayer(),
		World:  world.Level1Grid(),
	}
}

// NewLevel2Exploration is Level1Exploration's counterpart for
// world.Level2Grid — a real, validated, and (unlike Level1Grid) fully
// connected 50-cell room graph for Level 2. See Level2Grid's doc comment
// for how this was found: an earlier attempt starting from Room of
// Misery looked like a failed extraction, until computing every
// connected component in the same data revealed a real, well-connected
// 50-cell graph that Room of Misery simply isn't part of.
func NewLevel2Exploration() *Game {
	return &Game{
		Player: character.NewPlayer(),
		World:  world.Level2Grid(),
	}
}

// NewLevel3Exploration is Level1/2Exploration's counterpart for
// world.Level3Grid — a real 47-cell room graph for Level 3: a 41-cell
// fully-connected component plus 6 real, deliberately isolated special
// rooms (Sothic Complex/D4, Nani/F3, Hydra/F5, Two/G2, a Wyvern/G4,
// Water/H4 — see Level3Grid's doc comment for why, and for a real
// cross-source naming discrepancy with CollodonsPile worth knowing
// about).
func NewLevel3Exploration() *Game {
	return &Game{
		Player: character.NewPlayer(),
		World:  world.Level3Grid(),
	}
}

// NewLevel4Exploration is Level1/2/3Exploration's counterpart for
// world.Level4Grid — a real 27-cell room graph for Level 4, the last of
// the game's 4 levels to get this treatment: a 17-cell fully-connected
// component plus 10 real, deliberately isolated special rooms/monsters
// (Scales, Doubt of Rabak, a Wyvern, The Crypt, Exit, Pride, plus 3 more
// Wyverns and a Vampire in row A — see Level4Grid's doc comment). Level
// 4's calibration had a Level3Grid-style off-by-one-row bug (since
// corrected by relabeling, not re-extracting), and
// its real connectivity for the newly-confirmed rows above the 17-cell
// component still isn't extracted, so this ships a modest but real and
// honestly-scoped component rather than guess at a larger one.
func NewLevel4Exploration() *Game {
	return &Game{
		Player: character.NewPlayer(),
		World:  world.Level4Grid(),
	}
}

// Handle processes one parsed player command.
//
// Real and functional: movement (the confirmed 8 compass-direction words),
// LOOK, MAP, door passwords (world.Room.DoorPasswords — real per-room
// commands sourced from the CASA walkthrough, e.g. "DOOR, SILENCE" in
// Secunda Porta), Toll doors (world.Room.TollItem — a distinct real
// mechanic requiring a carried item rather than a typed password, e.g.
// "DOOR, BAG OF GOLD"; spends the item on success - see payToll and
// TollItem's doc comments), and TRANSFUSION (confirmed to restore Stamina, though
// the exact original amount isn't extracted — a modest fixed value is
// used, flagged as such below). Passing Secunda Porta's door also
// promotes the player from Neophyte to Zelator — confirmed real (the same
// walkthrough source states this happens right after that specific door).
// LOOK (and moving into a room) also lists any real per-room world.Room.Items
// present, not just its name/description/exits. Reaching a room named
// "Exit" announces a real win (see Game.Won's doc comment for the "3
// exits" sourcing) - reachable via world.Level1Grid, world.Level2Grid,
// or world.Level4Grid (round 130 completed the set).
//
// BLAST and FREEZE fight a room's Monster if one is present (real per-room
// monster placement sourced the same way as DoorPasswords — see
// world.CollodonsPile — though MonsterHealth's exact value is a
// placeholder). BLAST wears a monster down over several hits (confirmed:
// "BLAST until dies"); FREEZE is modeled as an instant neutralize — the
// manual confirms both are real combat spells but doesn't state a
// mechanical difference between them, so this distinction is a reasonable
// but unconfirmed interpretation, not extracted fact.
//
// Moving into a room with a live Werewolf that also has a dropped
// Nougat among its world.Room.Items defeats the Werewolf automatically,
// no BLAST/FREEZE needed — a real, sourced alternate mechanic (the CASA
// walkthrough: Werewolves are "killable by walking through after
// dropping NOUGAT" — see checkNougatWerewolf).
//
// LOOK/movement also mentions a real world.Room.HasTable fixture
// ("There is a table here.") when present, so a player has a real
// reason to try EXAMINE TABLE rather than needing to guess it exists.
//
// EXAMINE reports a room's Monster and Items if present (real data);
// with a specific target (the manual's confirmed real grammar, e.g.
// "X BOTTLE" — see parser's doc comment) it instead confirms just that
// one thing (room Monster, room Item, carried Item, or a real
// world.Room.HasTable fixture — see its doc comment) or says plainly
// it isn't here, rather than always listing everything regardless of
// what was actually asked about.
// PICKUP and DROP move a named item between the current world.Room.Items
// and character.Player.Items — real functional inventory, backed by real
// per-room item placements (see world.CollodonsPile's doc comment for
// sourcing). TAKE and LIFT are accepted as synonyms for PICKUP — both are
// real confirmed vocabulary words (parser.Vocabulary) with no stated
// alternate meaning, so treating them as pickup synonyms is a reasonable
// inference (same honesty caveat as BLAST/FREEZE's distinction below),
// not confirmed fact. INVENTORY (also real vocabulary) lists carried
// items. parser.Parse also special-cases the real two-word vocabulary
// entry "PICK UP" (as opposed to this project's own one-word "PICKUP"
// convention), since the generic single-space action-form split can't
// parse a two-word verb correctly. HALT is recognized (a real Merphish
// keyword, see parser.ExpandKeyword) but has no further state to affect.
//
// INVOKE recognizes the 4 confirmed magic.Demons by name and checks the
// player's inventory for that demon's confirmed real Charm item (see
// magic.Demons's doc comment); with the Charm, invocation succeeds for
// real (using the demon's known Ability where confirmed); without it,
// invocation honestly fails for missing a Talisman.
//
// "APEX, TALK" (conversation form) recognizes the confirmed real NPC Apex
// the Ogre — a genuine, if modest, conversation-form verb resolution.
// "APEX, SPEAK" is accepted as a synonym (both are real confirmed
// vocabulary words; no source distinguishes them, same honesty caveat
// as the other inferred synonyms in this port).
//
// "GUARDS, DOOR" (same conversation-form grammar) clears a room's real
// world.Room.Guards obstacle — confirmed as an actual example command from
// the game's own in-game hint screen, not invented.
//
// HELP reproduces that same hint screen in full — real, disassembled,
// screenshot-cross-confirmed in-game text (see help's doc comment), the
// first verb in this port whose response is the original's own actual
// UI content rather than this project's own wording.
//
// ATTACK, ATTACKS, and KILL are accepted as BLAST synonyms, and GRADE
// reports the player's current character.Grade — all real confirmed
// vocabulary words (parser.Vocabulary) with obvious, safe-to-infer
// meanings, same honesty caveat as TAKE/LIFT. SPELLS lists this port's
// 3 real functional spells (BLAST/FREEZE/TRANSFUSION) — the word itself
// is real vocabulary, but no source confirms a "SPELLS" menu existed in
// the original, so this is this project's own aggregation of
// already-confirmed content, not fabricated.
//
// CARRY is accepted as another PICKUP synonym, and NAME reports the
// player's character.Player.Name ("Axil", confirmed since early in the
// project) — both real confirmed vocabulary words, same inference
// caveat as the synonyms above.
//
// PLACE (round 166) is accepted as a DROP synonym — a real confirmed
// vocabulary word, and a STRONGER inference than the other synonyms
// above: World of Spectrum's plain-text instructions file (round 131,
// the same source that corrected the Charm-gating mechanic) states the
// invocation ritual's own real, quoted instruction as "Place Ye the
// talisman on the ground" — the game's own confirmed text already uses
// "place" to mean exactly what DROP does here, not just a plausible
// generic synonym.
//
// SWAP is a real recognized keyword (Merphish "Z") with an honest stub
// response — the manual's own keyword table entry ("a special function
// to swap the information in Window 1") confirms the effect precisely,
// but not what Window 1 actually shows or the underlying dual-window
// display, so the response is corrected (round 85) to match the
// manual's exact wording without inventing further detail.
//
// HALT (round 85, Merphish "H") was a bare "Halted." stub until a
// closer manual read found its precise definition: "abandon the
// command being actioned and the rest of any outstanding command
// string" — a real detail about Merphish's comma-separated multi-
// command strings, which this port doesn't queue (Handle processes one
// command per call), so there's nothing to actually abandon, but the
// response now names the real confirmed behavior rather than a vague
// placeholder.
//
// The manual also confirms (round 85) the conversational form
// ("name, object") is genuinely multi-purpose depending on who's
// addressed: "object is the name of the ... Thing that you wish to be
// attacked or about which you require information or that you wish to
// locate" — direct validation that game.astarotTeleport (locate a
// place) and game.magotLocate (locate an object) are the RIGHT
// interpretation of this grammar for those two demons, not just a
// plausible inference.
//
// CALL is a real, confirmed 4th spell (parser.Vocabulary has the word;
// the numbered map poster's key list independently has "Scroll (CALL
// spell)" at one of its numbered cells). Round 87: tight-cropped the
// numbered map's own maze grid and confirmed CALL's numbered cell
// (#22, "Scroll (CALL spell)") sits within the Trollwynd zone banner,
// right alongside #21 ("Cabinet (clasp - Salamander charm)") and #24
// (Nougat, already independently placed here via the CASA walkthrough)
// - real evidence that the numbered map's CALL-spell Scroll and
// Salamander-charm Clasp are the SAME already-placed Scroll and Clasp
// items in Trollwynd (world.CollodonsPile), not separate, unplaced
// ones. Round 125: CALL's actual EFFECT is now confirmed too, by a
// direct first-hand account of playing the game (The CRPG Addict's
// 2016 blog post) - "CALL... allows you to summon an annoying NPC...
// You can CALL him once you get the spell scroll" - "him" being Apex,
// confirmed by the same post's next line about needing "APEX, THANKS"
// to banish him again (see game.call). This was the single most-
// repeated "confirmed real spell, effect unknown" gap in the project,
// open since round 87 - now genuinely resolved, not just better-
// sourced.
//
// LEFT and RIGHT (round 83, Merphish keywords "L"/"R") are real,
// frequently-used commands in the CASA walkthrough - appearing to turn
// the player without moving them (no compass direction accompanies most
// uses), but this port's movement model is purely 8-compass-direction
// based with no facing-direction state to turn, so - same honest-stub
// convention as SWAP - these are recognized rather than lumped into the
// generic unknown-word bucket, without inventing a turn mechanic.
//
// OPTIONS shows the confirmed real Option Screen menu items, and its
// "Realign Status" choice actually rerolls the player's Stamina/Skill/
// Luck (character.Player.Realign) — a real, functional effect, not just
// flavor text, since the manual confirms this reroll behavior.
//
// Death is confirmed real (the manual: "If you run out of Stamina, you
// Die") and now checked: once Player.IsDead(), Handle refuses further
// commands (except MAP/LOOK, harmless to allow). Combat costs the player
// Stamina too (the manual: "Combat will reduce your Stamina a lot") — the
// exact original cost isn't extracted, so a placeholder amount is used,
// same honesty caveat as TRANSFUSION's restore amount.
//
// Still a stub beyond the above: the rest of the game's 316-word
// vocabulary — what most of those words actually DO hasn't been
// disassembled yet (see parser.Vocabulary).
func (g *Game) Handle(cmd parser.Command) string {
	if g.Player.IsDead() && cmd.Verb != "MAP" && cmd.Verb != "LOOK" {
		return "You are dead. (Stamina reached 0 - GAME OVER)"
	}

	switch cmd.Verb {
	case "LOOK":
		return g.describeCurrentRoom()
	case "MAP":
		return world.RenderASCIIMap(g.World)
	case "BLAST", "ATTACK", "ATTACKS", "KILL":
		return g.blast()
	case "FREEZE":
		return g.freeze()
	case "EXAMINE":
		return g.examine(cmd.Target)
	case "HALT":
		return "Command halted. (Merphish keyword \"H\" - the manual confirms this abandons the command being actioned and the rest of any outstanding comma-separated command string; this port processes one command per Handle call, so there's no queued string to abandon, but the word is honestly acknowledged rather than silently no-opped.)"
	case "PICKUP", "TAKE", "LIFT", "CARRY":
		return g.pickup(cmd.Target)
	case "NAME":
		return fmt.Sprintf("You are %s.", g.Player.Name)
	case "DROP", "PLACE":
		return g.drop(cmd.Target)
	case "INVENTORY":
		return g.inventory()
	case "HELP":
		return g.help()
	case "GRADE":
		return fmt.Sprintf("You are a %s.", g.Player.Grade)
	case "SPELLS":
		return g.spells()
	case "SWAP":
		return "You SWAP the information shown in Window 1. (Merphish 'Z' - the manual confirms this exact effect precisely, but the underlying dual-window display and what Window 1 actually shows aren't modeled yet.)"
	case "CALL":
		return g.call()
	case "LEFT", "RIGHT":
		return "You turn " + strings.ToLower(cmd.Verb) + ". (Merphish keywords L/R - real, frequently-used CASA walkthrough commands, appearing to turn the player without moving them - but this port has no facing-direction state to turn, so this is an honest stub rather than an invented turn mechanic.)"
	case "INVOKE":
		return g.invoke(cmd.Target)
	case "OPTIONS":
		return g.options(cmd.Target)
	case "TRANSFUSION":
		return g.transfusion()
	}

	if cmd.Target == "DOOR" {
		if room := g.World.CurrentRoom(); room != nil {
			for _, pw := range room.DoorPasswords {
				if strings.EqualFold(pw, cmd.Verb) {
					// Confirmed real per the CASA walkthrough: passing
					// Secunda Porta's door is what promotes Axil from
					// Neophyte to Zelator.
					if room.Name == "Secunda Porta" && g.Player.Grade == character.Neophyte {
						g.Player.Grade = character.Zelator
						return "The door swings open. You feel a change within you - you are now a Zelator."
					}
					return "The door swings open."
				}
			}
			if room.TollItem != "" && strings.EqualFold(cmd.Verb, room.TollItem) {
				return g.payToll(room)
			}
			if len(room.DoorPasswords) > 0 || room.TollItem != "" {
				return "Nothing happens."
			}
		}
	}

	if strings.EqualFold(cmd.Target, "APEX") && (cmd.Verb == "TALK" || cmd.Verb == "SPEAK") {
		return g.talkToApex()
	}

	if strings.EqualFold(cmd.Target, "APEX") && strings.EqualFold(cmd.Verb, "DOOR") {
		return g.apexDoorHint()
	}

	if strings.EqualFold(cmd.Target, "APEX") && cmd.Verb == "THANKS" {
		return "You thank Apex. He grunts and returns to his business."
	}

	if strings.EqualFold(cmd.Target, "ASTAROT") && cmd.Verb != "" {
		return g.astarotTeleport(cmd.Verb)
	}

	if strings.EqualFold(cmd.Target, "MAGOT") && cmd.Verb != "" {
		return g.magotLocate(cmd.Verb)
	}

	if strings.EqualFold(cmd.Target, "ASMODEE") && cmd.Verb != "" {
		return g.asmodeeDestroy(cmd.Verb)
	}

	if strings.EqualFold(cmd.Target, "BELEZBAR") && cmd.Verb != "" {
		return g.belezbarReveal(cmd.Verb)
	}

	if strings.EqualFold(cmd.Target, "GUARDS") && cmd.Verb == "DOOR" {
		return g.passGuards()
	}

	if strings.EqualFold(cmd.Target, "WATER") && strings.EqualFold(cmd.Verb, "FALL") {
		return g.passWater()
	}

	if strings.EqualFold(cmd.Target, "NEST") && strings.EqualFold(cmd.Verb, "PHOENIX") {
		return g.nestPhoenix()
	}

	if strings.EqualFold(cmd.Target, "CAULDRON") && strings.EqualFold(cmd.Verb, "ACHAD") {
		return g.cauldronAchad()
	}

	if dir, ok := world.ParseDirection(cmd.Verb); ok {
		return g.move(dir)
	}

	if !parser.KnownWord(cmd.Verb) && !parser.KnownWord(cmd.Target) {
		return "I don't understand that word."
	}
	return "I recognize that word, but don't know what it does yet (verb resolution not yet implemented)."
}

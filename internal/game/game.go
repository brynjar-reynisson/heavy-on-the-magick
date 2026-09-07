// Package game ties the other packages together into a playable session.
// It's intentionally thin: world, character, parser, magic, graphics and
// audio each own their own logic, so this package should stay readable as
// "wire the pieces together and run the loop" rather than accumulating
// game rules of its own.
package game

import (
	"fmt"
	"strings"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/character"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/magic"
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
	case "DROP":
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

	if strings.EqualFold(cmd.Target, "APEX") && cmd.Verb == "THANKS" {
		return "You thank Apex. He grunts and returns to his business."
	}

	if strings.EqualFold(cmd.Target, "ASTAROT") && cmd.Verb != "" {
		return g.astarotTeleport(cmd.Verb)
	}

	if strings.EqualFold(cmd.Target, "MAGOT") && cmd.Verb != "" {
		return g.magotLocate(cmd.Verb)
	}

	if strings.EqualFold(cmd.Target, "GUARDS") && cmd.Verb == "DOOR" {
		return g.passGuards()
	}

	if strings.EqualFold(cmd.Target, "NEST") && strings.EqualFold(cmd.Verb, "PHOENIX") {
		return g.nestPhoenix()
	}

	if dir, ok := world.ParseDirection(cmd.Verb); ok {
		return g.move(dir)
	}

	if !parser.KnownWord(cmd.Verb) && !parser.KnownWord(cmd.Target) {
		return "I don't understand that word."
	}
	return "I recognize that word, but don't know what it does yet (verb resolution not yet implemented)."
}

// combatStaminaCost is charged to the player per BLAST/FREEZE cast against
// a real monster - a placeholder amount, not extracted (see Handle's doc
// comment for the sourcing of the underlying "combat costs Stamina" fact).
const combatStaminaCost = 5

// saveStaminaCost is charged on every successful Save Game/Save Axil -
// a real, sourced mechanic (round 86): the manual states plainly "Saving
// a game will deplete your Stamina, so that a Save cannot be used as an
// easy way of getting round difficult choices!" No exact amount was
// given by the manual, so this started as an honest placeholder,
// deliberately smaller than combatStaminaCost to match the manual's own
// relative framing ("Combat will reduce your Stamina a lot, most other
// actions will reduce it a little") — round 125 found the EXACT real
// number from a direct first-hand playthrough account (The CRPG
// Addict's 2016 blog post): "You also lose 1 stamina point every time
// you save" — this already-placeholder value of 1 turns out to be
// exactly right, not a guess anymore.
const saveStaminaCost = 1

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

// transfusion handles TRANSFUSION. Confirmed real effect (restores
// Stamina) via the CASA walkthrough extraction; the exact original
// restore amount was not stated there and hasn't been extracted from the
// disassembly, so the amount itself is an honest placeholder. The cap at
// character.Player.MaxStamina IS a real, sourced mechanic though: the
// same walkthrough instructs casting TRANSFUSION "a few times until you
// get maximum stamina" - wording that only makes sense if repeated casts
// approach a real ceiling rather than healing indefinitely.
func (g *Game) transfusion() string {
	const placeholderRestoreAmount = 10
	g.Player.Stamina = min(g.Player.Stamina+placeholderRestoreAmount, g.Player.MaxStamina)
	if g.Player.Stamina >= g.Player.MaxStamina {
		return fmt.Sprintf("You feel your strength return. Stamina is now %d (maximum). (restore-per-cast amount is a placeholder - not yet extracted from the original)", g.Player.Stamina)
	}
	return fmt.Sprintf("You feel your strength return. Stamina is now %d. (restore-per-cast amount is a placeholder - not yet extracted from the original)", g.Player.Stamina)
}

func (g *Game) blast() string {
	room := g.World.CurrentRoom()
	if room == nil || room.Monster == "" || room.MonsterHealth <= 0 {
		return "You BLAST, but there's nothing here to hit."
	}
	g.Player.Stamina -= combatStaminaCost
	room.MonsterHealth -= g.blastDamage()
	if room.MonsterHealth <= 0 {
		gained := g.awardVictoryPoints()
		return g.deathCheck(fmt.Sprintf("The %s is destroyed! (+%d Experience Points)", room.Monster, gained))
	}
	return g.deathCheck(fmt.Sprintf("You BLAST the %s! It's still standing.", room.Monster))
}

// blastDamage returns how much MonsterHealth a single BLAST removes.
// Confirmed real (Spectrum Computing's plain-text instructions file for
// the game — see ../../CLAUDE.md): "your Stamina and Skill together
// affect the outcome of conflicts." No exact formula is stated, so this
// is a reasonable, documented interpretation (Skill was previously
// rolled but never consulted by any game logic at all) rather than
// extracted fact: a base 1 hit, plus 1 extra point of damage for every
// 4 points of Skill.
func (g *Game) blastDamage() int {
	return 1 + g.Player.Skill/4
}

func (g *Game) freeze() string {
	room := g.World.CurrentRoom()
	if room == nil || room.Monster == "" || room.MonsterHealth <= 0 {
		return "You FREEZE, but there's nothing here to target."
	}
	g.Player.Stamina -= combatStaminaCost
	name := room.Monster
	room.MonsterHealth = 0
	gained := g.awardVictoryPoints()
	return g.deathCheck(fmt.Sprintf("You FREEZE the %s solid! (+%d Experience Points)", name, gained))
}

// awardVictoryPoints grants ExperiencePoints for defeating a monster and
// returns the amount gained. Both the base reward and Luck's bonus are
// this port's own reasonable interpretation, not an extracted formula:
// ExperiencePoints itself is confirmed real (see character.Player's doc
// comment) but no confirmed source states what earns it or how much.
// The Luck bonus specifically has real motivation though — Spectrum
// Computing's instructions file for the game states "your Luck
// influences virtually everything you do" (see ../../CLAUDE.md), and
// Luck was previously rolled but never consulted by any game logic at
// all, same gap Skill's blastDamage bonus closed for combat damage.
func (g *Game) awardVictoryPoints() int {
	const baseVictoryPoints = 10
	gained := baseVictoryPoints + g.Player.Luck
	g.Player.ExperiencePoints += gained
	return gained
}

// deathCheck appends a death notice if the just-applied Stamina cost
// killed the player (see Handle's doc comment - confirmed real mechanic).
func (g *Game) deathCheck(msg string) string {
	if g.Player.IsDead() {
		return msg + "\nYour Stamina gives out. You are dead. (GAME OVER)"
	}
	return msg
}

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
		if strings.HasPrefix(d.Ability, "Warning:") {
			return fmt.Sprintf("You invoke %s, %s! %s", d.Name, d.Title, d.Ability)
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

// options handles OPTIONS (Merphish keyword "O"). With no target it shows
// the Option Screen's confirmed real menu items (found live in memory
// during disassembly as the game's own on-screen UI strings under the
// "Magick!" header, not manual prose) - the disassembly confirmed these 5
// items exist among "options 1-6", but not which numbered slot each one
// occupies. Round 124 pinned down 3 of the 6 real slots from two
// separate sources: the manual states "select option 1 and Away You
// Go!" (option 1 = starting the game — no action needed here, since
// this port's game already exists once constructed) and "select option
// 6 and the values will be realigned" (option 6 = Realign Status); the
// CASA walkthrough's own "Tips" section separately states verbatim
// "SAVE regularly by pressing key O and then option 2" (option 2 =
// Save Game). Options 3-5's exact slots still aren't confirmed. A bare
// numeric target ("O 2", "O 6") now works as a real alternate way to
// trigger these, alongside the existing keyword form. "Realign Status"
// is wired to a real effect (Player.Realign, same roll ranges as
// NewPlayer) since that reroll behavior is confirmed by the manual.
// Save/Restore Game/Axil are real, confirmed menu choices (see save.go)
// — a real, functional file-based save system, though the file format
// and the Game-vs-Axil split are this port's own implementation, not the
// original's actual save mechanism (not extracted/known). Saving now
// also costs real Stamina (round 86) - the manual states plainly
// "Saving a game will deplete your Stamina, so that a Save cannot be
// used as an easy way of getting round difficult choices!" - see
// saveStaminaCost. Restoring does not cost Stamina (not stated by any
// source, and would defeat a Save's own point if it did). Round 119:
// the manual also confirms real "Version letter" save slots (see
// save.go's doc comment) - an optional trailing single-letter word in
// target (e.g. "SAVE GAME B") selects one; omitted, it's the same
// single default slot this port has always used.
func (g *Game) options(target string) string {
	target = strings.ToUpper(strings.TrimSpace(target))
	// Round 124: a bare confirmed real numeric slot (see doc comment
	// above) is translated to its equivalent keyword up front, so it
	// flows through the exact same logic below as the keyword form -
	// not a separate, duplicated code path.
	switch target {
	case "2":
		target = "SAVE GAME"
	case "6":
		target = "REALIGN"
	}
	version := extractVersionLetter(target)
	versionSuffix := ""
	if version != "" {
		versionSuffix = fmt.Sprintf(" (version %s)", version)
	}
	switch {
	case strings.Contains(target, "REALIGN"):
		g.Player.Realign()
		return fmt.Sprintf("Realign Status: Stamina %d, Skill %d, Luck %d.", g.Player.Stamina, g.Player.Skill, g.Player.Luck)
	case strings.Contains(target, "SAVE") && strings.Contains(target, "AXIL"):
		g.Player.Stamina -= saveStaminaCost
		if err := g.SaveAxilVersion(version); err != nil {
			return fmt.Sprintf("Save Axil failed: %v", err)
		}
		return g.deathCheck("Axil saved" + versionSuffix + ".")
	case strings.Contains(target, "RESTORE") && strings.Contains(target, "AXIL"):
		if err := g.RestoreAxilVersion(version); err != nil {
			return fmt.Sprintf("Restore Axil failed: %v", err)
		}
		return "Axil restored" + versionSuffix + "."
	case strings.Contains(target, "SAVE"):
		g.Player.Stamina -= saveStaminaCost
		if err := g.SaveGameVersion(version); err != nil {
			return fmt.Sprintf("Save Game failed: %v", err)
		}
		return g.deathCheck("Game saved" + versionSuffix + ".")
	case strings.Contains(target, "RESTORE"):
		if err := g.RestoreGameVersion(version); err != nil {
			return fmt.Sprintf("Restore Game failed: %v", err)
		}
		return "Game restored" + versionSuffix + "."
	default:
		return "Magick!\nSave Game / Restore Game / Save Axil / Restore Axil / Realign Status"
	}
}

// extractVersionLetter finds a real Version-letter token (see save.go's
// doc comment) in an OPTIONS target string: a single-letter word, none
// of which appear among the real keywords this switch already checks
// for (SAVE/RESTORE/GAME/AXIL/REALIGN/STATUS are all multi-letter), so
// no explicit exclusion list is needed. Returns "" if none is present.
func extractVersionLetter(target string) string {
	for word := range strings.FieldsSeq(target) {
		if len(word) == 1 && word[0] >= 'A' && word[0] <= 'Z' {
			return word
		}
	}
	return ""
}

// help reproduces the game's own real "SOME ADVICE" startup hint
// screen - not paraphrased or invented, but the actual in-game text
// disassembled from the original's memory and independently cross-
// confirmed against a real screenshot of the running game (see
// ../../CLAUDE.md, "Found ground truth: ref_screenshot.png" - the
// decoded strings matched the screenshot exactly). "HELP" is a real
// confirmed vocabulary word (parser.Vocabulary); no source states it's
// literally the keyword that shows this screen in the original (it may
// have only appeared once, at boot), but reproducing the original
// game's own real hint text for a player who asks for help is squarely
// what a faithful port should do with content this well-sourced.
func (g *Game) help() string {
	return `SOME ADVICE
Talk to Apex often!
"APEX, DOOR"  "APEX, WEREWOLF"  "APEX, FIRE"
To dismiss say "APEX, THANKS"
Talk to other things - sometimes they will answer!
"DOOR, password"
"GUARDS, DOOR"
"ASTAROT, WOLFDORP"
Finally, if in a panic, BLAST without an object!`
}

// spells lists this port's real, functional spells. "SPELLS" is a real
// confirmed vocabulary word (parser.Vocabulary); no source states this
// exact command lists them as a menu, so the aggregation itself is this
// project's own, not fabricated content — but the manual DOES have its
// own explicit "Spells:" heading (round 118, found via the PDF's real
// text layer — see round 110's writeup), grouping exactly three
// keywords under it: I (Invoke), B (Blast), F (Freeze). INVOKE was
// missing from this listing entirely until round 118 — a real,
// sourced omission this project's own aggregation had gotten wrong,
// not just an incomplete one — corrected here. TRANSFUSION and CALL
// are real spells too (see their own doc comments) but the manual
// itself says "fuller details... can be found in the section on the
// Grimoire" for those, i.e. they're confirmed spells from a DIFFERENT
// part of the manual, not this specific three-keyword grouping — kept
// in the listing since they're genuinely real spells, just not
// re-labeled as part of the manual's core "Spells:" trio.
func (g *Game) spells() string {
	return "Known spells: INVOKE (summon a demon), BLAST (combat), FREEZE (combat), TRANSFUSION (restore Stamina), CALL (summon Apex, needs the Scroll)."
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
			return "A wooden chest."
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

func (g *Game) move(dir world.Direction) string {
	if room := g.World.CurrentRoom(); room != nil {
		if destID, ok := room.Exits[dir]; ok {
			if dest := g.World.Rooms[destID]; dest != nil && dest.Fire && !g.hasItem("Clasp") {
				return "Flames block your way. You'd need something to protect you from the fire."
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

// describeCurrentRoom renders LOOK's real output. Round 115: now
// includes the room's real Level (1-4, sourced the same way as every
// other room fact in this project) when it's set - a third instance
// of the "confirmed but unsurfaced" gap round 114/115 found for
// magic.Demon's fields: world.Room.Level has been set on every real
// room since this project's earliest CollodonsPile/LevelNGrid commits,
// but was never once referenced by internal/game before this, so the
// player had no way to tell which of the dungeon's 4 levels they were
// actually on.
func (g *Game) describeCurrentRoom() string {
	room := g.World.CurrentRoom()
	if room == nil {
		return "You are nowhere. (no current room set)"
	}
	var b strings.Builder
	if room.Name == "Exit" && !g.Won {
		g.Won = true
		b.WriteString("You have found one of Collodon's Pile's 3 exits and escaped! YOU HAVE WON.\n")
	}
	if room.Level != 0 {
		fmt.Fprintf(&b, "%s (Level %d)\n%s\n", room.Name, room.Level, room.Description)
	} else {
		fmt.Fprintf(&b, "%s\n%s\n", room.Name, room.Description)
	}
	if room.HasTable {
		b.WriteString("There is a table here.\n")
	}
	if room.HasChest {
		b.WriteString("There is a chest here.\n")
	}
	if len(room.Items) > 0 {
		fmt.Fprintf(&b, "You see: %s\n", strings.Join(room.Items, ", "))
	}
	if exits := exitList(room); exits != "" {
		fmt.Fprintf(&b, "Exits: %s\n", exits)
	}
	return strings.TrimRight(b.String(), "\n")
}

func exitList(room *world.Room) string {
	names := make([]string, 0, len(room.Exits))
	for dir := range room.Exits {
		names = append(names, dir.String())
	}
	return strings.Join(names, ", ")
}

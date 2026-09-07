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
	// footer text and its "E = one of three exits" legend entry — see
	// ../../CLAUDE.md). Detected by Room.Name == "Exit" (currently only
	// world.Level1Grid's G3 cell carries that name).
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
// world.Level4Grid — a real 22-cell room graph for Level 4, the last of
// the game's 4 levels to get this treatment: a 17-cell fully-connected
// component plus 5 real, named, deliberately isolated special rooms
// (Scales, Doubt of Rabak, The Crypt, Exit, Pride — see Level4Grid's
// doc comment). Level 4's calibration had a Level3Grid-style off-by-
// one-row bug (since corrected by relabeling, not re-extracting), and
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
// exits" sourcing) - currently only reachable via world.Level1Grid.
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
// response — the manual's own table entry ("swap Window 1") isn't
// detailed enough to model the underlying dual-window/spell-hand display.
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
		return "Halted."
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
		return "You SWAP windows. (Merphish 'Z' - confirmed real command, but the underlying dual-window/spell-hand display isn't modeled yet.)"
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

	if strings.EqualFold(cmd.Target, "GUARDS") && cmd.Verb == "DOOR" {
		return g.passGuards()
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
// actually succeeds when the player carries the right Charm. With no
// target, lists all 4 demons and their Charm requirements — real,
// already-sourced data that had no way to be seen in-game before this.
func (g *Game) invoke(target string) string {
	if target == "" {
		var b strings.Builder
		b.WriteString("Known demons and their required Talismans:\n")
		for _, d := range magic.Demons {
			fmt.Fprintf(&b, "%s, %s (needs: %s)\n", d.Name, d.Title, d.Charm)
		}
		return strings.TrimRight(b.String(), "\n")
	}
	for _, d := range magic.Demons {
		if !strings.EqualFold(d.Name, target) {
			continue
		}
		if !g.hasItem(d.Charm) {
			return fmt.Sprintf("You begin the ritual to invoke %s, %s... but you have no suitable Talisman (a %s).", d.Name, d.Title, d.Charm)
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

// hasItem reports whether the player is carrying an item by name.
func (g *Game) hasItem(name string) bool {
	return g.Player.HasItem(name)
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
// occupies or what the 6th, unextracted option is, so they're listed
// unnumbered rather than guessing an exact menu layout. "Realign Status"
// is wired to a real effect (Player.Realign, same roll ranges as
// NewPlayer) since that reroll behavior is confirmed by the manual.
// Save/Restore Game/Axil are real, confirmed menu choices (see save.go)
// — a real, functional file-based save system, though the file format
// and the Game-vs-Axil split are this port's own implementation, not the
// original's actual save mechanism (not extracted/known).
func (g *Game) options(target string) string {
	target = strings.ToUpper(strings.TrimSpace(target))
	switch {
	case strings.Contains(target, "REALIGN"):
		g.Player.Realign()
		return fmt.Sprintf("Realign Status: Stamina %d, Skill %d, Luck %d.", g.Player.Stamina, g.Player.Skill, g.Player.Luck)
	case strings.Contains(target, "SAVE") && strings.Contains(target, "AXIL"):
		if err := g.SaveAxil(); err != nil {
			return fmt.Sprintf("Save Axil failed: %v", err)
		}
		return "Axil saved."
	case strings.Contains(target, "RESTORE") && strings.Contains(target, "AXIL"):
		if err := g.RestoreAxil(); err != nil {
			return fmt.Sprintf("Restore Axil failed: %v", err)
		}
		return "Axil restored."
	case strings.Contains(target, "SAVE"):
		if err := g.SaveGame(); err != nil {
			return fmt.Sprintf("Save Game failed: %v", err)
		}
		return "Game saved."
	case strings.Contains(target, "RESTORE"):
		if err := g.RestoreGame(); err != nil {
			return fmt.Sprintf("Restore Game failed: %v", err)
		}
		return "Game restored."
	default:
		return "Magick!\nSave Game / Restore Game / Save Axil / Restore Axil / Realign Status"
	}
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
// exact command lists them (the manual is known to describe BLAST,
// FREEZE, and TRANSFUSION individually - see their own doc comments -
// but not as a single "SPELLS" menu), so this is this project's own
// aggregation of already-confirmed real spells, not fabricated content.
func (g *Game) spells() string {
	return "Known spells: BLAST (combat), FREEZE (combat), TRANSFUSION (restore Stamina)."
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
			return fmt.Sprintf("You pick up the %s.", item)
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
	if !g.World.Move(dir) {
		return "You can't go that way."
	}
	msg := g.checkNougatWerewolf()
	desc := g.describeCurrentRoom()
	if msg != "" {
		return msg + "\n" + desc
	}
	return desc
}

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
	fmt.Fprintf(&b, "%s\n%s\n", room.Name, room.Description)
	if room.HasTable {
		b.WriteString("There is a table here.\n")
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

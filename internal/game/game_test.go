package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/character"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/world"
)

func TestHandleUnknownWord(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("GOLANG, COMPILE"))
	if !strings.Contains(got, "don't understand") {
		t.Errorf("Handle(unknown word) = %q, want a not-understood response", got)
	}
}

func TestHandleKnownWordNoBehaviorYet(t *testing.T) {
	g := New()
	// TALK to MONSTER (as opposed to APEX, which now has real behavior -
	// see TestHandleTalkToApex) is real confirmed vocabulary but has no
	// verb resolution yet.
	got := g.Handle(parser.Parse("MONSTER, TALK"))
	if !strings.Contains(got, "don't know what it does") {
		t.Errorf("Handle(known word, unimplemented verb) = %q, want the not-yet-implemented response", got)
	}
}

func TestHandleTalkToApex(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("APEX, TALK"))
	if !strings.Contains(got, "Apex") {
		t.Errorf("Handle(APEX, TALK) = %q, want it to recognize the real NPC Apex the Ogre", got)
	}
}

func TestHandleSpeakToApexIsSynonymForTalk(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("APEX, SPEAK"))
	if !strings.Contains(got, "Apex") {
		t.Errorf("Handle(APEX, SPEAK) = %q, want the same recognition as APEX, TALK", got)
	}
}

func TestHandleCarryIsSynonymForPickup(t *testing.T) {
	g := New() // Room of Misery has a real sourced item: Grimoire
	got := g.Handle(parser.Parse("CARRY GRIMOIRE"))
	if !strings.Contains(got, "Grimoire") || !g.hasItem("Grimoire") {
		t.Errorf("Handle(CARRY GRIMOIRE) = %q, want it to pick up the Grimoire", got)
	}
}

func TestHandleNameReportsPlayerName(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("NAME"))
	if !strings.Contains(got, "Axil") {
		t.Errorf("Handle(NAME) = %q, want it to mention Axil", got)
	}
}

// TestPassGuardsClearsRealObstacle covers the confirmed real command
// "GUARDS, DOOR" (see world.Room.Guards and passGuards doc comments) —
// the same TARGET-comma-VERB grammar as "APEX, TALK".
func TestPassGuardsClearsRealObstacle(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Guards = true
	got := g.Handle(parser.Parse("GUARDS, DOOR"))
	if !strings.Contains(got, "let you pass") {
		t.Errorf("Handle(GUARDS, DOOR) with real guards present = %q, want them to let the player pass", got)
	}
	if g.World.CurrentRoom().Guards {
		t.Error("Guards should be cleared after a successful GUARDS, DOOR")
	}
}

func TestPassGuardsWithNoGuardsPresent(t *testing.T) {
	g := New() // Room of Misery has no Guards
	got := g.Handle(parser.Parse("GUARDS, DOOR"))
	if !strings.Contains(got, "no guards") {
		t.Errorf("Handle(GUARDS, DOOR) with no guards present = %q, want it to say there are none", got)
	}
}

// TestHandleInvokeSucceedsWithCharm pins round 131's correction: the
// real mechanic (a genuinely new source, World of Spectrum's separate
// plain-text instructions file) states "Place Ye the talisman on the
// ground and proceed with thy invocation from a distance" - the Charm
// must be dropped in the room, not merely carried.
func TestHandleInvokeSucceedsWithCharm(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Sunflower")
	got := g.Handle(parser.Parse("I MAGOT")) // Magot's confirmed Charm is Sunflower
	if strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(I MAGOT) with Sunflower on the ground = %q, want the invocation to succeed", got)
	}
	if !strings.Contains(got, "MAGOT") {
		t.Errorf("Handle(I MAGOT) = %q, want it to name Magot", got)
	}
}

// TestHandleInvokeCarriedNotDroppedCharmFails pins the other half of
// round 131's correction: merely CARRYING the Charm isn't enough, and
// gets a distinct, honest hint (not the furnace-room punishment, which
// is reserved for not having the Charm at all).
func TestHandleInvokeCarriedNotDroppedCharmFails(t *testing.T) {
	g := New()
	g.Player.Items = append(g.Player.Items, "Sunflower")
	got := g.Handle(parser.Parse("I MAGOT"))
	if !strings.Contains(got, "place it on the ground") {
		t.Errorf("Handle(I MAGOT) with Sunflower only carried = %q, want a place-it-on-the-ground hint", got)
	}
	if strings.Contains(got, "furnace room") {
		t.Errorf("Handle(I MAGOT) with the Charm carried (not dropped) = %q, want it to NOT trigger the no-Charm-at-all punishment", got)
	}
}

// TestHandleInvokeWithoutCharmTeleportsToFurnaceRoom pins round 126's
// real, sourced punishment mechanic (The CRPG Addict's first-hand
// playthrough account, the same source round 125 used to resolve
// CALL's effect): invoking a demon without its Charm doesn't just
// reject the command, it actually teleports the player to the real
// Furnace Room. Verifies both the message and the actual World state
// change, not just text.
func TestHandleInvokeWithoutCharmTeleportsToFurnaceRoom(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("I MAGOT")) // no Sunflower carried
	if !strings.Contains(got, "furnace room") {
		t.Errorf("Handle(I MAGOT) with no Charm = %q, want it to mention being sent to the furnace room", got)
	}
	if room := g.World.CurrentRoom(); room == nil || room.Name != "Furnace Room" {
		t.Errorf("current room after a failed INVOKE = %+v, want Furnace Room", room)
	}
}

// TestHandleApexThanksDismisses pins the hint screen's own confirmed
// dismiss phrase ("To dismiss say \"APEX, THANKS\"" — see game.help's
// verbatim text), wired for real for the first time this round.
func TestHandleApexThanksDismisses(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("APEX, THANKS"))
	if !strings.Contains(got, "thank Apex") {
		t.Errorf("Handle(APEX, THANKS) = %q, want a real dismiss response", got)
	}
}

func TestHandleAstarotTeleportRequiresSword(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("ASTAROT, WOLFDORP"))
	if !strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(ASTAROT, WOLFDORP) with no Sword carried = %q, want a Talisman rejection", got)
	}
}

// TestHandleAstarotTeleportSucceeds pins the real, hint-screen-confirmed
// example command ("ASTAROT, WOLFDORP" — see parser.Parse's package doc
// comment) actually teleporting the player, once they carry Astarot's
// confirmed Charm (Sword).
func TestHandleAstarotTeleportSucceeds(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Sword")
	got := g.Handle(parser.Parse("ASTAROT, WOLFDORP"))
	if strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(ASTAROT, WOLFDORP) with Sword carried = %q, want the teleport to succeed", got)
	}
	if !strings.Contains(got, "Wolfdorp") {
		t.Errorf("Handle(ASTAROT, WOLFDORP) = %q, want it to name Wolfdorp", got)
	}
	if got := g.World.CurrentRoom(); got == nil || got.Name != "Wolfdorp" {
		t.Errorf("after ASTAROT, WOLFDORP, current room = %+v, want Wolfdorp", got)
	}
}

func TestHandleAstarotTeleportUnknownLocation(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Sword")
	got := g.Handle(parser.Parse("ASTAROT, NARNIA"))
	if !strings.Contains(got, "doesn't recognize") {
		t.Errorf("Handle(ASTAROT, NARNIA) = %q, want an honest unknown-location rejection", got)
	}
	if got := g.World.CurrentRoom(); got == nil || got.Name != "Room of Misery" {
		t.Errorf("after an unknown-location ASTAROT command, current room = %+v, want unchanged (Room of Misery)", got)
	}
}

func TestHandleMagotLocateRequiresSunflower(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("MAGOT, GRIMOIRE"))
	if !strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(MAGOT, GRIMOIRE) with no Sunflower carried = %q, want a Talisman rejection", got)
	}
}

// TestHandleMagotLocateFindsRealItem pins the natural-inference "MAGOT,
// <object>" grammar (see magotLocate's doc comment) actually finding a
// real item's real room, once the player carries Magot's confirmed
// Charm (Sunflower). Grimoire is a real, already-shipped item in Room
// of Misery.
func TestHandleMagotLocateFindsRealItem(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Sunflower")
	got := g.Handle(parser.Parse("MAGOT, GRIMOIRE"))
	if strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(MAGOT, GRIMOIRE) with Sunflower carried = %q, want the locate to succeed", got)
	}
	if !strings.Contains(got, "Room of Misery") {
		t.Errorf("Handle(MAGOT, GRIMOIRE) = %q, want it to name Room of Misery", got)
	}
}

func TestHandleMagotLocateAlreadyCarried(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Sunflower")
	g.Player.Items = append(g.Player.Items, "Grimoire")
	got := g.Handle(parser.Parse("MAGOT, GRIMOIRE"))
	if !strings.Contains(got, "already carry") {
		t.Errorf("Handle(MAGOT, GRIMOIRE) while carrying it = %q, want an honest already-carried response", got)
	}
	// Regression: must echo the item's real stored casing ("Grimoire"),
	// not the player's raw uppercased typed target ("GRIMOIRE") - the
	// same casing bug already fixed once for examine().
	if !strings.Contains(got, "Grimoire") || strings.Contains(got, "the GRIMOIRE") {
		t.Errorf("Handle(MAGOT, GRIMOIRE) while carrying it = %q, want it to echo real casing \"Grimoire\"", got)
	}
}

func TestHandleMagotLocateUnknownObject(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Sunflower")
	got := g.Handle(parser.Parse("MAGOT, EXCALIBUR"))
	if !strings.Contains(got, "senses no such object") {
		t.Errorf("Handle(MAGOT, EXCALIBUR) = %q, want an honest not-found response", got)
	}
}

// TestHandleFireBlocksMovementWithoutClasp covers the real, sourced
// Fire mechanic (see world.Room.Fire's doc comment): the CASA
// walkthrough states the Clasp "enables you to walk through fire" -
// modeled as blocking movement into a Fire room without it. Uses a
// synthetic 2-room world (world.Level2Grid's real D6 Fire cell is
// currently isolated, unreachable via ordinary movement) so this real
// mechanic is exercised end-to-end even though it can't be in the
// shipped data yet.
func TestHandleFireBlocksMovementWithoutClasp(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Start", Exits: map[world.Direction]world.RoomID{world.North: 1}})
	w.AddRoom(&world.Room{ID: 1, Name: "Blaze", Fire: true})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("NORTH"))
	if !strings.Contains(got, "Flames block") {
		t.Errorf("Handle(NORTH) into a Fire room without Clasp = %q, want it blocked", got)
	}
	if g.World.CurrentRoom().Name != "Start" {
		t.Errorf("current room after a blocked Fire move = %q, want unchanged (Start)", g.World.CurrentRoom().Name)
	}

	g.Player.Items = append(g.Player.Items, "Clasp")
	got = g.Handle(parser.Parse("NORTH"))
	if strings.Contains(got, "Flames block") {
		t.Errorf("Handle(NORTH) into a Fire room WITH Clasp = %q, want it to succeed", got)
	}
	if g.World.CurrentRoom().Name != "Blaze" {
		t.Errorf("current room after a Clasp-protected Fire move = %q, want \"Blaze\"", g.World.CurrentRoom().Name)
	}
}

// TestHandleSwapItemRevealsRealItem covers round 132's real, sourced
// "protected item" mechanic (see world.Room.SwapItem's doc comment):
// dropping the room's real SwapItem reveals its RevealItem. Uses a
// synthetic room (the exact numbered-map cell this triad corresponds
// to hasn't been cross-referenced to a shipped room yet) so this real
// mechanic is exercised end-to-end even though it can't be in the
// shipped data yet - same pattern as TestHandleFireBlocksMovementWithoutClasp.
func TestHandleSwapItemRevealsRealItem(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Cache", SwapItem: "Ball", RevealItem: "Pellet"})
	g := &Game{Player: character.NewPlayer(), World: w}
	g.Player.Items = append(g.Player.Items, "Ball")

	got := g.Handle(parser.Parse("DROP BALL"))
	if !strings.Contains(got, "Pellet") {
		t.Errorf("Handle(DROP BALL) with SwapItem=Ball/RevealItem=Pellet = %q, want it to mention the Pellet", got)
	}
	room := g.World.CurrentRoom()
	found := false
	for _, item := range room.Items {
		if item == "Pellet" {
			found = true
		}
	}
	if !found {
		t.Errorf("room Items after the swap = %v, want it to include Pellet", room.Items)
	}
	if room.SwapItem != "" || room.RevealItem != "" {
		t.Errorf("SwapItem/RevealItem after the swap = %q/%q, want both cleared (one-time reveal)", room.SwapItem, room.RevealItem)
	}
}

// TestHandleSwapItemUnrelatedDropDoesNothing is a regression guard:
// dropping an item that ISN'T the room's real SwapItem must not
// trigger a reveal.
func TestHandleSwapItemUnrelatedDropDoesNothing(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Cache", SwapItem: "Ball", RevealItem: "Pellet"})
	g := &Game{Player: character.NewPlayer(), World: w}
	g.Player.Items = append(g.Player.Items, "Grimoire")

	g.Handle(parser.Parse("DROP GRIMOIRE"))
	room := g.World.CurrentRoom()
	for _, item := range room.Items {
		if item == "Pellet" {
			t.Errorf("dropping an unrelated item revealed Pellet early: room Items = %v", room.Items)
		}
	}
	if room.SwapItem != "Ball" || room.RevealItem != "Pellet" {
		t.Errorf("SwapItem/RevealItem after an unrelated drop = %q/%q, want both unchanged", room.SwapItem, room.RevealItem)
	}
}

func TestHandleMovementValidExit(t *testing.T) {
	g := New()
	// Room of Misery --East--> Secunda Porta, per the real (walkthrough-
	// sourced) room graph in world.CollodonsPile.
	got := g.Handle(parser.Parse("EAST"))
	if !strings.Contains(got, "Secunda Porta") {
		t.Errorf("Handle(EAST) from the starting room = %q, want it to describe Secunda Porta", got)
	}
}

func TestHandleMovementInvalidExit(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("SOUTH")) // Room of Misery has no South exit
	if got != "You can't go that way." {
		t.Errorf("Handle(SOUTH) with no such exit = %q, want rejection", got)
	}
}

func TestHandleLook(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "Room of Misery") {
		t.Errorf("Handle(LOOK) = %q, want the starting room's description", got)
	}
}

// TestHandleLookShowsLevel pins round 115: world.Room.Level (real data,
// set since this project's earliest room commits but never shown to
// the player before) now appears in LOOK's output. Room of Misery is
// real, confirmed Level 2.
func TestHandleLookShowsLevel(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "Level 2") {
		t.Errorf("Handle(LOOK) = %q, want it to show Room of Misery's real Level (2)", got)
	}
}

func TestHandleLookShowsRoomItems(t *testing.T) {
	g := New() // Room of Misery has a real sourced item: Grimoire
	got := g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "Grimoire") {
		t.Errorf("Handle(LOOK) = %q, want it to mention the Grimoire present in the room", got)
	}
}

func TestHandlePickUpTwoWordForm(t *testing.T) {
	g := New() // Room of Misery has a real sourced item: Grimoire
	got := g.Handle(parser.Parse("PICK UP GRIMOIRE"))
	if !strings.Contains(got, "Grimoire") {
		t.Errorf("Handle(PICK UP GRIMOIRE) = %q, want it to pick up the Grimoire", got)
	}
	if !g.hasItem("Grimoire") {
		t.Error("expected the player to be carrying the Grimoire after PICK UP GRIMOIRE")
	}
}

func TestHandleTakeAndLiftAreSynonymsForPickup(t *testing.T) {
	for _, verb := range []string{"TAKE", "LIFT"} {
		g := New()
		got := g.Handle(parser.Parse(verb + " GRIMOIRE"))
		if !strings.Contains(got, "Grimoire") || !g.hasItem("Grimoire") {
			t.Errorf("Handle(%s GRIMOIRE) = %q, want it to pick up the Grimoire", verb, got)
		}
	}
}

// TestHandleHelpShowsRealHintScreen covers the real, disassembled,
// screenshot-cross-confirmed "SOME ADVICE" in-game hint screen (see
// help's doc comment) — the game's own actual text, not invented.
func TestHandleHelpShowsRealHintScreen(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("HELP"))
	for _, want := range []string{"SOME ADVICE", "APEX, THANKS", "GUARDS, DOOR", "BLAST without an object"} {
		if !strings.Contains(got, want) {
			t.Errorf("Handle(HELP) = %q, want it to contain %q", got, want)
		}
	}
}

func TestHandleInventory(t *testing.T) {
	g := New()
	// Round 121: Axil starts with a real, confirmed Pouch (see
	// character.NewPlayer's doc comment) - INVENTORY should list it
	// from the very start, not claim nothing is carried.
	if got := g.Handle(parser.Parse("INVENTORY")); !strings.Contains(got, "Pouch") {
		t.Errorf("Handle(INVENTORY) at game start = %q, want it to list the real starting Pouch", got)
	}
	g.Handle(parser.Parse("PICKUP GRIMOIRE"))
	got := g.Handle(parser.Parse("INVENTORY"))
	if !strings.Contains(got, "Grimoire") {
		t.Errorf("Handle(INVENTORY) after picking up the Grimoire = %q, want it listed", got)
	}
}

// TestHandleInventoryEmptyWhenNothingCarried covers the actually-empty
// case (INVENTORY's other real branch) directly on a Player with no
// Items, since NewPlayer's real Pouch means New() itself never starts
// empty anymore.
func TestHandleInventoryEmptyWhenNothingCarried(t *testing.T) {
	g := New()
	g.Player.Items = nil
	if got := g.Handle(parser.Parse("INVENTORY")); !strings.Contains(got, "anything") {
		t.Errorf("Handle(INVENTORY) with nothing carried = %q, want it to say so", got)
	}
}

func TestHandleExamineShowsRoomItems(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("EXAMINE"))
	if !strings.Contains(got, "Grimoire") {
		t.Errorf("Handle(EXAMINE) = %q, want it to mention the Grimoire present in the room", got)
	}
}

// TestHandleExamineWithTargetConfirmsJustThatThing covers the real
// confirmed grammar "X BOTTLE" (see parser's doc comment): a targeted
// EXAMINE should confirm just the named thing, not always list
// everything in the room regardless of what was asked about.
func TestHandleExamineWithTargetConfirmsJustThatThing(t *testing.T) {
	g := New() // Room of Misery has real items: Grimoire, Poison-smeared book
	got := g.Handle(parser.Parse("X GRIMOIRE"))
	if !strings.Contains(got, "Grimoire") {
		t.Errorf("Handle(X GRIMOIRE) = %q, want it to mention the Grimoire", got)
	}
	if strings.Contains(got, "Poison-smeared book") {
		t.Errorf("Handle(X GRIMOIRE) = %q, want it to NOT mention the other item", got)
	}
}

// TestHandleLookMentionsTable covers the real HasTable fixture being
// surfaced in LOOK itself, not just discoverable by blindly guessing
// "EXAMINE TABLE".
func TestHandleLookMentionsTable(t *testing.T) {
	g := New() // Room of Misery has a real confirmed table
	got := g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "table") {
		t.Errorf("Handle(LOOK) in Room of Misery = %q, want it to mention the real table", got)
	}
}

// TestHandleExamineTable covers the real, sourced HasTable fixture (see
// world.Room.HasTable's doc comment) - the CASA walkthrough repeatedly
// uses "EXAMINE TABLE" in specific real rooms.
func TestHandleExamineTable(t *testing.T) {
	g := New() // Room of Misery has a real confirmed table
	got := g.Handle(parser.Parse("X TABLE"))
	if !strings.Contains(got, "table") {
		t.Errorf("Handle(X TABLE) in Room of Misery = %q, want it to acknowledge the real table", got)
	}
}

// walkToWolfdorp navigates a fresh Game from Room of Misery to Wolfdorp
// via real moves - the shared setup TestHandleDropPaysRealToll already
// uses, split out so other Wolfdorp-specific tests (like the HasChest
// ones below) don't have to repeat it.
func walkToWolfdorp(t *testing.T) *Game {
	t.Helper()
	g := New()
	g.Handle(parser.Parse("EAST"))          // Secunda Porta
	g.Handle(parser.Parse("DOOR, SILENCE")) // unlocks the door North
	g.Handle(parser.Parse("NORTH"))         // Trollwynd
	g.Handle(parser.Parse("SOUTH"))         // Sothic Complex
	g.Handle(parser.Parse("SOUTH"))         // Wolfdorp
	if room := g.World.CurrentRoom(); room == nil || room.Name != "Wolfdorp" {
		t.Fatalf("test setup bug: expected to be in Wolfdorp, got %+v", room)
	}
	return g
}

// TestHandleLookMentionsChest covers the real HasChest fixture (round
// 78, see world.Room.HasChest's doc comment) being surfaced in LOOK
// itself, not just discoverable by blindly guessing "EXAMINE CHEST".
func TestHandleLookMentionsChest(t *testing.T) {
	g := walkToWolfdorp(t)
	got := g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "chest") {
		t.Errorf("Handle(LOOK) in Wolfdorp = %q, want it to mention the real chest", got)
	}
}

// TestHandleExamineChest covers the real, sourced HasChest fixture - the
// CASA walkthrough uses "EXAMINE CHEST" before picking up Garlic here.
func TestHandleExamineChest(t *testing.T) {
	g := walkToWolfdorp(t)
	got := g.Handle(parser.Parse("X CHEST"))
	if !strings.Contains(got, "chest") {
		t.Errorf("Handle(X CHEST) in Wolfdorp = %q, want it to acknowledge the real chest", got)
	}
}

func TestHandleExamineWithTargetNotHereSaysSo(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("EXAMINE SWORD"))
	if !strings.Contains(got, "don't see") {
		t.Errorf("Handle(EXAMINE SWORD) with no sword present = %q, want it to say so", got)
	}
}

func TestHandleExamineWithTargetFindsCarriedItem(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("PICKUP GRIMOIRE"))
	got := g.Handle(parser.Parse("EXAMINE GRIMOIRE"))
	if !strings.Contains(got, "carrying") {
		t.Errorf("Handle(EXAMINE GRIMOIRE) after picking it up = %q, want it to recognize the carried item", got)
	}
}

func TestHandleMapTracksExploration(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST")) // visit Secunda Porta
	got := g.Handle(parser.Parse("MAP"))

	if !strings.Contains(got, "ROO") { // "Room of Misery"'s 3-letter label
		t.Errorf("Handle(MAP) = %q, want it to show the visited Room of Misery", got)
	}
	if !strings.Contains(got, "SEC") { // "Secunda Porta"'s 3-letter label
		t.Errorf("Handle(MAP) = %q, want it to show the visited Secunda Porta", got)
	}
	if strings.Contains(got, "TRO") { // Trollwynd was never visited
		t.Errorf("Handle(MAP) = %q, should not show the unvisited Trollwynd", got)
	}
}

func TestHandleDoorPasswordCorrect(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST")) // move to Secunda Porta, which has a door password
	got := g.Handle(parser.Parse("DOOR, SILENCE"))
	if !strings.Contains(got, "door swings open") {
		t.Errorf("Handle(DOOR, SILENCE) in Secunda Porta = %q, want it to succeed", got)
	}
}

func TestHandleDoorPasswordCaseInsensitive(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))
	got := g.Handle(parser.Parse("DOOR, silence"))
	if !strings.Contains(got, "door swings open") {
		t.Errorf("Handle(DOOR, silence) lowercase = %q, want it to still succeed", got)
	}
}

func TestHandleSecundaPortaDoorPromotesToZelator(t *testing.T) {
	g := New()
	if g.Player.Grade != character.Neophyte {
		t.Fatalf("test setup bug: new player Grade = %v, want Neophyte", g.Player.Grade)
	}
	g.Handle(parser.Parse("EAST")) // Secunda Porta
	got := g.Handle(parser.Parse("DOOR, SILENCE"))
	if g.Player.Grade != character.Zelator {
		t.Errorf("Grade after Secunda Porta's door = %v, want Zelator", g.Player.Grade)
	}
	if !strings.Contains(got, "Zelator") {
		t.Errorf("Handle(DOOR, SILENCE) = %q, want it to mention becoming a Zelator", got)
	}
}

func TestHandleDoorPasswordWrong(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST")) // Secunda Porta's real password is SILENCE, not THANKS
	got := g.Handle(parser.Parse("DOOR, THANKS"))
	if got != "Nothing happens." {
		t.Errorf("Handle(DOOR, THANKS) wrong password = %q, want rejection", got)
	}
}

func TestHandleDoorNoPasswordKnown(t *testing.T) {
	g := New() // Room of Misery has no known DoorPasswords
	got := g.Handle(parser.Parse("DOOR, ANYTHING"))
	if strings.Contains(got, "swings open") {
		t.Errorf("Handle(DOOR, ...) in a room with no known password = %q, should not succeed", got)
	}
}

// TestHandleTollDoorRequiresItem uses a synthetic TollItem room to
// isolate the "DOOR, <item>" form specifically (see
// TestHandleDropPaysRealToll below for the real, now-placed rooms and
// the "DROP <item>" form a fresh walkthrough re-read found to be the
// actual trigger phrase — see world.Room.TollItem's doc comment).
func TestHandleTollDoorRequiresItem(t *testing.T) {
	g := New()
	g.World.CurrentRoom().TollItem = "Bag of Gold"

	got := g.Handle(parser.Parse("DOOR, BAG OF GOLD"))
	if !strings.Contains(got, "don't have") {
		t.Errorf("Handle(DOOR, BAG OF GOLD) without the item = %q, want a don't-have response", got)
	}

	g.Player.Items = append(g.Player.Items, "Bag of Gold")
	got = g.Handle(parser.Parse("DOOR, BAG OF GOLD"))
	if !strings.Contains(got, "swings open") {
		t.Errorf("Handle(DOOR, BAG OF GOLD) carrying the item = %q, want it to succeed", got)
	}
	if g.hasItem("Bag of Gold") {
		t.Error("Bag of Gold should be spent (removed from inventory) after paying the toll")
	}
}

// TestHandleDropPaysRealToll covers the real, now-placed TollItem rooms
// (round 64) and confirms "DROP <item>" - the phrase a fresh CASA
// walkthrough re-read found is the actual trigger - opens the door.
func TestHandleDropPaysRealToll(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST")) // Secunda Porta
	g.Handle(parser.Parse("DOOR, SILENCE"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd
	g.Handle(parser.Parse("SOUTH")) // Sothic Complex
	g.Handle(parser.Parse("SOUTH")) // Wolfdorp
	room := g.World.CurrentRoom()
	if room.Name != "Wolfdorp" {
		t.Fatalf("test setup bug: expected to be in Wolfdorp, got %q", room.Name)
	}
	g.Handle(parser.Parse("PICKUP BAG"))
	g.Handle(parser.Parse("NORTH-WEST")) // Room of Stings
	if g.World.CurrentRoom().Name != "Room of Stings" {
		t.Fatalf("test setup bug: expected to be in Room of Stings, got %q", g.World.CurrentRoom().Name)
	}
	if got := g.Handle(parser.Parse("DROP BAG")); strings.Contains(got, "swings open") {
		t.Errorf("Handle(DROP BAG) at Room of Stings (needs a Key, not a Bag) = %q, should not open the door", got)
	}
	g.Handle(parser.Parse("PICKUP BAG")) // take it back before moving on

	g.Handle(parser.Parse("NORTH")) // Morfang, needs the Bag
	if g.World.CurrentRoom().Name != "Morfang" {
		t.Fatalf("test setup bug: expected to be in Morfang, got %q", g.World.CurrentRoom().Name)
	}
	got := g.Handle(parser.Parse("DROP BAG"))
	if !strings.Contains(got, "swings open") {
		t.Errorf("Handle(DROP BAG) at Morfang (needs a Bag) = %q, want it to open the door", got)
	}
	if g.hasItem("Bag") {
		t.Error("Bag should be spent (removed from inventory) after paying the real toll")
	}
}

func TestHandleAttackAndKillAreSynonymsForBlast(t *testing.T) {
	for _, verb := range []string{"ATTACK", "ATTACKS", "KILL"} {
		g := New() // Room of Misery has no known Monster
		got := g.Handle(parser.Parse(verb))
		if !strings.Contains(got, "nothing here") {
			t.Errorf("Handle(%s) with no monster = %q, want the same response as BLAST", verb, got)
		}
	}
}

func TestHandleGradeReportsCurrentGrade(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("GRADE"))
	if !strings.Contains(got, "Neophyte") {
		t.Errorf("Handle(GRADE) for a new player = %q, want it to mention Neophyte", got)
	}
}

// TestHandleSpellsListsRealSpells pins round 118: INVOKE joined this
// list because the manual's own explicit "Spells:" heading groups
// exactly I (Invoke), B (Blast), F (Freeze) - INVOKE was missing from
// this port's SPELLS listing entirely before, a real sourced omission,
// not just an incomplete aggregation.
func TestHandleSpellsListsRealSpells(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("SPELLS"))
	for _, want := range []string{"INVOKE", "BLAST", "FREEZE", "TRANSFUSION", "CALL"} {
		if !strings.Contains(got, want) {
			t.Errorf("Handle(SPELLS) = %q, want it to list %q", got, want)
		}
	}
}

// TestHandleCallWithoutScrollFails covers round 125's real, confirmed
// effect for CALL (see Handle's doc comment for the sourcing - The
// CRPG Addict's direct playthrough account): without the Scroll, it
// honestly fails rather than summoning Apex for free.
func TestHandleCallWithoutScrollFails(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("CALL"))
	if strings.Contains(got, "don't understand") {
		t.Errorf("Handle(CALL) = %q, want it recognized as a real spell, not an unknown word", got)
	}
	if strings.Contains(got, "Apex") {
		t.Errorf("Handle(CALL) without the Scroll = %q, want it to fail, not summon Apex", got)
	}
}

// TestHandleCallWithScrollSummonsApex covers CALL's real confirmed
// effect end-to-end: with the Scroll carried, it summons Apex.
func TestHandleCallWithScrollSummonsApex(t *testing.T) {
	g := New()
	g.Player.Items = append(g.Player.Items, "Scroll")
	got := g.Handle(parser.Parse("CALL"))
	if !strings.Contains(got, "Apex") {
		t.Errorf("Handle(CALL) with the Scroll carried = %q, want it to summon Apex", got)
	}
}

func TestHandleBlastWithNoMonster(t *testing.T) {
	g := New() // Room of Misery has no known Monster
	got := g.Handle(parser.Parse("BLAST"))
	if !strings.Contains(got, "nothing here") {
		t.Errorf("Handle(BLAST) with no monster = %q, want it to say there's nothing to hit", got)
	}
}

func TestHandleBlastDefeatsMonster(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST")) // Secunda Porta
	g.Handle(parser.Parse("DOOR, SILENCE"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd, which has a Monster

	room := g.World.CurrentRoom()
	if room.Name != "Trollwynd" || room.Monster == "" {
		t.Fatalf("test setup bug: expected to be in Trollwynd with a monster, got %+v", room)
	}
	startHealth := room.MonsterHealth

	// BLAST's per-hit damage now scales with Skill (see blastDamage), so
	// the monster may take fewer hits than its starting MonsterHealth to
	// defeat - loop until it's actually dead rather than assuming 1
	// damage per hit, with a generous safety cap against an infinite loop.
	var got string
	hits := 0
	for room.MonsterHealth > 0 && hits < startHealth+5 {
		got = g.Handle(parser.Parse("BLAST"))
		hits++
		if room.MonsterHealth > 0 && !strings.Contains(got, "still standing") {
			t.Errorf("Handle(BLAST) mid-fight = %q, want it to say the monster is still standing", got)
		}
	}
	if room.MonsterHealth > 0 {
		t.Fatalf("MonsterHealth after %d BLASTs = %d, want <= 0 (defeated)", hits, room.MonsterHealth)
	}
	if !strings.Contains(got, "destroyed") {
		t.Errorf("Handle(BLAST) on the killing blow = %q, want it to say destroyed", got)
	}

	// One more BLAST after defeat should report nothing left to hit.
	got = g.Handle(parser.Parse("BLAST"))
	if !strings.Contains(got, "nothing here") {
		t.Errorf("Handle(BLAST) after monster defeated = %q, want nothing-to-hit response", got)
	}
}

func TestHandleFreezeAwardsExperiencePoints(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd
	before := g.Player.ExperiencePoints

	got := g.Handle(parser.Parse("FREEZE"))
	want := 10 + g.Player.Luck
	if g.Player.ExperiencePoints != before+want {
		t.Errorf("ExperiencePoints after defeating a monster = %d, want %d (before %d + base 10 + Luck %d)", g.Player.ExperiencePoints, before+want, before, g.Player.Luck)
	}
	if !strings.Contains(got, "Experience Points") {
		t.Errorf("Handle(FREEZE) on a kill = %q, want it to mention Experience Points", got)
	}
}

func TestHandleBlastDefeatsMonsterAwardsExperiencePoints(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd
	room := g.World.CurrentRoom()
	before := g.Player.ExperiencePoints

	var got string
	hits := 0
	for room.MonsterHealth > 0 && hits < 10 {
		got = g.Handle(parser.Parse("BLAST"))
		hits++
	}
	if g.Player.ExperiencePoints <= before {
		t.Errorf("ExperiencePoints after defeating a monster with BLAST = %d, want more than before (%d)", g.Player.ExperiencePoints, before)
	}
	if !strings.Contains(got, "Experience Points") {
		t.Errorf("Handle(BLAST) on the killing blow = %q, want it to mention Experience Points", got)
	}
}

func TestHandleFreezeDefeatsMonsterInstantly(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd

	got := g.Handle(parser.Parse("FREEZE"))
	if !strings.Contains(got, "FREEZE") {
		t.Errorf("Handle(FREEZE) = %q, want it to mention freezing", got)
	}
	room := g.World.CurrentRoom()
	if room.MonsterHealth != 0 {
		t.Errorf("MonsterHealth after FREEZE = %d, want 0 (instant neutralize)", room.MonsterHealth)
	}
}

func TestBlastDamageScalesWithSkill(t *testing.T) {
	g := New()
	g.Player.Skill = 4 // bottom of the roll range
	if got, want := g.blastDamage(), 2; got != want {
		t.Errorf("blastDamage() with Skill 4 = %d, want %d", got, want)
	}
	g.Player.Skill = 12 // top of the roll range
	if got, want := g.blastDamage(), 4; got != want {
		t.Errorf("blastDamage() with Skill 12 = %d, want %d", got, want)
	}
}

func TestHandleBlastCostsStamina(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd, has a Monster
	before := g.Player.Stamina
	g.Handle(parser.Parse("BLAST"))
	if g.Player.Stamina != before-combatStaminaCost {
		t.Errorf("Stamina after BLAST = %d, want %d (before %d minus combatStaminaCost %d)", g.Player.Stamina, before-combatStaminaCost, before, combatStaminaCost)
	}
}

func TestHandleFreezeCostsStamina(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd, has a Monster
	before := g.Player.Stamina
	g.Handle(parser.Parse("FREEZE"))
	if g.Player.Stamina != before-combatStaminaCost {
		t.Errorf("Stamina after FREEZE = %d, want %d (before %d minus combatStaminaCost %d)", g.Player.Stamina, before-combatStaminaCost, before, combatStaminaCost)
	}
}

func TestHandleCombatCanKillPlayer(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd, has a Monster
	g.Player.Stamina = combatStaminaCost

	got := g.Handle(parser.Parse("BLAST"))
	if !g.Player.IsDead() {
		t.Fatalf("Player.Stamina = %d after BLAST, want <= 0 (dead)", g.Player.Stamina)
	}
	if !strings.Contains(got, "dead") {
		t.Errorf("Handle(BLAST) that reduces Stamina to 0 = %q, want it to mention death", got)
	}

	got = g.Handle(parser.Parse("EAST"))
	if got != "You are dead. (Stamina reached 0 - GAME OVER)" {
		t.Errorf("Handle(...) after death = %q, want the dead-rejection message", got)
	}
}

func TestHandleDeadPlayerCanStillLookAndMap(t *testing.T) {
	g := New()
	g.Player.Stamina = 0
	if got := g.Handle(parser.Parse("LOOK")); strings.Contains(got, "You are dead") {
		t.Errorf("Handle(LOOK) while dead = %q, want LOOK to still work", got)
	}
	if got := g.Handle(parser.Parse("MAP")); strings.Contains(got, "You are dead") {
		t.Errorf("Handle(MAP) while dead = %q, want MAP to still work", got)
	}
}

func TestHandleInvokeWithNoTargetListsDemons(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("INVOKE"))
	for _, want := range []string{"ASMODEE", "ASTAROT", "BELEZBAR", "MAGOT"} {
		if !strings.Contains(got, want) {
			t.Errorf("Handle(INVOKE) = %q, want it to list %s", got, want)
		}
	}
}

// TestHandleInvokeWithNoTargetShowsCorrespondences pins round 114:
// magic.Demon.Correspondences (extracted round 110, previously unused
// anywhere in-game) is now shown in the bare INVOKE listing.
func TestHandleInvokeWithNoTargetShowsCorrespondences(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("INVOKE"))
	if !strings.Contains(got, "Tourmaline") { // part of Astarot's real Correspondences
		t.Errorf("Handle(INVOKE) = %q, want it to include demon Correspondences (e.g. Astarot's gem, Tourmaline)", got)
	}
}

// TestHandleInvokeWithNoTargetShowsNumberSignAspect pins round 115:
// magic.Demon's Number/Sign/Aspect fields - real manual facts present
// since this project's first magic.Demons commit but never referenced
// by internal/game at all until now - are shown in the bare INVOKE
// listing too.
func TestHandleInvokeWithNoTargetShowsNumberSignAspect(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("INVOKE"))
	for _, want := range []string{"1376", "Sign of Gemini", "Legion"} { // Astarot's real Number/Sign/Aspect
		if !strings.Contains(got, want) {
			t.Errorf("Handle(INVOKE) = %q, want it to include %q (Astarot's Number/Sign/Aspect)", got, want)
		}
	}
}

func TestHandleInvokeRecognizesDemon(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("I ASTAROT")) // "I" expands to INVOKE
	if !strings.Contains(got, "ASTAROT") || !strings.Contains(got, "Talisman") {
		t.Errorf("Handle(I ASTAROT) = %q, want it to name Astarot and mention the missing Talisman", got)
	}
}

func TestHandleInvokeUnknownDemon(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("INVOKE MOTHRA"))
	if got != "There is no demon by that name." {
		t.Errorf("Handle(INVOKE MOTHRA) = %q, want the no-such-demon response", got)
	}
}

func TestHandlePickupMovesItemToInventory(t *testing.T) {
	g := New() // Room of Misery has a real sourced item: Grimoire
	got := g.Handle(parser.Parse("PICKUP GRIMOIRE"))
	if !strings.Contains(got, "Grimoire") {
		t.Errorf("Handle(PICKUP GRIMOIRE) = %q, want it to mention the Grimoire", got)
	}
	// Round 121: Axil's real starting Pouch (character.NewPlayer) means
	// Items is never empty before this pickup - check Grimoire was
	// added, not that it's the only item.
	found := false
	for _, item := range g.Player.Items {
		if item == "Grimoire" {
			found = true
		}
	}
	if !found {
		t.Errorf("Player.Items after pickup = %v, want it to include Grimoire", g.Player.Items)
	}
	room := g.World.CurrentRoom()
	for _, item := range room.Items {
		if item == "Grimoire" {
			t.Errorf("Room.Items after pickup still contains Grimoire: %v", room.Items)
		}
	}
}

// TestHandlePickupPoisonedItemCostsStamina covers round 128's real,
// sourced mechanic (a 1986 CRASH magazine review: "Poison damages
// Stamina upon contact") applied to Room of Misery's already-real
// "Poison-smeared book".
func TestHandlePickupPoisonedItemCostsStamina(t *testing.T) {
	g := New() // Room of Misery has a real sourced item: Poison-smeared book
	before := g.Player.Stamina
	got := g.Handle(parser.Parse("PICKUP POISON-SMEARED BOOK"))
	if !strings.Contains(got, "poisonous") {
		t.Errorf("Handle(PICKUP POISON-SMEARED BOOK) = %q, want it to mention the poison", got)
	}
	if want := before - poisonPickupStaminaCost; g.Player.Stamina != want {
		t.Errorf("Stamina after picking up a poisoned item = %d, want %d", g.Player.Stamina, want)
	}
}

// TestHandlePickupNonPoisonedItemDoesNotCostStamina is a regression
// guard: only items whose name actually contains "poison" should cost
// Stamina on pickup - a plain item like the Grimoire must not.
func TestHandlePickupNonPoisonedItemDoesNotCostStamina(t *testing.T) {
	g := New()
	before := g.Player.Stamina
	g.Handle(parser.Parse("PICKUP GRIMOIRE"))
	if g.Player.Stamina != before {
		t.Errorf("Stamina after picking up the Grimoire = %d, want unchanged %d", g.Player.Stamina, before)
	}
}

func TestHandlePickupMissingItem(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("PICKUP UNICORN"))
	if !strings.Contains(got, "no unicorn here") {
		t.Errorf("Handle(PICKUP UNICORN) = %q, want a not-here response", got)
	}
}

func TestHandleDropReturnsItemToRoom(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("PICKUP GRIMOIRE"))
	got := g.Handle(parser.Parse("DROP GRIMOIRE"))
	if !strings.Contains(got, "drop the Grimoire") {
		t.Errorf("Handle(DROP GRIMOIRE) = %q, want it to confirm the drop", got)
	}
	// Round 121: Axil's real starting Pouch (see character.NewPlayer)
	// is still carried after dropping the Grimoire - only the Grimoire
	// itself should be gone, not the whole inventory.
	for _, item := range g.Player.Items {
		if item == "Grimoire" {
			t.Errorf("Player.Items after drop = %v, want it to NOT include Grimoire", g.Player.Items)
		}
	}
	room := g.World.CurrentRoom()
	found := false
	for _, item := range room.Items {
		if item == "Grimoire" {
			found = true
		}
	}
	if !found {
		t.Errorf("Room.Items after drop = %v, want it to contain Grimoire again", room.Items)
	}
}

func TestHandleDropNotCarried(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("DROP UNICORN"))
	if !strings.Contains(got, "aren't carrying") {
		t.Errorf("Handle(DROP UNICORN) = %q, want a not-carrying response", got)
	}
}

func TestHandleSwapRecognized(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("Z")) // "Z" expands to SWAP
	if strings.Contains(got, "don't understand") {
		t.Errorf("Handle(Z) = %q, want a recognized (if stubbed) response, not the unknown-word rejection", got)
	}
}

// TestHandleLeftRightRecognized covers LEFT/RIGHT (Merphish keywords
// L/R) - real, frequently-used CASA walkthrough commands, recognized
// as real (not lumped in with "I don't understand that word") but
// honestly stubbed since this port has no facing-direction state.
func TestHandleLeftRightRecognized(t *testing.T) {
	for _, keyword := range []string{"L", "R"} {
		g := New()
		got := g.Handle(parser.Parse(keyword))
		if strings.Contains(got, "don't understand") {
			t.Errorf("Handle(%s) = %q, want it recognized as a real command, not an unknown word", keyword, got)
		}
	}
}

func TestHandleOptionsShowsMenu(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("OPTIONS"))
	if !strings.Contains(got, "Realign Status") {
		t.Errorf("Handle(OPTIONS) = %q, want it to list the confirmed menu including Realign Status", got)
	}
}

func TestHandleOptionsRealignRerollsStats(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("O REALIGN STATUS"))
	if !strings.Contains(got, "Realign Status") {
		t.Errorf("Handle(O REALIGN STATUS) = %q, want it to confirm the realign", got)
	}
	if g.Player.Stamina <= 0 || g.Player.Skill <= 0 || g.Player.Luck <= 0 {
		t.Errorf("Player after Realign = %+v, want all stats positive", g.Player)
	}
}

func TestHandleExamineReportsMonster(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd, has a Monster (Troll - see collodons_pile.go round 101)

	got := g.Handle(parser.Parse("EXAMINE"))
	if !strings.Contains(got, "Troll") {
		t.Errorf("Handle(EXAMINE) in a room with a monster = %q, want it mentioned", got)
	}
}

func TestHandleAbbreviatedMovement(t *testing.T) {
	g := New()
	// "E" should expand to EAST via parser.ExpandKeyword, same as typing
	// the full word - confirms the real Merphish abbreviation grammar
	// (see parser/keywords.go) works end to end through Handle.
	got := g.Handle(parser.Parse("E"))
	if !strings.Contains(got, "Secunda Porta") {
		t.Errorf("Handle(E) = %q, want it to move East to Secunda Porta same as Handle(EAST)", got)
	}
}

// TestHandleHalt covers HALT (Merphish "H") - round 85 corrected the
// response to reflect the manual's precise definition ("abandon the
// command being actioned and the rest of any outstanding command
// string"), not just a generic "Halted."
func TestHandleHalt(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("HALT"))
	if !strings.Contains(got, "halted") {
		t.Errorf("Handle(HALT) = %q, want it to acknowledge the command was halted", got)
	}
}

func TestHandleTransfusionRestoresStamina(t *testing.T) {
	g := New()
	g.Player.Stamina -= 5 // take damage first; a fresh player starts at MaxStamina already
	before := g.Player.Stamina
	got := g.Handle(parser.Parse("TRANSFUSION"))
	if g.Player.Stamina <= before {
		t.Errorf("Stamina after TRANSFUSION = %d, want more than before (%d)", g.Player.Stamina, before)
	}
	if !strings.Contains(got, "Stamina") {
		t.Errorf("Handle(TRANSFUSION) = %q, want it to report Stamina", got)
	}
}

func TestHandleTransfusionCapsAtMaxStamina(t *testing.T) {
	g := New() // a fresh player already starts at MaxStamina
	got := g.Handle(parser.Parse("TRANSFUSION"))
	if g.Player.Stamina != g.Player.MaxStamina {
		t.Errorf("Stamina after TRANSFUSION on a full-health player = %d, want it capped at MaxStamina %d", g.Player.Stamina, g.Player.MaxStamina)
	}
	if !strings.Contains(got, "maximum") {
		t.Errorf("Handle(TRANSFUSION) at full health = %q, want it to say maximum", got)
	}
}

func TestLevel1ExplorationStartRoomIsVisited(t *testing.T) {
	// Regression test: NewLevel1Exploration originally forgot to mark its
	// start room Visited (unlike New()/CollodonsPile), which crashed
	// Handle(MAP) with a nil pointer dereference the moment it was
	// actually run (caught by playing it live, not by unit tests alone -
	// see CLAUDE.md).
	g := NewLevel1Exploration()
	if !g.World.CurrentRoom().Visited {
		t.Fatal("NewLevel1Exploration's start room should be marked Visited immediately")
	}
	got := g.Handle(parser.Parse("MAP"))
	if got == "" {
		t.Error("Handle(MAP) on a fresh Level1Exploration game returned an empty string")
	}
}

func TestLevel2ExplorationStartRoomIsVisited(t *testing.T) {
	g := NewLevel2Exploration()
	if !g.World.CurrentRoom().Visited {
		t.Fatal("NewLevel2Exploration's start room should be marked Visited immediately")
	}
	got := g.Handle(parser.Parse("MAP"))
	if got == "" {
		t.Error("Handle(MAP) on a fresh Level2Exploration game returned an empty string")
	}
}

func TestLevel2ExplorationMovement(t *testing.T) {
	g := NewLevel2Exploration()
	// A1 --East--> A2 is a real, validated exit.
	got := g.Handle(parser.Parse("EAST"))
	if strings.Contains(got, "can't go that way") {
		t.Errorf("Handle(EAST) from A1 = %q, want a successful move to A2", got)
	}
	if g.World.CurrentRoom().Name != "A2" {
		t.Errorf("current room after EAST = %q, want A2", g.World.CurrentRoom().Name)
	}
}

// TestLevel2ExplorationStartRoomIsExit pins round 130's find:
// Level2Grid's own arbitrary starting anchor (A1) turned out to be a
// real, confirmed "Exit" cell (see level2_grid.go's doc comment) - so
// a real LOOK right at the start of a fresh -level2grid session
// genuinely announces a win, an honest quirk of A1 having been picked
// before its real name was known, not a bug.
func TestLevel2ExplorationStartRoomIsExit(t *testing.T) {
	g := NewLevel2Exploration()
	if g.Won {
		t.Fatal("Game.Won should be false immediately after construction, before any LOOK/move")
	}
	got := g.Handle(parser.Parse("LOOK"))
	if !g.Won {
		t.Fatalf("Game.Won = false after LOOK at the real Exit start room; output: %q", got)
	}
	if !strings.Contains(got, "YOU HAVE WON") {
		t.Errorf("Handle(LOOK) at Level2Exploration's start = %q, want the real win announcement", got)
	}
}

func TestLevel1ExplorationReachingExitWins(t *testing.T) {
	// Real, validated path from the start room (A1) to the Exit cell
	// (G3), found via BFS over Level1Grid's confirmed connectivity:
	// A1-B1-C1-C2-C3-C4-D4-E4-F4-G4-G3.
	g := NewLevel1Exploration()
	path := []string{"SOUTH", "SOUTH", "EAST", "EAST", "EAST", "SOUTH", "SOUTH", "SOUTH", "SOUTH", "WEST"}
	var last string
	for _, dir := range path {
		last = g.Handle(parser.Parse(dir))
	}
	if !g.Won {
		t.Fatalf("Game.Won = false after reaching the Exit room; last Handle() output: %q", last)
	}
	if !strings.Contains(last, "YOU HAVE WON") {
		t.Errorf("Handle() output on reaching Exit = %q, want a win announcement", last)
	}
}

// TestNougatDefeatsWerewolfOnDrop covers the real, sourced alternate
// mechanic (CASA walkthrough: Werewolves "killable by walking through
// after dropping NOUGAT" — see checkNougatWerewolf). Path to C2 (a real
// Werewolf, per Level1Grid): A1-South-B1-South-C1-East-C2.
func TestNougatDefeatsWerewolfOnDrop(t *testing.T) {
	g := NewLevel1Exploration()
	for _, dir := range []string{"SOUTH", "SOUTH", "EAST"} {
		g.Handle(parser.Parse(dir))
	}
	room := g.World.CurrentRoom()
	if room.Monster != "Werewolf" || room.MonsterHealth <= 0 {
		t.Fatalf("test setup bug: expected a live Werewolf at C2, got %+v", room)
	}
	g.Player.Items = append(g.Player.Items, "Nougat")
	got := g.Handle(parser.Parse("DROP NOUGAT"))
	if room.MonsterHealth > 0 {
		t.Errorf("Werewolf should be defeated after dropping Nougat, MonsterHealth = %d", room.MonsterHealth)
	}
	if !strings.Contains(got, "Nougat") {
		t.Errorf("Handle(DROP NOUGAT) with a live Werewolf present = %q, want it to mention the Nougat mechanic", got)
	}
}

// TestCollodonsPileNougatDefeatsWolfdorpWerewolf covers round 120's
// Wolfdorp addition end-to-end in the actual DEFAULT (CollodonsPile)
// game, not just the isolated Level1Grid mechanic test above: picks up
// the real, already-placed Nougat in Trollwynd, carries it to Wolfdorp
// (both real rooms on the real walkthrough path), and drops it there -
// confirming game.checkNougatWerewolf's real mechanic is now genuinely
// reachable in a normal playthrough, not just unit-testable in
// isolation.
func TestCollodonsPileNougatDefeatsWolfdorpWerewolf(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))          // Secunda Porta
	g.Handle(parser.Parse("DOOR, SILENCE")) // unlocks the door North
	g.Handle(parser.Parse("NORTH"))         // Trollwynd - has real Nougat
	if got := g.Handle(parser.Parse("PICKUP NOUGAT")); !strings.Contains(got, "Nougat") {
		t.Fatalf("test setup bug: PICKUP NOUGAT in Trollwynd = %q, want it to succeed", got)
	}
	g.Handle(parser.Parse("SOUTH")) // Sothic Complex
	g.Handle(parser.Parse("SOUTH")) // Wolfdorp - has the real Werewolf

	room := g.World.CurrentRoom()
	if room.Name != "Wolfdorp" || room.Monster != "Werewolf" || room.MonsterHealth <= 0 {
		t.Fatalf("test setup bug: expected a live Werewolf in Wolfdorp, got %+v", room)
	}
	g.Handle(parser.Parse("DROP NOUGAT"))
	if room.MonsterHealth > 0 {
		t.Errorf("Wolfdorp's Werewolf should be defeated after dropping the real Nougat carried from Trollwynd, MonsterHealth = %d", room.MonsterHealth)
	}
}

// TestCollodonsPileGarlicDefeatsMorfangVampire covers round 126's
// second CRPG Addict find, end-to-end in the actual DEFAULT
// (CollodonsPile) game: picks up the real, already-placed Garlic in
// Wolfdorp, carries it to Morfang (both real rooms on the real
// walkthrough path), and drops it there - confirming
// game.checkGarlicVampire's real mechanic is genuinely reachable in a
// normal playthrough, not just unit-testable in isolation.
func TestCollodonsPileGarlicDefeatsMorfangVampire(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))          // Secunda Porta
	g.Handle(parser.Parse("DOOR, SILENCE")) // unlocks the door North
	g.Handle(parser.Parse("NORTH"))         // Trollwynd
	g.Handle(parser.Parse("NORTH"))         // Agile Stair
	g.Handle(parser.Parse("SOUTH-EAST"))    // Methos
	g.Handle(parser.Parse("SOUTH"))         // Sothic Complex
	g.Handle(parser.Parse("SOUTH"))         // Wolfdorp - has the real Garlic
	if got := g.Handle(parser.Parse("EXAMINE CHEST")); !strings.Contains(got, "chest") {
		t.Fatalf("test setup bug: EXAMINE CHEST in Wolfdorp = %q, want it to acknowledge the chest", got)
	}
	if got := g.Handle(parser.Parse("PICKUP GARLIC")); !strings.Contains(got, "Garlic") {
		t.Fatalf("test setup bug: PICKUP GARLIC in Wolfdorp = %q, want it to succeed", got)
	}
	g.Handle(parser.Parse("NORTH-WEST")) // Room of Stings
	g.Handle(parser.Parse("NORTH"))      // Morfang - has the real Vampire

	room := g.World.CurrentRoom()
	if room.Name != "Morfang" || room.Monster != "Vampire" || room.MonsterHealth <= 0 {
		t.Fatalf("test setup bug: expected a live Vampire in Morfang, got %+v", room)
	}
	got := g.Handle(parser.Parse("DROP GARLIC"))
	if room.MonsterHealth > 0 {
		t.Errorf("Morfang's Vampire should be defeated after dropping the real Garlic carried from Wolfdorp, MonsterHealth = %d", room.MonsterHealth)
	}
	if !strings.Contains(got, "Garlic") {
		t.Errorf("Handle(DROP GARLIC) with a live Vampire present = %q, want it to mention the Garlic mechanic", got)
	}
}

func TestLevel1ExplorationMovementAndCombat(t *testing.T) {
	g := NewLevel1Exploration()
	// A1 (start) has a real Ghost, per world.Level1Grid's extracted data.
	room := g.World.CurrentRoom()
	if room.Monster != "Ghost" {
		t.Fatalf("test setup bug: expected the start room to have a Ghost, got %+v", room)
	}
	startHealth := room.MonsterHealth
	hits := 0
	for room.MonsterHealth > 0 && hits < startHealth+5 {
		g.Handle(parser.Parse("BLAST"))
		hits++
	}
	if room.MonsterHealth > 0 {
		t.Errorf("MonsterHealth after %d BLASTs = %d, want <= 0 (defeated)", hits, room.MonsterHealth)
	}

	// A1 --East--> A2 is a real, validated exit.
	got := g.Handle(parser.Parse("EAST"))
	if strings.Contains(got, "can't go that way") {
		t.Errorf("Handle(EAST) from A1 = %q, want a successful move to A2", got)
	}
	if g.World.CurrentRoom().Name != "A2" {
		t.Errorf("current room after EAST = %q, want A2", g.World.CurrentRoom().Name)
	}
}

func TestLevel3ExplorationStartRoomIsVisited(t *testing.T) {
	g := NewLevel3Exploration()
	if !g.World.CurrentRoom().Visited {
		t.Fatal("NewLevel3Exploration's start room should be marked Visited immediately")
	}
	got := g.Handle(parser.Parse("MAP"))
	if got == "" {
		t.Error("Handle(MAP) on a fresh Level3Exploration game returned an empty string")
	}
}

func TestLevel3ExplorationMovement(t *testing.T) {
	g := NewLevel3Exploration()
	got := g.Handle(parser.Parse("EAST"))
	if strings.Contains(got, "can't go that way") {
		t.Errorf("Handle(EAST) from A1 = %q, want a successful move to A2", got)
	}
	if g.World.CurrentRoom().Name != "A2" {
		t.Errorf("current room after EAST = %q, want A2", g.World.CurrentRoom().Name)
	}
}

func TestLevel4ExplorationStartRoomIsVisited(t *testing.T) {
	g := NewLevel4Exploration()
	if !g.World.CurrentRoom().Visited {
		t.Fatal("NewLevel4Exploration's start room should be marked Visited immediately")
	}
	got := g.Handle(parser.Parse("MAP"))
	if got == "" {
		t.Error("Handle(MAP) on a fresh Level4Exploration game returned an empty string")
	}
}

func TestLevel4ExplorationMovement(t *testing.T) {
	g := NewLevel4Exploration()
	got := g.Handle(parser.Parse("EAST"))
	if strings.Contains(got, "can't go that way") {
		t.Errorf("Handle(EAST) from the start room = %q, want a successful move", got)
	}
}

func TestNewGameStartsInRoomOfMisery(t *testing.T) {
	g := New()
	room := g.World.CurrentRoom()
	if room == nil || room.Name != "Room of Misery" {
		t.Errorf("starting room = %+v, want Room of Misery", room)
	}
	if !room.Visited {
		t.Error("starting room should be marked Visited immediately")
	}
}

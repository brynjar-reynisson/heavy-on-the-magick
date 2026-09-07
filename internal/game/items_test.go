package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/character"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/world"
)

func TestHandleCarryIsSynonymForPickup(t *testing.T) {
	g := New() // Room of Misery has a real sourced item: Grimoire
	got := g.Handle(parser.Parse("CARRY GRIMOIRE"))
	if !strings.Contains(got, "Grimoire") || !g.hasItem("Grimoire") {
		t.Errorf("Handle(CARRY GRIMOIRE) = %q, want it to pick up the Grimoire", got)
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

func TestHandleExamineReportsMonster(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd, has a Monster (Troll - see collodons_pile.go round 101)

	got := g.Handle(parser.Parse("EXAMINE"))
	if !strings.Contains(got, "Troll") {
		t.Errorf("Handle(EXAMINE) in a room with a monster = %q, want it mentioned", got)
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

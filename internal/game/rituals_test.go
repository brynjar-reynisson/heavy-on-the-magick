package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/character"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/world"
)

// TestHandleNestPhoenixRequiresRealNest covers round 136's real,
// sourced ritual command (see game.nestPhoenix's doc comment): saying
// "NEST, PHOENIX" anywhere that isn't really named "Nest of Phoenix"
// is an honest rejection, not a fabricated success.
func TestHandleNestPhoenixRequiresRealNest(t *testing.T) {
	g := New() // starts in Room of Misery, not the Nest of Phoenix
	got := g.Handle(parser.Parse("NEST, PHOENIX"))
	if !strings.Contains(got, "no phoenix nest here") {
		t.Errorf("Handle(NEST, PHOENIX) outside the real nest = %q, want an honest rejection", got)
	}
}

// TestHandleNestPhoenixFullRitual covers the real ritual succeeding
// once every confirmed requirement is met (real room name, carrying
// the Clasp/"Salamander charm", and an Egg already dropped there).
// Uses a synthetic room (no shipped World.Room is named "Nest of
// Phoenix" yet - real, scoped follow-up work) - same pattern as
// TestHandleFireBlocksMovementWithoutClasp/TestHandleSwapItemRevealsRealItem.
func TestHandleNestPhoenixFullRitual(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Nest of Phoenix"})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("NEST, PHOENIX"))
	if !strings.Contains(got, "Salamander charm") {
		t.Errorf("Handle(NEST, PHOENIX) with no Clasp carried = %q, want the Salamander-charm hint", got)
	}

	g.Player.Items = append(g.Player.Items, "Clasp")
	got = g.Handle(parser.Parse("NEST, PHOENIX"))
	if !strings.Contains(got, "drop an Egg") {
		t.Errorf("Handle(NEST, PHOENIX) with the Clasp but no Egg dropped = %q, want the drop-an-Egg hint", got)
	}

	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Egg")
	got = g.Handle(parser.Parse("NEST, PHOENIX"))
	if !strings.Contains(got, "ritual succeeds") {
		t.Errorf("Handle(NEST, PHOENIX) with every requirement met = %q, want the ritual to succeed", got)
	}
}

// TestHandleCauldronAchadRequiresRealCauldron covers round 139's real,
// sourced ritual command (see game.cauldronAchad's doc comment):
// saying "CAULDRON, ACHAD" anywhere that isn't really named "Kitchen
// of Ai" (corrected round 178 - round 177 briefly guessed "Room of
// Nani" instead) is an honest rejection.
func TestHandleCauldronAchadRequiresRealCauldron(t *testing.T) {
	g := New() // starts in Room of Misery, not Kitchen of Ai
	got := g.Handle(parser.Parse("CAULDRON, ACHAD"))
	if !strings.Contains(got, "no cauldron here") {
		t.Errorf("Handle(CAULDRON, ACHAD) outside the real cauldron = %q, want an honest rejection", got)
	}
}

// TestHandleCauldronAchadFullRitual covers the real ritual succeeding
// once every confirmed requirement is met: real room name, the
// Scroll removed, and Ulna/Thigh/Skull all dropped. Uses a synthetic
// room named "Kitchen of Ai" (the real name, per round 178 -
// Level3Grid's own H2 isn't in the same World as this synthetic test
// setup) - same pattern as TestHandleNestPhoenixFullRitual.
func TestHandleCauldronAchadFullRitual(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Kitchen of Ai", Items: []string{"Scroll"}})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("CAULDRON, ACHAD"))
	if !strings.Contains(got, "take the Scroll out") {
		t.Errorf("Handle(CAULDRON, ACHAD) with the Scroll still inside = %q, want the take-the-Scroll-out hint", got)
	}

	g.Handle(parser.Parse("PICKUP SCROLL"))
	got = g.Handle(parser.Parse("CAULDRON, ACHAD"))
	if !strings.Contains(got, "needs more bones") {
		t.Errorf("Handle(CAULDRON, ACHAD) with the Scroll out but no bones = %q, want the needs-more-bones hint", got)
	}

	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Ulna", "Thigh")
	got = g.Handle(parser.Parse("CAULDRON, ACHAD"))
	if !strings.Contains(got, "needs more bones") {
		t.Errorf("Handle(CAULDRON, ACHAD) with only 2 of 3 bones = %q, want it to still be missing bones", got)
	}

	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Skull")
	got = g.Handle(parser.Parse("CAULDRON, ACHAD"))
	if !strings.Contains(got, "ritual succeeds") {
		t.Errorf("Handle(CAULDRON, ACHAD) with every requirement met = %q, want the ritual to succeed", got)
	}
}

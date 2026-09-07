package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
)

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

// TestHandleAsmodeeDestroyRequiresErlstone covers round 162's real,
// sourced ability (Hardcore Gaming 101: "Asmodee destroys any object
// you ask of him") - the same Charm-on-the-ground gating already
// established for Astarot/Magot.
func TestHandleAsmodeeDestroyRequiresErlstone(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("ASMODEE, GRIMOIRE"))
	if !strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(ASMODEE, GRIMOIRE) with no Erlstone carried = %q, want a Talisman rejection", got)
	}
}

// TestHandleAsmodeeDestroysCarriedItem covers destroying an item the
// player is currently carrying.
func TestHandleAsmodeeDestroysCarriedItem(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Erlstone")
	g.Handle(parser.Parse("PICKUP GRIMOIRE"))
	got := g.Handle(parser.Parse("ASMODEE, GRIMOIRE"))
	if strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(ASMODEE, GRIMOIRE) with Erlstone on the ground = %q, want the destroy to succeed", got)
	}
	if !strings.Contains(got, "Grimoire") || !strings.Contains(got, "crumbles") {
		t.Errorf("Handle(ASMODEE, GRIMOIRE) = %q, want it to confirm the Grimoire was destroyed", got)
	}
	if g.hasItem("Grimoire") {
		t.Error("Grimoire should be gone from the player's inventory after Asmodee destroys it")
	}
}

// TestHandleAsmodeeDestroysRoomItem covers destroying an item lying in
// a room (not carried), mirroring magotLocate's own search order.
func TestHandleAsmodeeDestroysRoomItem(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Erlstone")
	got := g.Handle(parser.Parse("ASMODEE, GRIMOIRE"))
	if !strings.Contains(got, "Grimoire") || !strings.Contains(got, "crumbles") {
		t.Errorf("Handle(ASMODEE, GRIMOIRE) with the Grimoire in the room = %q, want it destroyed", got)
	}
	room := g.World.CurrentRoom()
	for _, item := range room.Items {
		if item == "Grimoire" {
			t.Errorf("room Items after Asmodee destroys the Grimoire = %v, want it gone", room.Items)
		}
	}
}

func TestHandleAsmodeeDestroyUnknownObject(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Erlstone")
	got := g.Handle(parser.Parse("ASMODEE, EXCALIBUR"))
	if !strings.Contains(got, "no such object") {
		t.Errorf("Handle(ASMODEE, EXCALIBUR) = %q, want an honest not-found response", got)
	}
}

package game

import (
	"strings"
	"testing"
	"time"

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
// round 131's correction: merely CARRYING the Charm isn't enough.
// ROUND 183 correction: this used to get a distinct, harmless hint
// instead of the furnace-room punishment - the user, having
// independently finished the same third gameplay video AND confirmed
// it live in SpecEmu, established that carrying-not-grounded fails the
// exact same lethal way as having no Charm at all. The message wording
// still distinguishes the 2 cases (see punishFailedInvoke), but both
// now lead to the same real furnace-room death.
func TestHandleInvokeCarriedNotDroppedCharmFails(t *testing.T) {
	g := New()
	g.Player.Items = append(g.Player.Items, "Sunflower")
	got := g.Handle(parser.Parse("I MAGOT"))
	if !strings.Contains(got, "must be on the ground") {
		t.Errorf("Handle(I MAGOT) with Sunflower only carried = %q, want the carrying-specific hint", got)
	}
	if !strings.Contains(got, "furnace room") {
		t.Errorf("Handle(I MAGOT) with the Charm carried (not dropped) = %q, want the same real furnace-room punishment as having no Charm at all", got)
	}
}

// TestHandleInvokeWithoutCharmTeleportsToFurnaceRoom pins round 126's
// real, sourced punishment mechanic (The CRPG Addict's first-hand
// playthrough account, the same source round 125 used to resolve
// CALL's effect): invoking a demon without its Charm doesn't just
// reject the command, it actually teleports the player to the real
// Furnace Room. Round 183: the user independently finished the same
// third gameplay video and separately confirmed live in SpecEmu that
// this is actually lethal - "sends Axil to the furnace room, where he
// dies horribly" - so this now also checks real death, not just the
// teleport.
func TestHandleInvokeWithoutCharmTeleportsToFurnaceRoom(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("I MAGOT")) // no Sunflower carried
	if !strings.Contains(got, "furnace room") {
		t.Errorf("Handle(I MAGOT) with no Charm = %q, want it to mention being sent to the furnace room", got)
	}
	if room := g.World.CurrentRoom(); room == nil || room.Name != "Furnace Room" {
		t.Errorf("current room after a failed INVOKE = %+v, want Furnace Room", room)
	}
	if !g.Player.IsDead() {
		t.Error("Player should be dead after being flung into the Furnace Room")
	}
	if !strings.Contains(got, "die horribly") {
		t.Errorf("Handle(I MAGOT) with no Charm = %q, want the real death message", got)
	}
}

// TestHandleAstarotTeleportRequiresSword covers round 183's correction:
// a failed Charm check now kills the player via the same real furnace-
// room punishment as bare INVOKE (see punishFailedInvoke's doc
// comment), not a harmless rejection.
func TestHandleAstarotTeleportRequiresSword(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("ASTAROT, WOLFDORP"))
	if !strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(ASTAROT, WOLFDORP) with no Sword carried = %q, want a Talisman rejection", got)
	}
	if !g.Player.IsDead() {
		t.Error("Player should be dead after a failed ASTAROT invocation with no Charm")
	}
}

// TestHandleAstarotTeleportSucceeds pins the real, hint-screen-confirmed
// example command ("ASTAROT, WOLFDORP" — see parser.Parse's package doc
// comment) actually teleporting the player, once they carry Astarot's
// confirmed Charm (Sword). The success wording ("Best place for you")
// is the game's own exact confirmed text (round 180 - seen 3 separate
// times in the same footage), replacing this port's own earlier
// invented "In an instant, you are transported to X" wording.
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
	if !strings.Contains(got, "Best place for you") {
		t.Errorf("Handle(ASTAROT, WOLFDORP) = %q, want the real confirmed success phrase \"Best place for you\"", got)
	}
	if got := g.World.CurrentRoom(); got == nil || got.Name != "Wolfdorp" {
		t.Errorf("after ASTAROT, WOLFDORP, current room = %+v, want Wolfdorp", got)
	}
}

// TestHandleAstarotTeleportUnknownLocation covers round 180's real,
// video-confirmed rejection text - "No such place." - replacing this
// port's own earlier invented wording.
func TestHandleAstarotTeleportUnknownLocation(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Sword")
	got := g.Handle(parser.Parse("ASTAROT, NARNIA"))
	if !strings.Contains(got, "No such place") {
		t.Errorf("Handle(ASTAROT, NARNIA) = %q, want the real confirmed rejection \"No such place.\"", got)
	}
	if got := g.World.CurrentRoom(); got == nil || got.Name != "Room of Misery" {
		t.Errorf("after an unknown-location ASTAROT command, current room = %+v, want unchanged (Room of Misery)", got)
	}
}

// TestHandleMagotLocateRequiresSunflower covers round 183's correction
// - see TestHandleAstarotTeleportRequiresSword's doc comment.
func TestHandleMagotLocateRequiresSunflower(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("MAGOT, GRIMOIRE"))
	if !strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(MAGOT, GRIMOIRE) with no Sunflower carried = %q, want a Talisman rejection", got)
	}
	if !g.Player.IsDead() {
		t.Error("Player should be dead after a failed MAGOT invocation with no Charm")
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
// TestHandleAsmodeeDestroyRequiresErlstone covers round 183's
// correction - see TestHandleAstarotTeleportRequiresSword's doc
// comment. Also independently confirmed live in SpecEmu by the user.
func TestHandleAsmodeeDestroyRequiresErlstone(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("ASMODEE, GRIMOIRE"))
	if !strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(ASMODEE, GRIMOIRE) with no Erlstone carried = %q, want a Talisman rejection", got)
	}
	if !g.Player.IsDead() {
		t.Error("Player should be dead after a failed ASMODEE invocation with no Charm")
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

// TestHandleAsmodeeDestroysLockedDoor covers round 183's real, user-
// recalled ability (see asmodeeDestroy's own doc comment): "ASMODEE,
// DOOR" clears a real locked-door obstacle on the current room -
// DoorPasswords, TollItem, and Guards are all structural room facts,
// not named Items the generic search would ever find.
func TestHandleAsmodeeDestroysLockedDoor(t *testing.T) {
	g := New()
	room := g.World.CurrentRoom()
	room.Items = append(room.Items, "Erlstone")
	room.DoorPasswords = []string{"SILENCE"}

	got := g.Handle(parser.Parse("ASMODEE, DOOR"))
	if !strings.Contains(got, "door crumbles") {
		t.Errorf("Handle(ASMODEE, DOOR) with a real locked door = %q, want it destroyed", got)
	}
	if len(room.DoorPasswords) > 0 {
		t.Error("DoorPasswords should be cleared after ASMODEE, DOOR destroys the lock")
	}
}

// TestHandleAsmodeeDestroyDoorWithNoLock is the honest negative: no
// real lock on the current room means nothing to destroy.
func TestHandleAsmodeeDestroyDoorWithNoLock(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Erlstone")
	got := g.Handle(parser.Parse("ASMODEE, DOOR"))
	if !strings.Contains(got, "no such object") {
		t.Errorf("Handle(ASMODEE, DOOR) with no real lock present = %q, want an honest not-found response", got)
	}
}

// TestHandleBelezbarRevealRequiresMantis covers round 164's real
// ability (completing all 4 demons - Astarot/Magot/Asmodee already had
// one) - the same Charm-on-the-ground gating as the other 3.
// TestHandleBelezbarRevealRequiresMantis covers round 183's correction
// - see TestHandleAstarotTeleportRequiresSword's doc comment.
func TestHandleBelezbarRevealRequiresMantis(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("BELEZBAR, PEBBLE"))
	if !strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(BELEZBAR, PEBBLE) with no Mantis carried = %q, want a Talisman rejection", got)
	}
	if !g.Player.IsDead() {
		t.Error("Player should be dead after a failed BELEZBAR invocation with no Charm")
	}
}

// TestHandleBelezbarRevealsRealDisguise covers the one confirmed,
// sourced disguise (numbered map poster #59: "Pebble (disguised
// Erlstone)") - see belezbarDisguises's doc comment.
func TestHandleBelezbarRevealsRealDisguise(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Mantis")
	got := g.Handle(parser.Parse("BELEZBAR, PEBBLE"))
	if strings.Contains(got, "no suitable Talisman") {
		t.Errorf("Handle(BELEZBAR, PEBBLE) with Mantis on the ground = %q, want the reveal to succeed", got)
	}
	if !strings.Contains(got, "Pebble") || !strings.Contains(got, "Erlstone") {
		t.Errorf("Handle(BELEZBAR, PEBBLE) = %q, want it to reveal the real disguise (Erlstone)", got)
	}
}

// TestHandleBelezbarRevealOrdinaryObject covers the honest negative:
// an object with no confirmed disguise appears to be exactly what it
// seems, not a fabricated secret identity.
func TestHandleBelezbarRevealOrdinaryObject(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Mantis")
	got := g.Handle(parser.Parse("BELEZBAR, GRIMOIRE"))
	if !strings.Contains(got, "exactly what it seems") {
		t.Errorf("Handle(BELEZBAR, GRIMOIRE) with no confirmed disguise = %q, want an honest ordinary-object response", got)
	}
}

// TestHandleDemonPatienceRunsOutAfterThreeNonsenseAttempts covers round
// 184's user-specified rule, modeling a real failure the user saw live
// in SpecEmu ("he didn't say anything worthy soon enough"): a demon
// sends Axil to the furnace after demonNonsenseLimit consecutive
// unrecognized attempts, the same lethal outcome as a missing Talisman.
func TestHandleDemonPatienceRunsOutAfterThreeNonsenseAttempts(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Sword")

	var got string
	for range demonNonsenseLimit {
		got = g.Handle(parser.Parse("ASTAROT, NARNIA"))
	}
	if !strings.Contains(got, "furnace room") {
		t.Errorf("Handle(ASTAROT, NARNIA) x%d = %q, want the demon to lose patience and send Axil to the furnace", demonNonsenseLimit, got)
	}
	if !g.Player.IsDead() {
		t.Error("Player should be dead after exhausting a demon's patience with nonsense")
	}
}

// TestHandleDemonPatienceSurvivesFewerThanLimitNonsenseAttempts is the
// regression guard: fewer than demonNonsenseLimit nonsense attempts
// doesn't trigger anything yet.
func TestHandleDemonPatienceSurvivesFewerThanLimitNonsenseAttempts(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Sword")

	for i := range demonNonsenseLimit - 1 {
		got := g.Handle(parser.Parse("ASTAROT, NARNIA"))
		if strings.Contains(got, "furnace room") {
			t.Errorf("Handle(ASTAROT, NARNIA) attempt %d = %q, want no punishment yet", i+1, got)
		}
	}
	if g.Player.IsDead() {
		t.Error("Player should still be alive before the nonsense limit is reached")
	}
}

// TestHandleDemonPatienceClearsOnSuccess covers that a worthy (real)
// attempt resets the nonsense counter entirely, matching
// invokeWithPatience's own doc comment - a later nonsense attempt for
// the same demon starts counting from zero again, not where the
// earlier, cleared session left off. Uses MAGOT rather than ASTAROT
// specifically because a successful ASTAROT teleports the player away
// from the room holding the Sword, which would fail the NEXT attempt
// for the wrong reason (missing Charm, not nonsense) - MAGOT's locate
// ability doesn't move the player.
func TestHandleDemonPatienceClearsOnSuccess(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Sunflower")

	g.Handle(parser.Parse("MAGOT, EXCALIBUR")) // 1 nonsense attempt
	g.Handle(parser.Parse("MAGOT, GRIMOIRE"))  // real, worthy - clears the session
	got := g.Handle(parser.Parse("MAGOT, EXCALIBUR"))
	if strings.Contains(got, "furnace room") {
		t.Errorf("Handle(MAGOT, EXCALIBUR) after an intervening success = %q, want the session reset, not punished", got)
	}
}

// TestHandleDemonPatienceTimesOut covers the other half of round 184's
// rule: exceeding demonPatienceLimit since a session began ends it too,
// even with 0 nonsense attempts recorded - modeling "soon enough" as a
// real time limit, not just a strike count.
func TestHandleDemonPatienceTimesOut(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Items = append(g.World.CurrentRoom().Items, "Sword")

	g.Handle(parser.Parse("ASTAROT, NARNIA")) // starts the session
	g.demonSession.since = time.Now().Add(-demonPatienceLimit - time.Second)

	got := g.Handle(parser.Parse("ASTAROT, NARNIA"))
	if !strings.Contains(got, "furnace room") {
		t.Errorf("Handle(ASTAROT, NARNIA) after the patience timeout = %q, want the demon to lose patience", got)
	}
	if !g.Player.IsDead() {
		t.Error("Player should be dead after a demon's patience timing out")
	}
}

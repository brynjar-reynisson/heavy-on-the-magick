package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/character"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/world"
)

// TestHandleChasmKillsPlayerWithoutFlask covers round 183's real,
// user-recalled hazard (see world.Room.Chasm's doc comment): unlike
// Fire, entering a Chasm room without a Flask is lethal, not just
// blocked.
func TestHandleChasmKillsPlayerWithoutFlask(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Start", Exits: map[world.Direction]world.RoomID{world.North: 1}})
	w.AddRoom(&world.Room{ID: 1, Name: "The Chasm", Chasm: true})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("NORTH"))
	if !strings.Contains(got, "GAME OVER") {
		t.Errorf("Handle(NORTH) into a Chasm without a Flask = %q, want death", got)
	}
	if !g.Player.IsDead() {
		t.Error("Player should be dead after falling into the Chasm without a Flask")
	}
}

// TestHandleChasmSafeWithFlask is the regression guard for the same
// mechanic: carrying a Flask bridges the Chasm safely.
func TestHandleChasmSafeWithFlask(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Start", Exits: map[world.Direction]world.RoomID{world.North: 1}})
	w.AddRoom(&world.Room{ID: 1, Name: "The Chasm", Chasm: true})
	g := &Game{Player: character.NewPlayer(), World: w}
	g.Player.Items = append(g.Player.Items, "Flask")

	got := g.Handle(parser.Parse("NORTH"))
	if strings.Contains(got, "GAME OVER") {
		t.Errorf("Handle(NORTH) into a Chasm WITH a Flask = %q, want it to succeed", got)
	}
	if g.World.CurrentRoom().Name != "The Chasm" {
		t.Errorf("current room after crossing the Chasm with a Flask = %q, want \"The Chasm\"", g.World.CurrentRoom().Name)
	}
}

// TestHandleEnteringMedusaRoomNeverKills covers round 185's correction:
// the user clarified that entering her room is NEVER lethal by itself
// (regardless of Mirror) - only trying to move PAST her (see
// TestHandleMedusaKillsOnMoveAttemptWithoutMirror below) or BLASTing
// her is.
func TestHandleEnteringMedusaRoomNeverKills(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Start", Exits: map[world.Direction]world.RoomID{world.North: 1}})
	w.AddRoom(&world.Room{ID: 1, Name: "The Pit", Monster: "Medusa", MonsterHealth: 3})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("NORTH"))
	if strings.Contains(got, "GAME OVER") {
		t.Errorf("Handle(NORTH) into Medusa's room without a Mirror = %q, want entry to succeed (never lethal by itself)", got)
	}
	if g.World.CurrentRoom().Name != "The Pit" {
		t.Errorf("current room after entering Medusa's room = %q, want \"The Pit\"", g.World.CurrentRoom().Name)
	}
}

// TestHandleMedusaKillsOnMoveAttemptWithoutMirror covers the real
// death trigger the user described: "only when he tries to walk into
// her" - i.e. any move attempt from within her room, without a Mirror,
// meets her gaze. Distinct from checkMirrorMedusa (round 178), which
// handles actually DROPPING the Mirror to destroy her outright.
func TestHandleMedusaKillsOnMoveAttemptWithoutMirror(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Start", Exits: map[world.Direction]world.RoomID{world.North: 1}})
	w.AddRoom(&world.Room{ID: 1, Name: "The Pit", Monster: "Medusa", MonsterHealth: 3, Exits: map[world.Direction]world.RoomID{world.South: 0}})
	g := &Game{Player: character.NewPlayer(), World: w}

	g.Handle(parser.Parse("NORTH")) // enters The Pit safely, no Mirror

	got := g.Handle(parser.Parse("SOUTH"))
	if !strings.Contains(got, "GAME OVER") {
		t.Errorf("Handle(SOUTH) trying to walk past Medusa without a Mirror = %q, want death", got)
	}
	if !g.Player.IsDead() {
		t.Error("Player should be dead after trying to walk past Medusa without a Mirror")
	}
}

// TestHandleMedusaSafeToLeaveWithMirror is the regression guard:
// carrying a Mirror makes moving past her safe outright - the same
// "carried item bypasses a hazard" pattern as Fire/Clasp, not
// something requiring her to be separately defeated first.
func TestHandleMedusaSafeToLeaveWithMirror(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Start", Exits: map[world.Direction]world.RoomID{world.North: 1}})
	w.AddRoom(&world.Room{ID: 1, Name: "The Pit", Monster: "Medusa", MonsterHealth: 3, Exits: map[world.Direction]world.RoomID{world.South: 0}})
	g := &Game{Player: character.NewPlayer(), World: w}
	g.Player.Items = append(g.Player.Items, "Mirror")

	g.Handle(parser.Parse("NORTH")) // enters The Pit safely

	got := g.Handle(parser.Parse("SOUTH"))
	if strings.Contains(got, "GAME OVER") {
		t.Errorf("Handle(SOUTH) with a Mirror carried = %q, want it to succeed", got)
	}
	if g.World.CurrentRoom().Name != "Start" {
		t.Errorf("current room after leaving Medusa's room with a Mirror = %q, want \"Start\"", g.World.CurrentRoom().Name)
	}
}

// TestHandleMedusaSafeToLeaveAfterDefeat covers the OTHER real way past
// her: checkMirrorMedusa (dropping the Mirror in her room) destroys her
// outright, after which leaving is safe even without still carrying a
// Mirror (MonsterHealth<=0 means the move-attempt check no longer
// applies at all).
func TestHandleMedusaSafeToLeaveAfterDefeat(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Start", Exits: map[world.Direction]world.RoomID{world.North: 1}})
	w.AddRoom(&world.Room{ID: 1, Name: "The Pit", Monster: "Medusa", MonsterHealth: 3, Exits: map[world.Direction]world.RoomID{world.South: 0}})
	g := &Game{Player: character.NewPlayer(), World: w}
	g.Player.Items = append(g.Player.Items, "Mirror")

	g.Handle(parser.Parse("NORTH"))       // enters The Pit safely
	g.Handle(parser.Parse("DROP MIRROR")) // triggers checkMirrorMedusa

	got := g.Handle(parser.Parse("SOUTH")) // no Mirror carried anymore
	if strings.Contains(got, "GAME OVER") {
		t.Errorf("Handle(SOUTH) after Medusa is defeated = %q, want it to succeed", got)
	}
	if g.World.CurrentRoom().Name != "Start" {
		t.Errorf("current room after leaving a defeated Medusa's room = %q, want \"Start\"", g.World.CurrentRoom().Name)
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

// TestHandleWaterBlocksMovementUntilCleared covers round 181's real
// fix: world.Room.Water (since round 169) was only ever checked by
// passWater ("WATER, FALL") and reported as a LOOK-time hint - it
// never actually blocked movement, so a player could walk straight
// past a real Water hazard without ever saying the confirmed command
// (flagged directly: "'WATER, FALL' should be a hindrance until
// that's spoken"). Unlike Fire (blocks entering the room ahead), Water
// blocks LEAVING the current room in any direction until cleared, since
// passWater's own real command is said from within the watery room
// itself - blocking entry would make it unreachable to clear at all.
func TestHandleWaterBlocksMovementUntilCleared(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Flooded", Water: true, Exits: map[world.Direction]world.RoomID{world.North: 1}})
	w.AddRoom(&world.Room{ID: 1, Name: "Dry Room"})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("NORTH"))
	if !strings.Contains(got, "water blocks") {
		t.Errorf("Handle(NORTH) out of a Water room before clearing it = %q, want it blocked", got)
	}
	if g.World.CurrentRoom().Name != "Flooded" {
		t.Errorf("current room after a blocked Water move = %q, want unchanged (Flooded)", g.World.CurrentRoom().Name)
	}

	g.Handle(parser.Parse("WATER, FALL"))
	got = g.Handle(parser.Parse("NORTH"))
	if strings.Contains(got, "water blocks") {
		t.Errorf("Handle(NORTH) after WATER, FALL = %q, want it to succeed", got)
	}
	if g.World.CurrentRoom().Name != "Dry Room" {
		t.Errorf("current room after a cleared Water move = %q, want \"Dry Room\"", g.World.CurrentRoom().Name)
	}
}

// TestHandleLookMentionsGuards covers round 152's real fix: before this,
// world.Room.Guards (a real, already-functional "GUARDS, DOOR" obstacle
// since round 19) had no LOOK-time hint at all - a player with no way to
// already know Guards were present would never think to try that exact
// command. Uses a synthetic room since Level1Grid/Level2Grid's real
// Guards placements aren't reachable from a fresh default-mode game.
func TestHandleLookMentionsGuards(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Gate", Guards: true})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "Guards") {
		t.Errorf("Handle(LOOK) with real Guards present = %q, want it mentioned", got)
	}
}

// TestHandleLookMentionsWater covers round 172's real fix: round 169's
// Water hazard (see world.Room.Water and game.passWater) never got the
// same LOOK-time hint round 152 already gave Guards/locked doors - a
// player standing in a real Water room had no hint "WATER, FALL" was
// even relevant there.
func TestHandleLookMentionsWater(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Pool", Water: true})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "water") {
		t.Errorf("Handle(LOOK) with real Water present = %q, want it mentioned", got)
	}
}

// TestHandleLookMentionsLockedDoor covers round 152's real fix for the
// other half of the same "confirmed but unsurfaced at LOOK-time" gap:
// world.Room.DoorPasswords/TollItem (real since rounds 9/64) never had
// any LOOK-time hint either - only discoverable by already guessing the
// right "DOOR, <word>" command. Deliberately checks the hint does NOT
// leak the actual password/toll item, since no source confirms a real
// in-game hint text - only the bare fact that a door exists.
func TestHandleLookMentionsLockedDoor(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Passworded", DoorPasswords: []string{"SILENCE"}})
	w.AddRoom(&world.Room{ID: 1, Name: "Tolled", TollItem: "Key"})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "locked door") {
		t.Errorf("Handle(LOOK) with real DoorPasswords present = %q, want a locked-door hint", got)
	}
	if strings.Contains(got, "SILENCE") {
		t.Errorf("Handle(LOOK) = %q, must not leak the real password", got)
	}

	g.World.Teleport(1)
	got = g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "locked door") {
		t.Errorf("Handle(LOOK) with real TollItem present = %q, want a locked-door hint", got)
	}
	if strings.Contains(got, "Key") {
		t.Errorf("Handle(LOOK) = %q, must not leak the real toll item", got)
	}
}

// TestHandleLookHintsAtAdjacentFire covers round 152's third real fix:
// move()'s own pre-move Fire check (round 80) has always blocked a
// Fire-hazard exit, but LOOK never hinted at it beforehand - a player
// only ever discovered it by trying to move there and getting rejected.
// Checks both the no-Clasp (hinted) and has-Clasp (no longer relevant,
// silent) cases, mirroring move's own exact behavior.
func TestHandleLookHintsAtAdjacentFire(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Start", Exits: map[world.Direction]world.RoomID{world.North: 1}})
	w.AddRoom(&world.Room{ID: 1, Name: "Blaze", Fire: true})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "Flames block the way North") {
		t.Errorf("Handle(LOOK) next to a Fire room without Clasp = %q, want a real hint", got)
	}

	g.Player.Items = append(g.Player.Items, "Clasp")
	got = g.Handle(parser.Parse("LOOK"))
	if strings.Contains(got, "Flames block") {
		t.Errorf("Handle(LOOK) next to a Fire room WITH Clasp = %q, want no hint (matches move's own no-longer-blocked behavior)", got)
	}
}

// TestHandleLookHintsAtNearbyMonster covers round 181's real,
// partially-confirmed mechanism: a "MONSTER NEARBY" panel appears in
// real gameplay footage just before entering a room with a live
// monster (letter codes not legible enough to decode fully) - modeled
// honestly as a proactive proximity warning, mirroring
// TestHandleLookHintsAtAdjacentFire's own pattern, not a fabricated
// random-encounter spawner.
func TestHandleLookHintsAtNearbyMonster(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Start", Exits: map[world.Direction]world.RoomID{world.North: 1}})
	w.AddRoom(&world.Room{ID: 1, Name: "Den", Monster: "Werewolf", MonsterHealth: 2})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "sense a monster nearby, to the North") {
		t.Errorf("Handle(LOOK) next to a live-monster room = %q, want a real proximity hint", got)
	}

	w.Rooms[1].MonsterHealth = 0
	got = g.Handle(parser.Parse("LOOK"))
	if strings.Contains(got, "sense a monster") {
		t.Errorf("Handle(LOOK) next to a DEFEATED monster's room = %q, want no hint", got)
	}
}

// TestHandleLookMarksLevelChangingExit covers round 181's response to
// a direct user question: does this port show the original's own
// special exit graphics for a level-changing direction? A frame-by-
// frame search found real, confirmed evidence the mechanic exists
// (Agile Stair's own status line reads "Level 3" then "Level 4" while
// nominally the same room) but never caught the original's exact
// glyph clearly enough to reproduce - this surfaces the same real,
// sourced fact (each room's own Level) with this port's own plain-text
// marker instead ("^" up, "v" down), honestly distinct from a
// same-level exit.
func TestHandleLookMarksLevelChangingExit(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Landing", Level: 3, Exits: map[world.Direction]world.RoomID{world.North: 1, world.East: 2}})
	w.AddRoom(&world.Room{ID: 1, Name: "Upstairs", Level: 4})
	w.AddRoom(&world.Room{ID: 2, Name: "SameFloor", Level: 3})
	g := &Game{Player: character.NewPlayer(), World: w}

	got := g.Handle(parser.Parse("LOOK"))
	if !strings.Contains(got, "North^") {
		t.Errorf("Handle(LOOK) with a level-4 exit North from level 3 = %q, want \"North^\"", got)
	}
	if strings.Contains(got, "East^") || strings.Contains(got, "Eastv") {
		t.Errorf("Handle(LOOK) with a same-level exit East = %q, want no level marker", got)
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

func TestLevel1ExplorationReachingExitWins(t *testing.T) {
	// Real, validated path from the start room (A1) to the Exit cell
	// (G3), found via BFS over Level1Grid's confirmed connectivity:
	// A1-B1-C1-C2-C3-C4-D4-E4-F4-G4-G3. D4 has a real Guards obstacle
	// (round 60) - since round 181's fix (Guards now actually blocks
	// movement, not just a LOOK-time hint - see game.move's doc
	// comment), "GUARDS, DOOR" is a real, required step on this path,
	// not just flavor.
	g := NewLevel1Exploration()
	path := []string{"SOUTH", "SOUTH", "EAST", "EAST", "EAST", "SOUTH", "GUARDS, DOOR", "SOUTH", "SOUTH", "SOUTH", "WEST"}
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

// TestReachingExitBelowPhilosophusStillWinsButNotesTheRealRequirement
// covers round 147's find (the official map poster's own footer
// banner: "TO LOCATE ALL 3 EXITS, AXIL MUST BECOME PHILOSOPHUS") -
// Won must still become true for a default (Neophyte) player, since
// this port has no confirmed path to Philosophus yet and shouldn't
// silently make its own win condition unreachable, but the real,
// sourced requirement is surfaced as an honest additional note.
func TestReachingExitBelowPhilosophusStillWinsButNotesTheRealRequirement(t *testing.T) {
	g := NewLevel1Exploration()
	path := []string{"SOUTH", "SOUTH", "EAST", "EAST", "EAST", "SOUTH", "GUARDS, DOOR", "SOUTH", "SOUTH", "SOUTH", "WEST"}
	var last string
	for _, dir := range path {
		last = g.Handle(parser.Parse(dir))
	}
	if !g.Won {
		t.Fatalf("Game.Won = false after reaching the Exit room as a Neophyte; last Handle() output: %q", last)
	}
	if !strings.Contains(last, "Philosophus") {
		t.Errorf("Handle() output on reaching Exit below Philosophus = %q, want it to note the real Grade requirement", last)
	}
}

// TestReachingExitAsPhilosophusOmitsTheNote is a regression guard: a
// player who already holds the Philosophus Grade (or higher) shouldn't
// see the below-Philosophus note, since the real requirement is met.
func TestReachingExitAsPhilosophusOmitsTheNote(t *testing.T) {
	g := NewLevel1Exploration()
	g.Player.Grade = character.Philosophus
	path := []string{"SOUTH", "SOUTH", "EAST", "EAST", "EAST", "SOUTH", "GUARDS, DOOR", "SOUTH", "SOUTH", "SOUTH", "WEST"}
	var last string
	for _, dir := range path {
		last = g.Handle(parser.Parse(dir))
	}
	if strings.Contains(last, "Philosophus") {
		t.Errorf("Handle() output on reaching Exit as Philosophus = %q, want no Grade-requirement note", last)
	}
}

// TestNougatDefeatsWerewolfOnDrop covers the real, sourced alternate
// mechanic (CASA walkthrough: Werewolves "killable by walking through
// after dropping NOUGAT" — see checkNougatWerewolf). Path to C2 (a real
// Werewolf, per Level1Grid): A1-South-B1-South-C1-East-C2.
// TestPelletDefeatsSlugOnDrop covers round 135's real, sourced
// alternate mechanic (World of Spectrum's plain-text instructions
// file: Slugs need "a Pellet") - see checkPelletSlug. Level2Grid's
// real Slug (C2) isn't on a simple path from the start room in this
// file's own connectivity data, so teleports there directly (the same
// mechanism astarotTeleport uses) rather than walking an unrelated
// path just to reach it.
func TestPelletDefeatsSlugOnDrop(t *testing.T) {
	g := NewLevel2Exploration()
	id, ok := g.World.FindRoomByName("C2")
	if !ok {
		t.Fatal("test setup bug: Level2Grid has no room named C2")
	}
	g.World.Teleport(id)
	room := g.World.CurrentRoom()
	if room.Monster != "Slug" || room.MonsterHealth <= 0 {
		t.Fatalf("test setup bug: expected a live Slug at C2, got %+v", room)
	}
	g.Player.Items = append(g.Player.Items, "Pellet")
	got := g.Handle(parser.Parse("DROP PELLET"))
	if room.MonsterHealth > 0 {
		t.Errorf("Slug should be defeated after dropping Pellet, MonsterHealth = %d", room.MonsterHealth)
	}
	if !strings.Contains(got, "Pellet") {
		t.Errorf("Handle(DROP PELLET) with a live Slug present = %q, want it to mention the Pellet mechanic", got)
	}
}

// TestSnakeWardsOffHydraOnDrop covers round 145's real, sourced ward-
// off mechanic (World of Spectrum's plain-text instructions file: "To
// pass the Hydras you need a Snake") - see checkSnakeHydra. No shipped
// World.Room carries Monster == "Hydra" yet (a real, honest scope gap
// - see checkSnakeHydra's doc comment), so this uses a synthetic room,
// same pattern as TestHandleFireBlocksMovementWithoutClasp.
func TestSnakeWardsOffHydraOnDrop(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Lair", Monster: "Hydra", MonsterHealth: 3})
	g := &Game{Player: character.NewPlayer(), World: w}
	g.Player.Items = append(g.Player.Items, "Snake")

	got := g.Handle(parser.Parse("DROP SNAKE"))
	room := g.World.CurrentRoom()
	if room.MonsterHealth > 0 {
		t.Errorf("Hydra should be warded off after dropping Snake, MonsterHealth = %d", room.MonsterHealth)
	}
	if !strings.Contains(got, "Snake") {
		t.Errorf("Handle(DROP SNAKE) with a live Hydra present = %q, want it to mention the Snake mechanic", got)
	}
}

// TestMirrorDestroysMedusaOnDrop covers round 178's real, video-
// confirmed instant-kill mechanic ("THE MIRROR DESTROYS MEDUSA") - see
// checkMirrorMedusa. Mirror (Trollwynd) and Medusa (Level4Grid H5) are
// 2 different, unmerged datasets, so this uses a synthetic room, same
// pattern as TestSnakeWardsOffHydraOnDrop.
func TestMirrorDestroysMedusaOnDrop(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "The Pit", Monster: "Medusa", MonsterHealth: 3})
	g := &Game{Player: character.NewPlayer(), World: w}
	g.Player.Items = append(g.Player.Items, "Mirror")

	got := g.Handle(parser.Parse("DROP MIRROR"))
	room := g.World.CurrentRoom()
	if room.MonsterHealth > 0 {
		t.Errorf("Medusa should be destroyed after dropping the Mirror, MonsterHealth = %d", room.MonsterHealth)
	}
	if !strings.Contains(got, "Mirror") {
		t.Errorf("Handle(DROP MIRROR) with a live Medusa present = %q, want it to mention the Mirror mechanic", got)
	}
}

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
	g.Handle(parser.Parse("PICKUP KEY"))    // real Room of Stings toll currency (round 180)
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
	g.Handle(parser.Parse("DOOR, WOLF")) // unlocks Wolfdorp's own door (round 181: now actually required)
	g.Handle(parser.Parse("NORTH-WEST")) // Room of Stings
	g.Handle(parser.Parse("DROP KEY"))   // pays Room of Stings' own real toll (round 181: now actually required)
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

// TestCollodonsPileSlatDefeatsNidusCyclops covers round 146's real,
// sourced instant-kill mechanic (World of Spectrum's plain-text
// instructions file: "the slat kills the Cyclops") end-to-end in the
// actual DEFAULT (CollodonsPile) game: picks up the real, already-
// placed Slat in Morfang, carries it East through Room of Arrows to
// Nidus (all real, already-connected rooms on the confirmed
// walkthrough path), and drops it on the real Cyclops there. Also
// exercises round 181's real locked-door fix along the whole path
// (Room of Stings/Morfang/Room of Arrows all now genuinely require
// paying their real TollItem to proceed) - including the Slat's own
// double duty (Room of Arrows' toll AND the Cyclops kill both need it,
// resolved by payToll leaving the paid item retrievable rather than
// deleting it, per its own doc comment).
func TestCollodonsPileSlatDefeatsNidusCyclops(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))          // Secunda Porta
	g.Handle(parser.Parse("DOOR, SILENCE")) // unlocks the door North
	g.Handle(parser.Parse("NORTH"))         // Trollwynd
	g.Handle(parser.Parse("PICKUP KEY"))    // real Room of Stings toll currency (round 180)
	g.Handle(parser.Parse("NORTH"))         // Agile Stair
	g.Handle(parser.Parse("SOUTH-EAST"))    // Methos
	g.Handle(parser.Parse("SOUTH"))         // Sothic Complex
	g.Handle(parser.Parse("SOUTH"))         // Wolfdorp
	g.Handle(parser.Parse("PICKUP BAG"))    // real Morfang toll currency (round 64)
	g.Handle(parser.Parse("DOOR, WOLF"))    // unlocks Wolfdorp's own door (round 181: now actually required)
	g.Handle(parser.Parse("NORTH-WEST"))    // Room of Stings
	g.Handle(parser.Parse("DROP KEY"))      // pays Room of Stings' own real toll (round 181: now actually required)
	g.Handle(parser.Parse("NORTH"))         // Morfang - has the real Slat
	if got := g.Handle(parser.Parse("PICKUP SLAT")); !strings.Contains(got, "Slat") {
		t.Fatalf("test setup bug: PICKUP SLAT in Morfang = %q, want it to succeed", got)
	}
	g.Handle(parser.Parse("DROP BAG"))  // pays Morfang's own real toll (round 181: now actually required)
	g.Handle(parser.Parse("EAST"))      // Room of Arrows
	g.Handle(parser.Parse("DROP SLAT")) // pays Room of Arrows' own real toll (round 181: now actually required)
	// ROUND 181: payToll leaves the paid item retrievable on the real
	// table it's placed on (see payToll's doc comment) rather than
	// deleting it outright - the SAME Slat that just paid this room's
	// toll is picked back up here and carried onward to Nidus, since
	// game.checkSlatCyclops's own real mechanic needs it there too. No
	// source suggests 2 separate Slats exist; this is the honest way
	// both already-verified mechanics stay reachable together.
	if got := g.Handle(parser.Parse("PICKUP SLAT")); !strings.Contains(got, "Slat") {
		t.Fatalf("test setup bug: PICKUP SLAT (back) in Room of Arrows = %q, want it to succeed", got)
	}
	g.Handle(parser.Parse("EAST")) // Nidus - has the real Cyclops

	room := g.World.CurrentRoom()
	if room.Name != "Nidus" || room.Monster != "Cyclops" || room.MonsterHealth <= 0 {
		t.Fatalf("test setup bug: expected a live Cyclops in Nidus, got %+v", room)
	}
	got := g.Handle(parser.Parse("DROP SLAT"))
	if room.MonsterHealth > 0 {
		t.Errorf("Nidus's Cyclops should be defeated after dropping the real Slat carried from Morfang, MonsterHealth = %d", room.MonsterHealth)
	}
	if !strings.Contains(got, "Slat") {
		t.Errorf("Handle(DROP SLAT) with a live Cyclops present = %q, want it to mention the Slat mechanic", got)
	}
}

// TestNuggetAlsoDefeatsWerewolfOnDrop covers round 146's extension:
// Nugget (not just Nougat) also wards off Werewolves, per 2
// independent sources (see checkNougatWerewolf's doc comment).
func TestNuggetAlsoDefeatsWerewolfOnDrop(t *testing.T) {
	g := NewLevel1Exploration()
	for _, dir := range []string{"SOUTH", "SOUTH", "EAST"} {
		g.Handle(parser.Parse(dir))
	}
	room := g.World.CurrentRoom()
	if room.Monster != "Werewolf" || room.MonsterHealth <= 0 {
		t.Fatalf("test setup bug: expected a live Werewolf at C2, got %+v", room)
	}
	g.Player.Items = append(g.Player.Items, "Nugget")
	got := g.Handle(parser.Parse("DROP NUGGET"))
	if room.MonsterHealth > 0 {
		t.Errorf("Werewolf should be defeated after dropping Nugget, MonsterHealth = %d", room.MonsterHealth)
	}
	if !strings.Contains(got, "Nugget") {
		t.Errorf("Handle(DROP NUGGET) with a live Werewolf present = %q, want it to mention the Nugget mechanic", got)
	}
}

// TestSilverNuggetAlsoDefeatsWerewolfOnDrop covers round 149's
// extension: "Silver Nugget" (not just bare "Nugget") also wards off
// Werewolves, per 2 independent sources both qualifying the item as
// silver (see checkNougatWerewolf's doc comment).
func TestSilverNuggetAlsoDefeatsWerewolfOnDrop(t *testing.T) {
	g := NewLevel1Exploration()
	for _, dir := range []string{"SOUTH", "SOUTH", "EAST"} {
		g.Handle(parser.Parse(dir))
	}
	room := g.World.CurrentRoom()
	if room.Monster != "Werewolf" || room.MonsterHealth <= 0 {
		t.Fatalf("test setup bug: expected a live Werewolf at C2, got %+v", room)
	}
	g.Player.Items = append(g.Player.Items, "Silver Nugget")
	got := g.Handle(parser.Parse("DROP SILVER NUGGET"))
	if room.MonsterHealth > 0 {
		t.Errorf("Werewolf should be defeated after dropping the Silver Nugget, MonsterHealth = %d", room.MonsterHealth)
	}
	if !strings.Contains(got, "Nugget") {
		t.Errorf("Handle(DROP SILVER NUGGET) with a live Werewolf present = %q, want it to mention the Nugget mechanic", got)
	}
}

func TestLevel1ExplorationMovementAndCombat(t *testing.T) {
	g := NewLevel1Exploration()
	withGrimoire(g) // BLAST now requires it - see spellRequiresItem
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

func TestLevel4ExplorationMovement(t *testing.T) {
	g := NewLevel4Exploration()
	got := g.Handle(parser.Parse("EAST"))
	if strings.Contains(got, "can't go that way") {
		t.Errorf("Handle(EAST) from the start room = %q, want a successful move", got)
	}
}

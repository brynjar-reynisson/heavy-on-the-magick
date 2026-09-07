package main

import (
	"image/color"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/game"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/world"
)

// TestMonsterGlyphColorMatchesLegend pins the exact letter+color pairs
// read from heavymap-grid-clean.gif's own legend (see monsterGlyphColor's
// doc comment) - a regression check against a future accidental edit,
// since these values are load-bearing facts, not arbitrary choices.
func TestMonsterGlyphColorMatchesLegend(t *testing.T) {
	want := map[string]string{
		"Troll": "t", "Cyclops": "c", "Ghost": "g", "Slug": "s",
		"Vampire": "w", "Medusa": "m", "Werewolf": "w", "Wyvern": "w",
	}
	for name, letter := range want {
		gc, ok := monsterGlyphColor[name]
		if !ok {
			t.Errorf("monsterGlyphColor is missing %q", name)
			continue
		}
		if gc.letter != letter {
			t.Errorf("monsterGlyphColor[%q].letter = %q, want %q", name, gc.letter, letter)
		}
	}
	// Vampire and Medusa share one exact color per the legend (only the
	// letter distinguishes them); Troll and Cyclops likewise.
	if monsterGlyphColor["Vampire"].c != monsterGlyphColor["Medusa"].c {
		t.Error("Vampire and Medusa should share the same confirmed color")
	}
	if monsterGlyphColor["Troll"].c != monsterGlyphColor["Cyclops"].c {
		t.Error("Troll and Cyclops should share the same confirmed color")
	}
	if monsterGlyphColor["Ghost"].c == monsterGlyphColor["Slug"].c {
		t.Error("Ghost and Slug are confirmed distinct shades of green, should not match")
	}
	if got := monsterGlyphColor["Wyvern"].c; got != (color.RGBA{0, 132, 255, 255}) {
		t.Errorf("Wyvern color = %v, want the confirmed rare blue RGB(0,132,255)", got)
	}
}

// TestGuardsColorMatchesLegend pins guardsColor's exact confirmed RGB
// against regression.
func TestGuardsColorMatchesLegend(t *testing.T) {
	if guardsColor != (color.RGBA{255, 0, 0, 255}) {
		t.Errorf("guardsColor = %v, want the confirmed bright red RGB(255,0,0)", guardsColor)
	}
}

// TestInvokeCommandForPicksCarriedCharm covers the I key's real
// target-selection logic (see invokeDemonForGroundedCharm's doc
// comment): with no text input, it scans a set of item names (the
// current room's real Items, since round 131's Charm-on-the-ground
// correction) for any demon's Charm.
func TestInvokeCommandForPicksCarriedCharm(t *testing.T) {
	if got := invokeCommandFor([]string{"Grimoire", "sword"}); got != "INVOKE ASTAROT" {
		t.Errorf("invokeCommandFor with a Sword present = %q, want \"INVOKE ASTAROT\"", got)
	}
	if got := invokeCommandFor([]string{"Grimoire"}); got != "INVOKE" {
		t.Errorf("invokeCommandFor with no Charm present = %q, want bare \"INVOKE\"", got)
	}
}

// TestApexPortraitShouldShow covers the K key's real portrait-display
// logic (see drawApexPortrait's doc comment): the portrait only appears
// right after a real "APEX, TALK" response, not for unrelated log lines.
func TestApexPortraitShouldShow(t *testing.T) {
	if apexPortraitShouldShow(nil) {
		t.Error("apexPortraitShouldShow(nil) = true, want false: no log yet")
	}
	if apexPortraitShouldShow([]string{"Room of Misery", "Exits: East"}) {
		t.Error("apexPortraitShouldShow with an unrelated last line = true, want false")
	}
	if !apexPortraitShouldShow([]string{"Room of Misery", "Apex the Ogre eyes you warily, then grunts. He might share what he knows, if you treat him with respect."}) {
		t.Error("apexPortraitShouldShow after a real talkToApex response = false, want true")
	}
}

// TestInvokedDemonPortraitName covers the demon-portrait priority logic
// (see currentPortraitName's doc comment): only a real "You invoke ..."
// success line should match, and it should correctly pick out which of
// the 4 confirmed demons was named.
func TestInvokedDemonPortraitName(t *testing.T) {
	name, ok := invokedDemonPortraitName([]string{"You invoke ASTAROT, the Spirit of Assemblage! In an instant, you are transported to Wolfdorp."})
	if !ok || name != "astarot" {
		t.Errorf("invokedDemonPortraitName(ASTAROT line) = (%q, %v), want (\"astarot\", true)", name, ok)
	}
	if _, ok := invokedDemonPortraitName([]string{"Room of Misery"}); ok {
		t.Error("invokedDemonPortraitName with an unrelated last line = true, want false")
	}
	if _, ok := invokedDemonPortraitName(nil); ok {
		t.Error("invokedDemonPortraitName(nil) = true, want false")
	}
}

// TestCurrentPortraitNamePriority covers currentPortraitName's real
// priority order: a just-happened conversation/invocation outranks the
// room's ambient monster, which itself only shows if a portrait for it
// actually exists.
func TestCurrentPortraitNamePriority(t *testing.T) {
	gui := &GUI{g: game.New(), portraits: map[string]*ebiten.Image{}}

	if _, ok := gui.currentPortraitName(); ok {
		t.Error("currentPortraitName with no log yet and no monster = true, want false")
	}

	gui.g.World.CurrentRoom().Monster = "Vampire"
	gui.g.World.CurrentRoom().MonsterHealth = 2
	// No "vampire" key in gui.portraits (map is empty) - should NOT claim
	// a monster portrait it doesn't actually have loaded.
	if _, ok := gui.currentPortraitName(); ok {
		t.Error("currentPortraitName claimed a monster portrait not present in gui.portraits")
	}

	gui.portraits["vampire"] = &ebiten.Image{}
	if name, ok := gui.currentPortraitName(); !ok || name != "vampire" {
		t.Errorf("currentPortraitName with a live Vampire and a loaded portrait = (%q, %v), want (\"vampire\", true)", name, ok)
	}

	gui.log = []string{"You invoke ASTAROT, the Spirit of Assemblage! In an instant, you are transported to Wolfdorp."}
	gui.portraits["astarot"] = &ebiten.Image{}
	if name, ok := gui.currentPortraitName(); !ok || name != "astarot" {
		t.Errorf("currentPortraitName should prioritize a just-invoked demon over the room's monster, got (%q, %v)", name, ok)
	}
}

// TestSubmitTypedCommandBlankDoesNothing covers round 174's real
// typed-command line (see updateTyping's doc comment): pressing ENTER
// with nothing typed should not resubmit anything or echo an empty
// prompt line.
func TestSubmitTypedCommandBlankDoesNothing(t *testing.T) {
	g := game.New()
	if echo, result := submitTypedCommand(g, "   "); echo != "" || result != "" {
		t.Errorf("submitTypedCommand(blank) = (%q, %q), want (\"\", \"\")", echo, result)
	}
}

// TestSubmitTypedCommandExpandsRealAbbreviation confirms a single typed
// letter ("W") goes through the exact same ExpandKeyword step the text
// frontend uses (parser.ExpandKeyword: "W" -> "WEST") and actually moves
// the player - the real fix for this round's request that keypresses
// behave like the original's own Merphish grammar, not this port's
// former WASD-as-North roguelike convention.
func TestSubmitTypedCommandExpandsRealAbbreviation(t *testing.T) {
	g := game.New()
	before := g.World.Current
	echo, result := submitTypedCommand(g, "e")
	if echo != "> E" {
		t.Errorf("submitTypedCommand(\"e\") echo = %q, want \"> E\"", echo)
	}
	if g.World.Current == before {
		t.Errorf("submitTypedCommand(\"e\") did not move the player; result = %q", result)
	}
}

// TestSubmitTypedCommandSupportsDiagonal confirms a real hyphenated
// compass word reaches the game exactly as the original expects - the
// specific gap this round's request named ("arrows... don't capture
// something like NORTH-EAST").
func TestSubmitTypedCommandSupportsDiagonal(t *testing.T) {
	g := game.New()
	_, result := submitTypedCommand(g, "NORTH-EAST")
	if strings.Contains(result, "don't understand") {
		t.Errorf("submitTypedCommand(\"NORTH-EAST\") = %q, want a real recognized direction response", result)
	}
}

// TestSubmitTypedCommandSupportsConversationForm confirms the real
// "NAME, OBJECT" grammar (e.g. talking to Apex) still works when typed
// through this new input line, not just the single-key K shortcut.
func TestSubmitTypedCommandSupportsConversationForm(t *testing.T) {
	g := game.New()
	_, result := submitTypedCommand(g, "APEX, TALK")
	if !strings.Contains(result, "Apex the Ogre") {
		t.Errorf("submitTypedCommand(\"APEX, TALK\") = %q, want a real Apex response", result)
	}
}

// TestWrapLineKeepsShortLinesWhole confirms wrapLine doesn't needlessly
// split a line that already fits.
func TestWrapLineKeepsShortLinesWhole(t *testing.T) {
	got := wrapLine("Room of Misery", 26)
	if len(got) != 1 || got[0] != "Room of Misery" {
		t.Errorf("wrapLine(short line) = %v, want [\"Room of Misery\"]", got)
	}
}

// TestWrapLineBreaksOnWordBoundaries covers the real bug this function
// fixes (see its own doc comment): a long game.Handle response line -
// this exact one bled across the middle panel's border into the stats
// panel before wrapping existed - must break into multiple lines, none
// exceeding maxChars, and never splitting a word in half.
func TestWrapLineBreaksOnWordBoundaries(t *testing.T) {
	const maxChars = 26
	got := wrapLine("(room description not yet extracted from the original game)", maxChars)
	if len(got) < 2 {
		t.Fatalf("wrapLine(long line) = %v, want more than 1 line", got)
	}
	for _, line := range got {
		if len(line) > maxChars {
			t.Errorf("wrapLine line %q exceeds maxChars=%d", line, maxChars)
		}
	}
	if strings.Join(got, " ") != "(room description not yet extracted from the original game)" {
		t.Errorf("wrapLine(long line) = %v, lost or reordered words", got)
	}
}

// TestStatsPanelLines covers the right status-bar panel's real
// confirmed Stamina/Skill/Luck/XP content, matching the real reference
// SpecEmu screenshot's own green STAMINA/SKILL/LUCK panel (see
// statsPanelLines's doc comment for the honest XP-addition caveat).
func TestStatsPanelLines(t *testing.T) {
	g := game.New()
	g.Player.Stamina = 30
	g.Player.MaxStamina = 39
	g.Player.Skill = 8
	g.Player.Luck = 4
	g.Player.ExperiencePoints = 21

	got := strings.Join(statsPanelLines(g), "\n")
	for _, want := range []string{"STAMINA 30/39", "SKILL   8", "LUCK    4", "XP      21"} {
		if !strings.Contains(got, want) {
			t.Errorf("statsPanelLines() = %q, want it to contain %q", got, want)
		}
	}
}

// TestExitsPanelLinesMatchRealScreenshot pins the compass-grid layout
// against the real SpecEmu screenshot of Room of Misery, which showed
// "EXITS:" with W on the left and E on the right of the same box - the
// panel exitsPanelLines is meant to approximate.
func TestExitsPanelLinesMatchRealScreenshot(t *testing.T) {
	room := &world.Room{Exits: map[world.Direction]world.RoomID{
		world.West: 1,
		world.East: 2,
	}}
	lines := exitsPanelLines(room)
	if lines[0] != "EXITS:" {
		t.Errorf("exitsPanelLines[0] = %q, want \"EXITS:\"", lines[0])
	}
	// The W/E compass row (index 3, see exitsPanelLines: "EXITS:", "",
	// NW/N/NE row, W/E row, SW/S/SE row) should show W on the left and E
	// on the right, matching the real screenshot's own left/right layout.
	weRow := lines[3]
	wIdx := strings.Index(weRow, "W")
	eIdx := strings.Index(weRow, "E")
	if wIdx == -1 || eIdx == -1 || wIdx >= eIdx {
		t.Errorf("exitsPanelLines W/E row = %q, want W positioned before E (left/right compass layout)", weRow)
	}
}

// TestExitsPanelLinesNoExits confirms an exit-less room (e.g. an
// isolated named cell) still renders a stable, non-panicking panel.
func TestExitsPanelLinesNoExits(t *testing.T) {
	room := &world.Room{}
	lines := exitsPanelLines(room)
	if lines[0] != "EXITS:" {
		t.Errorf("exitsPanelLines(no exits)[0] = %q, want \"EXITS:\"", lines[0])
	}
}

// TestRoomStatusLines covers the left panel's OTHER real mode (see
// showRoomStatus's doc comment - a real observed SWAP/"Z" effect):
// showing the room name/level/grade instead of Exits.
func TestRoomStatusLines(t *testing.T) {
	g := game.New()
	room := g.World.CurrentRoom()
	room.Name = "Sothic Complex"
	room.Level = 2

	got := strings.Join(roomStatusLines(g), "\n")
	for _, want := range []string{"YOU ARE IN THE", "SOTHIC COMPLEX", "ON LEVEL 2", "YOUR GRADE IS", "NEOPHYTE"} {
		if !strings.Contains(got, want) {
			t.Errorf("roomStatusLines() = %q, want it to contain %q", got, want)
		}
	}
}

package main

import (
	"image/color"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/game"
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

// TestItemsColorIsLegible pins itemsColor away from black - see its
// doc comment for why the confirmed real "object" icon color (black)
// can't be used against this GUI's black background.
func TestItemsColorIsLegible(t *testing.T) {
	if itemsColor == (color.RGBA{0, 0, 0, 255}) {
		t.Error("itemsColor must not be black - it would be invisible against this GUI's black background")
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

// TestFixturesText covers round 151's real, previously-unsurfaced-in-
// the-GUI HasTable/HasChest content (see drawFixtures's doc comment).
func TestFixturesText(t *testing.T) {
	if got := fixturesText(false, false); got != "" {
		t.Errorf("fixturesText(false, false) = %q, want empty", got)
	}
	if got := fixturesText(true, false); got != "Table" {
		t.Errorf("fixturesText(true, false) = %q, want \"Table\"", got)
	}
	if got := fixturesText(false, true); got != "Chest" {
		t.Errorf("fixturesText(false, true) = %q, want \"Chest\"", got)
	}
	if got := fixturesText(true, true); got != "Table, Chest" {
		t.Errorf("fixturesText(true, true) = %q, want \"Table, Chest\"", got)
	}
}

// TestStatsLine constructs a bare *GUI directly (not via NewGUI, which
// touches ebiten's audio/image APIs and needs a real display/audio
// device) since statsLine only reads gui.g.Player - a pure formatting
// function worth testing despite living in package main.
func TestStatsLine(t *testing.T) {
	gui := &GUI{g: game.New()}
	gui.g.Player.Stamina = 30
	gui.g.Player.MaxStamina = 39
	gui.g.Player.Skill = 8
	gui.g.Player.Luck = 4
	gui.g.Player.ExperiencePoints = 21

	got := gui.statsLine()
	for _, want := range []string{"Axil", "Neophyte", "Stamina 30/39", "Skill 8", "Luck 4", "XP 21"} {
		if !strings.Contains(got, want) {
			t.Errorf("statsLine() = %q, want it to contain %q", got, want)
		}
	}
}

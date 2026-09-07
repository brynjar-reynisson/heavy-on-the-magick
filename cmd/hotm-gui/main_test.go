package main

import (
	"image/color"
	"strings"
	"testing"

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
// target-selection logic (see invokeCarriedDemon's doc comment): with
// no text input, it scans the player's items for any demon's Charm.
func TestInvokeCommandForPicksCarriedCharm(t *testing.T) {
	if got := invokeCommandFor([]string{"Grimoire", "sword"}); got != "INVOKE ASTAROT" {
		t.Errorf("invokeCommandFor with a carried Sword = %q, want \"INVOKE ASTAROT\"", got)
	}
	if got := invokeCommandFor([]string{"Grimoire"}); got != "INVOKE" {
		t.Errorf("invokeCommandFor with no Charm carried = %q, want bare \"INVOKE\"", got)
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

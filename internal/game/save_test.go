package game

import (
	"os"
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
)

// withTempSaveDir runs fn inside a temporary working directory so save
// tests don't write hotm-save.json/hotm-axil-save.json into the repo.
func withTempSaveDir(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
}

func TestSaveGameRoundTrip(t *testing.T) {
	withTempSaveDir(t)
	g := New()
	g.Handle(parser.Parse("EAST")) // move away from the start room
	g.Player.Stamina = 12

	if err := g.SaveGame(); err != nil {
		t.Fatalf("SaveGame() error = %v", err)
	}

	g2 := New()
	if err := g2.RestoreGame(); err != nil {
		t.Fatalf("RestoreGame() error = %v", err)
	}
	if g2.Player.Stamina != 12 {
		t.Errorf("Stamina after restore = %d, want 12", g2.Player.Stamina)
	}
	if g2.World.CurrentRoom().Name != "Secunda Porta" {
		t.Errorf("current room after restore = %q, want Secunda Porta", g2.World.CurrentRoom().Name)
	}
}

func TestSaveAxilRoundTrip(t *testing.T) {
	withTempSaveDir(t)
	g := New()
	g.Player.Stamina = 7
	g.Player.Items = append(g.Player.Items, "Grimoire")

	if err := g.SaveAxil(); err != nil {
		t.Fatalf("SaveAxil() error = %v", err)
	}

	g2 := New()
	if err := g2.RestoreAxil(); err != nil {
		t.Fatalf("RestoreAxil() error = %v", err)
	}
	if g2.Player.Stamina != 7 {
		t.Errorf("Stamina after RestoreAxil = %d, want 7", g2.Player.Stamina)
	}
	// Round 121: Axil's real starting Pouch (character.NewPlayer) is
	// already in Items before this test appends Grimoire, so the saved/
	// restored set is [Pouch, Grimoire], not just [Grimoire].
	want := []string{"Pouch", "Grimoire"}
	if len(g2.Player.Items) != len(want) {
		t.Fatalf("Items after RestoreAxil = %v, want %v", g2.Player.Items, want)
	}
	for i, w := range want {
		if g2.Player.Items[i] != w {
			t.Errorf("Items after RestoreAxil = %v, want %v", g2.Player.Items, want)
		}
	}
	// RestoreAxil should not touch World.
	if g2.World.CurrentRoom().Name != "Room of Misery" {
		t.Errorf("current room after RestoreAxil = %q, want unchanged Room of Misery", g2.World.CurrentRoom().Name)
	}
}

// TestHandleOptionsNumericSlots pins round 124: the CASA walkthrough's
// own "Tips" section ("SAVE regularly by pressing key O and then
// option 2") and the manual ("select option 6 and the values will be
// realigned") confirm real numbered Option Screen slots - "O 2" and
// "O 6" should now work exactly like the keyword forms.
func TestHandleOptionsNumericSlots(t *testing.T) {
	withTempSaveDir(t)
	g := New()
	if got := g.Handle(parser.Parse("O 2")); got != "Game saved." {
		t.Errorf(`Handle(O 2) = %q, want "Game saved." (option 2 = Save Game, per CASA's own Tips section)`, got)
	}
	got := g.Handle(parser.Parse("O 6"))
	if !strings.HasPrefix(got, "Realign Status:") {
		t.Errorf(`Handle(O 6) = %q, want it to start with "Realign Status:" (option 6, per the manual)`, got)
	}
}

func TestHandleOptionsSaveAndRestoreGame(t *testing.T) {
	withTempSaveDir(t)
	g := New()
	got := g.Handle(parser.Parse("O SAVE GAME"))
	if got != "Game saved." {
		t.Errorf("Handle(O SAVE GAME) = %q, want %q", got, "Game saved.")
	}
	got = g.Handle(parser.Parse("O RESTORE GAME"))
	if got != "Game restored." {
		t.Errorf("Handle(O RESTORE GAME) = %q, want %q", got, "Game restored.")
	}
}

// TestHandleSaveCostsStamina pins the real, sourced mechanic (round 86):
// the manual states "Saving a game will deplete your Stamina, so that a
// Save cannot be used as an easy way of getting round difficult
// choices!" Restoring should NOT cost Stamina (would defeat the point).
func TestHandleSaveCostsStamina(t *testing.T) {
	withTempSaveDir(t)
	g := New()
	before := g.Player.Stamina

	g.Handle(parser.Parse("O SAVE GAME"))
	if g.Player.Stamina != before-saveStaminaCost {
		t.Errorf("Stamina after Save Game = %d, want %d (before %d minus saveStaminaCost %d)", g.Player.Stamina, before-saveStaminaCost, before, saveStaminaCost)
	}

	afterSave := g.Player.Stamina
	g.Handle(parser.Parse("O RESTORE GAME"))
	if g.Player.Stamina != afterSave {
		t.Errorf("Stamina after Restore Game = %d, want unchanged %d (restoring shouldn't cost Stamina)", g.Player.Stamina, afterSave)
	}
}

func TestRestoreGameMissingFile(t *testing.T) {
	withTempSaveDir(t)
	g := New()
	if err := g.RestoreGame(); err == nil {
		t.Error("RestoreGame() with no save file present, want an error")
	}
}

// TestSaveGameVersionRoundTrip pins round 119's real "Version letter"
// save slots (see save.go's doc comment): two different versions are
// genuinely separate files, and restoring one doesn't touch the other.
func TestSaveGameVersionRoundTrip(t *testing.T) {
	withTempSaveDir(t)
	gA := New()
	gA.Handle(parser.Parse("EAST"))
	gA.Player.Stamina = 11
	if err := gA.SaveGameVersion("A"); err != nil {
		t.Fatalf("SaveGameVersion(A) error = %v", err)
	}

	gB := New()
	gB.Player.Stamina = 22
	if err := gB.SaveGameVersion("B"); err != nil {
		t.Fatalf("SaveGameVersion(B) error = %v", err)
	}

	restored := New()
	if err := restored.RestoreGameVersion("A"); err != nil {
		t.Fatalf("RestoreGameVersion(A) error = %v", err)
	}
	if restored.Player.Stamina != 11 {
		t.Errorf("Stamina after RestoreGameVersion(A) = %d, want 11 (version B's save shouldn't be touched)", restored.Player.Stamina)
	}
	if restored.World.CurrentRoom().Name != "Secunda Porta" {
		t.Errorf("current room after RestoreGameVersion(A) = %q, want Secunda Porta", restored.World.CurrentRoom().Name)
	}
}

// TestHandleOptionsSaveGameWithVersionLetter covers the real end-to-end
// Handle path: "O SAVE GAME B" saves to a version-lettered slot,
// distinct from the plain "O SAVE GAME" default slot.
func TestHandleOptionsSaveGameWithVersionLetter(t *testing.T) {
	withTempSaveDir(t)
	g := New()
	g.Player.Stamina = 30

	got := g.Handle(parser.Parse("O SAVE GAME B"))
	if got != "Game saved (version B)." {
		t.Errorf("Handle(O SAVE GAME B) = %q, want %q", got, "Game saved (version B).")
	}
	if _, err := os.Stat("hotm-save-B.json"); err != nil {
		t.Errorf("hotm-save-B.json should exist after O SAVE GAME B: %v", err)
	}
	if _, err := os.Stat("hotm-save.json"); err == nil {
		t.Error("hotm-save.json (the default slot) should NOT exist after only saving to version B")
	}

	got = g.Handle(parser.Parse("O RESTORE GAME B"))
	if got != "Game restored (version B)." {
		t.Errorf("Handle(O RESTORE GAME B) = %q, want %q", got, "Game restored (version B).")
	}
}

package game

import (
	"os"
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
	if len(g2.Player.Items) != 1 || g2.Player.Items[0] != "Grimoire" {
		t.Errorf("Items after RestoreAxil = %v, want [Grimoire]", g2.Player.Items)
	}
	// RestoreAxil should not touch World.
	if g2.World.CurrentRoom().Name != "Room of Misery" {
		t.Errorf("current room after RestoreAxil = %q, want unchanged Room of Misery", g2.World.CurrentRoom().Name)
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

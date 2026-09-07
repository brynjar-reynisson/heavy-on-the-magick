package game

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/character"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/world"
)

// Save/Restore file names. The manual's Option Screen (confirmed real —
// see game.go's OPTIONS doc comment) offers "Save Game"/"Restore Game"
// and "Save Axil"/"Restore Axil" as separate choices; modeled here as
// separate files, one holding the whole session (world + player), the
// other just the player character, matching that Game/Axil distinction
// (the manual doesn't state the original's actual save format or
// mechanism, so the file-based approach and these names are this port's
// own implementation choice, not extracted fact).
//
// Round 119: the manual DOES confirm one real detail about the
// original's save mechanism, though — "When Saving or Restoring a
// game, you will be asked for a Version letter — this is to ensure
// that the right game is restored, so keep a note of Version letters."
// This is real evidence the original supports multiple save slots per
// type (Game/Axil), identified by a letter — not just the single fixed
// slot this port had modeled since round 21. saveFileName below adds
// that: an empty version keeps the exact original filename (so
// existing saves/tests are unaffected), a non-empty one (any letter)
// gets its own separate file.
const (
	gameSaveFile = "hotm-save.json"
	axilSaveFile = "hotm-axil-save.json"
)

// saveFileName returns base unchanged for the empty (default) version,
// or a per-version filename (e.g. "hotm-save-B.json") otherwise.
func saveFileName(base, version string) string {
	if version == "" {
		return base
	}
	return strings.TrimSuffix(base, ".json") + "-" + strings.ToUpper(version) + ".json"
}

type gameSaveData struct {
	Player *character.Player
	World  *world.World
}

// SaveGame writes the full session (player and world state, including
// visited rooms and monster/item state) to gameSaveFile. Equivalent to
// SaveGameVersion("").
func (g *Game) SaveGame() error {
	return g.SaveGameVersion("")
}

// SaveGameVersion is SaveGame, but to a specific real Version letter
// slot (round 119 — see the const block's doc comment) rather than the
// single default file.
func (g *Game) SaveGameVersion(version string) error {
	data := gameSaveData{Player: g.Player, World: g.World}
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(saveFileName(gameSaveFile, version), bytes, 0o644)
}

// RestoreGame replaces the current session with what's in gameSaveFile.
// Equivalent to RestoreGameVersion("").
func (g *Game) RestoreGame() error {
	return g.RestoreGameVersion("")
}

// RestoreGameVersion is RestoreGame, but from a specific real Version
// letter slot rather than the single default file.
func (g *Game) RestoreGameVersion(version string) error {
	bytes, err := os.ReadFile(saveFileName(gameSaveFile, version))
	if err != nil {
		return err
	}
	var data gameSaveData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	g.Player = data.Player
	g.World = data.World
	return nil
}

// SaveAxil writes just the player character (not room/world state) to
// axilSaveFile. Equivalent to SaveAxilVersion("").
func (g *Game) SaveAxil() error {
	return g.SaveAxilVersion("")
}

// SaveAxilVersion is SaveAxil, but to a specific real Version letter
// slot rather than the single default file.
func (g *Game) SaveAxilVersion(version string) error {
	bytes, err := json.MarshalIndent(g.Player, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(saveFileName(axilSaveFile, version), bytes, 0o644)
}

// RestoreAxil replaces just the player character with what's in
// axilSaveFile, leaving the current world/room state untouched.
// Equivalent to RestoreAxilVersion("").
func (g *Game) RestoreAxil() error {
	return g.RestoreAxilVersion("")
}

// RestoreAxilVersion is RestoreAxil, but from a specific real Version
// letter slot rather than the single default file.
func (g *Game) RestoreAxilVersion(version string) error {
	bytes, err := os.ReadFile(saveFileName(axilSaveFile, version))
	if err != nil {
		return err
	}
	var p character.Player
	if err := json.Unmarshal(bytes, &p); err != nil {
		return err
	}
	g.Player = &p
	return nil
}

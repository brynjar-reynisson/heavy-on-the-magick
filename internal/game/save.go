package game

import (
	"encoding/json"
	"os"

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
const (
	gameSaveFile = "hotm-save.json"
	axilSaveFile = "hotm-axil-save.json"
)

type gameSaveData struct {
	Player *character.Player
	World  *world.World
}

// SaveGame writes the full session (player and world state, including
// visited rooms and monster/item state) to gameSaveFile.
func (g *Game) SaveGame() error {
	data := gameSaveData{Player: g.Player, World: g.World}
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(gameSaveFile, bytes, 0o644)
}

// RestoreGame replaces the current session with what's in gameSaveFile.
func (g *Game) RestoreGame() error {
	bytes, err := os.ReadFile(gameSaveFile)
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
// axilSaveFile.
func (g *Game) SaveAxil() error {
	bytes, err := json.MarshalIndent(g.Player, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(axilSaveFile, bytes, 0o644)
}

// RestoreAxil replaces just the player character with what's in
// axilSaveFile, leaving the current world/room state untouched.
func (g *Game) RestoreAxil() error {
	bytes, err := os.ReadFile(axilSaveFile)
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

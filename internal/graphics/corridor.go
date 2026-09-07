package graphics

import (
	"bytes"
	_ "embed"
	"image"
)

// assets/corridor_sample.png is a real, tight-cropped extract of an
// actual in-game room screenshot — sourced from the same
// maps.speccy.cz "Speccy Screenshot Maps" atlas as the demon/monster
// portraits (heavymap-speccy-screenshots.png in the repo root, credited
// to its creator Hippy Smith). Round 96: precisely located via the
// atlas's own printed row/column labels ("A" and "1", tight-crop-
// confirmed the same way every other extraction in this project is) —
// this is Level 2's cell A1, which is already this project's own
// confirmed real starting room for `game.NewLevel2Exploration()` (see
// level2_grid.go: `w := New(level2Room("A1"))`). Not yet wired into any
// live gameplay display (`cmd/hotm-gui` doesn't currently support the
// `-level2grid` exploration mode at all — a real, scoped follow-up),
// but this is the first actual extracted ROOM scene this project has
// (as opposed to the demon/monster/NPC portraits already wired in) —
// direct evidence of what the original game's real corridor rendering
// looked like, for future graphics-fidelity work to check against.
//
//go:embed assets/corridor_sample.png
var corridorSamplePNG []byte

// CorridorSample decodes the embedded real Level2Grid-A1 room
// screenshot. Panics on failure, matching how other embedded-asset-at-
// startup code in this package behaves.
func CorridorSample() image.Image {
	img, _, err := image.Decode(bytes.NewReader(corridorSamplePNG))
	if err != nil {
		panic("graphics: failed to decode embedded corridor_sample.png: " + err.Error())
	}
	return img
}

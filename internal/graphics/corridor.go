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
// level2_grid.go: `w := New(level2Room("A1"))`). Wired into live
// gameplay in round 97 (`cmd/hotm-gui -level2grid`) — this was the
// first actual extracted ROOM scene this project has (as opposed to
// the demon/monster/NPC portraits already wired in), and round 98
// added Level 1's equivalent (Level1CorridorSample, below) — direct
// evidence of what the original game's real corridor rendering looked
// like, for future graphics-fidelity work to check against.
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

// assets/level1_corridor_sample.png is Level 1's own real room screenshot,
// extracted the same way and from the same atlas as corridor_sample.png
// above (round 98) — Level 1's quadrant sits in the atlas's top-left,
// Level 2's in the top-right (confirmed via each quadrant's own printed
// "Level N" heading, not assumed from position alone). Row A's content
// starts at the same y≈484-490 in both quadrants (consistent with a
// shared row-grid layout across the atlas), while column 1 starts at
// x≈571 here vs. x≈5360 for Level 2 — located via the same "find the
// atlas's own printed row/column labels, then tight-crop" method as
// round 96, not pixel-darkness-gap scanning. This is Level 1's own
// confirmed real starting room (level1_grid.go: `w :=
// New(level1Room("A1"))`), the same well-grounded pairing as
// CorridorSample()'s Level 2 cell.
//
//go:embed assets/level1_corridor_sample.png
var level1CorridorSamplePNG []byte

// Level1CorridorSample decodes the embedded real Level1Grid-A1 room
// screenshot. Panics on failure, matching CorridorSample() above.
func Level1CorridorSample() image.Image {
	img, _, err := image.Decode(bytes.NewReader(level1CorridorSamplePNG))
	if err != nil {
		panic("graphics: failed to decode embedded level1_corridor_sample.png: " + err.Error())
	}
	return img
}

// assets/level3_corridor_sample.png is Level 3's own real room
// screenshot (round 103), extracted the same way as Level 1/2's above —
// this atlas's Level 3 quadrant sits bottom-left (confirmed via its own
// printed "Level 3" heading), with column 1 starting at x≈574 (matching
// Level 1's x≈571 almost exactly — the same shared column-grid layout
// carries down to this quadrant too, a good consistency check) and row
// A's content starting at y≈3026. This is Level 3's own confirmed real
// starting room (level3_grid.go: `w := New(level3Room("A1"))`), showing
// a distinctive altar/table-and-cauldron scene, unlike Level 1 and 2's
// plain corridor samples.
//
//go:embed assets/level3_corridor_sample.png
var level3CorridorSamplePNG []byte

// Level3CorridorSample decodes the embedded real Level3Grid-A1 room
// screenshot. Panics on failure, matching CorridorSample() above.
func Level3CorridorSample() image.Image {
	img, _, err := image.Decode(bytes.NewReader(level3CorridorSamplePNG))
	if err != nil {
		panic("graphics: failed to decode embedded level3_corridor_sample.png: " + err.Error())
	}
	return img
}

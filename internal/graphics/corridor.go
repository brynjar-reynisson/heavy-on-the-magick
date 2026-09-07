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

// assets/level4_corridor_sample.png is Level 4's own real room
// screenshot (round 104), extracted the same way as the other 3
// levels' — this atlas's Level 4 quadrant sits bottom-right (confirmed
// via its own printed "Level 4" heading), with column 1 starting at
// x≈5360 and row A at y≈3025 (both match Level 2's calibration almost
// exactly, as expected for the two right-column quadrants).
//
// UNLIKE the other 3 levels, this is NOT cell A1: world.Level4Grid's
// own real starting room is F2 (`w := New(level4Room("F2"))`), not A1
// — see Level4Grid's doc comment: its extraction never confirmed A1 as
// part of the real, playable 17-cell component the way Levels 1-3's A1
// starting rooms were confirmed reachable, so `game.NewLevel4Exploration`
// starts the player at F2 instead. Extracting A1 here would have been
// visually consistent with the other 3 but semantically wrong — it
// would show art for a room the player never actually starts in. F2's
// position was computed from row A/column 1's calibration (row F = +5
// row-heights, column 2 = +1 column-width) then pixel-verified directly
// (not trusted from arithmetic alone) — the initial calculated crop
// included a real black gap between F1 and F2 that isn't part of
// either room, caught and fixed by scanning for where the floor color
// is actually contiguous before finalizing the crop.
//
//go:embed assets/level4_corridor_sample.png
var level4CorridorSamplePNG []byte

// Level4CorridorSample decodes the embedded real Level4Grid-F2 room
// screenshot. Panics on failure, matching CorridorSample() above.
func Level4CorridorSample() image.Image {
	img, _, err := image.Decode(bytes.NewReader(level4CorridorSamplePNG))
	if err != nil {
		panic("graphics: failed to decode embedded level4_corridor_sample.png: " + err.Error())
	}
	return img
}

// assets/room_of_misery_sample.png is a real screenshot of Room of
// Misery — world.CollodonsPile's own real starting room, i.e. the room
// every DEFAULT (unflagged) cmd/hotm-gui session begins in, unlike the
// other 4 samples above which only ever show in a -levelNgrid mode.
// Round 105: extracted from the same atlas's Level 2 quadrant, cell F4
// — level2_grid.go's own doc comment already establishes F4 as Room of
// Misery (tight-crop-confirmed by name against a DIFFERENT source, the
// clean grid map), so this reuses that identification rather than
// re-deriving it, then locates F4 in THIS atlas via row/column offset
// arithmetic from the already-calibrated row-A/column-1 origin (row F
// = +5 rows, column 4 = +3 columns), the same approach round 104 used
// for Level 4's F2.
//
// Cross-checked before trusting it: the cell immediately adjacent
// (one column left, connected by a passage line) is a room showing the
// real "SATOR AREPO TENET OPERA ROTAS" magic word-square inscribed on
// a wall plaque — a strong, independent confirmation this is the right
// column, since level2_grid.go's F3 (immediately before F4) is named
// "Sign", and a wall-inscribed magic square is exactly what "Sign"
// would show. Depicts a robed figure standing between two small
// pedestal tables — consistent with Room of Misery's own already-
// confirmed HasTable: true.
//
//go:embed assets/room_of_misery_sample.png
var roomOfMiserySamplePNG []byte

// RoomOfMiserySample decodes the embedded real Room-of-Misery
// screenshot. Panics on failure, matching CorridorSample() above.
func RoomOfMiserySample() image.Image {
	img, _, err := image.Decode(bytes.NewReader(roomOfMiserySamplePNG))
	if err != nil {
		panic("graphics: failed to decode embedded room_of_misery_sample.png: " + err.Error())
	}
	return img
}

// assets/room_of_stings_sample.png and assets/room_of_arrows_sample.png
// are real screenshots of two more CollodonsPile rooms (round 108) —
// unlike Room of Misery (the world's starting room), these show up
// whenever the player is actually IN that specific room, wherever it
// is on the map, thanks to cmd/hotm-gui's round-108 generalization from
// "art tied to the start room" to "art tied to any room with a known
// screenshot" (see GUI.roomArt's doc comment).
//
// Both are level1_grid.go's own confirmed cell/name pairing (F3 "Room
// of Stings", F5 "Room of Arrows" — the same two cells whose real
// names bridge CollodonsPile and Level1Grid, see level1_grid.go's doc
// comment). Located in the atlas via the same row/column arithmetic as
// every other sample here (row F = +5 rows, columns 3/5 respectively
// from the already-calibrated row-A/column-1 origin), pixel-verified
// rather than trusted from arithmetic alone. Room of Arrows's crop
// shows an actual bow and arrow leaning against a pedestal — strong
// content confirmation the column count was right, the same kind of
// independent check round 105's Sator Square gave for Room of Misery.
//
//go:embed assets/room_of_stings_sample.png
var roomOfStingsSamplePNG []byte

//go:embed assets/room_of_arrows_sample.png
var roomOfArrowsSamplePNG []byte

// RoomOfStingsSample decodes the embedded real Room-of-Stings
// screenshot. Panics on failure, matching CorridorSample() above.
func RoomOfStingsSample() image.Image {
	img, _, err := image.Decode(bytes.NewReader(roomOfStingsSamplePNG))
	if err != nil {
		panic("graphics: failed to decode embedded room_of_stings_sample.png: " + err.Error())
	}
	return img
}

// RoomOfArrowsSample decodes the embedded real Room-of-Arrows
// screenshot. Panics on failure, matching CorridorSample() above.
func RoomOfArrowsSample() image.Image {
	img, _, err := image.Decode(bytes.NewReader(roomOfArrowsSamplePNG))
	if err != nil {
		panic("graphics: failed to decode embedded room_of_arrows_sample.png: " + err.Error())
	}
	return img
}

// assets/wolfdorp_sample.png (round 109) is a real screenshot from
// within CollodonsPile's Wolfdorp room, but with an honestly LOWER
// confidence level than the samples above: Room of Misery/Room of
// Stings/Room of Arrows each had an exact, independently-confirmed
// cell (level2_grid.go's F4, level1_grid.go's F3/F5) established
// BEFORE this project ever tried extracting art for them. Wolfdorp has
// no such single named cell — CollodonsPile's "Wolfdorp" is a zone
// abstraction spanning roughly 18-24 real per-cell rooms (heavymap-
// grid-clean.gif's own colored zone boundary, viewed directly at full
// resolution: the magenta area covering row A columns 1-6, all of rows
// B and C, and part of row D). Cell A2 was picked as a representative,
// unremarkable cell WITHIN that zone (not overlapping any of the
// zone's own specially-labeled sub-cells, and distinct in content from
// Level1Grid's own A1 art — a decorative wall rosette and pedestal,
// not A1's chest) — located in the atlas via the same row/column
// arithmetic as every other sample (row A, column 2), pixel-verified.
// Shown honestly as "representative Wolfdorp-zone art," not a claim
// this is THE definitive Wolfdorp screenshot the way the exact-cell
// samples above are.
//
//go:embed assets/wolfdorp_sample.png
var wolfdorpSamplePNG []byte

// WolfdorpSample decodes the embedded representative Wolfdorp-zone
// screenshot. Panics on failure, matching CorridorSample() above.
func WolfdorpSample() image.Image {
	img, _, err := image.Decode(bytes.NewReader(wolfdorpSamplePNG))
	if err != nil {
		panic("graphics: failed to decode embedded wolfdorp_sample.png: " + err.Error())
	}
	return img
}

// assets/nidus_sample.png (round 113) is a second zone-level-confidence
// sample, same standard as WolfdorpSample above — heavymap-grid-
// clean.gif's own colored zone boundary shows "Nidus" as a green area
// spanning roughly row E-H, columns 6-8, with no single cell singled
// out as "the" Nidus room. Cell F6 was picked as a representative,
// unremarkable cell within that zone (no monster/item marker at this
// position in the clean map), located in the atlas via the same row/
// column arithmetic as every other sample (row F, column 6),
// pixel-verified — shows a distinctive two-archway room with a
// stalagmite formation, visually distinct from every other sample
// already shipped. Shown honestly as "representative Nidus-zone art,"
// not exact-cell precision.
//
//go:embed assets/nidus_sample.png
var nidusSamplePNG []byte

// NidusSample decodes the embedded representative Nidus-zone
// screenshot. Panics on failure, matching CorridorSample() above.
func NidusSample() image.Image {
	img, _, err := image.Decode(bytes.NewReader(nidusSamplePNG))
	if err != nil {
		panic("graphics: failed to decode embedded nidus_sample.png: " + err.Error())
	}
	return img
}

// assets/trollwynd_sample.png (round 116) is a third zone-level-
// confidence sample, same standard as Wolfdorp/Nidus above — this time
// from LEVEL 3, not Level 1. heavymap-grid-clean.gif's own colored zone
// boundary shows "Trollwynd" as a green area spanning roughly rows B-D,
// columns 4-8 (matching zone_monsters.go's "Trollwynd: Troll x4"
// sighting, already cross-confirmed in round 71 against this same
// zone's 4 tight-crop-verified Troll icons). Cell B5 was picked as a
// representative, unmarked cell within that zone, located in the atlas
// via the same row/column arithmetic as every other Level 3 sample
// (row B, column 5, reusing Level3CorridorSample's row-A/column-1
// origin) — pixel-verified (left/top edges matched the calibration's
// prediction closely, a good consistency check), showing a doorway and
// a distinctive round shield/disc object on the floor.
//
//go:embed assets/trollwynd_sample.png
var trollwyndSamplePNG []byte

// TrollwyndSample decodes the embedded representative Trollwynd-zone
// screenshot. Panics on failure, matching CorridorSample() above.
func TrollwyndSample() image.Image {
	img, _, err := image.Decode(bytes.NewReader(trollwyndSamplePNG))
	if err != nil {
		panic("graphics: failed to decode embedded trollwynd_sample.png: " + err.Error())
	}
	return img
}

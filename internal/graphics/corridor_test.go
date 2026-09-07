package graphics

import (
	"image"
	"testing"
)

// TestCorridorSampleDecodesToRealArt pins the embedded real room
// screenshot (see corridor.go's doc comment for sourcing) actually
// decoding to plausible, non-trivial art — not a corrupt or placeholder
// image.
func TestCorridorSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, CorridorSample(), "CorridorSample()")
}

// TestLevel1CorridorSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 98's Level 1 equivalent.
func TestLevel1CorridorSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, Level1CorridorSample(), "Level1CorridorSample()")
}

// TestLevel3CorridorSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 103's Level 3 equivalent.
func TestLevel3CorridorSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, Level3CorridorSample(), "Level3CorridorSample()")
}

// TestLevel4CorridorSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 104's Level 4 equivalent
// (cell F2, not A1 - see Level4CorridorSample's doc comment for why).
func TestLevel4CorridorSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, Level4CorridorSample(), "Level4CorridorSample()")
}

// TestRoomOfMiserySampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 105's Room of Misery
// sample - the first of these 5 real screenshots that shows in
// cmd/hotm-gui's DEFAULT (unflagged, CollodonsPile) mode.
func TestRoomOfMiserySampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, RoomOfMiserySample(), "RoomOfMiserySample()")
}

// TestRoomOfStingsSampleDecodesToRealArt and
// TestRoomOfArrowsSampleDecodesToRealArt mirror
// TestCorridorSampleDecodesToRealArt for round 108's two additions.
func TestRoomOfStingsSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, RoomOfStingsSample(), "RoomOfStingsSample()")
}

func TestRoomOfArrowsSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, RoomOfArrowsSample(), "RoomOfArrowsSample()")
}

// TestAgileStairSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 137's Agile Stair
// addition - exact-cell confidence (level1_grid.go's own A7), the same
// tier as Room of Misery/Stings/Arrows above.
func TestAgileStairSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, AgileStairSample(), "AgileStairSample()")
}

// TestSothicComplexSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 144's Sothic Complex
// addition - zone-level confidence, same tier as WolfdorpSample above.
func TestSothicComplexSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, SothicComplexSample(), "SothicComplexSample()")
}

// TestMorfangSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 143's Morfang addition -
// zone-level confidence, same tier as WolfdorpSample above.
func TestMorfangSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, MorfangSample(), "MorfangSample()")
}

// TestMethosSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 142's Methos addition -
// zone-level confidence, same tier as WolfdorpSample above.
func TestMethosSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, MethosSample(), "MethosSample()")
}

// TestFurnaceRoomSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 141's Furnace Room
// addition - exact-cell confidence (level1_grid.go's own A8).
func TestFurnaceRoomSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, FurnaceRoomSample(), "FurnaceRoomSample()")
}

// TestWolfdorpSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 109's Wolfdorp-zone
// sample.
func TestWolfdorpSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, WolfdorpSample(), "WolfdorpSample()")
}

// TestNidusSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 113's Nidus-zone sample.
func TestNidusSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, NidusSample(), "NidusSample()")
}

// TestTrollwyndSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 116's Trollwynd-zone
// sample.
func TestTrollwyndSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, TrollwyndSample(), "TrollwyndSample()")
}

// TestPilefootSampleDecodesToRealArt mirrors
// TestCorridorSampleDecodesToRealArt for round 123's Pilefoot-zone
// sample.
func TestPilefootSampleDecodesToRealArt(t *testing.T) {
	assertRealArt(t, PilefootSample(), "PilefootSample()")
}

func assertRealArt(t *testing.T, img image.Image, label string) {
	t.Helper()
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w < 100 || h < 50 {
		t.Errorf("%s size = %dx%d, want a real room-scene-sized image (at least 100x50)", label, w, h)
	}

	first := img.At(bounds.Min.X, bounds.Min.Y)
	allSame := true
	for y := bounds.Min.Y; y < bounds.Max.Y && allSame; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.At(x, y) != first {
				allSame = false
				break
			}
		}
	}
	if allSame {
		t.Errorf("%s is a single flat color - want real varied room art", label)
	}
}

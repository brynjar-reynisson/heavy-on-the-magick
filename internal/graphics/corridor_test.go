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

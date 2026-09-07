package graphics

import "testing"

// TestCorridorSampleDecodesToRealArt pins the embedded real room
// screenshot (see corridor.go's doc comment for sourcing) actually
// decoding to plausible, non-trivial art — not a corrupt or placeholder
// image.
func TestCorridorSampleDecodesToRealArt(t *testing.T) {
	img := CorridorSample()
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w < 100 || h < 50 {
		t.Errorf("CorridorSample() size = %dx%d, want a real room-scene-sized image (at least 100x50)", w, h)
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
		t.Error("CorridorSample() is a single flat color - want real varied room art")
	}
}

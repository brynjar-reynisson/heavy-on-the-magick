package graphics

import "testing"

// TestApexPortraitDecodesToRealArt pins the embedded real portrait asset
// (see portraits.go's doc comment for sourcing) actually decoding to
// plausible, non-trivial art — not a corrupt or placeholder 1x1 image.
func TestApexPortraitDecodesToRealArt(t *testing.T) {
	img := ApexPortrait()
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w < 50 || h < 50 {
		t.Errorf("ApexPortrait() size = %dx%d, want a real portrait-sized image (at least 50x50)", w, h)
	}

	// The real portrait is monochrome black-and-white line art (see the
	// source screenshot) - confirm it isn't accidentally blank (all one
	// color), which would indicate a bad crop or corrupt asset.
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
		t.Error("ApexPortrait() is a single flat color - want real varied portrait art")
	}
}

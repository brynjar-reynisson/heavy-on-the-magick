package graphics

import "testing"

// TestPortraitDecodesToRealArt pins every embedded real portrait asset
// (see portraits.go's doc comment for sourcing) actually decoding to
// plausible, non-trivial art — not a corrupt or placeholder 1x1 image.
func TestPortraitDecodesToRealArt(t *testing.T) {
	for _, name := range PortraitNames {
		img := Portrait(name)
		bounds := img.Bounds()
		w, h := bounds.Dx(), bounds.Dy()
		if w < 30 || h < 30 {
			t.Errorf("Portrait(%q) size = %dx%d, want a real portrait-sized image (at least 30x30)", name, w, h)
		}

		// The real portraits are monochrome black-and-white line art (see
		// the source screenshots) - confirm each isn't accidentally blank
		// (all one color), which would indicate a bad crop or corrupt asset.
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
			t.Errorf("Portrait(%q) is a single flat color - want real varied portrait art", name)
		}
	}
}

// TestApexPortraitMatchesPortrait pins that the pre-existing convenience
// wrapper (kept for cmd/hotm-gui's one call site) is exactly equivalent
// to the generalized Portrait("apex") added in round 76.
func TestApexPortraitMatchesPortrait(t *testing.T) {
	want := Portrait("apex").Bounds()
	got := ApexPortrait().Bounds()
	if got != want {
		t.Errorf("ApexPortrait().Bounds() = %v, want %v (same as Portrait(\"apex\"))", got, want)
	}
}

// TestPortraitUnknownNamePanics pins the documented panic-on-unknown-name
// behavior rather than a silent nil/zero-value return.
func TestPortraitUnknownNamePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Portrait(\"nosuchcreature\") did not panic, want it to")
		}
	}()
	Portrait("nosuchcreature")
}

package graphics

import (
	"bytes"
	"image/png"
	"testing"
)

var _ Renderer = (*PNGRenderer)(nil) // compile-time interface check

func TestPNGRendererProducesValidPNG(t *testing.T) {
	r := NewPNGRenderer(4, 4, 4)
	r.Clear(Black)
	r.DrawGlyph(RuneGlyphs[0], 1, 1, White, false)

	var buf bytes.Buffer
	if err := r.WritePNG(&buf); err != nil {
		t.Fatalf("WritePNG: %v", err)
	}

	img, err := png.Decode(&buf)
	if err != nil {
		t.Fatalf("output is not a valid PNG: %v", err)
	}
	bounds := img.Bounds()
	wantW, wantH := 4*8*4, 4*8*4
	if bounds.Dx() != wantW || bounds.Dy() != wantH {
		t.Errorf("image size = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), wantW, wantH)
	}
}

func TestPNGRendererDrawsSetPixels(t *testing.T) {
	r := NewPNGRenderer(1, 1, 1)
	r.Clear(Black)
	// A glyph with the top-left pixel set.
	r.DrawGlyph(Glyph{0x80, 0, 0, 0, 0, 0, 0, 0}, 0, 0, White, false)

	got := r.img.RGBAAt(0, 0)
	want := palette[White]
	if got != want {
		t.Errorf("pixel (0,0) = %+v, want %+v", got, want)
	}
	// A pixel that should still be background (Black).
	got2 := r.img.RGBAAt(7, 7)
	want2 := palette[Black]
	if got2 != want2 {
		t.Errorf("pixel (7,7) = %+v, want %+v (unset bit)", got2, want2)
	}
}

// TestPNGRendererBrightUsesBrightPalette pins the real ZX Spectrum ULA
// BRIGHT attribute distinction (round 88) - a bright draw must use
// brightPalette's confirmed higher-intensity RGB, not palette's normal
// values, and must differ from the same color drawn non-bright.
func TestPNGRendererBrightUsesBrightPalette(t *testing.T) {
	r := NewPNGRenderer(2, 1, 1)
	r.Clear(Black)
	full := Glyph{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	r.DrawGlyph(full, 0, 0, Red, false)
	r.DrawGlyph(full, 1, 0, Red, true)

	normal := r.img.RGBAAt(0, 0)
	bright := r.img.RGBAAt(8, 0)
	if normal != palette[Red] {
		t.Errorf("non-bright Red pixel = %+v, want palette[Red] %+v", normal, palette[Red])
	}
	if bright != brightPalette[Red] {
		t.Errorf("bright Red pixel = %+v, want brightPalette[Red] %+v", bright, brightPalette[Red])
	}
	if normal == bright {
		t.Error("bright and non-bright Red rendered identically, want the confirmed higher-intensity RGB to differ")
	}
}

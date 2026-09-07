// Command render-glyphs draws every confirmed graphics asset extracted so
// far (see internal/graphics/font.go) to a PNG, as a visual sanity check
// of the Go rendering pipeline against what the disassembly work found.
package main

import (
	"flag"
	"os"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/graphics"
)

func main() {
	out := flag.String("out", "glyphs.png", "output PNG path")
	cellSize := flag.Int("cell-size", 16, "pixels per native ZX Spectrum pixel")
	flag.Parse()

	cols := len(graphics.RuneGlyphs) + 1 // +1 for the known-icon column
	r := graphics.NewPNGRenderer(cols, 1, *cellSize)
	r.Clear(graphics.Black)

	for i, g := range graphics.RuneGlyphs {
		r.DrawGlyph(g, i, 0, graphics.White, false)
	}
	// Real confirmed ULA attribute 0x43 = ink magenta BRIGHT (see
	// KnownIcons's doc comment) - not a cosmetic choice.
	r.DrawGlyph(graphics.KnownIcons["magenta-icon"], len(graphics.RuneGlyphs), 0, graphics.Magenta, true)

	f, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := r.WritePNG(f); err != nil {
		panic(err)
	}
}

package graphics

import (
	"image"
	"image/color"
	"image/png"
	"io"
)

// palette maps the ZX Spectrum's 8 colors (non-bright) to real RGB.
var palette = map[Color]color.RGBA{
	Black:   {0, 0, 0, 255},
	Blue:    {0, 0, 214, 255},
	Red:     {214, 0, 0, 255},
	Magenta: {214, 0, 214, 255},
	Green:   {0, 214, 0, 255},
	Cyan:    {0, 214, 214, 255},
	Yellow:  {214, 214, 0, 255},
	White:   {214, 214, 214, 255},
}

// brightPalette maps the same 8 colors to their real ZX Spectrum ULA
// BRIGHT-attribute RGB values — 255 per "on" channel instead of 214, the
// standard real hardware bright-intensity value (matching this project's
// own already-confirmed exact bright RGB facts, e.g. Ghost's bright green
// (0,255,0) and Wraith/Medusa's bright red (255,0,0) from the clean grid
// map's legend — see CLAUDE.md). Black has no bright variant on real
// hardware (BRIGHT only affects "on" bits), so it's identical to palette.
var brightPalette = map[Color]color.RGBA{
	Black:   {0, 0, 0, 255},
	Blue:    {0, 0, 255, 255},
	Red:     {255, 0, 0, 255},
	Magenta: {255, 0, 255, 255},
	Green:   {0, 255, 0, 255},
	Cyan:    {0, 255, 255, 255},
	Yellow:  {255, 255, 0, 255},
	White:   {255, 255, 255, 255},
}

// PNGRenderer implements Renderer by drawing into an in-memory image,
// exportable as a real PNG file. This is the first actual Renderer
// implementation (see screen.go's interface) — no live-display backend
// (e.g. ebiten) is wired up yet, but this proves the drawing model works
// and is useful on its own for rendering confirmed assets to check by eye
// (see cmd/render-glyphs).
type PNGRenderer struct {
	img *image.RGBA
	// CellSize is the on-screen pixel size of one ZX Spectrum bitmap pixel
	// (the original is 256x192 native pixels; CellSize lets output be
	// upscaled for visibility, matching how we rendered glyphs during the
	// disassembly work — see ../../CLAUDE.md's render_tile.py). Round 129:
	// a real 1986 CRASH magazine review confirms this upscaling isn't
	// just a convenient debugging choice — the original itself "is formed
	// in memory and blown up onto the screen... with the result that
	// individual pixels become conspicuous" — this project's own
	// CellSize>1 usage (cmd/hotm-gui's HUD renders at CellSize=8) is
	// therefore confirmed faithful to the real technique, not merely a
	// happy coincidence.
	CellSize int
}

// NewPNGRenderer creates a renderer for a canvas widthCols x heightRows
// 8x8 character cells, matching the original's 32x24-cell (256x192-pixel)
// screen when widthCols=32, heightRows=24.
func NewPNGRenderer(widthCols, heightRows, cellSize int) *PNGRenderer {
	if cellSize < 1 {
		cellSize = 1
	}
	w := widthCols * 8 * cellSize
	h := heightRows * 8 * cellSize
	return &PNGRenderer{
		img:      image.NewRGBA(image.Rect(0, 0, w, h)),
		CellSize: cellSize,
	}
}

func (r *PNGRenderer) Clear(bg Color) {
	c := palette[bg]
	bounds := r.img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r.img.Set(x, y, c)
		}
	}
}

// DrawGlyph draws an 8x8 Glyph at character-cell (col, row). bright
// selects the real ZX Spectrum ULA's BRIGHT attribute bit (round 88 -
// see screen.go's Renderer.DrawGlyph doc comment and brightPalette).
func (r *PNGRenderer) DrawGlyph(g Glyph, col, row int, fg Color, bright bool) {
	c := palette[fg]
	if bright {
		c = brightPalette[fg]
	}
	baseX := col * 8 * r.CellSize
	baseY := row * 8 * r.CellSize
	for py := range 8 {
		rowBits := g[py]
		for px := range 8 {
			if rowBits&(0x80>>px) == 0 {
				continue // unset pixel: leave background showing through
			}
			for sy := range r.CellSize {
				for sx := range r.CellSize {
					r.img.Set(baseX+px*r.CellSize+sx, baseY+py*r.CellSize+sy, c)
				}
			}
		}
	}
}

func (r *PNGRenderer) Present() {
	// No-op for a static image renderer: there's no live display to flip.
	// Call WritePNG to actually get the result out.
}

// Image exposes the underlying canvas directly (as a standard
// image.Image), so other frontends — e.g. cmd/hotm-gui's live ebiten
// display — can reuse this renderer's drawing logic instead of
// duplicating it, converting via their own toolkit's image-from-image
// constructor.
func (r *PNGRenderer) Image() image.Image {
	return r.img
}

// WritePNG encodes the current canvas as a PNG.
func (r *PNGRenderer) WritePNG(w io.Writer) error {
	return png.Encode(w, r.img)
}

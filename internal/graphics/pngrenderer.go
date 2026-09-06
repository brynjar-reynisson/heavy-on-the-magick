package graphics

import (
	"image"
	"image/color"
	"image/png"
	"io"
)

// palette maps the ZX Spectrum's 8 colors (non-bright) to real RGB. Real
// Spectrum hardware also has a "bright" variant of each color (roughly
// doubling the RGB values); PNGRenderer doesn't distinguish bright/normal
// yet since nothing built on it has needed that distinction so far — see
// the TODO on DrawGlyph.
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
	// disassembly work — see ../../CLAUDE.md's render_tile.py).
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

// DrawGlyph draws an 8x8 Glyph at character-cell (col, row).
//
// TODO: fg is currently always treated as non-bright (see palette); the
// original ZX Spectrum ULA also has a BRIGHT attribute bit per cell (used
// throughout the game — e.g. the confirmed magenta/green spell-icon
// squares in CLAUDE.md are both "bright" variants). Once bright rendering
// is needed, extend Color or add a separate bright flag here rather than
// guessing at RGB values not yet cross-checked against a real screenshot.
func (r *PNGRenderer) DrawGlyph(g Glyph, col, row int, fg Color) {
	c := palette[fg]
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

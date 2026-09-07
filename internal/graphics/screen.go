package graphics

// Renderer draws the game's visuals. No implementation exists yet — the
// plan is an ebiten-backed one (see ../../CLAUDE.md's "port" discussion:
// a Go port doesn't need to replicate the ZX Spectrum's bit-packed/
// attribute-clash screen format, just draw the same pictures with normal
// 2D primitives), but the interface is defined now so game logic can be
// written against it before the renderer exists.
type Renderer interface {
	// DrawGlyph draws an 8x8 Glyph at the given cell position (col, row in
	// 8x8-character-cell units, matching the original's coordinate system)
	// in the given foreground color, useful for spell sigils and text.
	// bright selects the ZX Spectrum ULA's BRIGHT attribute bit (roughly
	// doubling ink intensity) - a real, confirmed-in-hardware distinction
	// (see KnownIcons's "magenta-icon" doc comment: its real attribute
	// byte is 0x43, ink=magenta BRIGHT), not a cosmetic option.
	DrawGlyph(g Glyph, col, row int, fg Color, bright bool)

	// Clear wipes the screen to a background color.
	Clear(bg Color)

	// Present flips/flushes the frame to the display.
	Present()
}

// Color mirrors the ZX Spectrum's 8-color (×2 brightness) palette, since
// the game's original art was designed within that constraint. A modern
// renderer can map these to real RGB values however it likes.
type Color int

const (
	Black Color = iota
	Blue
	Red
	Magenta
	Green
	Cyan
	Yellow
	White
)

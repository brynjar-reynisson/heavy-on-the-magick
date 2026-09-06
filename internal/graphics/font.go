// Package graphics holds the game's visual assets as they get extracted
// from the disassembly, plus the future rendering layer (a real renderer —
// likely ebiten — hasn't been wired up yet; see screen.go).
package graphics

// Glyph is one 8x8, 1-bit-per-pixel character cell — the ZX Spectrum's
// native character/UDG format. Row 0 is the top pixel row; within a row,
// bit 7 is the leftmost pixel.
type Glyph [8]byte

// RuneGlyphs are confirmed real, correctly-decoded magic-symbol glyphs from
// the game's picture table at Z80 address 48054 (see ../../CLAUDE.md,
// "Graphics table at 48054 — status"). These render as clean, coherent
// rune/sigil shapes and are almost certainly the custom font used for
// spellcasting sigils — but WHICH symbol means what, and their in-game
// index/ordering, has not been confirmed yet (see 31932's "1 of 16 runes"
// selector in CLAUDE.md — these are 3 of that font, not necessarily in the
// right order).
//
// Byte source: table entries 26 and 27 (ptr=51944, ptr=51952), read
// directly from decompressed RAM page 5, count=1 row_count=1 (single
// 8-byte glyph, no picture-transform ambiguity — unlike the still-uncracked
// 120-byte "big picture" entries, these small ones were confirmed to
// render correctly as plain untransformed bytes).
var RuneGlyphs = []Glyph{
	{0x00, 0x00, 0x42, 0x5a, 0x52, 0x42, 0x00, 0x00}, // table entry 26
	{0x00, 0x02, 0x5a, 0x5e, 0x7e, 0x5a, 0x42, 0x00}, // table entry 27
}

// KnownIcons are on-screen sprites confirmed by reading live emulator
// memory directly (see ../../CLAUDE.md, "Got real ground-truth pixel data
// for the magenta icon"), but NOT found in the picture table at 48054 —
// their actual source location/format is still unknown. Recorded here as
// ground truth to check any future graphics-format theory against.
var KnownIcons = map[string]Glyph{
	// The magenta spell-selection icon square in the ritual-circle screen,
	// read from live bitmap memory at 18540 (+256*row), attribute 22892
	// (row 11, col 12, ink=magenta bright, paper=black).
	"magenta-icon": {0xf8, 0x88, 0xa8, 0xa9, 0x88, 0xf8, 0x91, 0x00},
}

// Package graphics holds the game's real visual assets and the renderer
// that draws them — PNGRenderer (screen.go/pngrenderer.go) has been the
// one real, working Renderer implementation for a long time; both the
// offline exporter (cmd/render-glyphs) and the live GUI (cmd/hotm-gui,
// via ebiten.NewImageFromImage) share the exact same drawing code rather
// than each having their own. This package's assets have grown well past
// the original confirmed rune-glyph font: 13 real extracted portraits
// (Portrait/PortraitNames — Apex, all 4 demons, all 8 monster types, from
// the real in-game screenshot atlas heavymap-speccy-screenshots.png) and
// a growing set of real per-room corridor screenshots (the *Sample
// functions — CorridorSample, Level1CorridorSample, RoomOfMiserySample,
// etc.), all live in cmd/hotm-gui's default and level-grid exploration
// modes. See ../../CLAUDE.md's round 99/108 notes for this package's own
// history of stale-doc-comment corrections — worth checking again
// whenever a doc comment here starts sounding out of date, since this
// package has grown substantially since most of its own comments were
// first written.
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

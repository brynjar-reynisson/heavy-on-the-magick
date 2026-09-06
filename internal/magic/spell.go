// Package magic models Axil's spellcasting system: Crowley/Golden-Dawn-
// themed "invocations" using a set of rune sigils.
//
// Confirmed from decompressed RAM string data (see ../../CLAUDE.md, RAM
// page 5 dump): vocabulary including INVOKE, OBJECT, SPELL, SPITTED,
// POISON, NIGHTSHADE, BLAST, MAGICK, and demon names such as BELEZBAR
// (= Beelzebub), ASMODEE, ASTAROT. The rune-selection mechanism (pick 1 of
// 16 glyphs, see graphics.RuneGlyphs and CLAUDE.md's trace of routine
// 31932) is understood structurally, but which rune means what, and how
// runes combine into an actual spell, is not yet known.
package magic

// Rune identifies one of the 16 selectable sigils (see routine 31932's
// `E & 15` glyph-index masking in the disassembly). Ordering/meaning is
// unconfirmed — these are placeholder names until the real ones are found.
type Rune int

const MaxRunes = 16

// Spell is a not-yet-understood combination of runes and/or a target. The
// real invocation grammar (how BLAST differs from a targeted spell, what
// role Grade plays in what's castable) hasn't been extracted yet.
type Spell struct {
	Runes []Rune
}

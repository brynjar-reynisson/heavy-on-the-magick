// Package magic models Axil's occult mechanics: the 4 confirmed demons
// (Demon, Demons — see demons.go) and their real, invocable abilities,
// plus the 12 confirmed zodiac sign-to-key pairings (ZodiacKey,
// ZodiacKeys — see zodiac_keys.go). All 4 demons now have real,
// functional conversation-form commands in internal/game (Astarot
// teleports, Magot locates, Asmodee destroys, Belezbar reveals — see
// game.astarotTeleport/magotLocate/asmodeeDestroy/belezbarReveal and
// ../../CLAUDE.md's rounds 75-164 for the full sourcing history).
//
// ROUND 168: this file used to describe a DIFFERENT, earlier avenue -
// a "rune sigil" spellcasting system (pick 1 of 16 glyphs, see
// graphics.RuneGlyphs and CLAUDE.md's trace of routine 31932) - as the
// package's main subject, with Demons/ZodiacKeys barely mentioned.
// That framing had gone stale: the rune-combination mechanism was
// never cracked (which rune means what, and how they'd combine into a
// spell, was never confirmed), and the game's REAL, functional spell
// system (BLAST/FREEZE/TRANSFUSION/CALL/INVOKE) turned out to be
// sourced entirely differently - from a published walkthrough, not
// from decoding these glyphs - and lives in internal/game, not here.
// The old placeholder Rune/Spell types (never consumed anywhere in
// this codebase, confirmed via a full grep before removal) are gone;
// the one real, sourced structural fact they existed to record - a
// confirmed 16-glyph selection space, from routine 31932's own `E &
// 15` masking - is preserved as a doc comment on graphics.RuneGlyphs
// instead, the package that actually holds real extracted glyph data.
package magic

package audio

// PitchTable is the game's real chromatic note-to-tone-period table,
// extracted byte-for-byte from Z80 memory address 64812 (see
// ../../CLAUDE.md, "Found the sound routine and a confirmed pitch table").
// Each entry is a delay-loop period value consumed by the beeper-toggle
// routine at Z80 address 64671/64733 (found via `OUT (254)` — the ZX
// Spectrum's only sound output, a single-bit speaker toggle).
//
// Confirmed genuinely musical: consecutive entries differ by an average
// ratio of 1.0606, essentially identical to the true 12-tone-equal-
// temperament semitone ratio 2^(1/12) = 1.0595 (within 0.1%). 53 real
// note entries spanning ~4.4 octaves, followed by a terminator/sentinel
// byte (value 1, not a real note — found by index 53's sharp discontinuity
// from the smooth geometric sequence).
//
// The routine that reads this (Z80 `64649`) does `noteIndex + 12` before
// indexing. Round 90 confirmed WHY: `noteIndex` is a SIGNED byte, not
// unsigned - see NoteIndex. The +12 gives headroom for negative offsets
// (down to -12) while still landing inside this table's unsigned 0-52
// index range, using a single shared table for a wider effective pitch
// range than 53 straight ascending entries alone would allow.
var PitchTable = []byte{
	255, 240, 227, 215, 203, 192, 180, 171, 161, 151,
	144, 136, 128, 121, 114, 108, 102, 96, 91, 86,
	81, 76, 72, 68, 64, 61, 57, 54, 51, 48,
	45, 43, 40, 38, 36, 34, 32, 30, 28, 27,
	25, 24, 23, 21, 20, 19, 18, 17, 16, 15,
	14, 13, 12,
}

// PitchTableTerminator is the sentinel value following PitchTable in the
// original data — not a playable note.
const PitchTableTerminator = 1

// StartupMelody is a raw note stream found immediately after PitchTable
// in memory (64812+54 onward), not yet confirmed to be THE startup
// jingle specifically, but positioned exactly where a melody would
// follow a pitch table, and shaped like one: runs of the same value 2-4
// times in a row (a held note, played across several timer ticks).
//
// CORRECTED (round 90): 254 and 252 were previously treated as special
// rest/silence markers (RestMarkerA/RestMarkerB), guessed because they
// fall outside PitchTable's unsigned 0-52 range. That guess was wrong.
// Directly parsing a .z80 memory snapshot (see SecondaryMelody below)
// found a second, clearly real, musically coherent note stream that
// uses the exact same encoding and contains several values in this same
// "looks out of range" territory (244/245/246/247) forming an obviously
// intentional, coherent melodic phrase once read as NoteIndex expects
// (SIGNED bytes) - not silence. Applying the same rule here: 254 (signed
// -2) and 252 (signed -4) decode to real, valid PitchTable indices (10
// and 8), not silence. See NoteIndex - RenderNotes now uses it for
// every entry, with no special-casing.
var StartupMelody = []byte{
	26, 41, 27, 41, 29, 41, 31, 41, 32, 41, 31, 41, 27, 41, 31, 41,
	26, 26, 27, 27, 29, 29, 31, 31, 32, 32, 31, 31, 27, 27, 31, 31,
	0, 0, 0, 0, 12, 12, 12, 12, 14, 14, 14, 14, 12, 12, 0, 0,
	10, 10, 10, 10, 12, 12, 12, 12, 10, 10, 10, 10, 254, 254, 10, 10,
	8, 8, 8, 8, 252, 252, 252, 252, 252, 252, 252, 252, 252, 254, 0, 3,
	2, 2, 2, 2, 5, 5, 5, 5, 7, 7, 7, 7, 2, 2, 2, 2,
	0, 0, 12, 0, 0, 0, 12, 0, 0, 0, 12, 0, 0, 0, 12, 14,
	15, 14, 15, 14, 15, 14, 12, 14, 15, 14, 15, 14, 17, 17, 17, 17,
	17, 17, 17, 17, 17, 19, 20, 17, 22, 20, 19, 20,
}

// SecondaryMelody is a SECOND, independent real note stream, found and
// confirmed this round (90) via disassembly + a direct .z80 memory
// snapshot read (not live emulator access - snapshot parsing was
// validated by re-deriving PitchTable and StartupMelody's own bytes
// from the same file and confirming an exact match against the
// already-shipped data above before trusting anything new from it).
//
// Source: the sound routine at Z80 64671 reads TWO independent,
// separately-advancing note streams per call, each via its own 2-byte
// pointer variable (64627 and 64631) chained forward one byte at a time
// by the shared reader routine at 64636/64663 - StartupMelody is the
// stream reached through the first pointer (initialized to 64865,
// i.e. one byte before this project's already-confirmed 64866 start);
// this is the stream reached through the SECOND pointer, initialized to
// 65155. Extracted from 65156 (the first byte after that pointer's own
// +1 initial advance) up to the first occurrence of the confirmed
// chain/loop marker value 64 (0x40, at address 65444) - the same
// terminator value routine 64636 checks for explicitly, so this is a
// real, well-defined stream boundary, not an arbitrary cutoff.
//
// Round 93 gave this a standalone live keybinding (cmd/hotm-gui's B
// key); round 99 went further and wired it into ACTUAL gameplay
// playback too — since the source routine reads both streams together
// on every call (see above), the GUI's startup sound now mixes this
// with StartupMelody (see audio.MixNotes) rather than requiring a
// separate manual press to ever hear it during real play.
//
// Round 111 actually TRACED what the two-stream combining produces on
// real hardware (see synth.go's tStatesPerPeriodUnit doc comment for
// the full cycle-count derivation): it's real bit-level XOR
// interleaving of two independently-clocked toggle counters on ONE
// shared speaker output bit (E's period from StartupMelody's current
// note, L's from SecondaryMelody's) — not simple alternation and not
// true multi-channel mixing, closer to a beat-frequency/interference
// pattern. MixNotes's sample-averaging remains this port's own
// simplification of that (reproducing the exact bit-interleave is a
// separate, not-yet-attempted task), but it's now a documented
// simplification of a KNOWN real mechanism, not a guess at an unknown
// one.
var SecondaryMelody = []byte{
	26, 24, 27, 24, 29, 24, 31, 24, 32, 24, 31, 24, 27, 24, 31, 24,
	26, 24, 27, 24, 29, 24, 31, 24, 32, 24, 31, 24, 27, 24, 31, 24,
	26, 24, 27, 24, 29, 24, 31, 24, 32, 24, 31, 24, 27, 24, 31, 24,
	26, 24, 27, 24, 29, 24, 31, 24, 32, 24, 31, 24, 27, 24, 31, 24,
	26, 24, 27, 24, 29, 24, 31, 24, 32, 24, 31, 24, 27, 24, 31, 24,
	26, 24, 27, 24, 29, 24, 31, 24, 32, 24, 31, 24, 27, 24, 31, 24,
	26, 24, 27, 24, 29, 24, 31, 24, 32, 24, 31, 24, 27, 24, 31, 24,
	26, 24, 27, 24, 29, 24, 31, 24, 32, 24, 31, 24, 27, 24, 31, 29,
	31, 29, 32, 29, 34, 29, 36, 29, 37, 29, 36, 29, 32, 29, 36, 29,
	31, 29, 32, 29, 34, 29, 36, 29, 37, 29, 36, 29, 32, 29, 36, 29,
	26, 24, 27, 24, 29, 24, 31, 24, 32, 24, 31, 24, 27, 24, 31, 24,
	26, 24, 27, 24, 29, 24, 31, 24, 32, 24, 31, 24, 27, 24, 31, 41,
	12, 14, 12, 14, 15, 14, 12, 14, 0, 2, 0, 2, 3, 2, 0, 2,
	10, 12, 10, 12, 14, 12, 10, 14, 12, 10, 12, 10, 12, 10, 9, 5,
	36, 24, 36, 24, 36, 24, 36, 24, 244, 19, 245, 24, 246, 19, 247, 24,
	0, 3, 0, 4, 0, 5, 0, 6, 0, 3, 0, 4, 12, 11, 10, 9,
	5, 7, 8, 5, 14, 14, 12, 12, 5, 7, 8, 5, 14, 14, 12, 12,
	19, 17, 19, 17, 19, 17, 15, 14, 12, 15, 19, 24, 27, 31, 36, 36,
}

// NoteIndex converts a raw note-stream byte (from StartupMelody or
// SecondaryMelody) into the real PitchTable index the game's own
// routine (Z80 64649) computes: the byte is interpreted as SIGNED, then
// offset by +12 (see PitchTable's doc comment). Confirmed (round 90) by
// checking this decodes SecondaryMelody's otherwise-out-of-range-looking
// values (244/245/246/247, i.e. signed -12/-11/-10/-9) into a musically
// coherent phrase — the same rule then correctly re-explains
// StartupMelody's 254/252 as real notes too, not silence. A result of
// 53 lands exactly on PitchTableTerminator's position (not a playable
// note); RenderNotes treats any index outside 0-52 as silence.
func NoteIndex(raw byte) int {
	return int(int8(raw)) + 12
}

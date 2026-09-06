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
// indexing — i.e. callers pass a note index that's offset by one octave
// (12 semitones) from this table's own start, though the exact reason
// (headroom for a transpose/octave parameter elsewhere) isn't confirmed.
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

// StartupMelody is a raw note-index stream found immediately after
// PitchTable in memory (64812+54 onward), not yet confirmed to be THE
// startup jingle specifically, but positioned exactly where a melody
// would follow a pitch table, and shaped like one: runs of the same
// value 2-4 times in a row (a held note, played across several timer
// ticks) mixed with two special values, 254 and 252, that fall well
// outside the valid 0-52 note range and are presumed to be
// rest/silence or duration-control markers — their exact meaning is NOT
// confirmed (the routine that reads this stream during real playback
// hasn't been traced yet).
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

// RestValues are StartupMelody entries outside the valid note range —
// treat as silence, not an out-of-bounds PitchTable lookup.
const (
	RestMarkerA = 254
	RestMarkerB = 252
)

// IsRest reports whether a StartupMelody entry is a silence marker rather
// than a playable note index.
func IsRest(noteIndex byte) bool {
	return noteIndex == RestMarkerA || noteIndex == RestMarkerB
}

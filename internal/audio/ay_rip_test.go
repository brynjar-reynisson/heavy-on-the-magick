package audio

import (
	"bytes"
	"os"
	"testing"
)

// TestPitchTableStartupMelodySecondaryMelodyMatchRealAYRip is round
// 157's real, permanent, automated proof of an independent third-party
// confirmation - not just a one-off manual check. testdata/
// HeavyOnTheMagick.ay is a genuine AY music rip of this exact game
// (ripped by Pawel Ochman, 8 Oct 2001, distributed by World of
// Spectrum/spectrumcomputing.co.uk - see PitchTable's doc comment for
// the full sourcing), built from a completely separate extraction
// effort than this project's own from-scratch disassembly. Confirms
// PitchTable, the corrected 288-byte StartupMelody, and SecondaryMelody
// all appear byte-for-byte, in that exact order, inside the real rip's
// own embedded song data - the strongest available cross-validation for
// this project's whole sound reconstruction.
func TestPitchTableStartupMelodySecondaryMelodyMatchRealAYRip(t *testing.T) {
	rip, err := os.ReadFile("testdata/HeavyOnTheMagick.ay")
	if err != nil {
		t.Fatalf("reading testdata/HeavyOnTheMagick.ay: %v", err)
	}

	pitchIdx := bytes.Index(rip, PitchTable)
	if pitchIdx < 0 {
		t.Fatal("PitchTable not found byte-for-byte in the real AY rip")
	}
	terminatorIdx := pitchIdx + len(PitchTable)
	if rip[terminatorIdx] != PitchTableTerminator {
		t.Errorf("byte after PitchTable in the AY rip = %d, want PitchTableTerminator (%d)", rip[terminatorIdx], PitchTableTerminator)
	}

	startupIdx := bytes.Index(rip, StartupMelody)
	if startupIdx < 0 {
		t.Fatal("StartupMelody not found byte-for-byte in the real AY rip")
	}
	if startupIdx != terminatorIdx+1 {
		t.Errorf("StartupMelody starts at AY rip offset %d, want immediately after PitchTable's terminator (%d)", startupIdx, terminatorIdx+1)
	}

	secondaryIdx := bytes.Index(rip, SecondaryMelody)
	if secondaryIdx < 0 {
		t.Fatal("SecondaryMelody not found byte-for-byte in the real AY rip")
	}
	if secondaryIdx <= startupIdx {
		t.Errorf("SecondaryMelody found at AY rip offset %d, want it to appear after StartupMelody (%d)", secondaryIdx, startupIdx)
	}

	if len(StartupMelody) != len(SecondaryMelody) {
		t.Errorf("len(StartupMelody)=%d != len(SecondaryMelody)=%d - the real AY rip confirms both real streams share one length (288), bounded by the same terminator convention", len(StartupMelody), len(SecondaryMelody))
	}
}

// TestAYRipContainsExactlyOneSong is round 158's real, second,
// INDEPENDENT confirmation (from an entirely different methodology
// than round 117's disassembly-based one) that the original 1986 game
// really does have exactly one piece of music, not an unported gap.
//
// Round 117 found the confirmed sound routine's only real entry point
// (Z80 CALL 64671) is called from exactly one place across all 8 of
// this repo's disassembly snapshots - direct, static evidence from
// this project's own reverse-engineering. This test checks a
// completely different, real-world source for the same fact: the AY
// file format's own header (see the AY-file-format spec; also
// self-verified against this exact file - see ay_rip_test.go's other
// test) has a dedicated "number of songs" field, since a real ripper
// examining an actual running game genuinely can, and often does, rip
// MULTIPLE distinct tunes into one file (many Spectrum games have a
// title tune, an in-game tune, and a game-over tune, each as its own
// numbered song). This particular rip - built by a real human in 2001,
// examining the real running game, independent of any disassembly -
// contains exactly ONE song. Two independent methods (a static code
// count, and a human ripper's real-world judgment) landing on the same
// answer is real, additive evidence, not proof by itself, but
// meaningfully stronger than either one alone.
func TestAYRipContainsExactlyOneSong(t *testing.T) {
	rip, err := os.ReadFile("testdata/HeavyOnTheMagick.ay")
	if err != nil {
		t.Fatalf("reading testdata/HeavyOnTheMagick.ay: %v", err)
	}
	if string(rip[0:8]) != "ZXAYEMUL" {
		t.Fatalf("testdata/HeavyOnTheMagick.ay does not start with the real AY magic bytes: %q", rip[0:8])
	}
	// Per the AY file format spec, byte 16 is "NumOfSongs - 1" (a real
	// ripped song count, not this project's own invented interpretation).
	const numSongsMinusOne = 16
	if got := rip[numSongsMinusOne]; got != 0 {
		t.Errorf("AY rip's NumOfSongs-1 field = %d (i.e. %d songs), want 0 (i.e. exactly 1 song)", got, int(got)+1)
	}
}

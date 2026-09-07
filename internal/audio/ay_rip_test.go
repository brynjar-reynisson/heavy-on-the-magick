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

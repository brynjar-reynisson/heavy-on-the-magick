package audio

import (
	"bytes"
	"testing"
)

func TestPitchTableIsMusical(t *testing.T) {
	// Each entry should represent a higher pitch (lower period) than the
	// last, confirming the table is a real descending/ascending scale, not
	// arbitrary data.
	for i := 1; i < len(PitchTable); i++ {
		if PitchTable[i] >= PitchTable[i-1] {
			t.Fatalf("PitchTable[%d]=%d is not less than PitchTable[%d]=%d; expected a monotonic scale",
				i, PitchTable[i], i-1, PitchTable[i-1])
		}
	}
}

func TestPeriodToFrequency(t *testing.T) {
	if PeriodToFrequency(0) != 0 {
		t.Error("PeriodToFrequency(0) should be silence (0 Hz)")
	}
	// Smaller period = higher frequency.
	low := PeriodToFrequency(PitchTable[0])
	high := PeriodToFrequency(PitchTable[len(PitchTable)-1])
	if high <= low {
		t.Errorf("expected the last (smallest-period) table entry to be a higher frequency: low=%.1f high=%.1f", low, high)
	}
}

func TestSquareWaveShape(t *testing.T) {
	samples := SquareWave(440, 0.01, SampleRate)
	if len(samples) != int(0.01*SampleRate) {
		t.Fatalf("len(samples) = %d, want %d", len(samples), int(0.01*SampleRate))
	}
	for _, s := range samples {
		if s != 1 && s != -1 {
			t.Fatalf("square wave sample %v is not +/-1", s)
		}
	}
}

func TestSquareWaveSilence(t *testing.T) {
	samples := SquareWave(0, 0.01, SampleRate)
	for _, s := range samples {
		if s != 0 {
			t.Fatalf("expected silence for 0 Hz, got sample %v", s)
		}
	}
}

// TestRenderNotesHandlesOutOfRangeIndex covers RenderNotes's silence
// fallback for a NoteIndex result outside PitchTable's 0-52 range - 41
// is a real, confirmed StartupMelody value that lands exactly on
// PitchTableTerminator's position (53), not a playable note.
func TestRenderNotesHandlesOutOfRangeIndex(t *testing.T) {
	samples := RenderNotes([]byte{0, 41, 5}, 0.01, SampleRate)
	want := 3 * int(0.01*SampleRate)
	if len(samples) != want {
		t.Fatalf("len(samples) = %d, want %d", len(samples), want)
	}
}

// TestNoteIndex pins the real, confirmed (round 90) signed-byte+12
// decoding rule - see NoteIndex's doc comment for the disassembly and
// cross-stream evidence.
func TestNoteIndex(t *testing.T) {
	cases := []struct {
		raw  byte
		want int
	}{
		{0, 12},   // plain positive value
		{244, 0},  // signed -12 -> the lowest valid index
		{245, 1},  // signed -11
		{254, 10}, // formerly (wrongly) treated as a rest marker
		{252, 8},  // formerly (wrongly) treated as a rest marker
		{41, 53},  // lands exactly on PitchTableTerminator's position
	}
	for _, c := range cases {
		if got := NoteIndex(c.raw); got != c.want {
			t.Errorf("NoteIndex(%d) = %d, want %d", c.raw, got, c.want)
		}
	}
}

// TestRenderNotesHeldNoteHasContinuousPhase covers the real fix for a
// repeated note run (see RenderNotes's doc comment): it should render as
// one continuous tone, identical to a single SquareWave call spanning
// the whole run's duration - not as separate re-triggered notes, which
// would restart the waveform's phase at each tick boundary and produce
// an audible click. note 0's real period (PitchTable[NoteIndex(0)]=128)
// doesn't evenly divide a 0.01s tick at 44100Hz, so a naive
// re-triggering implementation would produce different samples than
// this continuous-phase one - a real behavioral difference, not just a
// cosmetic one.
func TestRenderNotesHeldNoteHasContinuousPhase(t *testing.T) {
	const tick = 0.01
	held := RenderNotes([]byte{0, 0, 0}, tick, SampleRate)
	freq := PeriodToFrequency(PitchTable[NoteIndex(0)])
	continuous := SquareWave(freq, 3*tick, SampleRate)
	if len(held) != len(continuous) {
		t.Fatalf("len(held) = %d, want %d (matching one continuous SquareWave call)", len(held), len(continuous))
	}
	for i := range held {
		if held[i] != continuous[i] {
			t.Fatalf("held[%d] = %v, want %v (a held run must not restart the waveform's phase)", i, held[i], continuous[i])
		}
	}
}

// TestMixNotesAveragesBothStreams covers the basic mixing math: two
// identical single-note streams should average back to the exact same
// waveform (a==b means (a+b)/2==a for every sample), and the result
// should be exactly as long as RenderNotes would produce for that one
// note (both inputs the same length here, so no padding is exercised).
func TestMixNotesAveragesBothStreams(t *testing.T) {
	notes := []byte{5}
	solo := RenderNotes(notes, 0.01, SampleRate)
	mixed := MixNotes(notes, notes, 0.01, SampleRate)
	if len(mixed) != len(solo) {
		t.Fatalf("len(mixed) = %d, want %d", len(mixed), len(solo))
	}
	for i := range mixed {
		if mixed[i] != solo[i] {
			t.Fatalf("mixed[%d] = %v, want %v (averaging two identical streams should reproduce the same waveform)", i, mixed[i], solo[i])
		}
	}
}

// TestMixNotesPadsShorterStream covers the length-mismatch case real
// StartupMelody/SecondaryMelody hit (different lengths): the shorter
// stream must be silence-padded, not truncate the mix, so the longer
// stream's tail still plays (just alone, at half amplitude from the
// averaging).
func TestMixNotesPadsShorterStream(t *testing.T) {
	short := []byte{5}
	long := []byte{5, 5, 5}
	mixed := MixNotes(short, long, 0.01, SampleRate)
	wantLen := len(RenderNotes(long, 0.01, SampleRate))
	if len(mixed) != wantLen {
		t.Fatalf("len(mixed) = %d, want %d (padded to the longer stream)", len(mixed), wantLen)
	}
	longSolo := RenderNotes(long, 0.01, SampleRate)
	tailStart := len(RenderNotes(short, 0.01, SampleRate))
	for i := tailStart; i < len(mixed); i++ {
		want := longSolo[i] / 2
		if mixed[i] != want {
			t.Fatalf("mixed[%d] = %v, want %v (long stream's tail, halved, once the short stream has run out)", i, mixed[i], want)
		}
	}
}

func TestWriteWAVProducesValidHeader(t *testing.T) {
	samples := SquareWave(440, 0.001, SampleRate)
	var buf bytes.Buffer
	if err := WriteWAV(&buf, samples, SampleRate); err != nil {
		t.Fatalf("WriteWAV: %v", err)
	}
	data := buf.Bytes()
	if len(data) < 44 {
		t.Fatalf("WAV output too short: %d bytes", len(data))
	}
	if string(data[0:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		t.Fatalf("missing RIFF/WAVE header: %q", data[:12])
	}
}

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

// TestMixNotesPadsShorterStream covers the general length-mismatch
// case (two arbitrary streams of different lengths - StartupMelody and
// SecondaryMelody themselves turned out to share one real length after
// round 157's correction, so this uses synthetic data instead): the
// shorter stream must be silence-padded, not truncate the mix, so the
// longer stream's tail still plays (just alone, at half amplitude from
// the averaging).
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

// TestXORTickTogglesFirstIterationCancels pins the real traced
// behavior (see xorTickToggles's doc comment): both counters start at
// 1 every tick, so the very first iteration always wraps BOTH of them
// simultaneously - two toggles that cancel out, producing no audible
// edge, even though a real period is always eventually reached.
func TestXORTickTogglesFirstIterationCancels(t *testing.T) {
	toggles := xorTickToggles(5, 7, 96)
	if len(toggles) != 0 {
		t.Errorf("xorTickToggles(5, 7, 96) = %v, want no toggles (both counters wrap on iteration 1 and cancel)", toggles)
	}
}

// TestXORTickTogglesSinglePeriod hand-verifies the exact T-state
// offsets of every toggle for a single active counter (periodB=0,
// silence): iteration 1 wraps E (starts at 1) and reloads it to
// periodA=5, producing one toggle at t=96; the next wrap is 5
// iterations later (5*96=480 T-states after that), at t=576.
func TestXORTickTogglesSinglePeriod(t *testing.T) {
	toggles := xorTickToggles(5, 0, 600)
	want := []float64{96, 576}
	if len(toggles) != len(want) {
		t.Fatalf("xorTickToggles(5, 0, 600) = %v, want %v", toggles, want)
	}
	for i := range want {
		if toggles[i] != want[i] {
			t.Errorf("xorTickToggles(5, 0, 600)[%d] = %v, want %v", i, toggles[i], want[i])
		}
	}
}

// TestXORTickTogglesBothSilent covers periodA=periodB=0 (neither
// stream has a valid note this tick): no counter ever wraps, so no
// toggles occur at all - real silence, not just a near-DC low tone.
func TestXORTickTogglesBothSilent(t *testing.T) {
	toggles := xorTickToggles(0, 0, 1000)
	if len(toggles) != 0 {
		t.Errorf("xorTickToggles(0, 0, 1000) = %v, want no toggles", toggles)
	}
}

func TestTickPeriodPastEndIsSilence(t *testing.T) {
	if got := tickPeriod([]byte{5}, 1); got != 0 {
		t.Errorf("tickPeriod([]byte{5}, 1) = %d, want 0 (past the end - silence)", got)
	}
}

func TestTickPeriodMatchesPitchTable(t *testing.T) {
	got := tickPeriod([]byte{0}, 0)
	want := int(PitchTable[NoteIndex(0)])
	if got != want {
		t.Errorf("tickPeriod([]byte{0}, 0) = %d, want %d (PitchTable[NoteIndex(0)])", got, want)
	}
}

// TestTogglesToSamplesHoldsLevelBetweenToggles covers the basic PCM
// conversion: the level should stay constant between toggle events and
// flip exactly at each one.
func TestTogglesToSamplesHoldsLevelBetweenToggles(t *testing.T) {
	// 1 sample = cpuHz/sampleRate T-states; pick round numbers so each
	// sample corresponds to exactly 100 T-states.
	const sr = 35000 // cpuHz(3.5e6)/sr = 100 T-states per sample
	samples := togglesToSamples([]float64{250, 550}, 1000, sr)
	// sample i covers T-state i*100; toggle at 250 falls between
	// samples 2 (200T) and 3 (300T) -> flips starting sample 3; toggle
	// at 550 flips starting sample 6 (600T).
	want := []float32{1, 1, 1, -1, -1, -1, 1, 1, 1, 1}
	if len(samples) != len(want) {
		t.Fatalf("len(samples) = %d, want %d", len(samples), len(want))
	}
	for i := range want {
		if samples[i] != want[i] {
			t.Errorf("samples[%d] = %v, want %v (full: %v)", i, samples[i], want[i], samples)
		}
	}
}

// TestRenderXORInterleavedSilentWhenNoValidNotes covers total silence
// (both streams past their end / no valid notes anywhere): the output
// should be constant (no toggles ever occur), not an error or garbage.
func TestRenderXORInterleavedSilentWhenNoValidNotes(t *testing.T) {
	samples := RenderXORInterleaved([]byte{41}, []byte{41}, 0.01, SampleRate) // 41 = PitchTableTerminator position, invalid
	if len(samples) == 0 {
		t.Fatal("RenderXORInterleaved with no valid notes produced no samples")
	}
	for i, s := range samples {
		if s != 1 {
			t.Fatalf("samples[%d] = %v, want a constant 1 (no toggles - both streams silent)", i, s)
		}
	}
}

// TestRenderXORInterleavedProducesAudibleOutput covers the real,
// shipped melodies together: confirms real toggles actually occur (not
// silence) and the output is real length (StartupMelody and
// SecondaryMelody are the same real 288-tick length since round 157,
// so either one's tick count works as the basis here).
func TestRenderXORInterleavedProducesAudibleOutput(t *testing.T) {
	samples := RenderXORInterleaved(StartupMelody, SecondaryMelody, 0.15, SampleRate)
	wantLen := int(float64(len(SecondaryMelody)) * 0.15 * float64(SampleRate))
	if samples == nil || len(samples) < wantLen-SampleRate || len(samples) > wantLen+SampleRate {
		t.Errorf("len(samples) = %d, want approximately %d (SecondaryMelody's %d ticks * 0.15s * %d Hz)", len(samples), wantLen, len(SecondaryMelody), SampleRate)
	}
	allSame := true
	first := samples[0]
	for _, s := range samples {
		if s != first {
			allSame = false
			break
		}
	}
	if allSame {
		t.Error("RenderXORInterleaved(StartupMelody, SecondaryMelody, ...) produced constant output - want real audible toggling")
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

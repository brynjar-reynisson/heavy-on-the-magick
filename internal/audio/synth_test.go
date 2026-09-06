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

func TestRenderNotesHandlesRests(t *testing.T) {
	samples := RenderNotes([]byte{0, RestMarkerA, 5}, 0.01, SampleRate)
	want := 3 * int(0.01*SampleRate)
	if len(samples) != want {
		t.Fatalf("len(samples) = %d, want %d", len(samples), want)
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

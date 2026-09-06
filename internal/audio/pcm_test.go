package audio

import "testing"

func TestToStereo16Length(t *testing.T) {
	samples := SquareWave(440, 0.01, SampleRate)
	out := ToStereo16(samples)
	want := len(samples) * 4
	if len(out) != want {
		t.Fatalf("len(out) = %d, want %d (4 bytes per sample: 2 channels x 2 bytes)", len(out), want)
	}
}

func TestToStereo16ChannelsMatch(t *testing.T) {
	samples := []float32{1, -1, 0}
	out := ToStereo16(samples)
	for i := range samples {
		l := out[i*4 : i*4+2]
		r := out[i*4+2 : i*4+4]
		if string(l) != string(r) {
			t.Errorf("sample %d: left=%v right=%v, want equal (mono duplicated to stereo)", i, l, r)
		}
	}
}

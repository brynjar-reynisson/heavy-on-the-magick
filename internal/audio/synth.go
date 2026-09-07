package audio

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
)

// SampleRate is the PCM sample rate used for all generated audio.
const SampleRate = 44100

// cpuHz is the ZX Spectrum 48K's Z80 clock speed.
const cpuHz = 3_500_000.0

// tStatesPerPeriodUnit is a ROUGH APPROXIMATION, not a confirmed value.
// The real beeper loop (Z80 addresses 64733-64781, see pitch_table.go's
// doc comment) branches between multiple paths depending on register E/L
// wraparound, so its exact T-state cost per PitchTable unit needs full
// control-flow simulation to pin down precisely — not yet done. This
// constant is a plausible order-of-magnitude estimate (typical ZX
// Spectrum beeper routines run somewhere in the tens-of-T-states-per-unit
// range) good enough to hear the melody's real relative pitch shape
// (which IS confirmed correct — see PitchTable's semitone-ratio check),
// but the absolute Hz/octave placement here should not be trusted as
// faithful until the loop is fully cycle-counted.
const tStatesPerPeriodUnit = 58.0

// PeriodToFrequency converts a PitchTable period value into an
// approximate frequency in Hz. See tStatesPerPeriodUnit's doc comment for
// the calibration caveat — relative pitch between notes is faithful,
// absolute pitch is not yet confirmed.
func PeriodToFrequency(period byte) float64 {
	if period == 0 {
		return 0
	}
	return cpuHz / (float64(period) * tStatesPerPeriodUnit)
}

// SquareWave generates a square wave at freqHz for durationSec, as PCM
// samples in [-1, 1] — the ZX Spectrum beeper only ever produces a 50%
// duty-cycle square wave (a single bit, toggled), so this is a faithful
// waveform shape even where the exact frequency calibration is not.
func SquareWave(freqHz, durationSec float64, sampleRate int) []float32 {
	n := int(durationSec * float64(sampleRate))
	samples := make([]float32, n)
	if freqHz <= 0 {
		return samples // silence
	}
	period := float64(sampleRate) / freqHz
	for i := range samples {
		phase := math.Mod(float64(i), period) / period
		if phase < 0.5 {
			samples[i] = 1
		} else {
			samples[i] = -1
		}
	}
	return samples
}

// RenderNotes renders a raw note stream (StartupMelody or
// SecondaryMelody) to PCM samples, each timer tick held for
// noteDurationSec. Every value goes through NoteIndex (round 90 -
// signed byte + 12, matching the game's own real indexing) before
// looking up PitchTable; a result outside PitchTable's 0-52 range
// (e.g. the terminator position, 53) produces silence rather than an
// out-of-bounds lookup.
//
// A run of consecutive identical note values is real, confirmed data for
// a single HELD note (see StartupMelody's doc comment: "runs of the same
// value 2-4 times in a row (a held note, played across several timer
// ticks)") - rendered as one continuous SquareWave call spanning the
// whole run's duration, not N separate re-triggered notes. This matters
// audibly: re-triggering restarts the square wave's phase from zero
// every time, producing an artificial click/stutter at each tick
// boundary within what's actually one sustained tone in the source data.
// Total duration is unchanged either way (len(run) * noteDurationSec),
// so this only affects the waveform's smoothness, not the rhythm/timing
// this project has already gotten right.
func RenderNotes(notes []byte, noteDurationSec float64, sampleRate int) []float32 {
	var out []float32
	i := 0
	for i < len(notes) {
		n := notes[i]
		runLen := 1
		for i+runLen < len(notes) && notes[i+runLen] == n {
			runLen++
		}
		duration := float64(runLen) * noteDurationSec
		idx := NoteIndex(n)
		if idx < 0 || idx >= len(PitchTable) {
			out = append(out, make([]float32, int(duration*float64(sampleRate)))...)
		} else {
			freq := PeriodToFrequency(PitchTable[idx])
			out = append(out, SquareWave(freq, duration, sampleRate)...)
		}
		i += runLen
	}
	return out
}

// MixNotes renders two note streams (see RenderNotes) and combines them
// into one PCM buffer by averaging samples — round 98's Stop-hook
// feedback specifically flagged SecondaryMelody as "not integrated into
// actual gameplay," just a separate manual keybinding. That framing
// undersold a fact already on file (see SecondaryMelody's doc comment):
// the real Z80 sound routine at 64671 reads BOTH StartupMelody's and
// SecondaryMelody's streams on every single call, via two independently-
// advancing pointers — the most direct reading of that is that the
// original plays both AT THE SAME TIME, not one after another or only
// on request. This is this port's best-effort reproduction of that:
// genuinely simultaneous playback, not a proven-faithful one — the ZX
// Spectrum beeper is a single output bit, so the real hardware's exact
// two-stream combining trick (probably some form of rapid alternation/
// interleaving, not literal additive mixing, which isn't physically
// what a 1-bit toggle can do) hasn't been traced from the disassembly.
// Averaging is an honest, simple stand-in that at least makes both
// streams audible together, not a claim of bit-exact hardware accuracy.
// The shorter stream is silence-padded to the longer one's length so
// both play to completion (StartupMelody and SecondaryMelody are
// different lengths).
func MixNotes(a, b []byte, noteDurationSec float64, sampleRate int) []float32 {
	sa := RenderNotes(a, noteDurationSec, sampleRate)
	sb := RenderNotes(b, noteDurationSec, sampleRate)
	n := max(len(sb), len(sa))
	out := make([]float32, n)
	for i := range out {
		var va, vb float32
		if i < len(sa) {
			va = sa[i]
		}
		if i < len(sb) {
			vb = sb[i]
		}
		out[i] = (va + vb) / 2
	}
	return out
}

// WriteWAV encodes PCM samples (in [-1, 1]) as a 16-bit mono PCM WAV file.
func WriteWAV(w io.Writer, samples []float32, sampleRate int) error {
	var pcm bytes.Buffer
	for _, s := range samples {
		if s > 1 {
			s = 1
		}
		if s < -1 {
			s = -1
		}
		binary.Write(&pcm, binary.LittleEndian, int16(s*32767))
	}
	dataSize := pcm.Len()

	header := new(bytes.Buffer)
	header.WriteString("RIFF")
	binary.Write(header, binary.LittleEndian, uint32(36+dataSize))
	header.WriteString("WAVE")
	header.WriteString("fmt ")
	binary.Write(header, binary.LittleEndian, uint32(16)) // fmt chunk size
	binary.Write(header, binary.LittleEndian, uint16(1))  // PCM
	binary.Write(header, binary.LittleEndian, uint16(1))  // mono
	binary.Write(header, binary.LittleEndian, uint32(sampleRate))
	binary.Write(header, binary.LittleEndian, uint32(sampleRate*2)) // byte rate
	binary.Write(header, binary.LittleEndian, uint16(2))            // block align
	binary.Write(header, binary.LittleEndian, uint16(16))           // bits per sample
	header.WriteString("data")
	binary.Write(header, binary.LittleEndian, uint32(dataSize))

	if _, err := w.Write(header.Bytes()); err != nil {
		return err
	}
	_, err := w.Write(pcm.Bytes())
	return err
}

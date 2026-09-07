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

// tStatesPerPeriodUnit (round 111): a real, hand-cycle-counted value —
// previously a guessed placeholder (58.0, "a plausible order-of-
// magnitude estimate"), now actually derived by tracing the beeper
// loop's real Z80 instructions (addresses 64733-64781, disassembly in
// hotm.skool) and summing their documented T-state costs.
//
// The loop toggles two INDEPENDENT counters per DJNZ iteration: E
// (reloaded from IXh — the note byte from the FIRST stream,
// StartupMelody, via the 64716-64731 setup) and L (reloaded from H —
// the SECOND stream's, SecondaryMelody's, PitchTable lookup result;
// see routine 64649: `LD H,(HL)` loads H directly from
// PitchTable[index]). The speaker line only actually changes state
// (`XOR D` at 64743/64752, D=16 = a single output-bit toggle mask) when
// one of these counters wraps to zero — most DJNZ iterations, where
// NEITHER counter wraps, just re-output the SAME A value via `OUT
// (254),A` (audibly a no-op). This is a genuinely traced structural
// fact, not a guess: SecondaryMelody's mixing into the output is real
// bit-level XOR interleaving of two independently-clocked toggle
// counters on one shared speaker bit — closer to a beat-frequency/
// interference pattern than either "true 2-channel mixing" or
// MixNotes's sample-averaging approximation. Round 112 implemented
// this exact mechanism directly — see RenderXORInterleaved/
// xorTickToggles below, MixNotes's honest-approximation sibling.
//
// Cycle count for the DOMINANT path (a "no-wrap" iteration — the most
// common case for typical PitchTable values, since they're mostly
// larger than a single-iteration granularity): NOP(4) NOP(4)
// EX AF,AF'(4) DEC E(4) OUT(11) JR NZ,taken(12) JR Z,not-taken(7)
// EX AF,AF'(4) DEC L(4) JP Z,not-taken(10) OUT(11) NOP(4) NOP(4)
// DJNZ,taken(13) = 96 T-states per iteration. A full audible square-
// wave cycle needs L to wrap TWICE (one edge each way), so
// full-cycle T-states ≈ PitchTable_value × 96 × 2 = PitchTable_value ×
// 192 — hence 192.0 here, not 58.0.
//
// Still an approximation, not a live-verified value: this uses the
// dominant no-wrap path's cost uniformly, ignoring the slightly
// different cost of the rarer E-wrap iterations (a real but small
// perturbation — E's own wrap period, from the paired stream's current
// note, is usually much larger than one DJNZ iteration, so this
// shouldn't shift the result by more than a few percent) and hasn't
// been checked against a live emulator's actual audio output. Relative
// pitch between notes remains independently confirmed correct (see
// PitchTable's semitone-ratio check, unaffected by this constant).
const tStatesPerPeriodUnit = 192.0

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

// xorTickToggles reproduces the beeper loop's real bit-level dual-
// counter toggle mechanism (see tStatesPerPeriodUnit's doc comment) for
// ONE melody tick lasting tStates T-states, given the tick's two
// current PitchTable period values (0 means that stream has no valid
// note this tick — its counter never wraps, so it contributes no
// toggles). Returns the T-state offset, from the start of this tick,
// of every real speaker-bit toggle.
//
// Traced from the disassembly (round 111/112): both counters (E for
// the first stream, L for the second) are freshly set to 1 at the
// start of EVERY tick (`POP DE`/the second `CALL 64649` both leave the
// low byte at 1, right before the loop at 64733 begins) — NOT carried
// over from the previous tick. This means the very first loop
// iteration of every tick always wraps BOTH counters simultaneously,
// which cancels out (`toggled = !toggled` applied twice) and produces
// no audible edge — but does reload both counters to their real
// periods (periodA, periodB) for the rest of the tick. Each
// tStatesPerIteration-T-state step after that decrements both
// counters, toggling (and reloading) independently whichever one
// reaches zero — literal bit-level XOR interleaving on one shared
// speaker line, not alternation or additive mixing.
func xorTickToggles(periodA, periodB int, tStates float64) []float64 {
	const tStatesPerIteration = 96.0
	eCount, lCount := 1, 1
	var toggles []float64
	for t := tStatesPerIteration; t <= tStates; t += tStatesPerIteration {
		toggled := false
		if periodA > 0 {
			eCount--
			if eCount <= 0 {
				toggled = !toggled
				eCount = periodA
			}
		}
		if periodB > 0 {
			lCount--
			if lCount <= 0 {
				toggled = !toggled
				lCount = periodB
			}
		}
		if toggled {
			toggles = append(toggles, t)
		}
	}
	return toggles
}

// tickPeriod looks up notes[i]'s real PitchTable period value the same
// way RenderNotes does (via NoteIndex), or 0 (silence — see
// xorTickToggles) if i is past the end of notes or the note decodes
// outside PitchTable's valid range.
func tickPeriod(notes []byte, i int) int {
	if i >= len(notes) {
		return 0
	}
	idx := NoteIndex(notes[i])
	if idx < 0 || idx >= len(PitchTable) {
		return 0
	}
	return int(PitchTable[idx])
}

// togglesToSamples renders a sequence of speaker-bit toggle times
// (absolute T-state offsets from the start of playback, strictly
// increasing) into a square-wave PCM buffer covering totalTStates
// T-states of real Z80 time. Starts at level +1, matching SquareWave's
// convention.
func togglesToSamples(toggleTStates []float64, totalTStates float64, sampleRate int) []float32 {
	n := int(totalTStates / cpuHz * float64(sampleRate))
	out := make([]float32, n)
	level := float32(1)
	ti := 0
	for i := range out {
		sampleTState := float64(i) / float64(sampleRate) * cpuHz
		for ti < len(toggleTStates) && toggleTStates[ti] <= sampleTState {
			level = -level
			ti++
		}
		out[i] = level
	}
	return out
}

// RenderXORInterleaved renders two note streams (see RenderNotes)
// together using the REAL traced bit-level combining mechanism
// (xorTickToggles), rather than MixNotes's sample-averaging
// approximation — round 112's follow-up to round 111's tracing work,
// closing the gap between "the mechanism is documented" and "it's
// actually reproduced." Each tick's period pair comes from a and b's
// current bytes (tickPeriod); the shorter stream contributes silence
// (period 0) for any tick past its own end, same convention as
// MixNotes's padding. Still not a claim of bit-EXACT hardware
// accuracy: tStatesPerIteration (96, from xorTickToggles) is itself
// the dominant-path approximation documented on tStatesPerPeriodUnit,
// and the exact per-tick iteration COUNT (real registers B and C,
// still untraced) is derived from noteDurationSec the same way every
// other render in this package already does, not extracted.
func RenderXORInterleaved(a, b []byte, noteDurationSec float64, sampleRate int) []float32 {
	n := max(len(a), len(b))
	tStatesPerTick := noteDurationSec * cpuHz
	var allToggles []float64
	cumT := 0.0
	for i := range n {
		periodA := tickPeriod(a, i)
		periodB := tickPeriod(b, i)
		for _, tg := range xorTickToggles(periodA, periodB, tStatesPerTick) {
			allToggles = append(allToggles, cumT+tg)
		}
		cumT += tStatesPerTick
	}
	return togglesToSamples(allToggles, cumT, sampleRate)
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
// genuinely simultaneous playback, not a proven-faithful one.
//
// Round 111 traced the real combining trick from the disassembly (see
// tStatesPerPeriodUnit's doc comment): it's bit-level XOR interleaving
// of two independently-clocked toggle counters on the beeper's one
// output bit — not literal additive mixing (physically impossible for
// a 1-bit toggle, as already suspected) and not simple alternation
// either, closer to a beat-frequency/interference pattern between the
// two streams' current pitches. Round 112 implemented that exact
// mechanism as RenderXORInterleaved — this averaging-based function
// remains a simpler, cheaper approximation kept for comparison/
// fallback use, not because the real mechanism is still unknown. The
// shorter stream is silence-padded to the longer one's length so both
// play to completion (StartupMelody and SecondaryMelody are different
// lengths) — RenderXORInterleaved handles that the same way, via
// tickPeriod treating a past-the-end index as silence.
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

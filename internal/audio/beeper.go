// Package audio reproduces the game's sound. The ZX Spectrum 48K has no
// sound chip — all audio is a single "beeper" channel, produced by the CPU
// toggling I/O port 0xFE at precise timings to generate square waves.
//
// STATUS: found the sound routine (Z80 addresses 64671-64781, disassembled
// in ../../CLAUDE.md) and a confirmed, genuinely musical 53-note chromatic
// pitch table at address 64812 (see pitch_table.go — verified via its
// semitone ratio, not guessed). A square-wave PCM synthesizer
// (synth.go) can render that data to real audio (WAV export; no live
// playback backend wired up yet — see Player below). The exact period-to-
// Hz calibration is still an approximation (the beeper loop's T-state cost
// wasn't fully cycle-counted — see synth.go's tStatesPerPeriodUnit doc),
// so relative pitch is faithful but absolute pitch/octave placement isn't
// confirmed yet.
package audio

// Player will play back the game's beeper sound effects and music through
// a real audio backend. No implementation exists yet (WriteWAV in synth.go
// is the only working output path so far — useful for offline rendering/
// testing, not live playback).
type Player interface {
	// PlayTone plays a single square-wave tone.
	PlayTone(hz float64, durationMs int)
}

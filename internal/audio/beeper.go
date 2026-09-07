// Package audio reproduces the game's sound. The ZX Spectrum 48K has no
// sound chip — all audio is a single "beeper" channel, produced by the CPU
// toggling I/O port 0xFE at precise timings to generate square waves.
//
// STATUS: found the sound routine (Z80 addresses 64671-64781, disassembled
// in ../../CLAUDE.md) and a confirmed, genuinely musical 53-note chromatic
// pitch table at address 64812 (see pitch_table.go — verified via its
// semitone ratio, not guessed). A square-wave PCM synthesizer (synth.go)
// renders that data to real audio, and this package now has TWO real
// output paths, not just offline export: WriteWAV (synth.go) for offline
// rendering/testing, and RenderNotes + ToStereo16 (pcm.go) for live
// playback — cmd/hotm-gui feeds their output directly into ebiten's own
// audio player, so no separate Player interface/backend of this
// package's own was ever needed (an earlier round's forward-looking
// Player interface sat here unused and was removed once the real live-
// playback path turned out to be simpler than anticipated). The exact
// period-to-Hz calibration is still an approximation (the beeper loop's
// T-state cost wasn't fully cycle-counted — see synth.go's
// tStatesPerPeriodUnit doc), so relative pitch is faithful but absolute
// pitch/octave placement isn't confirmed yet.
package audio

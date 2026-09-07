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
// playback path turned out to be simpler than anticipated). Round 111
// actually cycle-counted the beeper loop's real Z80 T-state cost (see
// synth.go's tStatesPerPeriodUnit doc) instead of guessing it, and round
// 112 implemented the loop's real bit-level 2-stream combining mechanism
// (RenderXORInterleaved) rather than leaving it approximated.
//
// ONE MELODY IS LIKELY ALL THE ORIGINAL GAME'S AUDIO EVER WAS — not a
// porting gap. Searching every disassembly snapshot in this repo (8, as
// of round 117) for `CALL 64671` (the sound routine's only entry point)
// finds exactly ONE call site in each, always at address 64613 — the
// documented boot-sequence "play a tune while waiting for a keypress"
// loop. No other code anywhere calls this routine or any of its
// internal entry points (64733/64749/64764/64771 are only ever reached
// by falling through/jumping from within 64671 itself). A further scan
// of every standalone `OUT (254),A` outside this routine found 12 more,
// all genuinely border-color-setting code (port 254 controls border
// color AND the speaker simultaneously in ZX Spectrum hardware — an
// `OUT (254)` alone doesn't imply sound), not additional sound effects.
// This is real, reproducible evidence (not certainty) that the original
// 1986 game has exactly one piece of music, played once at boot — there
// is no "other rooms' ambient sound" or per-location audio to port,
// because the original very likely never had any either.
package audio

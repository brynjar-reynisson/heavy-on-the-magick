// Command render-melody renders one of the extracted note streams
// (StartupMelody or SecondaryMelody, see internal/audio/pitch_table.go),
// or both combined (round 99's audio.MixNotes sample-averaging
// approximation, or round 112's audio.RenderXORInterleaved — a real,
// traced bit-level reproduction of the actual Z80 combining mechanism;
// see either function's doc comment for what each honestly does and
// doesn't claim), to a WAV file so it can actually be listened to —
// useful for sanity-checking the pitch-table/melody extraction against
// ear, and for anyone wanting to hear the current state of the sound-
// porting work without wiring up a live audio backend.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/audio"
)

func main() {
	out := flag.String("out", "", "output WAV file path (default: <track>.wav)")
	track := flag.String("track", "xor", `which extracted note stream to render: "startup", "secondary", "mixed" (sample-averaging, round 99), or "xor" (real traced bit-level interleave, round 112 - this is what cmd/hotm-gui actually plays at startup as of round 112)`)
	noteSeconds := flag.Float64("note-seconds", 0.15, "seconds each note/rest is held for")
	flag.Parse()

	var samples []float32
	switch *track {
	case "startup":
		samples = audio.RenderNotes(audio.StartupMelody, *noteSeconds, audio.SampleRate)
	case "secondary":
		samples = audio.RenderNotes(audio.SecondaryMelody, *noteSeconds, audio.SampleRate)
	case "mixed":
		samples = audio.MixNotes(audio.StartupMelody, audio.SecondaryMelody, *noteSeconds, audio.SampleRate)
	case "xor":
		samples = audio.RenderXORInterleaved(audio.StartupMelody, audio.SecondaryMelody, *noteSeconds, audio.SampleRate)
	default:
		fmt.Fprintf(os.Stderr, "unknown -track %q, want \"startup\", \"secondary\", \"mixed\", or \"xor\"\n", *track)
		os.Exit(1)
	}

	outPath := *out
	if outPath == "" {
		outPath = *track + "_melody.wav"
	}

	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := audio.WriteWAV(f, samples, audio.SampleRate); err != nil {
		panic(err)
	}
}

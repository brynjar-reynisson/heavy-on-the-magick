// Command render-melody renders one of the extracted note streams
// (StartupMelody or SecondaryMelody, see internal/audio/pitch_table.go),
// or both mixed together (round 99's audio.MixNotes — see its doc
// comment for what "mixed" honestly does and doesn't claim), to a WAV
// file so it can actually be listened to — useful for sanity-checking
// the pitch-table/melody extraction against ear, and for anyone wanting
// to hear the current state of the sound-porting work without wiring up
// a live audio backend.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/audio"
)

func main() {
	out := flag.String("out", "", "output WAV file path (default: <track>.wav)")
	track := flag.String("track", "mixed", `which extracted note stream to render: "startup", "secondary", or "mixed" (both together, via audio.MixNotes - this is what cmd/hotm-gui actually plays at startup as of round 99; see round 90's CLAUDE.md writeup for what's confirmed vs. not about the second stream)`)
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
	default:
		fmt.Fprintf(os.Stderr, "unknown -track %q, want \"startup\", \"secondary\", or \"mixed\"\n", *track)
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

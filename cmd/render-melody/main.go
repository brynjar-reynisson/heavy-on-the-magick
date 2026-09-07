// Command render-melody renders one of the extracted note streams
// (StartupMelody or SecondaryMelody, see internal/audio/pitch_table.go)
// to a WAV file, so it can actually be listened to — useful for
// sanity-checking the pitch-table/melody extraction against ear, and for
// anyone wanting to hear the current state of the sound-porting work
// without wiring up a live audio backend.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/audio"
)

func main() {
	out := flag.String("out", "", "output WAV file path (default: <track>.wav)")
	track := flag.String("track", "startup", `which extracted note stream to render: "startup" or "secondary" (see round 90's CLAUDE.md writeup for what's confirmed vs. not about the second stream)`)
	noteSeconds := flag.Float64("note-seconds", 0.15, "seconds each note/rest is held for")
	flag.Parse()

	var notes []byte
	switch *track {
	case "startup":
		notes = audio.StartupMelody
	case "secondary":
		notes = audio.SecondaryMelody
	default:
		fmt.Fprintf(os.Stderr, "unknown -track %q, want \"startup\" or \"secondary\"\n", *track)
		os.Exit(1)
	}

	outPath := *out
	if outPath == "" {
		outPath = *track + "_melody.wav"
	}

	samples := audio.RenderNotes(notes, *noteSeconds, audio.SampleRate)

	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := audio.WriteWAV(f, samples, audio.SampleRate); err != nil {
		panic(err)
	}
}

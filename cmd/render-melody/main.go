// Command render-melody renders the extracted StartupMelody note data to a
// WAV file, so it can actually be listened to — useful for sanity-checking
// the pitch-table/melody extraction against ear, and for anyone wanting to
// hear the current state of the sound-porting work without wiring up a
// live audio backend.
package main

import (
	"flag"
	"os"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/audio"
)

func main() {
	out := flag.String("out", "startup_melody.wav", "output WAV file path")
	noteSeconds := flag.Float64("note-seconds", 0.15, "seconds each note/rest is held for")
	flag.Parse()

	samples := audio.RenderNotes(audio.StartupMelody, *noteSeconds, audio.SampleRate)

	f, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := audio.WriteWAV(f, samples, audio.SampleRate); err != nil {
		panic(err)
	}
}

# Heavy on the Magick — Go port

A from-scratch Go port of **Heavy on the Magick** (1986, Gargoyle Games), a
ZX Spectrum graphic adventure. Long-term goal: a faithful port of the
original's gameplay, graphics, and sound, structured so the Go code stays
easy to follow — and easy to extend with a new feature the original never
had: an in-game map of the rooms you've explored.

**Status: playable two ways — text and live graphical/audio.** Movement (8
compass directions) through **Collodon's Pile with real room names and
connections** — Room of Misery, Secunda Porta, Trollwynd, Agile Stair, and
more, sourced from a published walkthrough and cross-confirmed against
strings found in the game's own memory (per Hardcore Gaming 101, the full
dungeon has **255 rooms**; 13 are wired in so far) — plus exploration
tracking and the explored-map feature (`MAP` command) are real and
working:

```
go run ./cmd/hotm        # text frontend
go run ./cmd/hotm-gui     # live graphical window + real audio playback
```

`cmd/hotm-gui` is a real ebiten-backed window: it draws the confirmed rune
glyphs on screen, renders room text live, and plays a real square-wave tone
(synthesized from the game's own extracted pitch table) through your
speakers on every move — not just offline PNG/WAV export. `Alt+Enter`
toggles full screen, matching the usual emulator/game convention. Both
frontends call the exact same `internal/game` engine, so there's one
source of truth for game logic. `cmd/render-melody`/`cmd/render-glyphs`
remain for offline
asset inspection/export.

What's still missing: room *description* text (as opposed to names), most
of the dungeon (13/255 rooms), and most verb resolution (talking, combat,
spellcasting).

See [`CLAUDE.md`](./CLAUDE.md) for the full, detailed reverse-engineering
log — what's been confirmed about the original game's code and data,
what's still unknown, and the tools/techniques used to dig it out
(SkoolKit, a live Spectrum emulator's debugger, etc).

## Layout

- `cmd/hotm/` — text-mode entry point.
- `cmd/hotm-gui/` — live graphical/audio entry point (ebiten). Contains no
  game logic of its own — just input mapping and rendering, calling into
  `internal/game` like `cmd/hotm` does.
- `cmd/render-melody/`, `cmd/render-glyphs/` — offline asset export tools
  (WAV/PNG), useful for inspecting extracted data without running the game.
- `internal/world/` — the room/exit graph (`CollodonsPile`, the real
  walkthrough-sourced dungeon), movement, and the explored-map layout
  algorithm (`Layout`)/ASCII renderer (`RenderASCIIMap`).
- `internal/character/` — the player character (Axil) and the game's
  Golden-Dawn-themed grade system.
- `internal/parser/` — the "TARGET, VERB" command parser, plus the game's
  **real, extracted vocabulary** (`internal/parser/data/vocabulary.json`) —
  316 words pulled directly out of the original game's memory.
- `internal/magic/` — the spellcasting/rune system (early stub).
- `internal/graphics/` — visual assets extracted so far (confirmed rune
  glyphs) plus a real `Renderer` implementation (`PNGRenderer`) reused by
  both `cmd/render-glyphs` (offline) and `cmd/hotm-gui` (live).
- `internal/audio/` — the game's real extracted pitch table + a working
  square-wave synthesizer, reused by both `cmd/render-melody` (offline WAV)
  and `cmd/hotm-gui` (live playback via ebiten's audio player).
- `internal/game/` — wires the above packages together into one session;
  the single source of truth both frontends share.

Every package doc-comments which parts are confirmed-from-the-original vs.
placeholder/not-yet-reverse-engineered — check those before assuming
anything is "real" game behavior.

## Original game files

`hotm1.szx`/`hotm1.tzx`, `hotm-original/`, and the various `hotm*.z80`/
`*.skool`/`*.ctl` files are reverse-engineering artifacts (snapshots,
SkoolKit disassembly output) from the ongoing analysis — not part of the Go
port itself, but kept here since `CLAUDE.md` references them directly.

## Building

```
go build ./...
go test ./...
```

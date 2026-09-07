// Command hotm-gui is a live, on-screen frontend for the Go port — real
// graphics rendering (via ebiten, not just offline PNG export) and real
// audio playback (via ebiten's audio player, not just offline WAV export)
// during actual gameplay, wired to the same internal/game engine used by
// cmd/hotm's text frontend. No game logic lives in this file: it only
// translates keyboard input into parser.Command values and renders
// whatever internal/game.Game.Handle returns.
//
// Layout: redesigned to match a real SpecEmu screenshot of the actual
// original (Room of Misery, the default starting room) rather than this
// project's own earlier invented arrangement (a rune-glyph HUD strip plus
// a scrolling log, which doesn't resemble anything the original shows).
// The real screen is: one big room picture across the top, then a
// magenta-bordered 3-panel status bar below it - a left panel (EXITS, or
// - per a second reference screenshot showing the SAME slot after
// pressing Z/SWAP - the current room's name/level/grade; see
// showRoomStatus), a middle panel (message/command-echo text on a light
// background), and a right panel (STAMINA/SKILL/LUCK on green). This
// version approximates that structure with solid-color panels (no chain-
// link border texture asset exists) using the ZX Spectrum's own real,
// confirmed non-bright palette values.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	ebitenaudio "github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	etext "github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"

	hotmaudio "github.com/brynjar-reynisson/heavy-on-the-magick/internal/audio"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/game"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/graphics"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/magic"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/world"
)

// Layout constants. pictureHeight/statusBar* were sized from a real
// SpecEmu screenshot of Room of Misery: the room picture occupies the top
// ~58% of the screen, the 3-panel status bar the rest - 224/160 out of a
// 384-tall window matches that ratio while keeping this port's existing
// 512x384 window size unchanged.
const (
	screenWidth  = 512
	screenHeight = 384

	pictureHeight   = 224
	statusBarY      = pictureHeight
	statusBarHeight = screenHeight - pictureHeight // 160
	borderThickness = 4

	leftPanelWidth  = 148
	midPanelWidth   = 196
	rightPanelWidth = 148

	// maxLogLines: the message panel has real, limited vertical room
	// (statusBarHeight minus borders and the label row) - capped so the
	// log can never run past its own panel, the bug found and fixed in
	// this port's earlier single-column layout.
	maxLogLines = 6

	// midPanelMaxChars/midPanelMaxLines bound the message panel's real
	// wrapped text (see wrapLine/drawMidPanelText) - a long response
	// line (e.g. "(room description not yet extracted...)") would
	// otherwise be drawn at its full pixel width, bleeding across the
	// panel border into the stats panel to its right; game.Handle
	// responses are prose, not pre-wrapped to any particular width, so
	// this port's own renderer has to do it. basicfont.Face7x13 is a
	// fixed 7px-wide font, so char-count math is exact, not approximate.
	midPanelMaxChars = (midPanelWidth - 2*6) / 7
	midPanelMaxLines = (statusBarHeight - 2*borderThickness) / 16
)

var face = etext.NewGoXFace(basicfont.Face7x13)

// Colors matching the ZX Spectrum's real, confirmed non-bright palette
// (see internal/graphics.Color / pngrenderer.go's palette) - used here so
// the panel backgrounds match the same real values the offline renderer
// uses, not arbitrary RGB guesses.
var (
	black      = color.Black
	textOnLite = color.Black // the message/status panels have light backgrounds, so their text is black, matching the real screenshots
	zxCyan     = color.RGBA{0, 214, 214, 255}
	zxGreen    = color.RGBA{0, 214, 0, 255}
	zxMagenta  = color.RGBA{214, 0, 214, 255}
	zxWhite    = color.RGBA{214, 214, 214, 255} // the real ZX "white" (214, not 255) - the message panel's background
)

// keyDirections maps every real single-letter compass abbreviation (N/S/
// E/W - see parser/keywords.go's merphishKeywords) plus the 4 arrow keys
// to the 4 real cardinal directions. Diagonals (NE/SE/SW/NW) have no
// dedicated key at all - on a real Spectrum keyboard there's no single
// key for a 2-letter abbreviation either, you type both letters and
// press ENTER, which is exactly what this GUI's typed-command line
// (updateTyping) now does. Previously (see git history) this used a
// "roguelike numpad-on-letters" scheme where W meant North - directly
// wrong once compared against a real SpecEmu screenshot, since the
// original's own W means WEST.
var keyDirections = map[ebiten.Key]world.Direction{
	ebiten.KeyArrowUp:    world.North,
	ebiten.KeyN:          world.North,
	ebiten.KeyArrowRight: world.East,
	ebiten.KeyE:          world.East,
	ebiten.KeyArrowDown:  world.South,
	ebiten.KeyS:          world.South,
	ebiten.KeyArrowLeft:  world.West,
	ebiten.KeyW:          world.West,
}

type GUI struct {
	g         *game.Game
	log       []string
	audioCtx  *ebitenaudio.Context
	portraits map[string]*ebiten.Image
	roomArt   map[string]*ebiten.Image // real extracted room screenshots, keyed by world.Room.Name - see pictureImage
	// startupPlayer holds the looping startup-melody player (round 129)
	// so it isn't garbage-collected mid-loop; not otherwise read.
	startupPlayer *ebitenaudio.Player
	// typing/inputBuffer hold the real typed-command line (opened by
	// ENTER, see updateTyping) - lets any real Merphish word/abbreviation
	// or the conversation form ("NAME, OBJECT") reach game.Handle exactly
	// as the original expects, not just whichever subset has a dedicated
	// single-key shortcut below.
	typing      bool
	inputBuffer []rune
	// showRoomStatus toggles the left status-bar panel between EXITS
	// (false) and the room name/level/grade (true) - a real, observed
	// mechanic: a second reference screenshot of the same game, after
	// SWAP (Merphish "Z") was used, showed that exact panel's content
	// replaced by "YOU ARE IN THE <room>, ON LEVEL <n>, YOUR GRADE IS
	// <grade>" instead of "EXITS:". This is the first concrete evidence
	// of what SWAP's "Window 1" actually shows - previously an honest
	// stub with the display "not modeled yet" (see game.go's SWAP case).
	showRoomStatus bool
}

// NewGUI builds a live GUI session around g. Round 97: previously always
// hardcoded game.New() (CollodonsPile) — this GUI had never supported
// the 4 level-grid exploration modes cmd/hotm's text frontend has had
// for a long time, a real, previously-unaddressed gap between the two
// frontends. Passing the constructed *game.Game in (rather than
// building it internally) lets main's -levelNgrid flags select which
// world to play, the same way cmd/hotm's flags already do.
// roomArt maps real extracted room screenshots to the specific
// world.Room.Name they were extracted for (round 108 — previously a
// single image.Image tied to just the world's starting room; see
// selectGame for what's populated per mode).
func NewGUI(g *game.Game, roomArt map[string]image.Image) *GUI {
	gui := &GUI{
		g:        g,
		audioCtx: ebitenaudio.NewContext(hotmaudio.SampleRate),
	}
	if len(roomArt) > 0 {
		gui.roomArt = make(map[string]*ebiten.Image, len(roomArt))
		for name, img := range roomArt {
			gui.roomArt[name] = ebiten.NewImageFromImage(img)
		}
	}
	gui.appendLog(describeRoom(gui.g))
	gui.portraits = make(map[string]*ebiten.Image, len(graphics.PortraitNames))
	for _, name := range graphics.PortraitNames {
		gui.portraits[name] = ebiten.NewImageFromImage(graphics.Portrait(name))
	}
	gui.playStartupMelody()
	return gui
}

// playStartupMelody plays the real, extracted audio.StartupMelody
// combined with audio.SecondaryMelody. The disassembly (see
// SecondaryMelody's doc comment) found the real Z80 sound routine reads
// both note streams together on every call, via two independently-
// advancing pointers — the most direct reading of that fact is the
// original plays them AT THE SAME TIME, not one requiring a separate
// manual keypress to ever be heard (SecondaryMelody's Y-key binding,
// still available below, only ever exercises it in isolation).
//
// Uses audio.RenderXORInterleaved (not the simpler audio.MixNotes): the
// disassembly traced the real single-bit-speaker combining mechanism
// (bit-level XOR interleaving of two independently-clocked toggle
// counters, see tStatesPerPeriodUnit's doc comment), and
// RenderXORInterleaved actually reproduces it.
//
// Loops continuously (via ebiten's audio.InfiniteLoop): a real 1986
// CRASH magazine review confirms "Gargoyle have produced an intro tune
// which improves and becomes more complete the longer you leave it
// playing on the introduction screens" — direct confirmation the real
// game's tune loops repeatedly, not plays once and stops. This port has
// no separate "introduction screen" state to bound the loop to (unlike
// the original, gameplay starts immediately) - honestly simplified to
// loop for the life of the session.
func (gui *GUI) playStartupMelody() {
	samples := hotmaudio.RenderXORInterleaved(hotmaudio.StartupMelody, hotmaudio.SecondaryMelody, 0.15, hotmaudio.SampleRate)
	pcm := hotmaudio.ToStereo16(samples)
	loop := ebitenaudio.NewInfiniteLoop(bytes.NewReader(pcm), int64(len(pcm)))
	player, err := gui.audioCtx.NewPlayer(loop)
	if err != nil {
		return
	}
	gui.startupPlayer = player
	player.Play()
}

func describeRoom(g *game.Game) string {
	return g.Handle(parser.Parse("LOOK"))
}

func (gui *GUI) appendLog(s string) {
	for line := range strings.SplitSeq(s, "\n") {
		gui.log = append(gui.log, line)
	}
	if len(gui.log) > maxLogLines {
		gui.log = gui.log[len(gui.log)-maxLogLines:]
	}
}

// playBlip plays a short confirmed-pitch-table tone through ebiten's real,
// live audio player (internal/audio's synthesizer + PCM conversion do the
// actual sound generation; this just hands the bytes to ebiten to play
// through real speakers, unlike cmd/render-melody's offline-only WAV
// export).
func (gui *GUI) playBlip(noteIndex byte) {
	samples := hotmaudio.RenderNotes([]byte{noteIndex}, 0.08, hotmaudio.SampleRate)
	pcm := hotmaudio.ToStereo16(samples)
	player := gui.audioCtx.NewPlayerFromBytes(pcm)
	player.Play()
}

// Update handles all keyboard input. Single-key action shortcuts are now
// aligned to their REAL Merphish letter meaning (see parser/keywords.go's
// merphishKeywords) wherever one is defined, instead of this port's
// earlier ad hoc convenience bindings - e.g. D now really means DROP (was
// previously O, since D was tied up in the old roguelike movement
// scheme), X really means EXAMINE (was V), O really means OPTIONS (was
// DROP), H really means HALT (was this port's own HELP shortcut - HELP
// is still reachable by typing it in full via ENTER), L really means
// LEFT (was this port's own LOOK shortcut - same typing fallback), R
// really means RIGHT (was GRADE), Z really means SWAP (previously
// unbound), and N/S/E/W are real movement letters (see keyDirections)
// rather than this port's own NAME/SPELLS/etc. shortcuts. Letters with
// no real Merphish meaning (M=map, T=transfusion, G=pass guards, K=talk
// to Apex, J=inventory) are kept as this port's own reasonable
// conveniences, since there's no real letter they'd be overriding.
func (gui *GUI) Update() error {
	// Alt+Enter toggles full screen, matching the convention used by most
	// emulators and games (explicitly requested). Checked before the
	// typing-mode gate below so ALT+ENTER never gets swallowed as "submit
	// the empty command line".
	if ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
		return nil
	}

	if gui.typing {
		gui.updateTyping()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		gui.typing = true
		gui.inputBuffer = nil
		return nil
	}

	for key, dir := range keyDirections {
		if inpututil.IsKeyJustPressed(key) {
			// Goes through parser.Parse + world.ParseDirection exactly like
			// the text frontend (cmd/hotm) does — no GUI-only shortcut path.
			before := gui.g.World.Current
			result := gui.g.Handle(parser.Parse(directionWord(dir)))
			gui.appendLog(result)
			if gui.g.World.Current != before {
				gui.playBlip(byte(7 + int(dir)*3)) // varies pitch by direction, not gameplay-meaningful yet
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		gui.appendLog(gui.g.Handle(parser.Parse("EXAMINE")))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		gui.pickUpFirstItem()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		gui.dropFirstItem()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		gui.invokeDemonForGroundedCharm()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyB) {
		gui.handleAndPlay("BLAST")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		gui.handleAndPlay("FREEZE")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		gui.appendLog(gui.g.Handle(parser.Parse("LEFT")))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		gui.appendLog(gui.g.Handle(parser.Parse("RIGHT")))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		gui.appendLog(gui.g.Handle(parser.Parse("HALT")))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		gui.appendLog(gui.g.Handle(parser.Parse("OPTIONS")))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		gui.appendLog(gui.g.Handle(parser.Parse("SWAP")))
		gui.showRoomStatus = !gui.showRoomStatus
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		gui.handleAndPlay("TRANSFUSION")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		gui.handleAndPlay("GUARDS, DOOR")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		gui.appendLog(gui.g.Handle(parser.Parse("APEX, TALK")))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		gui.appendLog(gui.g.Handle(parser.Parse("MAP")))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		gui.appendLog(gui.g.Handle(parser.Parse("INVENTORY")))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyY) {
		gui.playSecondaryMelody()
	}
	return nil
}

// updateTyping handles keyboard input while the real typed-command line
// is open (entered via ENTER, see Update). This is what makes every real
// Merphish word/abbreviation - not just the subset with its own
// dedicated single-key shortcut above - genuinely reachable exactly as
// the original expects: typing "N" or "NE" and pressing ENTER goes
// through the exact same ExpandKeyword step the text frontend uses, and
// "ASTAROT, WOLFDORP" (the conversation form) works character-for-
// character the same way. ESCAPE cancels without submitting; BACKSPACE
// edits; any other typed character is captured via
// ebiten.AppendInputChars, the standard ebiten text-input API (works
// across keyboard layouts, unlike reading individual key codes).
func (gui *GUI) updateTyping() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		gui.typing = false
		gui.inputBuffer = nil
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) {
		raw := string(gui.inputBuffer)
		gui.typing = false
		gui.inputBuffer = nil
		if echo, result := submitTypedCommand(gui.g, raw); echo != "" {
			gui.appendLog(echo)
			gui.appendLog(result)
			// See showRoomStatus's doc comment: SWAP toggles the left
			// status panel, whether triggered by the Z quick-key (Update)
			// or typed out here in full.
			if parser.Parse(raw).Verb == "SWAP" {
				gui.showRoomStatus = !gui.showRoomStatus
			}
		}
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(gui.inputBuffer) > 0 {
		gui.inputBuffer = gui.inputBuffer[:len(gui.inputBuffer)-1]
	}
	gui.inputBuffer = ebiten.AppendInputChars(gui.inputBuffer)
}

// submitTypedCommand runs one real typed command through the exact same
// parser.Parse + game.Handle path as every other command in this file -
// split out from updateTyping so this behavior is testable without a
// real ebiten input context. A blank/whitespace-only command (ENTER
// pressed with nothing typed) submits nothing, matching the real game's
// own behavior of not re-running the last command on an empty line.
func submitTypedCommand(g *game.Game, raw string) (echo, result string) {
	cmd := strings.TrimSpace(raw)
	if cmd == "" {
		return "", ""
	}
	return "> " + strings.ToUpper(cmd), g.Handle(parser.Parse(cmd))
}

// playSecondaryMelody plays the real, extracted audio.SecondaryMelody
// (round 90's second discovered note stream) on its own, standalone —
// playStartupMelody already plays it MIXED with StartupMelody at actual
// startup, so this key is now for isolating/comparing the second voice
// alone, not the only way to ever hear it during real play. Bound to Y
// (not B) since B now sends the real Merphish BLAST command.
func (gui *GUI) playSecondaryMelody() {
	samples := hotmaudio.RenderNotes(hotmaudio.SecondaryMelody, 0.15, hotmaudio.SampleRate)
	pcm := hotmaudio.ToStereo16(samples)
	player := gui.audioCtx.NewPlayerFromBytes(pcm)
	player.Play()
}

// dropFirstItem handles the D key (real Merphish "D" = DROP - see
// Update's doc comment). Drops whichever item is first in the player's
// real character.Player.Items inventory, since this GUI has no free-text
// object entry for the instant-key path (the typed-command line does,
// via "DROP <name>").
func (gui *GUI) dropFirstItem() {
	if len(gui.g.Player.Items) == 0 {
		gui.appendLog(gui.g.Handle(parser.Parse("DROP")))
		return
	}
	gui.appendLog(gui.g.Handle(parser.Parse("DROP " + gui.g.Player.Items[0])))
}

// pickUpFirstItem handles the P key (real Merphish "P" = PICKUP). The
// instant-key path has no free-form target, so this targets whichever
// item is first in the current room's real, sourced world.Room.Items,
// going through the exact same game.Game.Handle("PICKUP ...") path
// either way, not a GUI-only shortcut that bypasses it.
func (gui *GUI) pickUpFirstItem() {
	room := gui.g.World.CurrentRoom()
	cmd := "PICKUP"
	if room != nil && len(room.Items) > 0 {
		cmd = "PICKUP " + room.Items[0]
	}
	gui.appendLog(gui.g.Handle(parser.Parse(cmd)))
}

// invokeDemonForGroundedCharm handles the I key (real Merphish "I" =
// INVOKE). The instant-key path has no free-form target, so this scans
// the CURRENT ROOM's real Items against each confirmed magic.Demons's
// Charm and invokes the first match, going through the exact same
// game.Game.Handle("INVOKE ...") path either way. With no matching Charm
// on the ground, falls back to bare INVOKE (lists the 4 demons and their
// requirements) rather than doing nothing.
//
// Scans the room, not the player's carried Items: a real source (World
// of Spectrum's plain-text instructions file) confirms the real
// mechanic requires the Charm to be dropped on the ground, not merely
// carried ("Place Ye the talisman on the ground and proceed with thy
// invocation from a distance" - see game.invoke's doc comment).
func (gui *GUI) invokeDemonForGroundedCharm() {
	room := gui.g.World.CurrentRoom()
	var items []string
	if room != nil {
		items = room.Items
	}
	gui.handleAndPlay(invokeCommandFor(items))
}

// invokeCommandFor picks which real game.Handle("INVOKE ...") command
// invokeDemonForGroundedCharm should send, given a set of item names
// (the current room's real Items) - split out so the target-selection
// logic is testable without needing a real audio context.
func invokeCommandFor(items []string) string {
	for _, d := range magic.Demons {
		for _, item := range items {
			if strings.EqualFold(item, d.Charm) {
				return "INVOKE " + d.Name
			}
		}
	}
	return "INVOKE"
}

// Event feedback notes. Only PitchTable itself is confirmed real data
// (extracted from the game's memory); which specific note plays for which
// game event here is this port's own choice, not an extracted mapping —
// same honesty distinction as the movement blip's pitch-by-direction
// formula above.
const (
	noteCombatHit    = 10 // BLAST/FREEZE connects, target still standing
	noteCombatDefeat = 30 // target destroyed/frozen solid
	noteHeal         = 35 // TRANSFUSION
	noteDeath        = 0  // player's Stamina reaches 0
	noteGuardsPass   = 20 // GUARDS, DOOR successfully clears a real obstacle
	noteInvokeOK     = 40 // INVOKE succeeds (the player carries the demon's Charm)
)

// handleAndPlay runs a combat/utility command through the exact same
// game.Game.Handle used everywhere else, then picks a feedback note by
// inspecting the real returned message for the outcome (defeated vs.
// still-standing vs. death) rather than guessing from the verb alone.
func (gui *GUI) handleAndPlay(verb string) {
	result := gui.g.Handle(parser.Parse(verb))
	gui.appendLog(result)
	switch {
	case strings.Contains(result, "GAME OVER") || strings.Contains(result, "You are dead"):
		gui.playBlip(noteDeath)
	case strings.Contains(result, "destroyed") || strings.Contains(result, "solid"):
		gui.playBlip(noteCombatDefeat)
	case strings.Contains(result, "strength return"):
		gui.playBlip(noteHeal)
	case strings.Contains(result, "still standing"):
		gui.playBlip(noteCombatHit)
	case strings.Contains(result, "let you pass"):
		gui.playBlip(noteGuardsPass)
	case strings.Contains(result, "You invoke"):
		gui.playBlip(noteInvokeOK)
	}
}

func directionWord(d world.Direction) string {
	switch d {
	case world.North:
		return "NORTH"
	case world.NorthEast:
		return "NORTH-EAST"
	case world.East:
		return "EAST"
	case world.SouthEast:
		return "SOUTH-EAST"
	case world.South:
		return "SOUTH"
	case world.SouthWest:
		return "SOUTH-WEST"
	case world.West:
		return "WEST"
	case world.NorthWest:
		return "NORTH-WEST"
	default:
		return ""
	}
}

// pictureImage picks whichever real extracted image belongs in the big
// top picture area this frame: a demon/NPC/monster portrait (see
// currentPortraitName) takes priority - the most immediately relevant
// feedback for what the player just did/is facing - falling back to the
// current room's own real extracted screenshot, if one exists. Neither
// existing is a real, honest state (most rooms have no extracted art
// yet): the caller just leaves the picture area black.
func (gui *GUI) pictureImage() (*ebiten.Image, bool) {
	if name, ok := gui.currentPortraitName(); ok {
		return gui.portraits[name], true
	}
	room := gui.g.World.CurrentRoom()
	if room == nil {
		return nil, false
	}
	img, ok := gui.roomArt[room.Name]
	return img, ok
}

// drawFitted scales img to fit entirely inside a (boxW x boxH) box,
// preserving its aspect ratio (uniform scale, never stretched), centered
// within the box - used for the big top picture area, whose real source
// images have varying, individually-cropped aspect ratios.
func drawFitted(screen *ebiten.Image, img *ebiten.Image, boxW, boxH float64) {
	iw := float64(img.Bounds().Dx())
	ih := float64(img.Bounds().Dy())
	if iw == 0 || ih == 0 {
		return
	}
	scale := boxW / iw
	if ih*scale > boxH {
		scale = boxH / ih
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((boxW-iw*scale)/2, (boxH-ih*scale)/2)
	screen.DrawImage(img, op)
}

// fillPanel fills a rectangular sub-region of screen with clr - the
// building block for the magenta-bordered 3-panel status bar (see
// Draw), since ebiten.Image.Fill only fills the whole image.
func fillPanel(screen *ebiten.Image, x, y, w, h int, clr color.Color) {
	screen.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image).Fill(clr)
}

func (gui *GUI) Draw(screen *ebiten.Image) {
	screen.Fill(black)

	if img, ok := gui.pictureImage(); ok {
		drawFitted(screen, img, screenWidth, pictureHeight)
	}
	gui.drawPictureBadges(screen)

	// The magenta-bordered 3-panel status bar (see the package doc
	// comment). Filling the whole bar magenta first, then the 3 inset
	// panels, leaves a uniform magenta border/gutter around and between
	// them - a solid-color approximation of the original's chain-link
	// border texture (no such texture asset exists in this port).
	fillPanel(screen, 0, statusBarY, screenWidth, statusBarHeight, zxMagenta)

	leftX := borderThickness
	midX := leftX + leftPanelWidth + borderThickness
	rightX := midX + midPanelWidth + borderThickness
	panelY := statusBarY + borderThickness
	panelH := statusBarHeight - 2*borderThickness

	leftBG := zxCyan
	if gui.showRoomStatus {
		leftBG = zxGreen // see showRoomStatus's doc comment - the real swapped panel is green, not cyan
	}
	fillPanel(screen, leftX, panelY, leftPanelWidth, panelH, leftBG)
	fillPanel(screen, midX, panelY, midPanelWidth, panelH, zxWhite)
	fillPanel(screen, rightX, panelY, rightPanelWidth, panelH, zxGreen)

	gui.drawLeftPanelText(screen, leftX, panelY)
	gui.drawMidPanelText(screen, midX, panelY)
	gui.drawRightPanelText(screen, rightX, panelY)
}

// drawPictureBadges overlays 2 small, real-colored indicators on the
// picture area when they apply: a live Monster with no confirmed
// portrait (monsterGlyphColor's letter+color, an honest "unconfirmed
// icon" fallback for whichever creature isn't one of the 13 extracted
// portraits) and a real, un-cleared Guards obstacle (guardsColor). Both
// are genuinely sourced facts (see their own doc comments) that predate
// this layout redesign - kept as small picture overlays rather than a
// separate HUD row, closer to how the original likely conveys in-room
// hazards pictorially rather than via a text sidebar.
func (gui *GUI) drawPictureBadges(screen *ebiten.Image) {
	room := gui.g.World.CurrentRoom()
	if room == nil {
		return
	}
	y := 4
	if room.Monster != "" && room.MonsterHealth > 0 {
		if _, hasPortrait := gui.currentPortraitName(); !hasPortrait {
			letter, c := "?", color.RGBA{255, 255, 255, 255}
			if gc, ok := monsterGlyphColor[room.Monster]; ok {
				letter, c = gc.letter, gc.c
			}
			opts := &etext.DrawOptions{}
			opts.GeoM.Translate(4, float64(y))
			opts.ColorScale.ScaleWithColor(c)
			etext.Draw(screen, letter+" "+room.Monster, face, opts)
			y += 16
		}
	}
	if room.Guards {
		opts := &etext.DrawOptions{}
		opts.GeoM.Translate(4, float64(y))
		opts.ColorScale.ScaleWithColor(guardsColor)
		etext.Draw(screen, "I Guards", face, opts)
	}
}

// exitLetter returns room's short compass abbreviation for d ("N", "NE",
// etc. - the same letters parser/keywords.go confirms as real Merphish
// input), or "" if room has no exit that way.
func exitLetter(room *world.Room, d world.Direction) string {
	if room == nil {
		return ""
	}
	if _, ok := room.Exits[d]; !ok {
		return ""
	}
	switch d {
	case world.North:
		return "N"
	case world.NorthEast:
		return "NE"
	case world.East:
		return "E"
	case world.SouthEast:
		return "SE"
	case world.South:
		return "S"
	case world.SouthWest:
		return "SW"
	case world.West:
		return "W"
	case world.NorthWest:
		return "NW"
	}
	return ""
}

// exitsPanelLines lays room's real Exits out compass-style (a 3x3 grid
// of text rows: NW/N/NE, W/·/E, SW/S/SE) approximating the real
// SpecEmu screenshot's own "EXITS:" panel, which showed each direction
// positioned roughly where that compass direction actually sits (W to
// the left, E to the right) rather than a plain comma list. Split out
// from drawLeftPanelText so this layout is testable without a real
// ebiten image.
func exitsPanelLines(room *world.Room) []string {
	return []string{
		"EXITS:",
		"",
		fmt.Sprintf("%-3s%-3s%-3s", exitLetter(room, world.NorthWest), exitLetter(room, world.North), exitLetter(room, world.NorthEast)),
		fmt.Sprintf("%-3s   %-3s", exitLetter(room, world.West), exitLetter(room, world.East)),
		fmt.Sprintf("%-3s%-3s%-3s", exitLetter(room, world.SouthWest), exitLetter(room, world.South), exitLetter(room, world.SouthEast)),
	}
}

// roomStatusLines is the left panel's OTHER real mode (see
// showRoomStatus's doc comment) - the room's name/level/grade, matching
// the real observed "YOU ARE IN THE <room> ON LEVEL <n> YOUR GRADE IS
// <grade>" screen text as closely as this port's own data allows. A
// further real detail on that screenshot ("1°=10°", presumably some
// grade-progress figure) isn't reproduced - its exact meaning isn't
// confirmed, so it's honestly omitted rather than guessed at.
func roomStatusLines(g *game.Game) []string {
	room := g.World.CurrentRoom()
	if room == nil {
		return []string{"YOU ARE IN THE", "VOID"}
	}
	lines := []string{"YOU ARE IN THE", strings.ToUpper(room.Name)}
	if room.Level != 0 {
		lines = append(lines, fmt.Sprintf("ON LEVEL %d", room.Level))
	}
	lines = append(lines, "YOUR GRADE IS", strings.ToUpper(g.Player.Grade.String()))
	return lines
}

func (gui *GUI) drawLeftPanelText(screen *ebiten.Image, x, y int) {
	var lines []string
	if gui.showRoomStatus {
		lines = roomStatusLines(gui.g)
	} else {
		lines = exitsPanelLines(gui.g.World.CurrentRoom())
	}
	opts := &etext.DrawOptions{}
	opts.GeoM.Translate(float64(x+6), float64(y+6))
	opts.LineSpacing = 16
	opts.ColorScale.ScaleWithColor(textOnLite)
	etext.Draw(screen, strings.Join(lines, "\n"), face, opts)
}

// statsPanelLines formats the player's real confirmed Stamina/Skill/
// Luck/XP stats (see character.Player's doc comment) one per line,
// matching the real screenshot's right-hand green panel (STAMINA/SKILL/
// LUCK) with XP added as a 4th line - real, sourced data this port
// already tracks that the reference screenshot doesn't happen to show in
// this exact panel, kept visible here rather than dropped.
func statsPanelLines(g *game.Game) []string {
	p := g.Player
	return []string{
		fmt.Sprintf("STAMINA %d/%d", p.Stamina, p.MaxStamina),
		fmt.Sprintf("SKILL   %d", p.Skill),
		fmt.Sprintf("LUCK    %d", p.Luck),
		fmt.Sprintf("XP      %d", p.ExperiencePoints),
	}
}

func (gui *GUI) drawRightPanelText(screen *ebiten.Image, x, y int) {
	opts := &etext.DrawOptions{}
	opts.GeoM.Translate(float64(x+6), float64(y+6))
	opts.LineSpacing = 16
	opts.ColorScale.ScaleWithColor(textOnLite)
	etext.Draw(screen, strings.Join(statsPanelLines(gui.g), "\n"), face, opts)
}

// wrapLine breaks s into whole-word lines no wider than maxChars -
// game.Handle's responses are ordinary prose, not pre-wrapped to any
// particular width, so the panel that displays them has to do this
// itself or long lines bleed across the panel's border into whatever's
// drawn next to it (a real bug this exact function fixes - caught live
// while first verifying this layout: "(room description not yet
// extracted from the original game)" bled visibly into the stats
// panel). A single word longer than maxChars is left on its own line
// rather than split mid-word.
func wrapLine(s string, maxChars int) []string {
	if maxChars <= 0 {
		return []string{s}
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{s}
	}
	lines := make([]string, 0, 1+len(s)/maxChars)
	cur := words[0]
	for _, w := range words[1:] {
		if len(cur)+1+len(w) > maxChars {
			lines = append(lines, cur)
			cur = w
			continue
		}
		cur += " " + w
	}
	return append(lines, cur)
}

// drawMidPanelText draws the recent command/response log, real-word-
// wrapped and capped to what actually fits in the panel (see
// midPanelMaxChars/midPanelMaxLines), and - while the typed-command
// line is open (see updateTyping) - the live input buffer as its own
// trailing line, replacing the log's own idle state. This is the panel
// a real reference screenshot showed holding exactly this kind of
// content (a scrolling command echo, e.g. "EAST,WEST,WEST"), on the
// same light background used here.
func (gui *GUI) drawMidPanelText(screen *ebiten.Image, x, y int) {
	raw := gui.log
	if gui.typing {
		raw = append(append([]string{}, raw...), "> "+string(gui.inputBuffer)+"_")
	}
	var wrapped []string
	for _, line := range raw {
		wrapped = append(wrapped, wrapLine(line, midPanelMaxChars)...)
	}
	if len(wrapped) > midPanelMaxLines {
		wrapped = wrapped[len(wrapped)-midPanelMaxLines:]
	}
	opts := &etext.DrawOptions{}
	opts.GeoM.Translate(float64(x+6), float64(y+6))
	opts.LineSpacing = 16
	opts.ColorScale.ScaleWithColor(textOnLite)
	etext.Draw(screen, strings.Join(wrapped, "\n"), face, opts)
}

// apexPortraitShouldShow reports whether the most recent log line is a
// talkToApex response ("APEX, TALK"/"APEX, SPEAK") — split out from
// currentPortraitName so this decision is testable without a real
// ebiten image. Checks the LAST line specifically (not the whole log) so
// the portrait only appears right after actually talking to Apex, not
// forever once it's scrolled into log history.
func apexPortraitShouldShow(log []string) bool {
	if len(log) == 0 {
		return false
	}
	return strings.Contains(log[len(log)-1], "Apex the Ogre")
}

// invokedDemonPortraitName reports which demon (if any) the most recent
// log line represents a successful "You invoke ..." response for (see
// game.invoke's doc comment) — split out for the same testability reason
// as apexPortraitShouldShow.
func invokedDemonPortraitName(log []string) (string, bool) {
	if len(log) == 0 {
		return "", false
	}
	last := log[len(log)-1]
	if !strings.HasPrefix(last, "You invoke ") {
		return "", false
	}
	for _, d := range magic.Demons {
		if strings.Contains(last, d.Name) {
			return strings.ToLower(d.Name), true
		}
	}
	return "", false
}

// currentPortraitName picks which real extracted portrait (see
// graphics.Portrait's doc comment for sourcing) belongs in the picture
// area this frame, given the room's live Monster and the most recent log
// line — a just-happened conversation/invocation takes priority over
// the room's ambient monster, since it's the more immediately relevant
// feedback for what the player just did.
func (gui *GUI) currentPortraitName() (string, bool) {
	if name, ok := invokedDemonPortraitName(gui.log); ok {
		return name, true
	}
	if apexPortraitShouldShow(gui.log) {
		return "apex", true
	}
	room := gui.g.World.CurrentRoom()
	if room != nil && room.Monster != "" && room.MonsterHealth > 0 {
		name := strings.ToLower(room.Monster)
		if _, ok := gui.portraits[name]; ok {
			return name, true
		}
	}
	return "", false
}

// monsterGlyphColor maps each confirmed monster name to the real
// letter+color icon it's drawn with on the game's own clean grid map
// (heavymap-grid-clean.gif — exact RGB values read directly from the
// image's own indexed palette and cross-checked against its printed
// legend, e.g. "w wraith" in bright red vs. "w werewolf" in magenta -
// see ../../CLAUDE.md's Level 1-3 monster-scan rounds). Used by
// drawPictureBadges as a fallback badge when no real portrait exists for
// the room's monster.
//
// "Wraith" renamed to "Vampire" (see Level1Grid's doc comment) - the
// map's own legend glossed this red "w" icon "wraith", but a second,
// more authoritative source (a real in-game creature-portrait
// screenshot) confirms the actual name is "Vampire".
var monsterGlyphColor = map[string]struct {
	letter string
	c      color.RGBA
}{
	"Troll":    {"t", color.RGBA{132, 132, 0, 255}},
	"Cyclops":  {"c", color.RGBA{132, 132, 0, 255}},
	"Ghost":    {"g", color.RGBA{0, 255, 0, 255}},
	"Slug":     {"s", color.RGBA{0, 132, 0, 255}},
	"Vampire":  {"w", color.RGBA{255, 0, 0, 255}},
	"Medusa":   {"m", color.RGBA{255, 0, 0, 255}},
	"Werewolf": {"w", color.RGBA{255, 0, 255, 255}},
	"Wyvern":   {"w", color.RGBA{0, 132, 255, 255}},
}

// guardsColor is the confirmed real color of the clean grid map's
// "guards" legend icon (a red "Ɪ" glyph) - see world.Room.Guards's doc
// comment for the sourcing. Used by drawPictureBadges.
var guardsColor = color.RGBA{255, 0, 0, 255}

func (gui *GUI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// selectGame picks which world to play from the same -levelNgrid flags
// cmd/hotm's text frontend already has (round 97 gave this GUI the same
// selection, previously CollodonsPile-only) - returns the constructed
// *game.Game plus a title suffix describing the mode, for the window
// title/log so it's clear which world is active, plus a map of real
// extracted room screenshots keyed by world.Room.Name (round 108
// generalized this from one image tied to the start room - see
// GUI.roomArt's doc comment).
func selectGame() (g *game.Game, modeTitle string, roomArt map[string]image.Image) {
	level1Grid := flag.Bool("level1grid", false, "play the extracted Level 1 grid (64 real cells) instead of CollodonsPile")
	level2Grid := flag.Bool("level2grid", false, "play the extracted Level 2 grid (50 real, fully-connected cells) instead of CollodonsPile")
	level3Grid := flag.Bool("level3grid", false, "play the extracted Level 3 grid (47 real cells: 41 fully-connected plus 6 isolated rooms) instead of CollodonsPile")
	level4Grid := flag.Bool("level4grid", false, "play the extracted Level 4 grid (27 real cells: 17 fully-connected plus 10 isolated rooms) instead of CollodonsPile")
	flag.Parse()

	switch {
	case *level1Grid:
		return game.NewLevel1Exploration(), " (Level 1 grid)", map[string]image.Image{
			"A1":             graphics.Level1CorridorSample(),
			"Agile Stair":    graphics.AgileStairSample(),
			"Furnace Room":   graphics.FurnaceRoomSample(),
			"Room of Stings": graphics.RoomOfStingsSample(),
			"Room of Arrows": graphics.RoomOfArrowsSample(),
		}
	case *level2Grid:
		return game.NewLevel2Exploration(), " (Level 2 grid)", map[string]image.Image{
			"A1":             graphics.CorridorSample(),
			"Room of Misery": graphics.RoomOfMiserySample(),
		}
	case *level3Grid:
		return game.NewLevel3Exploration(), " (Level 3 grid)", map[string]image.Image{"A1": graphics.Level3CorridorSample()}
	case *level4Grid:
		return game.NewLevel4Exploration(), " (Level 4 grid)", map[string]image.Image{"F2": graphics.Level4CorridorSample()}
	default:
		return game.New(), "", map[string]image.Image{
			"Room of Misery": graphics.RoomOfMiserySample(),
			"Room of Stings": graphics.RoomOfStingsSample(),
			"Room of Arrows": graphics.RoomOfArrowsSample(),
			"Agile Stair":    graphics.AgileStairSample(),
			"Furnace Room":   graphics.FurnaceRoomSample(),
			"Wolfdorp":       graphics.WolfdorpSample(),
			"Nidus":          graphics.NidusSample(),
			"Trollwynd":      graphics.TrollwyndSample(),
			"Pilefoot":       graphics.PilefootSample(),
			"Sothic Complex": graphics.SothicComplexSample(),
			"Methos":         graphics.MethosSample(),
			"Morfang":        graphics.MorfangSample(),
			"Sign":           graphics.SignSample(),
		}
	}
}

func main() {
	g, modeTitle, roomArt := selectGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Heavy on the Magick — Go port (live)" + modeTitle)
	if err := ebiten.RunGame(NewGUI(g, roomArt)); err != nil {
		log.Fatal(err)
	}
	fmt.Println("bye")
}

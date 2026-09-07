// Command hotm-gui is a live, on-screen frontend for the Go port — real
// graphics rendering (via ebiten, not just offline PNG export) and real
// audio playback (via ebiten's audio player, not just offline WAV export)
// during actual gameplay, wired to the same internal/game engine used by
// cmd/hotm's text frontend. No game logic lives in this file: it only
// translates keyboard input into parser.Command values and renders
// whatever internal/game.Game.Handle returns, plus the confirmed rune
// glyphs from internal/graphics as a HUD strip — reusing the existing
// PNGRenderer rather than duplicating glyph-drawing code.
package main

import (
	"fmt"
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

const (
	screenWidth  = 512
	screenHeight = 384
	maxLogLines  = 20
)

var face = etext.NewGoXFace(basicfont.Face7x13)

// Colors matching the ZX Spectrum's palette (see internal/graphics.Color)
// for visual consistency with the offline PNG renderer's output.
var (
	black = color.Black
	white = color.White
	cyan  = color.RGBA{0, 214, 214, 255}
	grey  = color.RGBA{128, 128, 128, 255}
)

// keyDirections maps keyboard keys to compass directions using the
// classic roguelike numpad-on-letters layout (Q/W/E / A-D / Z/X/C),
// plus arrow keys for the 4 cardinals — both confirmed-real directions
// from world.ParseDirection, not invented key names.
var keyDirections = map[ebiten.Key]world.Direction{
	ebiten.KeyW:          world.North,
	ebiten.KeyArrowUp:    world.North,
	ebiten.KeyE:          world.NorthEast,
	ebiten.KeyD:          world.East,
	ebiten.KeyArrowRight: world.East,
	ebiten.KeyC:          world.SouthEast,
	ebiten.KeyX:          world.South,
	ebiten.KeyArrowDown:  world.South,
	ebiten.KeyZ:          world.SouthWest,
	ebiten.KeyA:          world.West,
	ebiten.KeyArrowLeft:  world.West,
	ebiten.KeyQ:          world.NorthWest,
}

type GUI struct {
	g         *game.Game
	log       []string
	audioCtx  *ebitenaudio.Context
	hud       *ebiten.Image // confirmed rune glyphs, rendered once via graphics.PNGRenderer
	portraits map[string]*ebiten.Image
}

func NewGUI() *GUI {
	gui := &GUI{
		g:        game.New(),
		audioCtx: ebitenaudio.NewContext(hotmaudio.SampleRate),
	}
	gui.appendLog(describeRoom(gui.g))
	gui.hud = buildHUD()
	gui.portraits = make(map[string]*ebiten.Image, len(graphics.PortraitNames))
	for _, name := range graphics.PortraitNames {
		gui.portraits[name] = ebiten.NewImageFromImage(graphics.Portrait(name))
	}
	gui.playStartupMelody()
	return gui
}

// playStartupMelody plays the real, extracted audio.StartupMelody once
// when the GUI starts up (round 93) — previously this port's live GUI
// never played the melody at all, only individual event blips; the one
// piece of confirmed, real Z80 sound data this project has was
// completely unheard in the actual graphical game. Fire-and-forget,
// same as playBlip - doesn't block Update()/gameplay while it plays.
func (gui *GUI) playStartupMelody() {
	samples := hotmaudio.RenderNotes(hotmaudio.StartupMelody, 0.15, hotmaudio.SampleRate)
	pcm := hotmaudio.ToStereo16(samples)
	player := gui.audioCtx.NewPlayerFromBytes(pcm)
	player.Play()
}

// buildHUD renders the confirmed rune glyphs (internal/graphics.RuneGlyphs)
// using the SAME PNGRenderer built for offline PNG export (cmd/render-glyphs)
// — proving both frontends draw identical, confirmed assets rather than
// each having their own separate (and possibly diverging) drawing code.
func buildHUD() *ebiten.Image {
	r := graphics.NewPNGRenderer(len(graphics.RuneGlyphs)+1, 1, 8)
	r.Clear(graphics.Black)
	for i, gl := range graphics.RuneGlyphs {
		r.DrawGlyph(gl, i, 0, graphics.White, false)
	}
	// The magenta spell-icon's real confirmed ULA attribute is 0x43 -
	// ink=magenta BRIGHT (see KnownIcons's doc comment) - so this is the
	// first place this port renders a real confirmed bright color, not a
	// cosmetic choice.
	r.DrawGlyph(graphics.KnownIcons["magenta-icon"], len(graphics.RuneGlyphs), 0, graphics.Magenta, true)
	return ebiten.NewImageFromImage(r.Image())
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

func (gui *GUI) Update() error {
	// Alt+Enter toggles full screen, matching the convention used by most
	// emulators and games (explicitly requested).
	if ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	for key, dir := range keyDirections {
		if inpututil.IsKeyJustPressed(key) {
			// Goes through parser.Parse + world.ParseDirection exactly like
			// the text frontend (cmd/hotm) does — no GUI-only shortcut path.
			before := gui.g.World.Current
			result := gui.g.Handle(parser.Parse(directionWord(dir)))
			gui.appendLog(result)
			// Only play the movement blip if the move actually succeeded -
			// previously played unconditionally, so a blocked move ("You
			// can't go that way", or the new Fire-blocked rejection) sounded
			// identical to a real step. Comparing room IDs (not scanning
			// result text for specific rejection strings) stays correct
			// automatically as new kinds of blocked-movement messages are
			// added.
			if gui.g.World.Current != before {
				gui.playBlip(byte(7 + int(dir)*3)) // varies pitch by direction, not gameplay-meaningful yet
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		gui.appendLog(gui.g.Handle(parser.Parse("LOOK")))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		gui.appendLog(gui.g.Handle(parser.Parse("MAP")))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		gui.handleAndPlay("BLAST")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		gui.handleAndPlay("FREEZE")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		gui.handleAndPlay("TRANSFUSION")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		gui.pickUpFirstItem()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		gui.dropFirstItem()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		gui.appendLog(gui.g.Handle(parser.Parse("EXAMINE")))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		gui.invokeCarriedDemon()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		gui.handleAndPlay("GUARDS, DOOR")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		gui.appendLog(gui.g.Handle(parser.Parse("APEX, TALK")))
	}
	return nil
}

// dropFirstItem handles the O key (drOp — D and X, the more obvious
// letters, are already movement keys in this GUI's roguelike layout).
// The DROP counterpart to pickUpFirstItem: drops whichever item is first
// in the player's real character.Player.Items inventory.
func (gui *GUI) dropFirstItem() {
	if len(gui.g.Player.Items) == 0 {
		gui.appendLog(gui.g.Handle(parser.Parse("DROP")))
		return
	}
	gui.appendLog(gui.g.Handle(parser.Parse("DROP " + gui.g.Player.Items[0])))
}

// pickUpFirstItem handles the P key. The GUI has no text input, so
// PICKUP/DROP (which take a named object in the real command grammar,
// see parser.Parse) can't offer a free-form target the way the text
// frontend (cmd/hotm) can — this instead targets whichever item is first
// in the current room's real, sourced world.Room.Items, going through
// the exact same game.Game.Handle("PICKUP ...") path either way, not a
// GUI-only shortcut that bypasses it.
func (gui *GUI) pickUpFirstItem() {
	room := gui.g.World.CurrentRoom()
	cmd := "PICKUP"
	if room != nil && len(room.Items) > 0 {
		cmd = "PICKUP " + room.Items[0]
	}
	gui.appendLog(gui.g.Handle(parser.Parse(cmd)))
}

// invokeCarriedDemon handles the I key. The GUI has no text input, so
// INVOKE (which takes a specific demon name, see parser.Parse) can't
// offer a free-form target the way the text frontend can - this
// instead scans the player's real carried Items against each confirmed
// magic.Demons's Charm and invokes the first match, going through the
// exact same game.Game.Handle("INVOKE ...") path either way (same
// target-selection convention as pickUpFirstItem/dropFirstItem above).
// With no matching Charm carried, falls back to bare INVOKE (lists the
// 4 demons and their requirements) rather than doing nothing.
func (gui *GUI) invokeCarriedDemon() {
	gui.handleAndPlay(invokeCommandFor(gui.g.Player.Items))
}

// invokeCommandFor picks which real game.Handle("INVOKE ...") command
// invokeCarriedDemon should send, given the player's carried items -
// split out from invokeCarriedDemon so the target-selection logic is
// testable without needing a real audio context.
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

func (gui *GUI) Draw(screen *ebiten.Image) {
	screen.Fill(black)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(8, 8)
	screen.DrawImage(gui.hud, op)

	gui.drawMonster(screen)
	gui.drawGuards(screen)
	gui.drawItems(screen)
	gui.drawPortrait(screen)

	statsOpts := &etext.DrawOptions{}
	statsOpts.GeoM.Translate(8, 76) // just below the 64px-tall HUD row (drawn at y=8)
	statsOpts.ColorScale.ScaleWithColor(white)
	etext.Draw(screen, gui.statsLine(), face, statsOpts)

	drawOpts := &etext.DrawOptions{}
	drawOpts.GeoM.Translate(8, 100)
	drawOpts.LineSpacing = 16 // basicfont.Face7x13 is 13px tall; give lines breathing room
	drawOpts.ColorScale.ScaleWithColor(cyan)
	etext.Draw(screen, strings.Join(gui.log, "\n"), face, drawOpts)

	helpOpts := &etext.DrawOptions{}
	helpOpts.GeoM.Translate(8, screenHeight-20)
	helpOpts.ColorScale.ScaleWithColor(grey)
	etext.Draw(screen, "WASD/arrows+QEZC=move  L=look  M=map  V=examine  SPACE=blast  F=freeze  T=transfusion  P=pickup  O=drop  I=invoke  G=pass guards  K=talk to Apex  ALT+ENTER=fullscreen", face, helpOpts)
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
// log line represents a successful "You invoke <NAME>, ..." response
// for (see game.invoke's doc comment) — split out for the same
// testability reason as apexPortraitShouldShow.
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
// graphics.Portrait's doc comment for sourcing) belongs in the corner
// this frame, given the room's live Monster and the most recent log
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

// drawPortrait shows whichever real extracted portrait
// currentPortraitName picks - the first place this port displays actual
// extracted 1986 game art (demon/NPC/monster) instead of a custom-drawn
// approximation.
func (gui *GUI) drawPortrait(screen *ebiten.Image) {
	name, ok := gui.currentPortraitName()
	if !ok {
		return
	}
	img := gui.portraits[name]
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(screenWidth-img.Bounds().Dx()-8), 8)
	screen.DrawImage(img, op)
}

// monsterGlyphColor maps each confirmed monster name to the real
// letter+color icon it's drawn with on the game's own clean grid map
// (heavymap-grid-clean.gif — exact RGB values read directly from the
// image's own indexed palette and cross-checked against its printed
// legend, e.g. "w wraith" in bright red vs. "w werewolf" in magenta -
// see ../../CLAUDE.md's Level 1-3 monster-scan rounds). This is the
// first time this port has drawn a monster as anything but plain text -
// a real, sourced graphics improvement, not an invented sprite.
//
// "Wraith" renamed to "Vampire" in round 74 (see Level1Grid's doc
// comment) - the map's own legend glossed this red "w" icon "wraith",
// but a second, more authoritative source (a real in-game creature-
// portrait screenshot) confirms the actual name is "Vampire".
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

// drawMonster renders the current room's real Monster (if any and still
// alive) as its confirmed letter+color icon, next to the HUD row -
// previously only ever shown as plain log text ("You see: a Vampire").
// Monster names not in monsterGlyphColor (e.g. CollodonsPile's generic
// "monster") fall back to a plain white "?", honestly signaling an
// unconfirmed icon rather than guessing one.
func (gui *GUI) drawMonster(screen *ebiten.Image) {
	room := gui.g.World.CurrentRoom()
	if room == nil || room.Monster == "" || room.MonsterHealth <= 0 {
		return
	}
	letter, c := "?", color.RGBA{255, 255, 255, 255}
	if gc, ok := monsterGlyphColor[room.Monster]; ok {
		letter, c = gc.letter, gc.c
	}
	opts := &etext.DrawOptions{}
	opts.GeoM.Translate(280, 8)
	opts.ColorScale.ScaleWithColor(c)
	etext.Draw(screen, letter+" "+room.Monster, face, opts)
}

// guardsColor is the confirmed real color of the clean grid map's
// "guards" legend icon (a red "Ɪ" glyph) - see world.Room.Guards's doc
// comment for the sourcing. Rendered here as the plain ASCII letter "I"
// in that same confirmed color, the same "real color, plain-letter
// stand-in" convention monsterGlyphColor already uses.
var guardsColor = color.RGBA{255, 0, 0, 255}

// drawGuards renders a real, un-cleared world.Room.Guards obstacle in
// its confirmed color, next to the monster indicator - previously only
// ever implied by room text (Exits/description), never shown visually.
func (gui *GUI) drawGuards(screen *ebiten.Image) {
	room := gui.g.World.CurrentRoom()
	if room == nil || !room.Guards {
		return
	}
	opts := &etext.DrawOptions{}
	opts.GeoM.Translate(280, 26) // just below drawMonster's (280, 8) - verified visible live, unlike an earlier (400, 8) attempt that rendered nothing on screen for reasons not fully understood
	opts.ColorScale.ScaleWithColor(guardsColor)
	etext.Draw(screen, "I Guards", face, opts)
}

// itemsColor: the clean grid map's own legend draws its generic
// "object" icon in black ("xx") - but this GUI's background is also
// black, so rendering real Items in that confirmed color would be
// invisible. Unlike monsterGlyphColor/guardsColor, this is honestly NOT
// the confirmed icon color, just a legible stand-in (plain yellow,
// matching this project's existing HUD color conventions elsewhere).
var itemsColor = color.RGBA{255, 255, 0, 255}

// drawItems renders the current room's real Items (if any) as a plain
// list next to the monster/guards indicators - previously only ever
// shown as log text ("You see: Grimoire"), never in the live HUD area.
func (gui *GUI) drawItems(screen *ebiten.Image) {
	room := gui.g.World.CurrentRoom()
	if room == nil || len(room.Items) == 0 {
		return
	}
	opts := &etext.DrawOptions{}
	opts.GeoM.Translate(280, 44) // below drawGuards's (280, 26)
	opts.ColorScale.ScaleWithColor(itemsColor)
	etext.Draw(screen, strings.Join(room.Items, ", "), face, opts)
}

// statsLine renders the player's real confirmed stats (see
// character.Player's doc comment for which fields are confirmed real)
// as a single HUD line — previously not shown anywhere in this GUI at
// all, only in the text frontend's LOOK/EXAMINE output indirectly.
func (gui *GUI) statsLine() string {
	p := gui.g.Player
	return fmt.Sprintf("%s the %s | Stamina %d/%d | Skill %d | Luck %d | XP %d",
		p.Name, p.Grade, p.Stamina, p.MaxStamina, p.Skill, p.Luck, p.ExperiencePoints)
}

func (gui *GUI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Heavy on the Magick — Go port (live)")
	if err := ebiten.RunGame(NewGUI()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("bye")
}

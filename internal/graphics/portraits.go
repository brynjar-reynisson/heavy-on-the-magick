package graphics

import (
	"bytes"
	"embed"
	"image"
)

// assets/*.png are real, tight-cropped extracts of actual in-game
// portraits — sourced from maps.speccy.cz's "Speccy Screenshot Maps"
// atlas (heavymap-speccy-screenshots.png in the repo root, credited to
// its creator Hippy Smith), a composite built from real captured
// screenshots of the running 1986 game, not a fan redrawing. Every
// confirmed demon (Asmodee, Astarot, Belezbar, Magot), Apex the Ogre,
// and all 8 monster types (Ghost, Slug, Vampire, Wyvern, Cyclops,
// Medusa, Troll, Werewolf) has a portrait here — the same gallery whose
// real on-screen name resolved the "Wraith"/"Vampire" naming question
// (see ../../CLAUDE.md). Apex's was extracted first (round 75); the
// remaining 12 followed the identical bounding-box-scan-then-visually-
// verify method in round 76 - each crop was rendered into one review
// grid and individually checked before being trusted, the same
// discipline used for every other pixel-extraction in this project.
//
//go:embed assets/*.png
var portraitFS embed.FS

// Portrait decodes one embedded real creature/NPC/demon portrait by
// name (case-sensitive, matching the asset filename without extension,
// e.g. "Apex", "Vampire", "Belezbar" - see assets/ for the full list).
// Panics on an unknown name or corrupt embedded asset, matching how
// other embedded-asset-at-startup code in this codebase behaves (there
// is no reasonable runtime fallback for a build-time asset problem).
func Portrait(name string) image.Image {
	data, err := portraitFS.ReadFile("assets/" + name + ".png")
	if err != nil {
		panic("graphics: no embedded portrait named " + name + ": " + err.Error())
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic("graphics: failed to decode embedded " + name + ".png: " + err.Error())
	}
	return img
}

// ApexPortrait is a convenience wrapper for Portrait("apex") - kept for
// the one call site (cmd/hotm-gui) that existed before the round-76
// generalization to all 13 portraits, so that code didn't need to
// change its casing convention.
func ApexPortrait() image.Image {
	return Portrait("apex")
}

// PortraitNames are every embedded portrait's Portrait() key, in the
// same order as the source gallery (demons first, then Apex, then the
// 8 monsters) - useful for iterating all of them (e.g. a future round's
// regression test or a GUI asset-preview mode) without hardcoding the
// list twice.
var PortraitNames = []string{
	"asmodee", "astarot", "belezbar", "magot",
	"apex",
	"ghost", "slug", "vampire", "wyvern", "cyclops", "medusa", "troll", "werewolf",
}

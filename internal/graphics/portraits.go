package graphics

import (
	"bytes"
	_ "embed"
	"image"
)

// assets/apex.png is a real, tight-cropped extract of Apex the Ogre's
// actual in-game portrait — sourced from maps.speccy.cz's "Speccy
// Screenshot Maps" atlas (heavymap-speccy-screenshots.png in the repo
// root, credited to its creator Hippy Smith), a composite built from
// real captured screenshots of the running 1986 game, not a fan
// redrawing. This is the first place this port uses actual extracted
// game art instead of a custom-drawn approximation — see
// ../../CLAUDE.md for the fuller writeup of this source and how it also
// resolved the "Wraith"/"Vampire" naming question.
//
//go:embed assets/apex.png
var apexPortraitPNG []byte

// ApexPortrait decodes the embedded real Apex the Ogre portrait. Panics
// on failure, matching how other embedded-asset-at-startup packages in
// this codebase behave (there's no reasonable runtime fallback for a
// corrupt build-time asset).
func ApexPortrait() image.Image {
	img, _, err := image.Decode(bytes.NewReader(apexPortraitPNG))
	if err != nil {
		panic("graphics: failed to decode embedded apex.png: " + err.Error())
	}
	return img
}

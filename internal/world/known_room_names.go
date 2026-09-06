package world

// KnownUnplacedRoomNames are real dungeon room/zone names read off a
// second, much clearer fan-made map — a computer-rendered (not
// hand-drawn) 8x8-per-level grid with lettered/numbered cells (A1-H8),
// color-coded zones, and a monster-icon legend, also hosted at Spectrum
// Computing (HeavyOnTheMagick_5.gif — see ../../CLAUDE.md for the fuller
// write-up of what this map revealed). Precise enough to supersede a few
// earlier low-confidence guesses at these same names (documented in the
// git history / CLAUDE.md), but still NOT wired into CollodonsPile: full
// cell-by-cell connectivity across all 4 levels wasn't extracted (a
// large task on its own), so these are names only, not a connected graph.
//
// A key structural finding from this map, worth knowing before extending
// CollodonsPile further: real room granularity is per-maze-cell (~64
// cells/level x 4 levels =~ 256, matching the confirmed 255-room count),
// NOT per named-zone. Big color-coded zones like "Wolfdorp" or "Gorburg"
// span many cells and are more like neighborhoods; a handful of
// individually-labeled cells within them (e.g. "Room of Claws" inside
// the Nidus/Wolfdorp zone boundary on Level 1) are the actual distinct
// rooms. CollodonsPile's existing entries (sourced from a walkthrough)
// mix both kinds without distinguishing them - that's a known modeling
// simplification, not corrected here to avoid a larger, riskier
// restructuring within this pass.
var KnownUnplacedRoomNames = []string{
	// Individually-labeled cells (Level noted where confirmed):
	"Room of Claws",  // Level 1
	"Room of Horns",  // Level 2
	"Room of Icthys", // Level 2
	"Room of Flox",   // Level 2
	"Purity",         // Level 2
	"Room of Nani",   // Level 3 (corrected from an earlier "Mani" misread; "NANI" is the confirmed real vocabulary word, "MANI" doesn't appear in it)
	"Room of Two",    // Level 3
	"Room of Water",  // Level 3
	"Room of Rains",  // Level 3
	"Room of Scales", // Level 4
	"Doubt of Rabak", // Level 4
	"The Crypt",      // Level 4
	"The Chasm",      // Level 4
	"Room of Pride",  // Level 4
	"The Pit",        // Level 4 (zone-sized, but no smaller name found within it)

	// Larger color-coded zones (span many cells each):
	"Eye of Heaven", // Level 2
	"Wraithvale",    // Level 2
	"Slymole",       // Level 2
	"Quadra Porta",  // Level 2
	"Gorburg",       // Level 3
	"Rook of Hydra", // Level 3
	"Tertia Porta",  // Level 3
	"Kitchen of Ai", // Level 3
	"Wormring",      // Level 4
	"Lichgate",      // Level 4
}

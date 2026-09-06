package world

import "strings"

// directionWords maps the game's real confirmed compass vocabulary (see
// ../parser/vocabulary.go and ../../CLAUDE.md) to Direction. Spelling is
// exact: the diagonals are hyphenated ("NORTH-EAST", not "NORTHEAST"),
// confirmed byte-for-byte from the game's own word table.
var directionWords = map[string]Direction{
	"NORTH":      North,
	"NORTH-EAST": NorthEast,
	"EAST":       East,
	"SOUTH-EAST": SouthEast,
	"SOUTH":      South,
	"SOUTH-WEST": SouthWest,
	"WEST":       West,
	"NORTH-WEST": NorthWest,
}

// ParseDirection matches a player's word against the game's real compass
// vocabulary (case-insensitive). ok is false for anything else, including
// words the original parser recognizes for other purposes (most of its
// 300+ word vocabulary isn't about movement at all).
func ParseDirection(word string) (d Direction, ok bool) {
	d, ok = directionWords[strings.ToUpper(word)]
	return d, ok
}

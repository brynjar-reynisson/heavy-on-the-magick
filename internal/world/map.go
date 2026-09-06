package world

// Coord is a 2D grid position for laying out the explored map.
type Coord struct{ X, Y int }

// delta is the grid step a Direction represents when auto-laying-out rooms.
var delta = map[Direction]Coord{
	North:     {0, -1},
	NorthEast: {1, -1},
	East:      {1, 0},
	SouthEast: {1, 1},
	South:     {0, 1},
	SouthWest: {-1, 1},
	West:      {-1, 0},
	NorthWest: {-1, -1},
}

// Layout assigns every visited room a grid coordinate by walking the exit
// graph breadth-first from the world's starting room, stepping by each
// exit's compass direction. This is the "auto-map" a real Spectrum text
// adventure would produce if its dungeon exits are spatially consistent
// (walking East then West returns you to where you started).
//
// consistent is false if the same room was reached at two different
// coordinates (a non-Euclidean shortcut or loop in the original design) —
// in that case the caller should fall back to a node/edge graph rendering
// instead of trusting the grid positions, since they'll overlap or
// contradict each other. See the discussion in ../../CLAUDE.md.
func Layout(w *World, start RoomID) (positions map[RoomID]Coord, consistent bool) {
	positions = map[RoomID]Coord{start: {0, 0}}
	consistent = true

	queue := []RoomID{start}
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		room := w.Rooms[id]
		if room == nil || !room.Visited {
			continue
		}
		here := positions[id]
		for dir, destID := range room.Exits {
			dest := w.Rooms[destID]
			if dest == nil || !dest.Visited {
				continue
			}
			d := delta[dir]
			want := Coord{here.X + d.X, here.Y + d.Y}
			if got, seen := positions[destID]; seen {
				if got != want {
					consistent = false
				}
				continue
			}
			positions[destID] = want
			queue = append(queue, destID)
		}
	}
	return positions, consistent
}

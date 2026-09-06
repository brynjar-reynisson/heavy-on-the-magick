package world

// World holds the full room graph and the player's current position.
type World struct {
	Rooms   map[RoomID]*Room
	Current RoomID
}

// New creates an empty World. Rooms are added with AddRoom once the real
// room data has been extracted from the disassembly.
func New(start RoomID) *World {
	return &World{
		Rooms:   make(map[RoomID]*Room),
		Current: start,
	}
}

// AddRoom registers a room in the world.
func (w *World) AddRoom(r *Room) {
	w.Rooms[r.ID] = r
}

// CurrentRoom returns the room the player currently occupies.
func (w *World) CurrentRoom() *Room {
	return w.Rooms[w.Current]
}

// Move attempts to walk in direction d from the current room. Returns false
// if there's no exit that way (mirrors the original's "you can't go that
// way"-style rejection, once that response text is extracted).
func (w *World) Move(d Direction) bool {
	room := w.CurrentRoom()
	if room == nil {
		return false
	}
	dest, ok := room.Exits[d]
	if !ok {
		return false
	}
	w.Current = dest
	if next := w.Rooms[dest]; next != nil {
		next.Visited = true
	}
	return true
}

// VisitedRooms returns every room the player has discovered so far — the
// data set the explored-map feature renders.
func (w *World) VisitedRooms() []*Room {
	var visited []*Room
	for _, r := range w.Rooms {
		if r.Visited {
			visited = append(visited, r)
		}
	}
	return visited
}

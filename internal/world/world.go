package world

import "strings"

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

// AddRoom registers a room in the world, keeping its own independent
// deep copy rather than aliasing r.
//
// This matters because CollodonsPile/Level1-4Grid all build their World
// from the SAME package-level []*Room literal every time they're called
// (e.g. game.New() called twice, or two tests each constructing their
// own World from Level1Grid) - without cloning, every such World would
// share the exact same underlying *Room objects, so a runtime mutation
// in one World (Visited, PICKUP removing an Item, a cleared Guards
// obstacle, a defeated Monster) would silently leak into every other
// World built from that same source data, including ones created
// later. Found and fixed after a new test (checking a cleared Guards
// obstacle) broke an unrelated, already-passing test that assumed a
// fresh Level1Grid() still had its original Guards placement - real
// evidence this was an actual bug, not just a hypothetical one.
func (w *World) AddRoom(r *Room) {
	w.Rooms[r.ID] = r.clone()
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

// FindRoomByName looks up a room by its exact (case-insensitive) Name,
// for name-based mechanics like Astarot's "transports the player to a
// named location" ability (see game.astarotTeleport).
func (w *World) FindRoomByName(name string) (RoomID, bool) {
	for id, r := range w.Rooms {
		if strings.EqualFold(r.Name, name) {
			return id, true
		}
	}
	return 0, false
}

// Teleport moves the player directly to the given room, bypassing normal
// Exits - the real mechanic behind Astarot's confirmed location-transport
// ability, as opposed to Move's ordinary corridor-by-corridor walking.
func (w *World) Teleport(id RoomID) {
	w.Current = id
	if room := w.Rooms[id]; room != nil {
		room.Visited = true
	}
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

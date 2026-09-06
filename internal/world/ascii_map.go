package world

import (
	"fmt"
	"strings"
)

// RenderASCIIMap draws the rooms the player has explored so far as text —
// the simplest possible real implementation of the explored-map feature,
// exercising Layout end-to-end. Each visited room is drawn as a 3-char
// box using the first 3 letters of its Name, followed by a 1-char marker
// for real per-room state (see roomMarker); the current room is
// highlighted with brackets.
//
// If Layout reports the room graph isn't spatially consistent (see
// map.go's doc comment — a non-Euclidean shortcut in the room design),
// RenderASCIIMap falls back to a simple list instead of attempting a grid
// that would overlap or contradict itself.
func RenderASCIIMap(w *World) string {
	positions, consistent := Layout(w, w.startRoomForLayout())
	if !consistent || len(positions) == 0 {
		return renderRoomList(w)
	}

	minX, minY, maxX, maxY := 0, 0, 0, 0
	for _, p := range positions {
		minX, maxX = min(minX, p.X), max(maxX, p.X)
		minY, maxY = min(minY, p.Y), max(maxY, p.Y)
	}

	var b strings.Builder
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			id, ok := roomAt(positions, x, y)
			room := w.Rooms[id]
			if !ok || room == nil {
				b.WriteString("      ")
				continue
			}
			label := roomLabel(room)
			if id == w.Current {
				fmt.Fprintf(&b, "[%s]%s", label, roomMarker(room))
			} else {
				fmt.Fprintf(&b, " %s %s", label, roomMarker(room))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

// roomMarker is a real, at-a-glance indicator of what's actually in a
// room right now — "!" for a living Monster (real per-room data, see
// CollodonsPile/Level1Grid's doc comments for sourcing), "*" for one or
// more Items, or a blank if neither. A defeated Monster (MonsterHealth
// <= 0) no longer marks the room, so the map reflects real combat state
// as the player changes it, not just static room contents.
func roomMarker(r *Room) string {
	if r.Monster != "" && r.MonsterHealth > 0 {
		return "!"
	}
	if len(r.Items) > 0 {
		return "*"
	}
	return " "
}

func (w *World) startRoomForLayout() RoomID {
	// Layout needs a starting room to walk the graph from. The player's
	// starting room isn't tracked separately from Current (there's no
	// need to yet), so just use whichever visited room has the lowest ID
	// as a stable, deterministic anchor.
	var start RoomID
	first := true
	for id, r := range w.Rooms {
		if !r.Visited {
			continue
		}
		if first || id < start {
			start = id
			first = false
		}
	}
	return start
}

func roomAt(positions map[RoomID]Coord, x, y int) (RoomID, bool) {
	for id, p := range positions {
		if p.X == x && p.Y == y {
			return id, true
		}
	}
	return 0, false
}

func roomLabel(r *Room) string {
	name := strings.ToUpper(r.Name)
	if len(name) >= 3 {
		return name[:3]
	}
	return (name + "   ")[:3]
}

func renderRoomList(w *World) string {
	var b strings.Builder
	b.WriteString("(room layout not spatially consistent - showing visited rooms as a list)\n")
	for _, r := range w.VisitedRooms() {
		current := " "
		if r.ID == w.Current {
			current = ">"
		}
		fmt.Fprintf(&b, "%s %s %s\n", current, roomMarker(r), r.Name)
	}
	return b.String()
}

package world

import "testing"

func TestLevel4GridHas17Cells(t *testing.T) {
	if len(level4Cells) != 17 {
		t.Fatalf("len(level4Cells) = %d, want 17 (the validated main connected component)", len(level4Cells))
	}
}

func TestLevel4GridExitsAreReciprocal(t *testing.T) {
	w := Level4Grid()
	for id, room := range w.Rooms {
		for dir, destID := range room.Exits {
			dest, ok := w.Rooms[destID]
			if !ok {
				t.Fatalf("room %v exits %v to an unknown RoomID %v", id, dir, destID)
			}
			back, ok := dest.Exits[dir.Opposite()]
			if !ok || back != id {
				t.Errorf("room %q --%s--> %q has no matching reverse exit (%q --%s--> %q)", room.Name, dir, dest.Name, dest.Name, dir.Opposite(), room.Name)
			}
		}
	}
}

func TestLevel4GridIsFullyConnected(t *testing.T) {
	w := Level4Grid()
	visited := map[RoomID]bool{w.Current: true}
	queue := []RoomID{w.Current}
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		for _, dest := range w.Rooms[id].Exits {
			if !visited[dest] {
				visited[dest] = true
				queue = append(queue, dest)
			}
		}
	}
	if len(visited) != 17 {
		t.Errorf("only %d of 17 cells reachable from the start room, want all 17", len(visited))
	}
}

package world

// PlaceholderDungeon is a small, hand-built room graph used to exercise
// the movement/exploration engine end-to-end while the original game's
// real room/exit table is still being extracted from the disassembly
// (see ../../CLAUDE.md — the room table search has traced deep into the
// main game loop without landing on it yet).
//
// IMPORTANT: this is NOT extracted from the original game. Room names and
// descriptions here are invented placeholders, not "Collodon's Pile" — do
// not treat them as faithful content. What IS faithful: the movement
// mechanics themselves (8-direction compass exits, visited tracking, the
// map layout algorithm in map.go), which are built to the confirmed real
// rules and will work identically once real room data replaces this.
func PlaceholderDungeon() *World {
	w := New(0)

	w.AddRoom(&Room{
		ID:          0,
		Name:        "Entrance Hall",
		Description: "[placeholder] A damp stone hall. Passages lead off in several directions.",
		Exits: map[Direction]RoomID{
			North: 1,
			East:  2,
		},
		Visited: true,
	})
	w.AddRoom(&Room{
		ID:          1,
		Name:        "Narrow Passage",
		Description: "[placeholder] The passage narrows here, walls close on both sides.",
		Exits: map[Direction]RoomID{
			South:     0,
			NorthEast: 3,
		},
	})
	w.AddRoom(&Room{
		ID:          2,
		Name:        "Side Chamber",
		Description: "[placeholder] A small chamber, empty but for dust.",
		Exits: map[Direction]RoomID{
			West: 0,
		},
	})
	w.AddRoom(&Room{
		ID:          3,
		Name:        "Junction",
		Description: "[placeholder] Several passages meet here.",
		Exits: map[Direction]RoomID{
			SouthWest: 1,
		},
	})

	return w
}

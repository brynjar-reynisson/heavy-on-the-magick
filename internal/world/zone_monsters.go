package world

// ZoneMonsterSighting records a monster type seen within a named zone on
// the computer-rendered grid map (see known_room_names.go's doc comment
// for the source, HeavyOnTheMagick_5.gif). Recorded at zone granularity,
// not per-cell: the source map places each icon in one specific lettered
// cell (e.g. a Werewolf at C3), but which cell that corresponds to in
// CollodonsPile's coarser room graph isn't known, so only "this zone has
// this monster type, this many times" is captured — deliberately not
// attached to any specific world.Room to avoid overwriting or
// contradicting monster data already sourced from the CASA walkthrough.
//
// Cyclops in the Nidus zone is a good independent cross-confirmation: it
// matches CollodonsPile's Nidus room (likely the same zone - this room
// was itself corrected from "Midus" to "Nidus" per the game's own
// vocabulary, see collodons_pile.go), which already has a Cyclops from
// the walkthrough source.
type ZoneMonsterSighting struct {
	Zone    string
	Monster string
	Count   int
}

var ZoneMonsterSightings = []ZoneMonsterSighting{
	{Zone: "Wolfdorp", Monster: "Ghost", Count: 2},
	{Zone: "Wolfdorp", Monster: "Werewolf", Count: 2},
	{Zone: "Morfang", Monster: "Wraith", Count: 3},
	{Zone: "Nidus", Monster: "Cyclops", Count: 1},
	{Zone: "Wraithvale", Monster: "Wraith", Count: 1},
	{Zone: "Room of Icthys", Monster: "Slug", Count: 1},
	{Zone: "Sothic Complex", Monster: "Ghost", Count: 1},
	{Zone: "Gorburg", Monster: "Wyvern", Count: 1},
	{Zone: "Gorburg", Monster: "Ghost", Count: 2},
	{Zone: "Trollwynd", Monster: "Troll", Count: 4},
	{Zone: "Rook of Hydra", Monster: "Wyvern", Count: 1},
	{Zone: "Wormring", Monster: "Wyvern", Count: 4},
	{Zone: "Doubt of Rabak", Monster: "Wraith", Count: 1},
	{Zone: "Methos", Monster: "Wraith", Count: 1},
	{Zone: "The Pit", Monster: "Medusa", Count: 1},
}

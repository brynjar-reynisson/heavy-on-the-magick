package world

// NumberedRoomContent records what a fan-made numbered map/key poster for
// the game says is in dungeon room number N. Source: "Heavy on the
// Magick - The Map", compiled by A. Britton (Wakefield, Yorks), hosted at
// Spectrum Computing (spectrumcomputing.co.uk/pub/sinclair/games-maps/h/
// HeavyOnTheMagick_3.jpg and _4.jpg — see ../../CLAUDE.md for the fuller
// provenance/honesty notes). The poster numbers 102 maze cells across the
// dungeon's 4 levels and gives each one's contents in a printed (not
// hand-lettered) key list, which is why this data is high-confidence
// despite the maze artwork itself being hard to read precisely.
//
// What this does NOT capture: which numbered cell is which world.Room
// (that mapping isn't extracted — the numbers key into this poster's own
// maze grids, not into CollodonsPile's room graph), nor the maze
// connectivity between numbered cells (deliberately not attempted — see
// CLAUDE.md on why reading exact wall openings off a small hand-drawn
// image was judged too error-prone to responsibly claim as fact). This is
// therefore reference content for when/if that mapping is done, not yet
// playable data.
type NumberedRoomContent struct {
	Number   int
	Contents string
}

// NumberedRoomContents is the full 102-entry key list transcribed from
// the source poster, in room-number order. A few entries reference
// mechanics confirmed elsewhere too: #14's "Snake" wards Hydras and #59's
// "disguised Erlstone" is Asmodee's confirmed Charm (see
// internal/magic.Demons) — both good independent cross-confirmations.
// #96's "Nest of Phoenix" gets the same treatment (round 91): "PHOENIX"
// is independently a real, confirmed word in the game's own extracted
// 316-word parser vocabulary (parser.Vocabulary) — real evidence this
// entry describes actual in-game content, not just fan-map flavor text.
var NumberedRoomContents = []NumberedRoomContent{
	{1, "Grimoire"},
	{2, "Poison-smeared book"},
	{3, "Bag containing Night Shade"},
	{4, "Bag of gold (opens doors with coin pictures)"},
	{5, "Loaf of bread"},
	{6, "Sign - Leo, Key of Nickel"},
	{7, "Chest (Sunflower)"},
	{8, "Sign - Pisces, fishes, Key of Copper"},
	{9, "Chest (jar with hemlock)"},
	{10, "Cabinet"},
	{11, "Sign - Virgo, virgin, Key of Alum"},
	{12, "Egg - rock, protected"},
	{13, "Sign - Gemini, twins, Key of Lithic"},
	{14, "Snake (iron clasp) - inscribed with Undine"},
	{15, "Sign on the wall: means 'do all in order'"},
	{16, "Stalagmite"},
	{17, "Rock (snake, dead cold)"},
	{18, "Chest"},
	{19, "Sign - Capricornus, goat, Key of Magnak"},
	{20, "Rock"},
	{21, "Cabinet (clasp - Salamander charm)"},
	{22, "Scroll (CALL spell)"},
	{23, "Stalactite"},
	{24, "Nougat"},
	{25, "Sign - Aquarius, water carrier, Key of Cobalt"},
	{26, "Chest (mirror)"},
	{27, "Rock"},
	{28, "Rock (honey jar - food)"},
	{29, "Scroll (TRANSFUSION spell)"},
	{30, "Sign - Aries, ram, Key of Bronze"},
	{31, "Pellet - rock, protected"},
	{32, "Cabinet (Mantis)"},
	{33, "Chest"},
	{34, "Bone"},
	{35, "Bone"},
	{36, "Meat bone"},
	{37, "Rock (poison-smeared head)"},
	{38, "Skull"},
	{39, "Rock"},
	{40, "Sign - Taurus, bull, Key of Iron"},
	{41, "Bone"},
	{42, "Thigh"},
	{43, "Meat bone"},
	{44, "Rib"},
	{45, "Rock"},
	{46, "Rock"},
	{47, "Ulna"},
	{48, "Poison-smeared head"},
	{49, "Nugget (silver), rock, protected"},
	{50, "Cauldron of cold iron (scroll inside)"},
	{51, "Chest (leaf, bag of gold)"},
	{52, "Rock"},
	{53, "Rock"},
	{54, "Meat bone"},
	{55, "Ball of copper"},
	{56, "Sign - Libra, scales, Key of Brass"},
	{57, "Pebble"},
	{58, "Pebble"},
	{59, "Pebble (disguised Erlstone)"},
	{60, "Pebble"},
	{61, "Bag of gold"},
	{62, "Pebble"},
	{63, "Pebble"},
	{64, "Stalagmite and rock"},
	{65, "Rock, two stalagmites, stalactite, sword"},
	{66, "Stalagmite, rock"},
	{67, "Bag of gold"},
	{68, "Cabinet"},
	{69, "Sign - Sagittarius, archer, Key of Chrome"},
	{70, "Stalagmite"},
	{71, "Rock"},
	{72, "Stalagmite, stalactite, rock"},
	{73, "Bag of gold, stalagmite, stalactite"},
	{74, "Rock, jar of honey"},
	{75, "Stalagmite, two stalactites, rock"},
	{76, "Rock, stalactite"},
	{77, "Two stalagmites"},
	{78, "Two rocks"},
	{79, "Two stalagmites"},
	{80, "Rock"},
	{81, "Rock"},
	{82, "Stalagmite, loaf of bread"},
	{83, "Stalagmite"},
	{84, "Chest (bag of gold, garlic, foot)"},
	{85, "Sign - Scorpio, scorpion, Key of Zinc"},
	{86, "Stalagmite"},
	{87, "Stalagmite"},
	{88, "Two stalagmites"},
	{89, "Meat bone"},
	{90, "Sign - Cancer, crab, Key of Tin"},
	{91, "Shelf, jar of hemlock, jar of honey"},
	{92, "Stalagmite"},
	{93, "Rock"},
	{94, "Meat bone"},
	{95, "Cabinet"},
	{96, "Nest of Phoenix"},
	{97, "Chest (slat)"},
	{98, "Flask, rock cake"},
	{99, "Cabinet (poison-smeared rock)"},
	{100, "Meat bone"},
	{101, "Meat bone"},
	{102, "Ruby"},
}

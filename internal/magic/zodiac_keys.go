package magic

// ZodiacKey pairs a zodiac sign with the metal key it unlocks. Confirmed
// real via a fan-made numbered map/key poster for the game ("Heavy on the
// Magick" map, compiled by A. Britton, Wakefield, Yorks — hosted at
// Spectrum Computing, spectrumcomputing.co.uk/pub/sinclair/games-maps/h/
// HeavyOnTheMagick_3.jpg — see ../../CLAUDE.md), whose numbered key list
// states e.g. "Sign — Leo, Key of Nickel" for specific room numbers.
// Independently cross-confirmed: several of these metal names (Zinc, Tin,
// Alum, Lithic, Copper, Nickel, Chrome) also appear as key item labels on
// the separate official heavymap-levels1-2.jpg poster — two unrelated
// sources describing the same real game mechanic.
//
// Not yet wired to any gameplay: which room's door each key/sign actually
// unlocks isn't extracted (that would require precise numbered-room
// connectivity, which this project has deliberately not attempted to read
// off the hand-drawn maze grids - see CLAUDE.md's honesty notes on map
// decoding).
//
// Round 116: tried the same "numbered cell sits inside a zone banner"
// method that placed Sword/Sunflower/Erlstone/Mantis for two of these
// 12 Sign entries (#19 "Capricornus", #6 "Leo") - both landed in large,
// genuinely unbannered connector areas between named zones (near Agile
// Stair/Trollwynd for #19, near Room of Misery's immediate surroundings
// for #6), unlike Erlstone's #59 which had exactly one nearby banner.
// Left unplaced rather than guessed — a real, checked negative result,
// not an oversight — see CLAUDE.md for the full writeup. Most of this
// poster's cell numbering appears to cover such connector cells, which
// may mean the remaining 10 Signs are similarly hard to zone-attribute
// without a cleaner source; worth knowing before re-attempting this
// same method on the others.
type ZodiacKey struct {
	Sign   string
	Symbol string // the sign's stated emblem/creature
	Metal  string
}

// ZodiacKeys are the 12 confirmed sign-to-key pairings, in the order
// their numbered rooms appear on the source map.
var ZodiacKeys = []ZodiacKey{
	{Sign: "Leo", Symbol: "lion", Metal: "Nickel"},
	{Sign: "Pisces", Symbol: "fishes", Metal: "Copper"},
	{Sign: "Virgo", Symbol: "virgin", Metal: "Alum"},
	{Sign: "Gemini", Symbol: "twins", Metal: "Lithic"},
	{Sign: "Capricornus", Symbol: "goat", Metal: "Magnak"},
	{Sign: "Aquarius", Symbol: "water carrier", Metal: "Cobalt"},
	{Sign: "Aries", Symbol: "ram", Metal: "Bronze"},
	{Sign: "Taurus", Symbol: "bull", Metal: "Iron"},
	{Sign: "Libra", Symbol: "scales", Metal: "Brass"},
	{Sign: "Scorpio", Symbol: "scorpion", Metal: "Zinc"},
	{Sign: "Sagittarius", Symbol: "archer", Metal: "Chrome"},
	{Sign: "Cancer", Symbol: "crab", Metal: "Tin"},
}

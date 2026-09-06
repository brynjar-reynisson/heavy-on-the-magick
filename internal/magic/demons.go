package magic

// Demon is one of the four Princes documented in the official Gargoyle
// Games instruction manual's "grimoire" section ("Heavy on the Magick",
// 1986, Carter Follis Software — see ../../CLAUDE.md for how this was
// sourced). Fields are the manual's own structured attributes for each
// Demon (title, number, astrological sign, etc.) — factual game data, not
// a copy of the manual's prose (the manual carries an explicit copyright
// notice; only the attribute facts are captured here, not its wording).
//
// These are confirmed real, invocable via the INVOKE spell (Merphish
// keyword "I" — see parser.ExpandKeyword). Invocation requires a
// "suitable Talisman" per the manual; a second source (Spectrum
// Computing's plain-text instructions file for the game, distinct from
// the manual PDF — see ../../CLAUDE.md) confirms each demon's specific
// Charm item, so this requirement is now checked for real in
// internal/game rather than always failing.
type Demon struct {
	Name    string
	Title   string
	Number  int
	Sign    string // astrological/domain association
	Aspect  string // the form the demon takes/is associated with
	Ability string // what invoking the demon does for the player, if known
	Charm   string // the Talisman item needed to invoke this demon
}

// Demons are the 4 confirmed from the manual's grimoire section. Order
// matches the manual's own listing. Ability is filled in for all 4 from a
// second source (a "THE DEMONS"/"...ONIC PRINCES...INVOCATION..." panel
// spanning both heavymap-levels1-2.jpg and a second poster page,
// HeavyOnTheMagick_2.jpg from Spectrum Computing — see ../../CLAUDE.md),
// paraphrased not quoted. Asmodee's is a warning ("Be careful with
// Asmodee"), not a positive effect — fitting his "Great Destroyer" title,
// unlike the other three's helpful abilities. Charm comes from a third
// source (Spectrum Computing's HeavyOnTheMagick.txt instructions file) —
// "Sunflower" independently cross-confirms the "Sun-Flower" item already
// found on the map poster (internal/world.LevelTwoItems), and "Erlstone"
// cross-confirms entry #59 ("Pebble, disguised Erlstone") in
// internal/world.NumberedRoomContents — good evidence all these sources
// describe the same real game.
var Demons = []Demon{
	{Name: "ASMODEE", Title: "the Great Destroyer", Number: 122, Sign: "House of Mars", Aspect: "Basilisk", Ability: "Warning: be careful with Asmodee (no confirmed positive effect)", Charm: "Erlstone"},
	{Name: "ASTAROT", Title: "the Spirit of Assemblage", Number: 1376, Sign: "Sign of Gemini", Aspect: "Legion", Ability: "Transports the player to a named location, if its name is known", Charm: "Sword"},
	{Name: "BELEZBAR", Title: "the Master of Flies", Number: 20, Sign: "Firmament of Stars", Aspect: "Deceit", Ability: "Reveals the true nature of objects", Charm: "Mantis"},
	{Name: "MAGOT", Title: "the Diviner", Number: 443, Sign: "Realm of Air", Aspect: "Baboon", Ability: "Reveals the whereabouts of any named object", Charm: "Sunflower"},
}

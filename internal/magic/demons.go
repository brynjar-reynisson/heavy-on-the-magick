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

	// Correspondences is each Prince's real occult correspondence set
	// (colour/plant/perfume/gem etc., in the tradition this game's
	// theming is drawn from — see round 91's Golden Dawn note) from the
	// manual's grimoire section (round 110, extracted via the PDF's own
	// text layer — see ../../CLAUDE.md — rather than the visual page-
	// image reading earlier rounds used, a much more precise method).
	// Condensed to bare facts, not the manual's own copyrighted prose.
	// No confirmed mechanical use (same as Sign/Aspect) - real, sourced
	// flavor data, not wired into any gameplay effect.
	Correspondences string
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
	{Name: "ASMODEE", Title: "the Great Destroyer", Number: 122, Sign: "House of Mars", Aspect: "Basilisk", Ability: "Warning: be careful with Asmodee (no confirmed positive effect)", Charm: "Erlstone", Correspondences: "Colour green; plant Nettle; bows to red gems (no single named gem, unlike the other 3 Princes)"},
	{Name: "ASTAROT", Title: "the Spirit of Assemblage", Number: 1376, Sign: "Sign of Gemini", Aspect: "Legion", Ability: "Transports the player to a named location, if its name is known", Charm: "Sword", Correspondences: "Perfume Wormwood; favours Orchid and Magpie; gem Tourmaline"},
	{Name: "BELEZBAR", Title: "the Master of Flies", Number: 20, Sign: "Firmament of Stars", Aspect: "Deceit", Ability: "Reveals the true nature of objects", Charm: "Mantis", Correspondences: "Reveres Amaranth, Musk, and Locust; gem Turquoise"},
	{Name: "MAGOT", Title: "the Diviner", Number: 443, Sign: "Realm of Air", Aspect: "Baboon", Ability: "Reveals the whereabouts of any named object", Charm: "Sunflower", Correspondences: "Colour yellow; scent Galbanum; gems Topaz and Chalcedony"},
}

// Two more real numbers from the same grimoire section (round 110),
// not tied to any specific Demon: "the number of Magick is 11; but the
// number of the Great Abyss is 24" (paraphrased from the manual's
// prose). Not given a home here since no source found so far ties
// either number to a specific game mechanic, item, or room — recorded
// in this comment rather than invented a use for them, so a future
// round chasing either number knows it's real, sourced content, not
// something still to find in the manual.

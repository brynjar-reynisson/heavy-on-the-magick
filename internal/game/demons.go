package game

import (
	"fmt"
	"strings"
	"time"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/magic"
)

// demonPatienceLimit is how long (real time) a demon will wait, once
// correctly summoned (its Charm on the ground), for the player to name
// something "worthy" before losing patience - a real, user-recalled
// detail from watching a failed Asmodee invocation live in SpecEmu:
// "he didn't say anything worthy soon enough" (round 184). No source
// states an exact number, so 2 minutes (the user's own estimate, given
// alongside the nonsense-count rule below in the same request) is used
// as an honest, user-specified placeholder rather than an extracted
// fact.
const demonPatienceLimit = 2 * time.Minute

// demonNonsenseLimit is the other half of the same round-184 rule: how
// many consecutive unrecognized ("nonsense") objects/locations a
// summoned demon tolerates before losing patience regardless of how
// much time has passed.
const demonNonsenseLimit = 3

// demonSession tracks the real-time "patience" state for whichever
// demon is currently being addressed via a ground-charm-gated
// conversation command (ASTAROT/MAGOT/ASMODEE - see invokeWithPatience's
// own doc comment for why BELEZBAR isn't included).
type demonSession struct {
	demon    string
	since    time.Time
	nonsense int
}

// invokeWithPatience wraps a demon ability call (already confirmed to
// have its Charm correctly grounded) with round 184's real-time/
// nonsense-count patience rule. worthy reports whether THIS attempt
// named something the demon actually recognized (a real location/
// object it acted on); an unworthy ("nonsense") attempt counts toward
// demonNonsenseLimit, and either limit being reached sends the player
// to the furnace the same lethal way punishFailedInvoke does for a
// missing Talisman - the user's own words describe this as the same
// kind of consequence ("sends Axil to the Furnace"), not a milder one.
// A worthy attempt clears the session entirely, so a later invocation
// (of any demon) starts a fresh 2-minute/3-strike allowance.
//
// Only used for Astarot/Magot/Asmodee, each of which already has a
// real, existing "I don't recognize that" failure branch to treat as
// nonsense. Belezbar's own ability (BELEZBAR, <object>) always
// "succeeds" in some form (either a real disguise reveal or an honest
// "appears exactly what it seems") - there's no equivalent failure
// branch to call nonsense without inventing one, so it's left out of
// this mechanic rather than guessed at.
func (g *Game) invokeWithPatience(d magic.Demon, worthy bool, msg string) string {
	now := time.Now()
	session := g.demonSession
	if session == nil || session.demon != d.Name {
		session = &demonSession{demon: d.Name, since: now}
		g.demonSession = session
	} else if now.Sub(session.since) > demonPatienceLimit {
		g.demonSession = nil
		return g.sendToFurnace(fmt.Sprintf("%s's patience runs out - you took too long to say anything worthy.", d.Name))
	}
	if worthy {
		g.demonSession = nil
		return msg
	}
	session.nonsense++
	if session.nonsense >= demonNonsenseLimit {
		g.demonSession = nil
		return g.sendToFurnace(fmt.Sprintf("%s has heard enough nonsense.", d.Name))
	}
	return msg
}

// invoke handles INVOKE (Merphish keyword "I"), given the demon name the
// player named as cmd.Target. Recognizing the 4 confirmed magic.Demons,
// requiring each one's specific Charm, and each one's Ability are all
// real, sourced facts (see magic.Demons's doc comment) — invocation now
// actually succeeds when the right Charm is present. With no target,
// lists all 4 demons, their Charm requirements, and (round 114) their
// real occult Correspondences (magic.Demon.Correspondences, extracted
// round 110 from the manual's grimoire section but never actually shown
// anywhere in-game until now) — real, already-sourced data that had no
// way to be seen in-game before this. Round 115 closes the same gap for
// Number/Sign/Aspect, which were ALSO real, confirmed manual facts
// (present on every Demon since this project's very first magic.Demons
// commit) but — unlike Ability/Charm/Correspondences — had never been
// referenced by internal/game at all, found the same way round 114 was:
// checking each struct field against what actually reaches the player,
// not assuming "it's a field, so it must be used somewhere."
//
// Round 131 CORRECTION: the Charm gate previously checked g.hasItem
// (carried in inventory) — but a genuinely new source (World of
// Spectrum's separate plain-text instructions file, distinct from the
// PDF manual this project had already mined heavily) states the real
// mechanic precisely: "Place Ye the talisman on the ground and proceed
// with thy invocation from a distance." The Charm must be DROPPED in
// the room, not merely carried. Fixed to check g.roomHasItem instead;
// a player who IS carrying the Charm but hasn't dropped it gets an
// honest, specific hint (not punishFailedInvoke's furnace-room
// punishment, which is reserved for not having the Charm at all).
func (g *Game) invoke(target string) string {
	if target == "" {
		var b strings.Builder
		b.WriteString("Known demons and their required Talismans:\n")
		for _, d := range magic.Demons {
			fmt.Fprintf(&b, "%s, %s (Number %d, %s, Aspect: %s; needs: %s) - %s\n", d.Name, d.Title, d.Number, d.Sign, d.Aspect, d.Charm, d.Correspondences)
		}
		return strings.TrimRight(b.String(), "\n")
	}
	for _, d := range magic.Demons {
		if !strings.EqualFold(d.Name, target) {
			continue
		}
		if !g.roomHasItem(d.Charm) {
			return g.punishFailedInvoke(d)
		}
		if d.Ability != "" {
			return fmt.Sprintf("You invoke %s, %s! %s.", d.Name, d.Title, d.Ability)
		}
		return fmt.Sprintf("You invoke %s, %s! The ritual succeeds, though its exact effect isn't modeled yet.", d.Name, d.Title)
	}
	return "There is no demon by that name."
}

// punishFailedInvoke handles invoking a demon without its Charm
// properly in place - a real, confirmed punishment, not just a
// rejection message. The CRPG Addict's first-hand playthrough account
// (the same source that confirmed CALL's effect, round 125): "When you
// INVOKE them, you have to be holding their particular talisman--found
// within the dungeon--or they send you to a furnace room with no
// exits." If the active World has a real Furnace Room
// (world.CollodonsPile does, added round 126; so does world.Level1Grid,
// independently, as its own isolated A8 cell), the player is genuinely
// teleported there, same mechanism as game.astarotTeleport. Worlds
// without one (Level2-4Grid) fall back to the plain rejection message
// rather than fail - an honest scope limit, not a fabricated
// destination.
//
// ROUND 183 correction: the user, having independently finished the
// same third gameplay video, recalled that a failed invocation isn't
// just a banishment - "sends Axil to the furnace room, where he dies
// horribly" (the video itself never shows this, since walkthrough
// footage only shows success paths). Death is now real, tied to
// actually REACHING the furnace room - if the active World has none
// (the honest scope limit above), there's no confirmed mechanism for
// death either, so the plain rejection message stays non-lethal. This
// now covers BOTH real failure shapes uniformly: having no Charm at
// all, and carrying it without grounding it (round 131's own "must be
// on the ground, not merely carried" correction) - the user separately
// confirmed the OTHER 3 demons kill Axil the same way for either case,
// not just Asmodee. carrying distinguishes the message's own wording
// (still accurate either way) without changing the outcome.
func (g *Game) punishFailedInvoke(d magic.Demon) string {
	carrying := g.hasItem(d.Charm)
	var msg string
	if carrying {
		msg = fmt.Sprintf("You begin the ritual to invoke %s, %s... but the %s must be on the ground, not in your hand. The ritual backfires!", d.Name, d.Title, d.Charm)
	} else {
		msg = fmt.Sprintf("You begin the ritual to invoke %s, %s... but you have no suitable Talisman (a %s). The ritual backfires!", d.Name, d.Title, d.Charm)
	}
	return g.sendToFurnace(msg)
}

// sendToFurnace implements the shared real, confirmed "furnace room
// with no exits" punishment (The CRPG Addict's blog, round 126;
// extended to a real death by the user, round 183) - factored out of
// punishFailedInvoke so invokeWithPatience's own real-time/nonsense-
// count punishment (round 184) can reuse the exact same mechanism with
// its own, different prefix rather than a fabricated separate outcome.
func (g *Game) sendToFurnace(prefix string) string {
	if id, ok := g.World.FindRoomByName("Furnace Room"); ok {
		g.World.Teleport(id)
		g.Player.Stamina = 0
		return prefix + " You are flung into a furnace room with no exits. You die horribly! (GAME OVER)"
	}
	return prefix + " You are flung into a furnace room with no exits."
}

// astarotTeleport handles the confirmed real conversation-form command
// "ASTAROT, <location>" - the in-game hint screen's own literal example
// is "ASTAROT, WOLFDORP" (see game.help's doc comment and
// parser.Parse's package doc comment, which has carried this exact
// example since the two-grammar-form writeup). magic.Demons's Astarot
// entry already confirms the underlying ability ("Transports the player
// to a named location, if its name is known") and Charm ("Sword") - this
// is the first place that ability is actually implemented, using the
// same Charm-gating convention as bare INVOKE (round 131: the Charm
// must be on the ground, per World of Spectrum's plain-text
// instructions file - see invoke's doc comment). The location name is
// looked up against the CURRENT world's real Room names (Wolfdorp itself
// is a real, already-shipped CollodonsPile room), so this only reaches
// places that genuinely exist in whichever world is active - no
// fabricated destinations.
//
// ROUND 183: a failed Charm check now routes through
// punishFailedInvoke - same real, lethal furnace-room punishment as
// bare INVOKE, per the user's own direct recollection (see
// punishFailedInvoke's doc comment) - replacing this port's own
// earlier, harmless rejection/hint messages.
func (g *Game) astarotTeleport(location string) string {
	const astarotName, astarotTitle, astarotCharm = "Astarot", "the Spirit of Assemblage", "Sword"
	d := magic.Demon{Name: astarotName, Title: astarotTitle, Charm: astarotCharm}
	if !g.roomHasItem(astarotCharm) {
		return g.punishFailedInvoke(d)
	}
	id, ok := g.World.FindRoomByName(location)
	if !ok {
		// "No such place." is the game's own exact rejection text for
		// an unrecognized ASTAROT destination, confirmed via a direct
		// frame-by-frame review (2-second sampling interval) of a
		// fourth pass over the third gameplay video (round 180) -
		// replacing this port's own earlier invented wording. Round
		// 184: an unrecognized location is exactly the kind of
		// "nonsense" invokeWithPatience's own patience rule counts.
		return g.invokeWithPatience(d, false, "No such place.")
	}
	dest := g.World.Rooms[id]
	g.World.Teleport(id)
	// "Best place for you" is the game's own exact confirmed success
	// text - seen 3 separate times in the same footage (round 180:
	// "ASTAROT, SLYMOLE", "ASTAROT, LICHGATE", and the original "ASTAROT,
	// WOLFDORP" sequence), each ending in exactly this phrase -
	// replacing this port's own earlier invented "In an instant, you
	// are transported to X" wording. dest.Name is kept in the message
	// too (not part of the confirmed text, but useful, real player
	// feedback this port already relied on elsewhere and no source
	// contradicts including).
	msg := fmt.Sprintf("You invoke %s, %s! Best place for you: %s.", astarotName, astarotTitle, dest.Name)
	return g.invokeWithPatience(d, true, msg)
}

// magotLocate handles the conversation-form command "MAGOT, <object>".
// No source gives a literal "MAGOT, X" example the way the hint screen
// gives "ASTAROT, WOLFDORP" - but the manual's own confirmed grammar
// ("name, object") is general, not restricted to the 2 demons it
// happens to illustrate, and magic.Demons's Magot entry already
// confirms both the underlying ability ("Reveals the whereabouts of any
// named object") and Charm ("Sunflower"). Applying the same
// TARGET-comma-VERB grammar and Charm-gating convention already
// established for Astarot's teleport is a natural, honestly-flagged
// inference, not a fabricated mechanic - the same confidence level this
// project already applies to vocabulary synonyms like TAKE/LIFT.
// Checks the player's own inventory first (an item they're already
// carrying isn't "located" anywhere else), then every room in the
// CURRENT world - so, like Astarot's teleport, this only ever reports
// real placements that genuinely exist in whichever world is active.
// Echoes the item's real stored casing in both branches, not the
// player's raw uppercased typed target - the same casing-honesty fix
// already applied once before to examine().
//
// ROUND 183: a failed Charm check now routes through
// punishFailedInvoke - see astarotTeleport's own doc comment for why.
func (g *Game) magotLocate(object string) string {
	const magotName, magotTitle, magotCharm = "Magot", "the Diviner", "Sunflower"
	d := magic.Demon{Name: magotName, Title: magotTitle, Charm: magotCharm}
	if !g.roomHasItem(magotCharm) {
		return g.punishFailedInvoke(d)
	}
	for _, item := range g.Player.Items {
		if strings.EqualFold(item, object) {
			msg := fmt.Sprintf("You invoke %s, %s! No need - you already carry the %s yourself.", magotName, magotTitle, item)
			return g.invokeWithPatience(d, true, msg)
		}
	}
	for _, room := range g.World.Rooms {
		for _, item := range room.Items {
			if strings.EqualFold(item, object) {
				msg := fmt.Sprintf("You invoke %s, %s! The %s lies in %s.", magotName, magotTitle, item, room.Name)
				return g.invokeWithPatience(d, true, msg)
			}
		}
	}
	// Round 184: an object Magot can't locate anywhere is exactly the
	// kind of "nonsense" invokeWithPatience's own patience rule counts.
	msg := fmt.Sprintf("You invoke %s, %s! %s senses no such object anywhere nearby.", magotName, magotTitle, magotName)
	return g.invokeWithPatience(d, false, msg)
}

// asmodeeDestroy handles the conversation-form command "ASMODEE,
// <object>" (round 162). Asmodee's Ability sat as just a warning for
// many rounds ("Warning: be careful with Asmodee... no confirmed
// positive effect") — a 4th, genuinely different source (Hardcore
// Gaming 101's article on this game) states plainly "Asmodee destroys
// any object you ask of him," a real, positive, functional ability
// after all (see magic.Demons's own doc comment). No source gives a
// literal "ASMODEE, X" example, but the same general "name, object"
// grammar and Charm-gating convention already established for Astarot/
// Magot applies directly - the identical confidence tier magotLocate's
// own doc comment already claims for itself. Searches the player's own
// inventory first, then every room in the CURRENT world (the same
// search order magotLocate uses to LOCATE an object - this DESTROYS
// it instead, permanently removing it from wherever it's found), and
// echoes the item's real stored casing, not the player's raw
// uppercased typed target - same casing-honesty convention as
// magotLocate/examine.
//
// ROUND 183: a failed Charm check now routes through
// punishFailedInvoke - see astarotTeleport's own doc comment for why.
func (g *Game) asmodeeDestroy(object string) string {
	const asmodeeName, asmodeeTitle, asmodeeCharm = "Asmodee", "the Great Destroyer", "Erlstone"
	d := magic.Demon{Name: asmodeeName, Title: asmodeeTitle, Charm: asmodeeCharm}
	if !g.roomHasItem(asmodeeCharm) {
		return g.punishFailedInvoke(d)
	}
	// ROUND 183: "ASMODEE, DOOR" - user-recalled directly from finishing
	// the same third gameplay video independently ("Asmodee opening the
	// door"), matching round 180's own earlier frame ("ASMODEE, DOOR...
	// ASMODEE DESTROYS... THE DOOR TO THE TOMB"). A locked door is a
	// structural room fact (DoorPasswords/TollItem/Guards), not a named
	// Item the generic search below would ever find - destroying it
	// means clearing whatever real lock is on the CURRENT room, the
	// same state game.move's own locked-door check (round 181) gates
	// on.
	//
	// ROUND 184 clarification: the user's own reading of the video is
	// that this only worked "because the ruby stone is present (that I
	// assume is the talisman for Asmodee)... all other locations would
	// be without the talisman" - i.e. the door-destroy is gated on the
	// Charm exactly like every other Asmodee action, not restricted to
	// one specific unconfirmed room name ("the Tomb" is never confirmed
	// as a real, placed room anywhere in this project's data, only
	// mentioned in passing in the video's own flavor text about what
	// was BEHIND the door). The `roomHasItem(asmodeeCharm)` check above
	// already enforces exactly this - any room without Erlstone
	// grounded in it already fails before reaching this line - so no
	// further restriction is added here, only this clarifying note.
	if strings.EqualFold(object, "DOOR") {
		room := g.World.CurrentRoom()
		if room != nil && (len(room.DoorPasswords) > 0 || room.TollItem != "" || room.Guards) {
			room.DoorPasswords = nil
			room.TollItem = ""
			room.Guards = false
			msg := fmt.Sprintf("You invoke %s, %s! The door crumbles to nothing.", asmodeeName, asmodeeTitle)
			return g.invokeWithPatience(d, true, msg)
		}
		msg := fmt.Sprintf("You invoke %s, %s! %s finds no such object to destroy.", asmodeeName, asmodeeTitle, asmodeeName)
		return g.invokeWithPatience(d, false, msg)
	}
	for i, item := range g.Player.Items {
		if strings.EqualFold(item, object) {
			g.Player.Items = append(g.Player.Items[:i], g.Player.Items[i+1:]...)
			msg := fmt.Sprintf("You invoke %s, %s! The %s you carried crumbles to nothing.", asmodeeName, asmodeeTitle, item)
			return g.invokeWithPatience(d, true, msg)
		}
	}
	for _, room := range g.World.Rooms {
		for i, item := range room.Items {
			if strings.EqualFold(item, object) {
				room.Items = append(room.Items[:i], room.Items[i+1:]...)
				msg := fmt.Sprintf("You invoke %s, %s! The %s in %s crumbles to nothing.", asmodeeName, asmodeeTitle, item, room.Name)
				return g.invokeWithPatience(d, true, msg)
			}
		}
	}
	// Round 184: an object Asmodee can't find anywhere is exactly the
	// kind of "nonsense" invokeWithPatience's own patience rule counts.
	msg := fmt.Sprintf("You invoke %s, %s! %s finds no such object to destroy.", asmodeeName, asmodeeTitle, asmodeeName)
	return g.invokeWithPatience(d, false, msg)
}

// belezbarDisguises maps a real, sourced "disguised" item name to its
// true identity. Currently one confirmed entry: the fan-made numbered
// map poster's own key list (internal/world.NumberedRoomContents)
// gives entry #59 as "Pebble (disguised Erlstone)" - genuinely distinct
// from the OTHER, plain Pebbles at neighboring numbered cells (#57/#58/
// #60/#62/#63, none of which carry a "disguised" qualifier) - exactly
// the kind of real, confirmed unmasking Belezbar's Ability ("Reveals
// the true nature of objects") describes. Erlstone itself is already
// placed directly (not disguised) at Methos (round 107) - that
// placement is unaffected; this map exists for a room that might one
// day place the disguised "Pebble" form specifically, the same "real
// mechanic, not yet reachable in shipped data" pattern already used
// for TollItem/Fire/Guards/SwapItem before their own first placements.
//
// ROUND 178 - a real, unresolved second data point: a frame-by-frame
// review of a third gameplay video shows a DIFFERENT real Pebble
// (elsewhere in the dungeon, room not identified) whose real Belezbar
// reveal reads "IT'S INSCRIBED WITH THE WORD LICHGATE" - not Erlstone.
// "Lichgate" is itself a real, separately-confirmed zone name
// (known_room_names.go). Since Pebble is evidently a generic item name
// with multiple distinct real instances/rooms (the numbered map's own
// key list already shows several plain, undisguised Pebbles alongside
// #59's disguised one), this project's current single
// name-to-disguise mapping can't represent both real reveals for the
// same generic name without an exact room to distinguish them - left
// as an honest, documented conflict rather than arbitrarily picking
// one over the other or guessing a room.
var belezbarDisguises = map[string]string{
	"Pebble": "Erlstone",
}

// belezbarReveal handles the conversation-form command "BELEZBAR,
// <object>" (round 164) - the last of the 4 confirmed demons to get a
// real, functional ability wired up (Astarot teleports, Magot locates,
// Asmodee destroys - see their own doc comments). No source gives a
// literal "BELEZBAR, X" example, but the same general "name, object"
// grammar and Charm-gating convention already established for the
// other 3 applies directly. Checks belezbarDisguises for a real,
// sourced disguise; anything else gets an honest "appears to be
// exactly what it seems" - not a fabricated secret identity for every
// object, only the one genuinely confirmed case. The success wording
// ("<Object> is inscribed with the word <Word>.") matches the game's
// own real, confirmed phrasing - a direct frame-by-frame review of a
// third gameplay video shows the exact live text "BELEZBAR: PEBBLE" /
// "IT'S INSCRIBED WITH THE WORD LICHGATE" (round 178) - replacing this
// port's own earlier invented "it is really a X" wording.
//
// Unlike Astarot/Magot/Asmodee's Charms (Sword/Sunflower/Erlstone, all
// placed in world.CollodonsPile itself), Belezbar's Charm (Mantis) is
// ONLY placed in world.Level3Grid (round 103) - a real, pre-existing
// scope limit this project has flagged since round 108 ("the only one
// of the 4 demons' Charms still unreachable in default-mode play"),
// not something this round changes. So this command is only reachable
// in a real playthrough via `go run ./cmd/hotm -level3grid`, not the
// default game.New() - verified live that way, not in default mode.
//
// ROUND 183: a failed Charm check now routes through
// punishFailedInvoke - see astarotTeleport's own doc comment for why.
//
// ROUND 184: deliberately NOT wrapped in invokeWithPatience the way
// Astarot/Magot/Asmodee now are - see that function's own doc comment.
// Belezbar has no real "I don't recognize that" failure branch to
// treat as nonsense (every object gets SOME real response), so
// applying the same time/nonsense-count rule here would mean inventing
// a new failure mode this project has no source for.
func (g *Game) belezbarReveal(object string) string {
	const belezbarName, belezbarTitle, belezbarCharm = "Belezbar", "the Master of Flies", "Mantis"
	if !g.roomHasItem(belezbarCharm) {
		return g.punishFailedInvoke(magic.Demon{Name: belezbarName, Title: belezbarTitle, Charm: belezbarCharm})
	}
	for disguised, real := range belezbarDisguises {
		if strings.EqualFold(disguised, object) {
			return fmt.Sprintf("You invoke %s, %s! %s is inscribed with the word %s.", belezbarName, belezbarTitle, disguised, real)
		}
	}
	return fmt.Sprintf("You invoke %s, %s! The %s appears to be exactly what it seems.", belezbarName, belezbarTitle, object)
}

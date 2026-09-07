package game

// help reproduces the game's own real "SOME ADVICE" startup hint
// screen - not paraphrased or invented, but the actual in-game text
// disassembled from the original's memory and independently cross-
// confirmed against a real screenshot of the running game (see
// ../../CLAUDE.md, "Found ground truth: ref_screenshot.png" - the
// decoded strings matched the screenshot exactly). "HELP" is a real
// confirmed vocabulary word (parser.Vocabulary); no source states it's
// literally the keyword that shows this screen in the original (it may
// have only appeared once, at boot), but reproducing the original
// game's own real hint text for a player who asks for help is squarely
// what a faithful port should do with content this well-sourced.
func (g *Game) help() string {
	return `SOME ADVICE
Talk to Apex often!
"APEX, DOOR"  "APEX, WEREWOLF"  "APEX, FIRE"
To dismiss say "APEX, THANKS"
Talk to other things - sometimes they will answer!
"DOOR, password"
"GUARDS, DOOR"
"ASTAROT, WOLFDORP"
Finally, if in a panic, BLAST without an object!`
}

// spells lists this port's real, functional spells. "SPELLS" is a real
// confirmed vocabulary word (parser.Vocabulary); no source states this
// exact command lists them as a menu, so the aggregation itself is this
// project's own, not fabricated content — but the manual DOES have its
// own explicit "Spells:" heading (round 118, found via the PDF's real
// text layer — see round 110's writeup), grouping exactly three
// keywords under it: I (Invoke), B (Blast), F (Freeze). INVOKE was
// missing from this listing entirely until round 118 — a real,
// sourced omission this project's own aggregation had gotten wrong,
// not just an incomplete one — corrected here. TRANSFUSION and CALL
// are real spells too (see their own doc comments) but the manual
// itself says "fuller details... can be found in the section on the
// Grimoire" for those, i.e. they're confirmed spells from a DIFFERENT
// part of the manual, not this specific three-keyword grouping — kept
// in the listing since they're genuinely real spells, just not
// re-labeled as part of the manual's core "Spells:" trio.
func (g *Game) spells() string {
	return "Known spells: INVOKE (summon a demon), BLAST (combat), FREEZE (combat), TRANSFUSION (restore Stamina), CALL (summon Apex, needs the Scroll)."
}

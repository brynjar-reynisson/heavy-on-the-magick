// Command hotm is the entry point for the Go port of "Heavy on the Magick".
//
// Movement, real room names/connections (world.CollodonsPile, sourced from
// a published walkthrough and cross-confirmed against strings found in the
// game's own memory — see its doc comment), and the explored-map feature
// are all playable today. Room description TEXT (as opposed to names) and
// most verb resolution (combat, spellcasting, NPC talk) are not yet ported
// — see CLAUDE.md at the repo root for the full reverse-engineering status.
//
// The -level1grid flag switches to game.NewLevel1Exploration instead: a
// real, validated 64-cell grid extracted from the actual Level 1 map (see
// world.Level1Grid's doc comment) - far more rooms than CollodonsPile's
// walkthrough-sourced set, though not yet merged with it (different,
// not-yet-reconciled addressing schemes - see CLAUDE.md). -level2grid is
// its Level 2 counterpart (game.NewLevel2Exploration, world.Level2Grid) -
// a real, fully-connected 50-cell grid.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/game"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
)

func main() {
	level1Grid := flag.Bool("level1grid", false, "play the extracted Level 1 grid (64 real cells) instead of CollodonsPile")
	level2Grid := flag.Bool("level2grid", false, "play the extracted Level 2 grid (50 real, fully-connected cells) instead of CollodonsPile")
	level3Grid := flag.Bool("level3grid", false, "play the extracted Level 3 grid (47 real cells: 41 fully-connected plus 6 isolated rooms) instead of CollodonsPile")
	level4Grid := flag.Bool("level4grid", false, "play the extracted Level 4 grid (23 real cells: 17 fully-connected plus 6 isolated rooms) instead of CollodonsPile")
	flag.Parse()

	var g *game.Game
	switch {
	case *level1Grid:
		fmt.Println("Heavy on the Magick — Go port (Level 1 grid exploration mode)")
		fmt.Println("Real, validated 64-cell Level 1 map (world.Level1Grid) — not merged with CollodonsPile yet.")
		g = game.NewLevel1Exploration()
	case *level2Grid:
		fmt.Println("Heavy on the Magick — Go port (Level 2 grid exploration mode)")
		fmt.Println("Real, validated, fully-connected 50-cell Level 2 map (world.Level2Grid) — not merged with CollodonsPile yet.")
		g = game.NewLevel2Exploration()
	case *level3Grid:
		fmt.Println("Heavy on the Magick — Go port (Level 3 grid exploration mode)")
		fmt.Println("Real 47-cell Level 3 map (world.Level3Grid): 41 fully-connected plus 6 isolated rooms — not merged with CollodonsPile yet.")
		g = game.NewLevel3Exploration()
	case *level4Grid:
		fmt.Println("Heavy on the Magick — Go port (Level 4 grid exploration mode)")
		fmt.Println("Real 23-cell Level 4 map (world.Level4Grid): 17 fully-connected plus 6 isolated rooms — not merged with CollodonsPile yet.")
		g = game.NewLevel4Exploration()
	default:
		fmt.Println("Heavy on the Magick — Go port")
		fmt.Println("Real room names/map from Collodon's Pile; room description text and most verbs not yet ported.")
		g = game.New()
	}
	fmt.Println("Commands: LOOK, MAP, a compass direction (e.g. NORTH, SOUTH-EAST), BLAST, or QUIT.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "QUIT" {
			break
		}
		cmd := parser.Parse(line)
		fmt.Println(g.Handle(cmd))
	}
}

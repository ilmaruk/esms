package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ilmaruk/esms/internal"
	"github.com/ilmaruk/esms/internal/models"
	"github.com/ilmaruk/esms/internal/plugins/persistence/file"
)

type RosterFetcher interface {
	Fetch(uuid.UUID) (models.Roster, error)
}

type TeamsheetStorer interface {
	Store(models.Teamsheet) error
}

// // ESMS - Electronic Soccer Management Simulator
// // Copyright (C) <1998-2005>  Eli Bendersky
// //
// // This program is free software, licensed with the GPL (www.fsf.org)
// //
// #include <cstdio>
// #include <cstdlib>
// #include <cstring>
// #include <set>
// #include <cctype>
// #include <ctime>
// #include <cassert>
// #include "tsc.h"
// #include "rosterplayer.h"
// #include "util.h"
// #include "config.h"

type Config struct {
}

func main() {
	teamsheet := models.Teamsheet{
		ID:      uuid.New(),
		Players: make([]models.TeamsheetPlayer, 25),
	}

	// theConfig := Config{}

	// Either we get no arguments, and then we ask to enter
	// the filename and formation manually, or we get 2
	// arguments - filename and formation
	if len(os.Args) < 3 {
		fmt.Println("Usage:\n\ntsc [<rosterID> <formation & tactic> [0]]")
		internal.MyExit(false, 0)
	}
	rosterID := os.Args[1]
	formation := os.Args[2]

	fetcher := file.NewJSONRosterFetcher("./data")
	roster, err := fetcher.Fetch(uuid.MustParse(rosterID))
	if err != nil {
		panic(err)
	}

	// num_subs := the_config().get_int_config("NUM_SUBS", 7);
	numSubs := 7

	// The number of subs is not constant, therefore there is
	// a need for some smart assignment. The following array
	// sets the positions of thr first 5 subs, and then iterates
	// cyclicly. For example, if there are 2 subs allowed,
	// their positions will be GK (mandatory 1st !) and MF
	// If 7: GK, DF, MF, DF, FW, MF, DF
	//                              ^
	//                              cyclic repetition begins
	//
	// const char* sub_position[] = {"DFC", "MFC", "DFC", "FWC", "MFC"};
	subPosition := []string{"DFC", "MFC", "DFC", "FWC", "MFC"}

	// Iterates (cyclicly) over positions of subs,
	//
	var subPosIter = 0

	if len(roster.Players) < 11+numSubs {
		panic("Error: not enough players in roster")
	}

	var dfs, mfs, fws int
	var tactic string

	parseFormation(formation, &dfs, &mfs, &fws, &tactic)

	// Calculate indices of the last defender and the last midfielder
	//
	lastDF := dfs + 1
	lastMF := dfs + mfs + 1

	// Pick the players
	//
	// First, the best shot stopper is picked as a GK, then
	// others are picker according to the schedule of sub_position
	// as described above
	//

	// This will keep us from picking the same players more than once
	//
	chosenPlayers := make(map[uuid.UUID]bool)

	for i := 1; i <= 11; i++ {
		if i == 1 {
			teamsheet.Players[i].Pos = "GK"
		} else if i >= 2 && i <= lastDF {
			teamsheet.Players[i].Pos = "DFC"
		} else if i > lastDF && i <= lastMF {
			teamsheet.Players[i].Pos = "MFC"
		} else if i > lastMF && i <= 11 {
			teamsheet.Players[i].Pos = "FWC"
		}
	}

	// set the best GK for N.1 position
	//
	teamsheet.Players[1].Player = chooseBestPlayer(roster.Players, chosenPlayers, stGetter)
	chosenPlayers[teamsheet.Players[1].Player.ID] = true

	// From now on, j is the index for players in the teamsheet
	//

	// Set the starting defenders
	//
	for j := 2; j <= lastDF; j++ {
		teamsheet.Players[j].Player = chooseBestPlayer(roster.Players, chosenPlayers, tkGetter)
		chosenPlayers[teamsheet.Players[j].Player.ID] = true
	}

	// Set the starting midfielders
	//
	for j := lastDF + 1; j <= lastMF; j++ {
		teamsheet.Players[j].Player = chooseBestPlayer(roster.Players, chosenPlayers, psGetter)
		chosenPlayers[teamsheet.Players[j].Player.ID] = true
	}

	// Set the starting forwards
	//
	for j := lastMF + 1; j <= 11; j++ {
		teamsheet.Players[j].Player = chooseBestPlayer(roster.Players, chosenPlayers, shGetter)
		chosenPlayers[teamsheet.Players[j].Player.ID] = true
	}

	// Set the substitute GK
	//
	teamsheet.Players[12].Player = chooseBestPlayer(roster.Players, chosenPlayers, stGetter)
	teamsheet.Players[12].Pos = "GK"
	chosenPlayers[teamsheet.Players[12].Player.ID] = true

	var theBest models.RosterPlayer

	for j := 13; j <= numSubs+11; j++ {
		// What position should the current sub be on ?
		//
		if subPosition[subPosIter] == "DFC" {
			theBest = chooseBestPlayer(roster.Players, chosenPlayers, tkGetter)
		} else if subPosition[subPosIter] == "MFC" {
			theBest = chooseBestPlayer(roster.Players, chosenPlayers, psGetter)
		} else if subPosition[subPosIter] == "FWC" {
			theBest = chooseBestPlayer(roster.Players, chosenPlayers, shGetter)
		} else {
			panic("0")
		}

		teamsheet.Players[j].Player = theBest
		teamsheet.Players[j].Pos = subPosition[subPosIter]
		chosenPlayers[teamsheet.Players[j].Player.ID] = true
		subPosIter = (subPosIter + 1) % 5
	}

	storer := file.NewJSONTeamsheetStorer("./data")
	if err := storer.Store(teamsheet); err != nil {
		panic(err)
	}
}

func stGetter(player models.RosterPlayer) float64 {
	return float64(player.St*player.Fitness) / 100
}

func tkGetter(player models.RosterPlayer) float64 {
	return float64(player.Tk*player.Fitness) / 100
}

func psGetter(player models.RosterPlayer) float64 {
	return float64(player.Ps*player.Fitness) / 100
}

func shGetter(player models.RosterPlayer) float64 {
	return float64(player.Sh*player.Fitness) / 100
}

// / Gets the best player on some position from an array of roster players.
// /
// / players 		- the array of players
// / chosen_players 	- a set of already chosen players (those won't be chosen again)
// / skill 			- pointer to a function receiving a player and returning the skill by
// / 				  which "best" is judged.
// /
// / Returns the chosen player's name. Note: chosen_players is not modified !
// /
func chooseBestPlayer(players []models.RosterPlayer, chosenPlayers map[uuid.UUID]bool, skill func(models.RosterPlayer) float64) models.RosterPlayer {
	var bestSkill float64 = -1
	var theBest models.RosterPlayer

	for _, player := range players {
		if _, ok := chosenPlayers[player.ID]; !ok {
			if skill(player) > bestSkill && player.Injury == 0 && player.Suspension == 0 {
				bestSkill = skill(player)
				theBest = player
			}
		}
	}

	// if theBest == nil {
	// 	panic("theBest not set")
	// }

	return theBest
}

//     sprintf(teamsheetname, "%ssht.txt", teamname);

//     teamsheetfile = fopen(teamsheetname, "w");

//     // Start filling the team sheet with the roster name and the
//     // tactic
//     //
//     fprintf(teamsheetfile, "%s\n", teamname);
//     fprintf(teamsheetfile, "%s\n", tactic);

//     /* Print all the players and their position */
//     for (i = 1; i <= 11 + num_subs; i++)
//     {
//         fprintf(teamsheetfile, "\n%s %s", t_player[i].pos.c_str(), t_player[i].name.c_str());

//         if (i == 11)
//             fprintf(teamsheetfile, "\n");
//     }

//     /* Print the penalty kick taker (player number last_mf + 1) */
//     fprintf(teamsheetfile, "\n\nPK: %s\n\n", t_player[last_mf + 1].name.c_str());

//     printf("%s created successfully\n", teamsheetname);

//     fclose(teamsheetfile);

//     MY_EXIT(0);

//     return 0;
// }

// // Remove trailing newline
// //
// void chomp(char* str)
// {
//     int len = strlen(str);

//     if (str[len-1] == '\n')
//         str[len-1] = '\0';
// }

// Parses the formation line, finds out how many defenders,
// midfielders and forwards to pick, and the tactic to use,
// performs error checking
//
// For example: 442N means 4 DFs, 4 MFs, 2 FWs, playing N
func parseFormation(formation string, dfs, mfs, fws *int, tactic *string) {
	if len(formation) != 4 {
		fmt.Println("The formation string must be exactly 4 characters long")
		fmt.Println("For example: 442N")
		internal.MyExit(false, 0)
	}

	// Random formation ?
	//
	if strings.HasPrefix(formation, "rnd") {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

		// between 3 and 5
		*dfs = int(3 + rnd.Uint64()%3)

		// if there are 5 dfs, max of 4 mfs
		if *dfs == 5 {
			*mfs = int(1 + rnd.Uint64()%4)
		} else {
			// 5 mfs is also possible
			*mfs = int(1 + rnd.Uint64()%5)
		}

		*fws = 10 - *dfs - *mfs

		*tactic = formation[3:]

		return
	}

	*dfs, _ = strconv.Atoi(string(formation[0]))
	*mfs, _ = strconv.Atoi(string(formation[1]))
	*fws, _ = strconv.Atoi(string(formation[2]))

	*tactic = formation[3:]

	verifyPositionRange(*dfs)
	verifyPositionRange(*mfs)
	verifyPositionRange(*fws)

	if (*dfs + *mfs + *fws) != 10 {
		fmt.Println("The number of players on all positions added together must be 10")
		fmt.Println("For example: 442N")
		internal.MyExit(false, 0)
	}
}

func verifyPositionRange(n int) {
	if n < 1 || n > 8 {
		fmt.Println("The number of players on each position must be between 1 and 8")
		fmt.Println("For example: 442N")
		internal.MyExit(false, 0)
	}
}

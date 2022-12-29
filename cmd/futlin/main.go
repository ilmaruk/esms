package main

import (
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/ilmaruk/esms/internal"
	"github.com/ilmaruk/esms/internal/esms"
	"github.com/ilmaruk/esms/internal/models"
)

const HALF_LENGTH = 45
const NUM_PLAYERS = 11

const (
	DID_SHOT int = iota
	DID_FOUL
	DID_TACKLE
	DID_ASSIST
)

var teamIndices [2]int = [2]int{0, 1}

func defendingTeamIndex(off int) int {
	if off == 0 {
		return 1
	}
	return 0
}

func main() {
	rand.Seed(time.Now().UnixMicro())

	teamStatsTotalEnabled := false

	var teams [2]models.Team
	var err error

	initTeamsData()
	teams[0], err = esms.ReadTeam("./bin/foo.txt", "./bin/foosht.txt")
	if err != nil {
		panic(err)
	}
	teams[1], err = esms.ReadTeam("./bin/bar.txt", "./bin/barsht.txt")
	if err != nil {
		panic(err)
	}

	// For each half
	//
	for halfStart := 1; halfStart < 2*HALF_LENGTH; halfStart += HALF_LENGTH {
		var half int
		if halfStart == 1 {
			half = 1
		} else {
			half = 2
		}
		lastMinuteOfHalf := halfStart + HALF_LENGTH - 1
		inInjTime := false

		// Play the game minutes of this half
		//
		// last_minute_of_half will be increased by inj_time_length in
		// the end of the half
		//
		formalMinute := halfStart
		for minute := halfStart; minute <= lastMinuteOfHalf; minute++ {
			cleanInjCardIndicators()
			recalculateTeamsData()

			// For each team
			//
			for _, att := range teamIndices {
				def := defendingTeamIndex(att)

				// Calculate different events
				//
				ifShot(teams[att], teams[def])
				ifFoul(att)
				randomInjury(att)

				// scoreDiff := team[j].score - team[!j].score
				checkConditionals(att)
			}

			// fixme ?
			if teamStatsTotalEnabled {
				if minute == 1 || minute%10 == 0 {
					addTeamStatsTotal()
				}
			}

			if !inInjTime {
				formalMinute++

				updatePlayersMinuteCount()
			}

			if minute == lastMinuteOfHalf && !inInjTime {
				inInjTime = true

				// shouldn't have been increased, but we only know about
				// this now
				formalMinute--

				injTimeLength := howMuchInjTime()
				lastMinuteOfHalf += injTimeLength

				internal.PrintCommentary(os.Stdout, internal.COMM_INJURYTIME, injTimeLength)
			}
		}

		inInjTime = false

		if half == 1 {
			internal.PrintCommentary(os.Stdout, internal.COMM_HALFTIME)
		} else if half == 2 {
			internal.PrintCommentary(os.Stdout, internal.COMM_FULLTIME)
		}
	}

	calcAbility()
}

func initTeamsData() {

}

func cleanInjCardIndicators()   {}
func recalculateTeamsData()     {}
func ifFoul(w int)              {}
func randomInjury(w int)        {}
func checkConditionals(w int)   {}
func addTeamStatsTotal()        {}
func updatePlayersMinuteCount() {}
func calcAbility()              {}

func randomp(p int) bool {
	return rand.Intn(1000) < p
}

// Called on each minute to handle a scoring chance of team
// a for this minute.
//
func ifShot(att, def models.Team) {
	// int shooter, assister, tackler;
	// int chance_tackled;
	// int chance_assisted = 0;
	var assister, shooter int
	// var chanceAssisted bool

	// Did a scoring chance occur ?
	//
	if randomp(int(att.ShotProb)) {
		// There's a 0.75 probability that a chance was assisted, and
		// 0.25 that it's a solo
		//
		if randomp(7500) {
			assister = whoDidIt(att, DID_ASSIST)
			// chanceAssisted = true

			// shooter = who_got_assist(a, assister);

			internal.PrintCommentary(os.Stdout, internal.ASSISTEDCHANCE, "[MIN]", att.Name, att.Players[assister].Name, att.Players[shooter].Name)
			// team[a].player[assister].keypasses++;
		} else {
			shooter = whoDidIt(att, DID_SHOT)
			// chanceAssisted = false

			internal.PrintCommentary(os.Stdout, internal.CHANCE, "[MIN]", att.Name, att.Players[shooter].Name)
		}

		// chance_tackled = (int) (4000.0*((team[!a].team_tackling*3.0)/(team[a].team_passing*2.0+team[a].team_shooting)));

		// /* If the chance was tackled */
		// if (randomp(chance_tackled))
		// {
		//     tackler = who_did_it(!a, DID_TACKLE);
		//     team[!a].player[tackler].tackles++;

		//     fprintf(comm, "%s", the_commentary().rand_comment("TACKLE", team[!a].player[tackler].name).c_str());
		// }
		// else /* Chance was not tackled, it will be a shot on goal */
		// {
		//     fprintf(comm, "%s", the_commentary().rand_comment("SHOT", team[a].player[shooter].name).c_str());
		//     team[a].player[shooter].shots++;

		//     if (if_ontarget(a, shooter))
		//     {
		//         team[a].finalshots_on++;
		//         team[a].player[shooter].shots_on++;

		//         if (if_goal(a, shooter))
		//         {
		//             fprintf(comm, "%s", the_commentary().rand_comment("GOAL").c_str());

		//             if (!is_goal_cancelled())
		//             {
		//                 team[a].score++;

		//                 // If the assister was the shooter, there was no
		//                 // assist, but a simple goal.
		//                 //
		//                 if (chance_assisted && (assister != shooter))
		//                     team[a].player[assister].assists++; /* For final stats */

		//                 team[a].player[shooter].goals++;
		//                 team[!a].player[team[!a].current_gk].conceded++;

		//                 fprintf(comm, "\n          ...  %s %d-%d %s ...",
		//                         team[0].name,
		//                         team[0].score,
		//                         team[1].score,
		//                         team[1].name);

		//                 report_event* an_event = new report_event_goal(team[a].player[shooter].name,
		//                                          team[a].name, formal_minute_str().c_str());

		//                 report_vec.push_back(an_event);
		//             }
		//         }
		//         else
		//         {
		//             fprintf(comm, "%s", the_commentary().rand_comment("SAVE",
		//                     team[!a].player[team[!a].current_gk].name).c_str());
		//             team[!a].player[team[!a].current_gk].saves++;
		//         }
		//     }
		//     else
		//     {
		//         team[a].player[shooter].shots_off++;
		//         fprintf(comm, "%s", the_commentary().rand_comment("OFFTARGET").c_str());
		//         team[a].finalshots_off++;
		//     }
	}
}

// Given a team and an event (eg. SHOT)
// picks one player at (weighted) random
// that performed this event.
//
// For example, for SHOT, pick a player
// at weighted random according to his
// shooting skill
//
func whoDidIt(att models.Team, event int) int {
	// var total float64 = 0
	// var weight float64 = 0
	// var ar []float64 = make([]float64, NUM_PLAYERS+1)

	// Employs the weighted random algorithm
	// A player's chance to DO_IT is his
	// contribution relative to the team's total
	// contribution
	//

	// for k := 1; k <= NUM_PLAYERS; k++ {
	// 	switch event {
	// 	case DID_SHOT:
	// 		weight += att.Players[k].ShContrib * 100.0
	// 		total = att.TeamShooting * 100.0
	// 		break
	// 	case DID_FOUL:
	// 		weight += att.Players[k].Ag
	// 		total = att.Aggression
	// 		break
	// 	case DID_TACKLE:
	// 		weight += att.Players[k].TkContrib * 100.0
	// 		total = att.TeamTackling * 100.0
	// 		break
	// 	case DID_ASSIST:
	// 		weight += att.Players[k].PsContrib * 100.0
	// 		total = att.TeamPassing * 100.0
	// 		break
	// 	}

	// 	ar[k] = weight
	// }

	// randValue := total * rand.Float64()

	var k int
	// for k = 2; ar[k] <= randValue; k++ {
	// }

	return k
}

// Just random, for now
//
func howMuchInjTime() int {
	substitutions := float64(rand.Intn(6))
	injuries := float64(rand.Intn(2))
	fouls := float64(rand.Intn(10))

	return int(math.Ceil(substitutions*0.5 + injuries*0.5 + fouls*0.5))
}

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
			recalculateTeamData(&teams[0], &teams[1], true)
			recalculateTeamData(&teams[1], &teams[0], false)

			// b, _ := json.Marshal(teams[0])
			// fmt.Println(string(b))

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

func cleanInjCardIndicators() {}

func calcAggression(t *models.Team) {
	t.Aggression = 0

	for i := 0; i < NUM_PLAYERS; i++ {
		if t.Players[i].Active != 1 {
			t.Players[i].Ag = 0
		}

		t.Aggression += t.Players[i].Ag
	}
}

func calcPlayerContributions(p *models.Player) {
	if p.Active != 1 {
		p.TkContrib = 0
		p.PsContrib = 0
		p.ShContrib = 0

		return
	}

	// double tk_mult = tact_manager().get_mult(team[a].tactic, team[!a].tactic,
	// 	team[a].player[b].pos, "TK");
	// double ps_mult = tact_manager().get_mult(team[a].tactic, team[!a].tactic,
	// 		team[a].player[b].pos, "PS");
	// double sh_mult = tact_manager().get_mult(team[a].tactic, team[!a].tactic,
	// 		team[a].player[b].pos, "SH");

	var tkMult float64 = 1
	var psMult float64 = 1
	var shMult float64 = 1

	sideFactor := .75
	if (p.Side == "R" && p.LikesRight) || (p.Side == "L" && p.LikesLeft) || (p.Side == "C" && p.LikesCenter) {
		sideFactor = 1.
	}

	p.TkContrib = tkMult * float64(p.Tk) * sideFactor * p.Fatigue
	p.PsContrib = psMult * float64(p.Ps) * sideFactor * p.Fatigue
	p.ShContrib = shMult * float64(p.Sh) * sideFactor * p.Fatigue
}

// Adjusts players' total contributions, taking into account the
// side balance on each position
//
func adjustContribWithSideBalance(t *models.Team) {
	// The side balance:
	// For each position (w/o side), keep a vector of 3 elements
	// to specify the number of players playing R [0], L [1], C [2] on this position
	//
	balance := make(map[string][]int)

	// Init the side balance for all positions
	//
	for _, pos := range []string{"DF", "DM", "MF", "AM", "FW"} {
		balance[pos] = []int{0, 0, 0}
	}

	// Go over the team's players and record on what side they play,
	// updating the side balance
	//
	for b := 1; b < NUM_PLAYERS; b++ {
		if t.Players[b].Active == 1 && t.Players[b].Pos != "GK" {
			switch t.Players[b].Side {
			case "R":
				balance[t.Players[b].Pos][0]++
			case "L":
				balance[t.Players[b].Pos][1]++
			case "C":
				balance[t.Players[b].Pos][2]++
			}
		}
	}

	// For all positions, check if the side balance is equal for R and L
	// If it isn't, penalize the contributions of the players on those positions
	//
	// Additionally, penalize teams who play with more than 3 C players on
	// some position without R and L
	//
	for _, pos := range []string{"DF", "DM", "MF", "AM", "FW"} {
		onPosRight := balance[pos][0]
		onPosLeft := balance[pos][0]
		onPosCenter := balance[pos][0]

		var taxedMultiplier float64 = 1

		if onPosLeft != onPosRight {
			taxRatio := 0.25 * math.Abs(float64(onPosRight-onPosLeft)) / float64(onPosRight+onPosLeft)
			taxedMultiplier = 1 - taxRatio
		} else if onPosLeft == 0 && onPosRight == 0 && onPosCenter > 3 {
			taxedMultiplier = .87
		}

		if taxedMultiplier != 1 {
			for b := 2; b < NUM_PLAYERS; b++ {
				if t.Players[b].Active == 1 && t.Players[b].Pos == pos {
					t.Players[b].TkContrib *= taxedMultiplier
					t.Players[b].PsContrib *= taxedMultiplier
					t.Players[b].ShContrib *= taxedMultiplier
				}
			}
		}
	}
}

func calcTeamContributionsTotal(t *models.Team) {
	for b := 2; b < NUM_PLAYERS; b++ {
		t.TeamTackling += t.Players[b].TkContrib
		t.TeamPassing += t.Players[b].PsContrib
		t.TeamShooting += t.Players[b].ShContrib
	}
}

func calcShotProb(att, def *models.Team, home bool) {
	// Note: 1.0 is added to tackling, to avoid singularity when the
	// team tackling is 0
	//
	att.ShotProb = 1.8 * (float64(att.Aggression)/50.0 + 800.0*
		math.Pow(((1.0/3.0*att.TeamShooting+2.0/3.0*att.TeamPassing)/(def.TeamTackling+1.0)), 2))

	// If it is the home team, add home bonus
	//
	if home {
		att.ShotProb += 0 // homeBonus
	}
}

func recalculateTeamData(att, def *models.Team, home bool) {
	att.TeamTackling = 0
	att.TeamPassing = 0
	att.TeamShooting = 0
	calcAggression(att)

	for b := 1; b < NUM_PLAYERS; b++ {
		if att.Players[b].Active == 1 {
			fatigueDeduction := att.Players[b].NominalFatiguePerMinute
			mrnd := rand.Intn(100)
			fatigueDeduction += float64(mrnd-50) / 50.0 * 0.003

			att.Players[b].Fatigue -= fatigueDeduction

			if att.Players[b].Fatigue < 0.10 {
				att.Players[b].Fatigue = .10
			}
		}
	}

	for b := 1; b < NUM_PLAYERS; b++ {
		if b != att.CurrentGk {
			calcPlayerContributions(&att.Players[b])
		}
	}

	adjustContribWithSideBalance(att)
	calcTeamContributionsTotal(att)

	calcShotProb(att, def, home)
}
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
			att.Players[assister].KeyPasses++
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
	var total float64 = 0
	var weight float64 = 0
	var ar []float64 = make([]float64, NUM_PLAYERS)

	// Employs the weighted random algorithm
	// A player's chance to DO_IT is his
	// contribution relative to the team's total
	// contribution
	//

	for k := 0; k < NUM_PLAYERS; k++ {
		switch event {
		case DID_SHOT:
			weight += att.Players[k].ShContrib * 100.0
			total = att.TeamShooting * 100.0
			break
		case DID_FOUL:
			weight += float64(att.Players[k].Ag)
			total = float64(att.Aggression)
			break
		case DID_TACKLE:
			weight += att.Players[k].TkContrib * 100.0
			total = att.TeamTackling * 100.0
			break
		case DID_ASSIST:
			weight += att.Players[k].PsContrib * 100.0
			total = att.TeamPassing * 100.0
			break
		}

		ar[k] = weight
	}

	randValue := total * rand.Float64()

	var k int
	for k = 1; ar[k] <= randValue; k++ {
	}

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

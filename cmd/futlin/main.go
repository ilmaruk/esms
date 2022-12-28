package main

import (
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/ilmaruk/esms/internal"
)

const halfLength = 45

func main() {
	rand.Seed(time.Now().UnixMicro())

	teamStatsTotalEnabled := false

	// For each half
	//
	for halfStart := 1; halfStart < 2*halfLength; halfStart += halfLength {
		var half int
		if halfStart == 1 {
			half = 1
		} else {
			half = 2
		}
		lastMinuteOfHalf := halfStart + halfLength - 1
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
			for j := 0; j <= 1; j++ {
				// Calculate different events
				//
				ifShot(j)
				ifFoul(j)
				randomInjury(j)

				// scoreDiff := team[j].score - team[!j].score
				checkConditionals(j)
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

func cleanInjCardIndicators()   {}
func recalculateTeamsData()     {}
func ifShot(w int)              {}
func ifFoul(w int)              {}
func randomInjury(w int)        {}
func checkConditionals(w int)   {}
func addTeamStatsTotal()        {}
func updatePlayersMinuteCount() {}

// Just random, for now
//
func howMuchInjTime() int {
	substitutions := float64(rand.Intn(6))
	injuries := float64(rand.Intn(2))
	fouls := float64(rand.Intn(10))

	return int(math.Ceil(substitutions*0.5 + injuries*0.5 + fouls*0.5))
}

func calcAbility() {}

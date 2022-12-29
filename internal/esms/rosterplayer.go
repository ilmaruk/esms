package esms

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ilmaruk/esms/internal/models"
)

const (
	NUM_COLUMNS_IN_ROSTER = 25
	NUM_PLAYERS           = 16
)

func ReadTeam(rosterPath, tsPath string) (models.Team, error) {
	var team models.Team

	// First the roster players
	//
	rosterPlayers, err := ReadRoster(rosterPath)
	if err != nil {
		return team, err
	}
	team.RosterPlayers = rosterPlayers

	// Now for the rest
	//
	readFile, err := os.Open(tsPath)
	if err != nil {
		return team, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Team name and tactic
	//
	fileScanner.Scan()
	team.Name = fileScanner.Text()
	fileScanner.Scan()
	team.Tactic = fileScanner.Text()

	// Empty line
	fileScanner.Scan()

	// Players
	//
	team.Players = make([]models.Player, NUM_PLAYERS)
	for i := 0; i < NUM_PLAYERS; {
		fileScanner.Scan()
		row := fileScanner.Text()
		if len(row) == 0 {
			// Empty line
			continue
		}

		// Read players's position and name
		parts := strings.Fields(row)
		fullPos := parts[0]
		team.Players[i].Name = parts[1]

		// For GKs, just copy the position as is
		//
		if fullPos == "GK" {
			team.Players[i].Pos = "GK"
		} else {
			if !isLegalPosition(fullPos) {
				return team, fmt.Errorf("illegal position %s of %s in %s's teamsheet", fullPos, team.Players[i].Name, team.Name)
			}

			team.Players[i].Pos = fullPosToPosition(fullPos)
			team.Players[i].Side = fullPosToSide(fullPos)
		}

		// The first specified player must be a GK
		//
		if i == 0 && team.Players[i].Pos != "GK" {
			return team, fmt.Errorf("the first player in %s's teamsheet must be a GK", team.Name)
		}

		if team.Players[i].Pos == "PK:" {
			return team, fmt.Errorf("PK: where player %d was expected (%s)", i, team.Name)
		}

		found := false
		for _, player := range team.RosterPlayers {
			if player.Name != team.Players[i].Name {
				// Not the one
				continue
			}
			found = true

			if player.Injury > 0 {
				return team, fmt.Errorf("player %s (%s) is injured", player.Name, team.Name)
			}

			if player.Suspension > 0 {
				return team, fmt.Errorf("player %s (%s) is suspended", player.Name, team.Name)
			}

			team.Players[i].LikesLeft = strings.Contains(team.Players[i].PrefSide, "L")
			team.Players[i].LikesRight = strings.Contains(team.Players[i].PrefSide, "R")
			team.Players[i].LikesCenter = strings.Contains(team.Players[i].PrefSide, "C")

			team.Players[i].St = player.St
			team.Players[i].Tk = player.Tk
			team.Players[i].Ps = player.Ps
			team.Players[i].Sh = player.Sh
			team.Players[i].Stamina = player.Stamina

			// Each player has a nominal_fatigue_per_minute rating that's
			// calculated once, based on his stamina.
			//
			// I'd like the average rating be 0.031 - so that an average player
			// (stamina = 50) will lose 30 fitness points during a full game.
			//
			// The range is approximately 50 - 10 points, and the stamina range
			// is 1-99. So, first the ratio is normalized and then subtracted
			// from the average 0.031 (which, times 90 minutes, is 0.279).
			// The formula for each player is:
			//
			// fatigue            stamina - 50
			// ------- = 0.0031 - ------------  * 0.0022
			//  minute                 50
			//
			//
			// This gives (approximately) 30 lost fitness points for average players,
			// 50 for the worse stamina and 10 for the best stamina.
			//
			// A small random factor is added each minute, so the exact numbers are
			// not deterministic.
			//
			normalizedStaminaRatio := float64(team.Players[i].Stamina-50) / 50.0
			team.Players[i].NominalFatiguePerMinute = 0.0031 - normalizedStaminaRatio*0.0022

			team.Players[i].Ag = player.Ag
			team.Players[i].Fatigue = float64(player.Fitness) / 100.0
		}

		if !found {
			return team, fmt.Errorf("player %s (%s) doesn't exist in the roster file", team.Players[i].Name, team.Name)
		}

		i++
	}

	// Empty line
	//
	fileScanner.Scan()

	// There's an optional "PK: <Name>" line.
	// If it exists, the <Name> must be listed in the teamsheet.
	//
	fileScanner.Scan()
	row := fileScanner.Text()
	team.PenaltyTaker = -1

	if strings.HasPrefix(row, "PK:") {
		parts := strings.Fields(row)
		name := parts[1]

		for i := 0; i < NUM_PLAYERS; i++ {
			if name == team.Players[i].Name {
				team.PenaltyTaker = i
				break
			}
		}

		if team.PenaltyTaker == -1 {
			return team, fmt.Errorf("error in penalty kick taker of %s, player %s not listed", team.Name, name)
		}
	}

	if err := ensureNoDuplicateNames(team); err != nil {
		return team, err
	}

	// read_conditionals(teamsheet) // @todo: no conditionals for now

	// Set active flags
	team.Substitutions = 0
	team.Injuries = 0

	for i := 0; i < NUM_PLAYERS; i++ {
		if i < NUM_PLAYERS {
			team.Players[i].Active = 1
		} else {
			team.Players[i].Active = 2
		}
	}

	// In the beginning, player n.1 is always the GK
	//
	team.CurrentGk = 0

	// Data initialisation
	// @todo: unecessary setting to 0
	//
	team.Score = 0
	team.FinalShotsOn = 0
	team.FinalShotsOff = 0
	team.FinalFouls = 0
	team.TeamTackling = 0
	team.TeamPassing = 0
	team.TeamShooting = 0

	for i := 0; i < NUM_PLAYERS; i++ {
		team.Players[i].TkContrib = 0
		team.Players[i].PsContrib = 0
		team.Players[i].ShContrib = 0

		team.Players[i].YellowCards = 0
		team.Players[i].RedCards = 0
		team.Players[i].Injured = 0
		team.Players[i].TkAb = 0
		team.Players[i].PsAb = 0
		team.Players[i].ShAb = 0
		team.Players[i].StAb = 0

		// final stats initialization
		team.Players[i].Minutes = 0
		team.Players[i].Shots = 0
		team.Players[i].Goals = 0
		team.Players[i].Saves = 0
		team.Players[i].Assists = 0
		team.Players[i].Tackles = 0
		team.Players[i].KeyPasses = 0
		team.Players[i].Fouls = 0
		team.Players[i].RedCards = 0
		team.Players[i].YellowCards = 0
		team.Players[i].Conceded = 0
		team.Players[i].ShotsOn = 0
		team.Players[i].ShotsOff = 0
	}

	return team, nil
}

func ensureNoDuplicateNames(t models.Team) error {
	for i := 0; i < NUM_PLAYERS; i++ {
		for k := 0; k < NUM_PLAYERS; k++ {
			if k != i && t.Players[i].Name == t.Players[k].Name {
				return fmt.Errorf("player %s (%s) is named twice in the team sheet", t.Players[i].Name, t.Name)
			}
		}
	}
	return nil
}

func ReadRoster(path string) ([]models.RosterPlayer, error) {
	var players = make([]models.RosterPlayer, 0)

	readFile, err := os.Open(path)
	if err != nil {
		return players, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// two dummy reads, to read in the header
	//
	fileScanner.Scan()
	fileScanner.Scan()

	// read all players from the roster
	//
	for lineNum := 3; ; lineNum++ {
		if !fileScanner.Scan() {
			break
		}

		columns := strings.Fields(fileScanner.Text())

		// Empty lines are skipped
		//
		if len(columns) == 0 {
			continue
		}

		// If a non-empty line contains an incorrect amount of columns, it's
		// an error
		//
		if len(columns) != NUM_COLUMNS_IN_ROSTER {
			return players, fmt.Errorf("In roster %s, line %d: has %d columns (must be %d)", path, lineNum, len(columns), NUM_COLUMNS_IN_ROSTER)
		}

		player := models.RosterPlayer{
			Name:        columns[0],
			Age:         atoiNoError(columns[1]),
			Nationality: columns[2],
			PrefSide:    columns[3],
			St:          atoiNoError(columns[4]),
			Tk:          atoiNoError(columns[5]),
			Ps:          atoiNoError(columns[6]),
			Sh:          atoiNoError(columns[7]),
			Stamina:     atoiNoError(columns[8]),
			Ag:          atoiNoError(columns[9]),
			StAb:        atoiNoError(columns[10]),
			TkAb:        atoiNoError(columns[11]),
			PsAb:        atoiNoError(columns[12]),
			ShAb:        atoiNoError(columns[13]),
			Games:       atoiNoError(columns[14]),
			Saves:       atoiNoError(columns[15]),
			Tackles:     atoiNoError(columns[16]),
			KeyPasses:   atoiNoError(columns[17]),
			Shots:       atoiNoError(columns[18]),
			Goals:       atoiNoError(columns[19]),
			Assists:     atoiNoError(columns[20]),
			DP:          atoiNoError(columns[21]),
			Injury:      atoiNoError(columns[22]),
			Suspension:  atoiNoError(columns[23]),
			Fitness:     atoiNoError(columns[24]),
		}
		players = append(players, player)
	}

	for fileScanner.Scan() {
		fmt.Println(fileScanner.Text())
	}

	return players, nil
}

func atoiNoError(a string) int {
	i, _ := strconv.Atoi(a)
	return i
}

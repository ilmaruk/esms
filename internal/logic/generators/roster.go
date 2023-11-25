package generators

import (
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/ilmaruk/esms/internal/models"
	"github.com/ilmaruk/esms/internal/random"
	"github.com/spf13/viper"
)

var nationalities = []string{
	"arg", "aus", "bra", "bul",
	"cam", "cro", "den", "eng",
	"fra", "ger", "hol", "ire",
	"isr", "ita", "jap", "nig",
	"nor", "saf", "spa", "usa",
}

func CreateRoster(rnd random.Randomiser, numGK, numDF, numDM, numMF, numAM, numFW, mainSkill, midSkill, secSkill int) models.Roster {
	roster := models.Roster{
		ID:   uuid.New(),
		Name: GenerateTeamName(rnd),
	}

	nPlayers := numGK + numDF + numDM + numMF + numAM + numFW

	playersArr := make([]models.RosterPlayer, 0, nPlayers)
	for plCount := 1; plCount <= nPlayers; plCount++ {
		player := createRosterPlayer(rnd, plCount, numGK, numDF, numDM, numMF, numAM, numFW, mainSkill, midSkill, secSkill)
		playersArr = append(playersArr, player)
	}

	slices.SortFunc(playersArr[0:numGK], models.MoreSt)
	slices.SortFunc(playersArr[numGK:numGK+numDF+numDM], models.MoreTk)
	slices.SortFunc(playersArr[numGK+numDF+numDM:numGK+numDF+numDM+numMF+numAM], models.MorePs)
	slices.SortFunc(playersArr[numGK+numDF+numDM+numMF+numAM:], models.MoreSh)

	roster.Players = playersArr

	return roster
}

func createRosterPlayer(rnd random.Randomiser, plCount, numGK, numDF, numDM, numMF, numAM, numFW, mainSkill, midSkill, secSkill int) models.RosterPlayer {
	player := models.RosterPlayer{
		ID: uuid.New(),
	}

	halfSecSkill := secSkill / 2

	tempRand := 0

	// Name: empty, or generated, depends on flag in configuration file
	//
	player.Name = genRandomName(rnd)

	// Nationality: randomly chosen from 20 possibilities
	//
	tempRand = rnd.UniformRandom(19)
	// assert(temp_rand >= 0 && temp_rand <= 19);
	player.Nationality = nationalities[tempRand]

	// Age: Varies between 16 and 30
	//
	player.Age = rnd.AveragedRandom(23, 7)

	// Preferred side: preset probability for each
	//
	tempRand = rnd.UniformRandom(150)

	var tempSide string

	if tempRand <= 8 {
		tempSide = "RLC"
	} else if tempRand <= 13 {
		tempSide = "RL"
	} else if tempRand <= 23 {
		tempSide = "RC"
	} else if tempRand <= 33 {
		tempSide = "LC"
	} else if tempRand <= 73 {
		tempSide = "R"
	} else if tempRand <= 103 {
		tempSide = "L"
	} else {
		tempSide = "C"
	}

	player.PrefSide = tempSide

	// Skills: Depends on the position, first n_goalkeepers
	// will get the highest skill in St, and so on...
	//
	if plCount <= numGK {
		player.St = rnd.AveragedRandomPartDev(mainSkill, 3)
		player.Tk = rnd.AveragedRandomPartDev(halfSecSkill, 2)
		player.Ps = rnd.AveragedRandomPartDev(halfSecSkill, 2)
		player.Sh = rnd.AveragedRandomPartDev(halfSecSkill, 2)
	} else if plCount <= numGK+numDF {
		player.Tk = rnd.AveragedRandomPartDev(mainSkill, 3)
		player.St = rnd.AveragedRandomPartDev(halfSecSkill, 2)
		player.Ps = rnd.AveragedRandomPartDev(secSkill, 2)
		player.Sh = rnd.AveragedRandomPartDev(secSkill, 2)
	} else if plCount <= numGK+numDF+numDM {
		player.Ps = rnd.AveragedRandomPartDev(midSkill, 3)
		player.Tk = rnd.AveragedRandomPartDev(midSkill, 3)
		player.St = rnd.AveragedRandomPartDev(halfSecSkill, 2)
		player.Sh = rnd.AveragedRandomPartDev(secSkill, 2)
	} else if plCount <= numGK+numDF+numDM+numMF {
		player.Ps = rnd.AveragedRandomPartDev(mainSkill, 3)
		player.St = rnd.AveragedRandomPartDev(halfSecSkill, 2)
		player.Tk = rnd.AveragedRandomPartDev(secSkill, 2)
		player.Sh = rnd.AveragedRandomPartDev(secSkill, 2)
	} else if plCount <= numGK+numDF+numDM+numMF+numAM {
		player.Ps = rnd.AveragedRandomPartDev(midSkill, 3)
		player.Sh = rnd.AveragedRandomPartDev(midSkill, 3)
		player.Tk = rnd.AveragedRandomPartDev(secSkill, 2)
		player.St = rnd.AveragedRandomPartDev(halfSecSkill, 2)
	} else {
		player.Sh = rnd.AveragedRandomPartDev(mainSkill, 3)
		player.St = rnd.AveragedRandomPartDev(halfSecSkill, 2)
		player.Tk = rnd.AveragedRandomPartDev(secSkill, 2)
		player.Ps = rnd.AveragedRandomPartDev(secSkill, 2)
	}

	// Stamina
	//
	player.Stamina = rnd.AveragedRandomPartDev(viper.GetInt("averageStamina"), 2)

	// Aggression
	//
	player.Ag = rnd.AveragedRandomPartDev(viper.GetInt("averageAggression"), 3)

	// Abilities: set all to 300
	//
	player.StAb = 300
	player.TkAb = 300
	player.PsAb = 300
	player.ShAb = 300

	// Other stats
	//
	player.Games = 0
	player.Saves = 0
	player.Tackles = 0
	player.KeyPasses = 0
	player.Shots = 0
	player.Goals = 0
	player.Assists = 0
	player.Dp = 0
	player.Injury = 0
	player.Suspension = 0
	player.Fitness = 100

	return player
}

// A very rudimentary random name generator
func genRandomName(rnd random.Randomiser) string {
	const (
		vowelish             = "a,o,e,i,u"
		vowelishNotBegin     = "ew,ow,oo,oa,oi,oe,ae,ua"
		consonantish         = "b,c,d,f,g,h,j,k,l,m,n,p,r,s,t,v,y,z,br,cl,gr,st,jh,tr,ty,dr,kr,ry,bt,sh,ch,pr"
		consonantishNotBegin = "mn,nh,rt,rs,rst,dn,nd,ds,bt,bs,bl,sk,vr,ks,sy,ny,vr,sht,ck"
	)

	firstNameAbbr := string(int('A') + rnd.UniformRandom(25))

	lastWasVowel := false
	result := ""

	// Generate first name (a letter + "_")
	result += firstNameAbbr
	result += "_"

	// Generate beginning
	// Capitalize the first letter of the surename
	//
	if rnd.ThrowWithProb(50) {
		result += strings.ToUpper(rnd.RandElem(vowelish))
		lastWasVowel = true
	} else {
		result += strings.Title(rnd.RandElem(consonantish))
		lastWasVowel = false
	}

	howManyProceed := 2 + rnd.UniformRandom(3)

	for i := 0; i < howManyProceed; i++ {
		if lastWasVowel {
			if rnd.ThrowWithProb(50) {
				result += rnd.RandElem(consonantish)
			} else {
				result += rnd.RandElem(consonantishNotBegin)
			}
		} else {
			if rnd.ThrowWithProb(75) {
				result += rnd.RandElem(vowelish)
			} else {
				result += rnd.RandElem(vowelishNotBegin)
			}
		}

		lastWasVowel = !lastWasVowel
	}

	if len(result) > 12 {
		result = result[0 : 9+rnd.UniformRandom(3)]
	}

	return result
}

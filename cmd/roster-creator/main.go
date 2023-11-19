package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/ilmaruk/esms/internal"
	"github.com/ilmaruk/esms/internal/models"
	"github.com/spf13/viper"
)

// // ESMS - Electronic Soccer Management Simulator
// // Copyright (C) <1998-2005>  Eli Bendersky
// //
// // This program is free software, licensed with the GPL (www.fsf.org)
// //

var (
	waitFlag bool
	seed     int64
)

var nationalities = []string{
	"arg", "aus", "bra", "bul",
	"cam", "cro", "den", "eng",
	"fra", "ger", "hol", "ire",
	"isr", "ita", "jap", "nig",
	"nor", "saf", "spa", "usa",
}

const nGauss = 1000

var gaussianVars = make([]float64, nGauss)

var rnd *rand.Rand

func parseConfig() error {
	viper.SetConfigName("roster-creator-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Setting up some default values for the
	// configuration data variables
	viper.SetDefault("rosterNamePrefix", "aaa")
	viper.SetDefault("numRosters", 10)
	viper.SetDefault("numGK", 3)
	viper.SetDefault("numDF", 8)
	viper.SetDefault("numDM", 3)
	viper.SetDefault("numMF", 8)
	viper.SetDefault("numAM", 3)
	viper.SetDefault("numFW", 8)
	viper.SetDefault("cfgAverageStamina", 60)
	viper.SetDefault("cfgAverageAggression", 30)
	viper.SetDefault("averageMainSkill", 14)
	viper.SetDefault("averageMidSkill", 11)
	viper.SetDefault("averageSecondarySkill", 7)

	if err := viper.ReadInConfig(); err != nil {
		// It's fine if there's no config file: we have defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return nil
}

func main() {
	if err := parseConfig(); err != nil {
		panic(err)
	}

	// handling/parsing command line arguments
	flag.BoolVar(&waitFlag, "wait-on-exit", false, "whether to show a prompt before exiting")
	flag.Int64Var(&seed, "seed", time.Now().UnixMicro(), "the seed for the randomiser")
	flag.Parse()

	rnd = rand.New(rand.NewSource(seed))

	fillGaussianVarsArr(gaussianVars, nGauss)

	cfgNGk := viper.GetInt("numGK")
	cfgNDf := viper.GetInt("numDF")
	cfgNDm := viper.GetInt("numDM")
	cfgNMf := viper.GetInt("numMF")
	cfgNAm := viper.GetInt("numAM")
	cfgNFw := viper.GetInt("numFW")

	cfgAverageMainSkill := viper.GetInt("averageMainSkill")
	cfgAverageMidSkill := viper.GetInt("averageMidSkill")
	cfgAverageSecondarySkill := viper.GetInt("averageSecondarySkill")

	nPlayers := cfgNGk + cfgNDf + cfgNDm + cfgNMf + cfgNAm + cfgNFw

	for rosterCount := 1; rosterCount <= viper.GetInt("numRosters"); rosterCount++ {
		playersArr := make([]models.RosterPlayer, 0, nPlayers)
		for plCount := 1; plCount <= nPlayers; plCount++ {
			var player models.RosterPlayer
			tempRand := 0

			// Name: empty, or generated, depends on flag in configuration file
			//
			player.Name = genRandomName()

			// Nationality: randomly chosen from 20 possibilities
			//
			tempRand = uniformRandom(19)
			// assert(temp_rand >= 0 && temp_rand <= 19);
			player.Nationality = nationalities[tempRand]

			// Age: Varies between 16 and 30
			//
			player.Age = averagedRandom(23, 7)

			// Preferred side: preset probability for each
			//
			tempRand = uniformRandom(150)

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

			halfAverageSecondarySkill := cfgAverageSecondarySkill / 2

			// Skills: Depends on the position, first n_goalkeepers
			// will get the highest skill in St, and so on...
			//
			if plCount <= cfgNGk {
				player.St = averagedRandomPartDev(cfgAverageMainSkill, 3)
				player.Tk = averagedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Ps = averagedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Sh = averagedRandomPartDev(halfAverageSecondarySkill, 2)
			} else if plCount <= cfgNGk+cfgNDf {
				player.Tk = averagedRandomPartDev(cfgAverageMainSkill, 3)
				player.St = averagedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Ps = averagedRandomPartDev(cfgAverageSecondarySkill, 2)
				player.Sh = averagedRandomPartDev(cfgAverageSecondarySkill, 2)
			} else if plCount <= cfgNGk+cfgNDf+cfgNDm {
				player.Ps = averagedRandomPartDev(cfgAverageMidSkill, 3)
				player.Tk = averagedRandomPartDev(cfgAverageMidSkill, 3)
				player.St = averagedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Sh = averagedRandomPartDev(cfgAverageSecondarySkill, 2)
			} else if plCount <= cfgNGk+cfgNDf+cfgNDm+cfgNMf {
				player.Ps = averagedRandomPartDev(cfgAverageMainSkill, 3)
				player.St = averagedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Tk = averagedRandomPartDev(cfgAverageSecondarySkill, 2)
				player.Sh = averagedRandomPartDev(cfgAverageSecondarySkill, 2)
			} else if plCount <= cfgNGk+cfgNDf+cfgNDm+cfgNMf+cfgNAm {
				player.Ps = averagedRandomPartDev(cfgAverageMidSkill, 3)
				player.Sh = averagedRandomPartDev(cfgAverageMidSkill, 3)
				player.Tk = averagedRandomPartDev(cfgAverageSecondarySkill, 2)
				player.St = averagedRandomPartDev(halfAverageSecondarySkill, 2)
			} else {
				player.Sh = averagedRandomPartDev(cfgAverageMainSkill, 3)
				player.St = averagedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Tk = averagedRandomPartDev(cfgAverageSecondarySkill, 2)
				player.Ps = averagedRandomPartDev(cfgAverageSecondarySkill, 2)
			}

			// Stamina
			//
			player.Stamina = averagedRandomPartDev(viper.GetInt("averageStamina"), 2)

			// Aggression
			//
			player.Ag = averagedRandomPartDev(viper.GetInt("averageAggression"), 3)

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

			playersArr = append(playersArr, player)
		}

		slices.SortFunc(playersArr[0:cfgNGk], models.MoreSt)
		slices.SortFunc(playersArr[cfgNGk:cfgNGk+cfgNDf+cfgNDm], models.MoreTk)
		slices.SortFunc(playersArr[cfgNGk+cfgNDf+cfgNDm:cfgNGk+cfgNDf+cfgNDm+cfgNMf+cfgNAm], models.MorePs)
		slices.SortFunc(playersArr[cfgNGk+cfgNDf+cfgNDm+cfgNMf+cfgNAm:], models.MoreSh)

		filename := fmt.Sprintf("%s%d.txt", viper.GetString("rosterNamePrefix"), rosterCount)
		if err := models.WriteRosterPlayers(filename, playersArr); err != nil {
			panic(err)
		}
	}

	internal.MyExit(waitFlag, 0)
}

// Return a pseudo-random integer uniformly distributed
// between 0 and max
func uniformRandom(max int) int {
	return rnd.Intn(max + 1)
}

func averagedRandomPartDev(average, div int) int {
	return averagedRandom(average, average/div)
}

func averagedRandom(average int, maxDeviation int) int {
	randGaussian := gaussianVars[uniformRandom(nGauss-1)]
	deviation := float64(maxDeviation) * randGaussian

	return average + int(deviation)
}

func fillGaussianVarsArr(arr []float64, amount uint) {
	for i := uint(0); i < amount; i++ {
		var (
			s  float64
			v1 float64
			v2 float64
			x  float64
		)

		for {
			for {
				u1 := rnd.Float64()
				u2 := rnd.Float64()

				v1 = 2*u1 - 1
				v2 = 2*u2 - 1

				s = v1*v1 + v2*v2

				if s < 1.0 {
					break
				}
			}

			x = v1 * math.Sqrt(-2*math.Log(s)/s)

			if !(x >= 1.0 || x <= -1.0) {
				break
			}
		}

		arr[i] = x
	}
}

// Given a string with comma separated values (like "a,cd,k")
// returns a random value.
func randElem(csv string) string {
	elems := strings.Split(csv, ",")
	return elems[uniformRandom(len(elems)-1)]
}

// Throws a bet with probability prob of success. Returns
// true iff succeeded.
func throwWithProb(prob int) bool {
	aThrow := 1 + uniformRandom(99)
	return prob >= aThrow
}

// A very rudimentary random name generator
func genRandomName() string {
	const (
		vowelish             = "a,o,e,i,u"
		vowelishNotBegin     = "ew,ow,oo,oa,oi,oe,ae,ua"
		consonantish         = "b,c,d,f,g,h,j,k,l,m,n,p,r,s,t,v,y,z,br,cl,gr,st,jh,tr,ty,dr,kr,ry,bt,sh,ch,pr"
		consonantishNotBegin = "mn,nh,rt,rs,rst,dn,nd,ds,bt,bs,bl,sk,vr,ks,sy,ny,vr,sht,ck"
	)

	firstNameAbbr := string(int('A') + uniformRandom(25))

	lastWasVowel := false
	result := ""

	// Generate first name (a letter + "_")
	result += firstNameAbbr
	result += "_"

	// Generate beginning
	// Capitalize the first letter of the surename
	//
	if throwWithProb(50) {
		result += strings.ToUpper(randElem(vowelish))
		lastWasVowel = true
	} else {
		result += strings.ToUpper(randElem(consonantish))
		lastWasVowel = false
	}

	howManyProceed := 2 + uniformRandom(3)

	for i := 0; i < howManyProceed; i++ {
		if lastWasVowel {
			if throwWithProb(50) {
				result += randElem(consonantish)
			} else {
				result += randElem(consonantishNotBegin)
			}
		} else {
			if throwWithProb(75) {
				result += randElem(vowelish)
			} else {
				result += randElem(vowelishNotBegin)
			}
		}

		lastWasVowel = !lastWasVowel
	}

	if len(result) > 12 {
		result = result[0 : 9+uniformRandom(3)]
	}

	return result
}

func printRoster(roster []models.RosterPlayer) {
	b, _ := json.MarshalIndent(roster, "", "  ")
	fmt.Println("=============")
	fmt.Println(string(b))
}

// ESMS - Electronic Soccer Management Simulator
// Copyright (C) <1998-2005>  Eli Bendersky
//
// This program is free software, licensed with the GPL (www.fsf.org)
package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/exp/slices"

	"github.com/ilmaruk/esms/internal"
	"github.com/ilmaruk/esms/internal/models"
	"github.com/ilmaruk/esms/internal/random"
)

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

	rnd := random.NewEsmsRandomiser()

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

			halfAverageSecondarySkill := cfgAverageSecondarySkill / 2

			// Skills: Depends on the position, first n_goalkeepers
			// will get the highest skill in St, and so on...
			//
			if plCount <= cfgNGk {
				player.St = rnd.AveragedRandomPartDev(cfgAverageMainSkill, 3)
				player.Tk = rnd.AveragedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Ps = rnd.AveragedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Sh = rnd.AveragedRandomPartDev(halfAverageSecondarySkill, 2)
			} else if plCount <= cfgNGk+cfgNDf {
				player.Tk = rnd.AveragedRandomPartDev(cfgAverageMainSkill, 3)
				player.St = rnd.AveragedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Ps = rnd.AveragedRandomPartDev(cfgAverageSecondarySkill, 2)
				player.Sh = rnd.AveragedRandomPartDev(cfgAverageSecondarySkill, 2)
			} else if plCount <= cfgNGk+cfgNDf+cfgNDm {
				player.Ps = rnd.AveragedRandomPartDev(cfgAverageMidSkill, 3)
				player.Tk = rnd.AveragedRandomPartDev(cfgAverageMidSkill, 3)
				player.St = rnd.AveragedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Sh = rnd.AveragedRandomPartDev(cfgAverageSecondarySkill, 2)
			} else if plCount <= cfgNGk+cfgNDf+cfgNDm+cfgNMf {
				player.Ps = rnd.AveragedRandomPartDev(cfgAverageMainSkill, 3)
				player.St = rnd.AveragedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Tk = rnd.AveragedRandomPartDev(cfgAverageSecondarySkill, 2)
				player.Sh = rnd.AveragedRandomPartDev(cfgAverageSecondarySkill, 2)
			} else if plCount <= cfgNGk+cfgNDf+cfgNDm+cfgNMf+cfgNAm {
				player.Ps = rnd.AveragedRandomPartDev(cfgAverageMidSkill, 3)
				player.Sh = rnd.AveragedRandomPartDev(cfgAverageMidSkill, 3)
				player.Tk = rnd.AveragedRandomPartDev(cfgAverageSecondarySkill, 2)
				player.St = rnd.AveragedRandomPartDev(halfAverageSecondarySkill, 2)
			} else {
				player.Sh = rnd.AveragedRandomPartDev(cfgAverageMainSkill, 3)
				player.St = rnd.AveragedRandomPartDev(halfAverageSecondarySkill, 2)
				player.Tk = rnd.AveragedRandomPartDev(cfgAverageSecondarySkill, 2)
				player.Ps = rnd.AveragedRandomPartDev(cfgAverageSecondarySkill, 2)
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

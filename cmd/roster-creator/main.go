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
)

// // ESMS - Electronic Soccer Management Simulator
// // Copyright (C) <1998-2005>  Eli Bendersky
// //
// // This program is free software, licensed with the GPL (www.fsf.org)
// //
// #include <cstdio>
// #include <cstdlib>
// #include <sstream>
// #include <cstring>
// #include <iostream>
// #include <algorithm>
// #include <ctime>
// #include <cassert>
// #include <cmath>
// #include <climits>
// #include <cctype>

// using namespace std;

// #include "util.h"
// #include "config.h"
// #include "rosterplayer.h"
// #include "anyoption.h"

// // whether there is a wait on exit
// //
// bool waitflag = true;
var (
	waitFlag bool
	seed     int64
)

// char nationalities[20][4] = {"arg", "aus", "bra", "bul",
//
//	"cam", "cro", "den", "eng",
//	"fra", "ger", "hol", "ire",
//	"isr", "ita", "jap", "nig",
//	"nor", "saf", "spa", "usa"};
var nationalities = []string{
	"arg", "aus", "bra", "bul",
	"cam", "cro", "den", "eng",
	"fra", "ger", "hol", "ire",
	"isr", "ita", "jap", "nig",
	"nor", "saf", "spa", "usa",
}

// const unsigned N_GAUSS = 1000;
const nGauss = 1000

// double gaussian_vars[N_GAUSS] = {0};
var gaussianVars = make([]float64, nGauss)

// inline unsigned uniform_random(unsigned max);
// inline unsigned averaged_random_part_dev(unsigned average, unsigned div);
// inline unsigned averaged_random(unsigned average, unsigned max_deviation);
// void fill_gaussian_vars_arr(double *arr, unsigned amount);
// string gen_random_name(void);

func moreSt(p1, p2 models.RosterPlayer) int {
	if p1.St < p2.St {
		return 1
	}
	return -1
}

func moreTk(p1, p2 models.RosterPlayer) int {
	if p1.Tk < p2.Tk {
		return 1
	}
	return -1
}

func morePs(p1, p2 models.RosterPlayer) int {
	if p1.Ps < p2.Ps {
		return 1
	}
	return -1
}

func moreSh(p1, p2 models.RosterPlayer) int {
	if p1.Sh < p2.Sh {
		return 1
	}
	return -1
}

var rnd *rand.Rand

func main() {
	// handling/parsing command line arguments
	flag.BoolVar(&waitFlag, "wait-on-exit", false, "whether to show a prompt before exiting")
	flag.Int64Var(&seed, "seed", time.Now().UnixMicro(), "the seed for the randomiser")
	flag.Parse()

	rnd = rand.New(rand.NewSource(seed))

	fillGaussianVarsArr(gaussianVars, nGauss)

	//     the_config().load_config_file("roster_creator_cfg.txt");

	// Setting up some default values for the
	// configuration data variables
	//
	//     int cfg_n_rosters = the_config().get_int_config("N_ROSTERS", 10);
	//     int cfg_n_gk = the_config().get_int_config("N_GK", 3);
	//     int cfg_n_df = the_config().get_int_config("N_DF", 8);
	//     int cfg_n_dm = the_config().get_int_config("N_DM", 3);
	//     int cfg_n_mf = the_config().get_int_config("N_MF", 8);
	//     int cfg_n_am = the_config().get_int_config("N_AM", 3);
	//     int cfg_n_fw = the_config().get_int_config("N_FW", 5);
	//     int cfg_average_stamina = the_config().get_int_config("AVERAGE_STAMINA", 60);
	//     int cfg_average_aggression = the_config().get_int_config("AVERAGE_AGGRESSION", 30);
	//     int cfg_average_main_skill = the_config().get_int_config("AVERAGE_MAIN_SKILL", 14);
	//     int cfg_average_mid_skill = the_config().get_int_config("AVERAGE_MID_SKILL", 11);
	//     int cfg_average_secondary_skill = the_config().get_int_config("AVERAGE_SECONDARY_SKILL", 7);
	//     string cfg_roster_name_prefix = the_config().get_config_value("ROSTER_NAME_PREFIX");

	//     if (cfg_roster_name_prefix == "")
	//         cfg_roster_name_prefix = "aaa";
	cfgNRosters := 10
	cfgNGk := 3
	cfgNDf := 8
	cfgNDm := 3
	cfgNMf := 8
	cfgNAm := 3
	cfgNFw := 5
	cfgAverageStamina := 60
	cfgAverageAggression := 30

	cfgAverageMainSkill := 14
	cfgAverageMidSkill := 11
	cfgAverageSecondarySkill := 7

	cfgRosterNamePrefix := "aaa"

	nPlayers := cfgNGk + cfgNDf + cfgNDm + cfgNMf + cfgNAm + cfgNFw

	for rosterCount := 1; rosterCount <= cfgNRosters; rosterCount++ {
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
			player.Stamina = averagedRandomPartDev(cfgAverageStamina, 2)

			// Aggression
			//
			player.Ag = averagedRandomPartDev(cfgAverageAggression, 3)

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

		slices.SortFunc(playersArr[0:cfgNGk], moreSt)
		slices.SortFunc(playersArr[cfgNGk:cfgNGk+cfgNDf+cfgNDm], moreTk)
		slices.SortFunc(playersArr[cfgNGk+cfgNDf+cfgNDm:cfgNGk+cfgNDf+cfgNDm+cfgNMf+cfgNAm], morePs)
		slices.SortFunc(playersArr[cfgNGk+cfgNDf+cfgNDm+cfgNMf+cfgNAm:], moreSh)

		filename := fmt.Sprintf("%s%d.txt", cfgRosterNamePrefix, rosterCount)
		if err := models.WriteRosterPlayers(filename, playersArr); err != nil {
			panic(err)
		}

		printRoster(playersArr)
	}

	internal.MyExit(waitFlag, 0)
}

// int main(int argc, char *argv[])
// {

//     the_config().load_config_file("roster_creator_cfg.txt");

//     // Setting up some default values for the
//     // configuration data variables
//     //

//     for (int roster_count = 1; roster_count <= cfg_n_rosters; ++roster_count)
//     {
//     }

//     MY_EXIT(0);
//     return 0;
// }

// // Return a pseudo-random integer uniformly distributed
// // between 0 and max
// //
// inline unsigned uniform_random(unsigned max)
// {
//     double d = rand() / (double)RAND_MAX;
//     unsigned u = (unsigned)(d * (max + 1));

//	    return (u == max + 1 ? max - 1 : u);
//	}
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

// string gen_random_name(void)
// {

//     return result;
// }

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

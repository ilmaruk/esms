// ESMS - Electronic Soccer Management Simulator
// Copyright (C) <1998-2005>  Eli Bendersky
//
// This program is free software, licensed with the GPL (www.fsf.org)
package main

import (
	"flag"
	"time"

	"github.com/spf13/viper"

	"github.com/ilmaruk/esms/internal"
	"github.com/ilmaruk/esms/internal/logic"
	"github.com/ilmaruk/esms/internal/models"
	"github.com/ilmaruk/esms/internal/plugins/persistence/file"
	"github.com/ilmaruk/esms/internal/random"
)

type RosterStorer interface {
	Store(roster models.Roster) error
}

var (
	waitFlag bool
	seed     int64
)

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

	storer := file.NewJSONRosterStorer("./data")

	for rosterCount := 1; rosterCount <= viper.GetInt("numRosters"); rosterCount++ {
		roster := logic.CreateRoster(rnd, cfgNGk, cfgNDf, cfgNDm, cfgNMf, cfgNAm, cfgNFw, cfgAverageMainSkill, cfgAverageMidSkill, cfgAverageSecondarySkill)
		storer.Store(roster)
	}

	internal.MyExit(waitFlag, 0)
}

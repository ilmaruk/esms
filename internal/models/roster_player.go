package models

import (
	"fmt"
	"os"
)

// // ESMS - Electronic Soccer Management Simulator
// // Copyright (C) <1998-2005>  Eli Bendersky
// //
// // This program is free software, licensed with the GPL (www.fsf.org)
// //
// #include <cstdio>
// #include <cstdlib>
// #include <cstring>
// #include <fstream>

// #include "rosterplayer.h"
// #include "util.h"

type RosterPlayer struct {
	Name        string
	Nationality string
	Age         int
	PrefSide    string
	St          int
	Tk          int
	Ps          int
	Sh          int
	Stamina     int
	Ag          int
	StAb        int
	TkAb        int
	PsAb        int
	ShAb        int
	Saves       int
	Games       int
	Tackles     int
	KeyPasses   int
	Shots       int
	Goals       int
	Assists     int
	Dp          int
	Injury      int
	Suspension  int
	Fitness     int
}

// string read_roster_players(string roster_filename, RosterPlayerArray& players_arr)
// {
//     ifstream rosterfile(roster_filename.c_str());

//     if (!rosterfile)
//         return format_str("Failed to open roster %s", roster_filename.c_str());

//     string line;

//     // two dummy reads, to read in the header
//     //
//     getline(rosterfile, line);
//     getline(rosterfile, line);

//     // read all players from the roster
//     //
//     for (int line_num = 3;; ++line_num)
//     {
//         if (!getline(rosterfile, line))
//             break;

//         vector<string> columns = tokenize(line);

//         // Empty lines are skipped
//         //
//         if (columns.size() == 0)
//             continue;

//         // If a non-empty line contains an incorrect amount of columns, it's
//         // an error
//         //
//         if (columns.size() != NUM_COLUMNS_IN_ROSTER)
//             return format_str(	"In roster %s, line %d: has %d columns (must be %d)",
// 								roster_filename.c_str(), line_num, columns.size(), NUM_COLUMNS_IN_ROSTER);

// 		RosterPlayer player;

//         // Populate the player's data
//         // Not much error checking done, since rosters are all machine-generated,
//         // so str_atoi is used (0 is good enough for wrong numeric columns)
//         //
// 		player.name 		= columns[0];
//         player.age 			= str_atoi(columns[1]);
// 		player.nationality 	= columns[2];
// 		player.pref_side 	= columns[3];
//         player.st 			= str_atoi(columns[4]);
//         player.tk 			= str_atoi(columns[5]);
//         player.ps 			= str_atoi(columns[6]);
//         player.sh 			= str_atoi(columns[7]);
//         player.stamina 		= str_atoi(columns[8]);
//         player.ag 			= str_atoi(columns[9]);
//         player.st_ab 		= str_atoi(columns[10]);
//         player.tk_ab 		= str_atoi(columns[11]);
//         player.ps_ab 		= str_atoi(columns[12]);
//         player.sh_ab 		= str_atoi(columns[13]);
//         player.games		= str_atoi(columns[14]);
//         player.saves 		= str_atoi(columns[15]);
//         player.tackles 		= str_atoi(columns[16]);
//         player.keypasses 	= str_atoi(columns[17]);
//         player.shots 		= str_atoi(columns[18]);
//         player.goals 		= str_atoi(columns[19]);
//         player.assists 		= str_atoi(columns[20]);
//         player.dp 			= str_atoi(columns[21]);
//         player.injury 		= str_atoi(columns[22]);
//         player.suspension 	= str_atoi(columns[23]);
//         player.fitness 		= str_atoi(columns[24]);

// 		players_arr.push_back(player);
//     }

//     return "";
// }

func WriteRosterPlayers(rosterFilename string, playersArr []RosterPlayer) error {
	rosterFile, err := os.Create(rosterFilename)
	if err != nil {
		panic(err)
	}
	defer rosterFile.Close()

	fmt.Fprintln(rosterFile, "Name         Age Nat Prs St Tk Ps Sh Sm Ag KAb TAb PAb SAb Gam Sav Ktk Kps Sht Gls Ass  DP Inj Sus Fit")
	fmt.Fprintln(rosterFile, "------------------------------------------------------------------------------------------------------")

	for _, player := range playersArr {
		fmt.Fprintf(rosterFile, "%-13s%3d%4s%4s%3d%3d%3d%3d%3d%3d%4d%4d%4d%4d%4d%4d%4d%4d%4d%4d%4d%4d%4d%4d%4d\n",
			player.Name,
			player.Age,
			player.Nationality,
			player.PrefSide,
			player.St,
			player.Tk,
			player.Ps,
			player.Sh,
			player.Stamina,
			player.Ag,
			player.StAb,
			player.TkAb,
			player.PsAb,
			player.ShAb,
			player.Games,
			player.Saves,
			player.Tackles,
			player.KeyPasses,
			player.Shots,
			player.Goals,
			player.Assists,
			player.Dp,
			player.Injury,
			player.Suspension,
			player.Fitness)
	}

	return nil
}

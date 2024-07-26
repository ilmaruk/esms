// ESMS - Electronic Soccer Management Simulator
// Copyright (C) <1998-2005>  Eli Bendersky
//
// This program is free software, licensed with the GPL (www.fsf.org)
//
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <fstream>

#include "rosterplayer.h"
#include "util.h"

#include "nlohmann/json.hpp"

using json = nlohmann::json;

string
read_roster(string roster_filename, RosterPlayerArray &players_arr)
{
    ifstream rosterfile(roster_filename.c_str());

    if (!rosterfile)
        return format_str("Failed to open roster %s", roster_filename.c_str());

    json j;
    rosterfile >> j;
    rosterfile.close();

    Roster roster = j.get<Roster>();
    players_arr = roster.players;

    return "";
}

string write_json_roster(string roster_filename, Roster r)
{
    ofstream rosterfile(roster_filename.c_str());

    if (!rosterfile)
        return format_str("Failed to open roster %s", roster_filename.c_str());

    json j = r;
    rosterfile << j << endl;

    return "";
}

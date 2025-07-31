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

bool read_roster(string roster_filename, vector<RosterPlayer> &players_arr)
{
    ifstream rosterfile(roster_filename.c_str());

    if (!rosterfile)
        return false;

    json j;
    rosterfile >> j;
    rosterfile.close();

    Roster roster = j.get<Roster>();
    players_arr = roster.players;

    return true;
}

bool write_roster(string filename, Roster r)
{
    ofstream fh(filename.c_str());

    if (!fh)
        return false;

    json j = r;
    fh << j << endl;

    return true;
}

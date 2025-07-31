// ESMS - Electronic Soccer Management Simulator
// Copyright (C) <1998-2005>  Eli Bendersky
//
// This program is free software, licensed with the GPL (www.fsf.org)
//
#ifndef ROSTERPLAYER_H_DEFINED
#define ROSTERPLAYER_H_DEFINED

#include <string>
#include <vector>
using namespace std;

#include "models.h"

#include "nlohmann/json.hpp"
using json = nlohmann::json;

const unsigned NUM_COLUMNS_IN_ROSTER = 25;

// Reads a roster from a JSON file into the vector of RosterPlayers.
// Returns true on success, false otherwise.
//
bool read_roster(string filename, RosterPlayerArray &players_arr);

// Writes a vector of RosterPlayers into a JSON file.
// Returns true on success, false otherwise.
//
bool write_roster(string filename, Roster r);

#endif // ROSTERPLAYER_H_DEFINED

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

#include "nlohmann/json.hpp"
using json = nlohmann::json;

/// Represents player information as read from a roster
///
struct RosterPlayer
{
    string name;
    string team;
    string nationality;
    string pref_side;
    int age;

    int st;
    int tk;
    int ps;
    int sh;
    int ag;
    int stamina;

    int st_ab;
    int tk_ab;
    int ps_ab;
    int sh_ab;

    int games;
    int saves;
    int tackles;
    int keypasses;
    int shots;
    int goals;
    int assists;
    int dp;

    int injury;
    int suspension;
    int fitness;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(RosterPlayer, name, team, nationality, pref_side, age,
                                   st, tk, ps, sh, ag, stamina, st_ab, tk_ab, ps_ab, sh_ab,
                                   games, saves, tackles, keypasses, shots, goals, assists, dp,
                                   injury, suspension, fitness);
};

const unsigned NUM_COLUMNS_IN_ROSTER = 25;

typedef vector<RosterPlayer> RosterPlayerArray;
typedef vector<RosterPlayer>::iterator RosterPlayerIterator;
typedef vector<RosterPlayer>::const_iterator RosterPlayerConstIterator;

struct Roster
{
    string team_name;
    RosterPlayerArray players;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Roster, team_name, players);
};

/// Reads a roster from a JSON file into the vector of RosterPlayers.
/// Returns "" on success, and an error message if something went wrong.
///
string read_roster(string roster_filename, RosterPlayerArray &players_arr);

/// Writes a vector of RosterPlayers into a JSON file.
/// Returns "" on success, and an error message if something went wrong.
///
string write_json_roster(string roster_filename, Roster r);

#endif // ROSTERPLAYER_H_DEFINED

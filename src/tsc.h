// ESMS - Electronic Soccer Management Simulator
// Copyright (C) <1998-2005>  Eli Bendersky
//
// This program is free software, licensed with the GPL (www.fsf.org)
//
#ifndef TSC_H
#define TSC_H

#include <string>

using namespace std;

#include "nlohmann/json.hpp"
using json = nlohmann::json;

struct TeamsheetPlayer
{
    string pos;
    string name;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(TeamsheetPlayer, pos, name);
};

struct Teamsheet
{
    string team_name;
    string tactic;
    std::vector<TeamsheetPlayer> field;
    std::vector<TeamsheetPlayer> bench;
    string pk;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Teamsheet, team_name, tactic, field, bench, pk);
};

void EXIT(int rc);
void chomp(char *str);
bool parse_formation(const char *formation, int &dfs, int &mfs, int &fws, char *tactic);
void verify_position_range(int n);

#endif /* TSC_H */

#ifndef _MODELS_H
#define _MODELS_H

#include <string>
#include <vector>
using namespace std;

#include "nlohmann/json.hpp"
using json = nlohmann::json;

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

typedef vector<RosterPlayer> RosterPlayerArray;
typedef vector<RosterPlayer>::iterator RosterPlayerIterator;
typedef vector<RosterPlayer>::const_iterator RosterPlayerConstIterator;

inline int st_getter(RosterPlayerConstIterator player)
{
    return player->st * player->fitness / 100;
}

inline int tk_getter(RosterPlayerConstIterator player)
{
    return player->tk * player->fitness / 100;
}

inline int ps_getter(RosterPlayerConstIterator player)
{
    return player->ps * player->fitness / 100;
}

inline int sh_getter(RosterPlayerConstIterator player)
{
    return player->sh * player->fitness / 100;
}

struct Roster
{
    string team_name;
    RosterPlayerArray players;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Roster, team_name, players);
};

struct TeamsheetPlayer
{
    string pos;
    string name;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(TeamsheetPlayer, pos, name);
};

typedef vector<TeamsheetPlayer> TeamsheetPlayerArray;

struct Teamsheet
{
    string team_name;
    string tactic;
    TeamsheetPlayerArray field;
    TeamsheetPlayerArray bench;
    string pk;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Teamsheet, team_name, tactic, field, bench, pk);
};

/**
 * Fixtures
 */

struct Fixture
{
    string home_team;
    string away_team;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Fixture, home_team, away_team);
};

typedef vector<Fixture> FixturesWeek;
typedef vector<FixturesWeek> FixturesArray;

struct Fixtures
{
    FixturesArray weeks;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Fixtures, weeks);
};

#endif // _MODELS_H
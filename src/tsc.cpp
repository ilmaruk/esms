// ESMS - Electronic Soccer Management Simulator
// Copyright (C) <1998-2005>  Eli Bendersky
//
// This program is free software, licensed with the GPL (www.fsf.org)
//
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <set>
#include <cctype>
#include <ctime>
#include <cassert>
#include "tsc.h"
#include "rosterplayer.h"
#include "util.h"
#include "config.h"
#include "anyoption.h"

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

/// Gets the best player on some position from an array of roster players.
///
/// players 		- the array of players
/// chosen_players 	- a set of already chosen players (those won't be chosen again)
/// skill 			- pointer to a function receiving a player and returning the skill by
/// 				  which "best" is judged.
///
/// Returns the chosen player's name. Note: chosen_players is not modified !
///
string choose_best_player(const RosterPlayerArray &players,
                          const set<string> &chosen_players,
                          int (*skill)(RosterPlayerConstIterator player))
{
    int best_skill = -1;
    string name_of_best = "";

    for (RosterPlayerConstIterator player = players.begin(); player != players.end(); ++player)
    {
        if (chosen_players.find(player->name) == chosen_players.end() &&
            skill(player) > best_skill && !player->injury && !player->suspension)
        {
            best_skill = skill(player);
            name_of_best = player->name;
        }
    }

    assert(name_of_best != "");
    return name_of_best;
}

const char *WORK_DIR_ARG = "work-dir";
const char *TEAM_NAME_ARG = "team-name";
const char *FORMATION_ARG = "formation";

int main(int argc, char **argv)
{
    AnyOption *opt = new AnyOption();
    opt->noPOSIX();

    opt->setOption(WORK_DIR_ARG);
    opt->setOption(TEAM_NAME_ARG);
    opt->setOption(FORMATION_ARG);
    opt->processCommandArgs(argc, argv);

    string work_dir, config_file;
    if (opt->getValue(WORK_DIR_ARG))
    {
        work_dir = opt->getValue(WORK_DIR_ARG);
        config_file = work_dir + "/league.dat";
    }
    else
    {
        cout << "Usage: tsc --work-dir <work-dir> --roster-file <roster-file> --formation <formation>\n";
        MY_EXIT(1);
    }

    string team_name, roster_file;
    if (opt->getValue(TEAM_NAME_ARG))
    {
        team_name = opt->getValue(TEAM_NAME_ARG);
        roster_file = work_dir + "/" + team_name + ".json";
    }
    else
    {
        cout << "Usage: tsc --work-dir <work-dir> --roster-file <roster-file> --formation <formation>\n";
        MY_EXIT(1);
    }

    string formation;
    if (opt->getValue(FORMATION_ARG))
    {
        formation = opt->getValue(FORMATION_ARG);
    }
    else
    {
        cout << "Usage: tsc --work-dir <work-dir> --roster-file <roster-file> --formation <formation>\n";
        MY_EXIT(1);
    }

    FILE *teamsheetfile;

    int i, j;

    TeamsheetPlayer t_player[25];

    the_config().load_config_file(config_file);

    int num_subs = the_config().get_int_config("NUM_SUBS", 7);

    // The number of subs is not constant, therefore there is
    // a need for some smart assignment. The following array
    // sets the positions of thr first 5 subs, and then iterates
    // cyclicly. For example, if there are 2 subs allowed,
    // their positions will be GK (mandatory 1st !) and MF
    // If 7: GK, DF, MF, DF, FW, MF, DF
    //                              ^
    //                              cyclic repetition begins
    //
    const char *sub_position[] = {"DFC", "MFC", "DFC", "FWC", "MFC"};

    // Iterates (cyclicly) over positions of subs,
    //
    int sub_pos_iter = 0;

    RosterPlayerArray players;
    string msg = read_roster(roster_file, players);

    if (msg != "")
        die(msg.c_str());

    if (static_cast<int>(players.size()) < 11 + num_subs)
        die("Error: not enough players in roster\n");

    int dfs, mfs, fws;
    char tactic[2];

    if (!parse_formation(formation.c_str(), dfs, mfs, fws, tactic))
    {
        MY_EXIT(1);
    }

    // Calculate indices of the last defender and the last midfielder
    //
    int last_df = dfs + 1;
    int last_mf = dfs + mfs + 1;

    // Pick the players
    //
    // First, the best shot stopper is picked as a GK, then
    // others are picker according to the schedule of sub_position
    // as described above
    //

    // This will keep us from picking the same players more than once
    //
    set<string> chosen_players;

    for (i = 1; i <= 11; i++)
    {
        if (i == 1)
            t_player[i].pos = "GK";
        else if (i >= 2 && i <= last_df)
            t_player[i].pos = "DFC";
        else if (i > last_df && i <= last_mf)
            t_player[i].pos = "MFC";
        else if (i > last_mf && i <= 11)
            t_player[i].pos = "FWC";
    }

    // set the best GK for N.1 position
    //
    t_player[1].name = choose_best_player(players, chosen_players, st_getter);
    chosen_players.insert(t_player[1].name);

    // From now on, j is the index for players in the teamsheet
    //

    // Set the starting defenders
    //
    for (j = 2; j <= last_df; j++)
    {
        t_player[j].name = choose_best_player(players, chosen_players, tk_getter);
        chosen_players.insert(t_player[j].name);
    }

    // Set the starting midfielders
    //
    for (j = last_df + 1; j <= last_mf; j++)
    {
        t_player[j].name = choose_best_player(players, chosen_players, ps_getter);
        chosen_players.insert(t_player[j].name);
    }

    // Set the starting forwards
    //
    for (j = last_mf + 1; j <= 11; j++)
    {
        t_player[j].name = choose_best_player(players, chosen_players, sh_getter);
        chosen_players.insert(t_player[j].name);
    }

    // Set the substitute GK
    //
    t_player[12].name = choose_best_player(players, chosen_players, st_getter);
    t_player[12].pos = "GK";
    chosen_players.insert(t_player[12].name);

    string name_of_best = "";

    for (j = 13; j <= num_subs + 11; ++j)
    {
        // What position should the current sub be on ?
        //
        if (!strcmp(sub_position[sub_pos_iter], "DFC"))
            name_of_best = choose_best_player(players, chosen_players, tk_getter);
        else if (!strcmp(sub_position[sub_pos_iter], "MFC"))
            name_of_best = choose_best_player(players, chosen_players, ps_getter);
        else if (!strcmp(sub_position[sub_pos_iter], "FWC"))
            name_of_best = choose_best_player(players, chosen_players, sh_getter);
        else
            assert(0);

        t_player[j].name = name_of_best;
        t_player[j].pos = sub_position[sub_pos_iter];
        chosen_players.insert(t_player[j].name);
        sub_pos_iter = (sub_pos_iter + 1) % 5;
    }

    string teamsheet_filename = work_dir + "/" + team_name + "_sht.txt";

    teamsheetfile = fopen(teamsheet_filename.c_str(), "w");

    // Start filling the team sheet with the roster name and the
    // tactic
    //
    fprintf(teamsheetfile, "%s\n", team_name.c_str());
    fprintf(teamsheetfile, "%s\n", tactic);

    /* Print all the players and their position */
    for (i = 1; i <= 11 + num_subs; i++)
    {
        fprintf(teamsheetfile, "\n%s %s", t_player[i].pos.c_str(), t_player[i].name.c_str());

        if (i == 11)
            fprintf(teamsheetfile, "\n");
    }

    /* Print the penalty kick taker (player number last_mf + 1) */
    fprintf(teamsheetfile, "\n\nPK: %s\n\n", t_player[last_mf + 1].name.c_str());

    printf("%s created successfully\n", teamsheet_filename.c_str());

    fclose(teamsheetfile);

    return 0;
}

// Remove trailing newline
//
void chomp(char *str)
{
    int len = strlen(str);

    if (str[len - 1] == '\n')
        str[len - 1] = '\0';
}

// Parses the formation line, finds out how many defenders,
// midfielders and forwards to pick, and the tactic to use,
// performs error checking
//
// For example: 442N means 4 DFs, 4 MFs, 2 FWs, playing N
//
bool parse_formation(const char *formation, int &dfs, int &mfs,
                     int &fws, char *tactic)
{
    if (strlen(formation) != 4)
    {
        printf("The formation string must be exactly 4 characters long\n");
        printf("For example: 442N\n");
        return false;
    }

    // Random formation ?
    //
    if (!strncmp(formation, "rnd", 3))
    {
        srand(time(NULL));

        // between 3 and 5
        dfs = 3 + rand() % 3;

        // if there are 5 dfs, max of 4 mfs
        if (dfs == 5)
        {
            mfs = 1 + rand() % 4;
        }
        else // 5 mfs is also possible
        {
            mfs = 1 + rand() % 5;
        }

        fws = 10 - dfs - mfs;

        tactic[0] = formation[3];
        tactic[1] = '\0';

        return true;
    }

    dfs = formation[0] - '0';
    mfs = formation[1] - '0';
    fws = formation[2] - '0';

    tactic[0] = formation[3];
    tactic[1] = '\0';

    verify_position_range(dfs);
    verify_position_range(mfs);
    verify_position_range(fws);

    if (dfs + mfs + fws != 10)
    {
        printf("The number of players on all positions added together must be 10\n");
        printf("For example: 442N\n");
        return false;
    }

    return true;
}

void verify_position_range(int n)
{
    if (n < 1 || n > 8)
    {
        printf("The number of players on each position must be between 1 and 8\n");
        printf("For example: 442N\n");
        MY_EXIT(0);
    }
}

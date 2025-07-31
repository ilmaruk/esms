#include <cstdio>
#include <fstream>

#include "teamsheet.h"
#include "util.h"

#include "nlohmann/json.hpp"

using json = nlohmann::json;

bool write_teamsheet_as_text(const char *filename, const char *team_name, const char *tactic,
                             const TeamsheetPlayer *t_player, int num_subs, int last_mf)
{
    int i;

    FILE *fh = fopen(filename, "w");

    // Start filling the team sheet with the roster name and the
    // tactic
    //
    fprintf(fh, "%s\n", team_name);
    fprintf(fh, "%s\n", tactic);

    /* Print all the players and their position */
    for (i = 1; i <= 11 + num_subs; i++)
    {
        fprintf(fh, "\n%s %s", t_player[i].pos.c_str(), t_player[i].name.c_str());

        if (i == 11)
            fprintf(fh, "\n");
    }

    /* Print the penalty kick taker (player number last_mf + 1) */
    fprintf(fh, "\n\nPK: %s\n\n", t_player[last_mf + 1].name.c_str());

    fclose(fh);
    return true;
}

bool write_teamsheet(string filename, const Teamsheet ts)
{
    ofstream fh(filename.c_str());

    if (!fh)
        return false;

    json j = ts;
    fh << j << endl;

    return true;
}

bool read_teamsheet(string filename, Teamsheet &ts)
{
    ifstream fh(filename.c_str());

    if (!fh)
        return false;

    json j;
    fh >> j;
    fh.close();

    ts = j.get<Teamsheet>();

    return true;
}
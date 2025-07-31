#ifndef _TEAMSHEET_H
#define _TEAMSHEET_H

#include "tsc.h"
#include "models.h"

bool write_teamsheet_as_text(const char *filename, const char *team_name, const char *tactic,
                             const TeamsheetPlayer *t_player, int num_subs, int last_mf);

string write_teamsheet(const char *filename, const Teamsheet ts);

#endif // _TEAMSHEET_H
// ESMS - Electronic Soccer Management Simulator
// Copyright (C) <1998-2005>  Eli Bendersky
//
// This program is free software, licensed with the GPL (www.fsf.org)
//
#ifndef TSC_H
#define TSC_H

#include <string>

void EXIT(int rc);
void chomp(char *str);
bool parse_formation(const char *formation, int &dfs, int &mfs, int &fws, char *tactic);
void verify_position_range(int n);

#endif /* TSC_H */

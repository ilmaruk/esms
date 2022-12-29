package models

type Team struct {
	Score         int
	FinalShotsOn  int
	FinalShotsOff int
	FinalFouls    int

	Substitutions int // Number of substitutions made by a team
	Injuries      int
	Aggression    int
	ShotProb      float64

	TeamTackling float64
	TeamPassing  float64
	TeamShooting float64

	Name     string
	FullName string
	Tactic   string

	// If this is -1, the team has no preselected PK taker (the best shooter will
	// take the penalties). Otherwise, this is the number of the PK taker as
	// specified in the teamsheet.
	//
	PenaltyTaker int

	CurrentGk int
	Players   []Player // @todo: should use pointers to RosterPlayer

	RosterPlayers []RosterPlayer

	// vector<cond*> conds;
}

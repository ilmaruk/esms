package models

type Player struct {
	Name string

	// contains the 2-char position (w/o side)
	Pos string

	// 1-char side (L, R, C)
	Side string

	PrefSide   string
	St         int
	Tk         int
	Ps         int
	Sh         int
	Ag         int
	Stamina    int
	Injury     int
	Suspension int

	// we don't want to calculate it every time...
	LikesLeft   bool
	LikesRight  bool
	LikesCenter bool

	// These are used only in the game running phase
	TkContrib               float64
	PsContrib               float64
	ShContrib               float64
	NominalFatiguePerMinute float64
	Fatigue                 float64
	Injured                 int // 0 - no; 1 - yes (For the updater)

	Active int // Status: 0 - unavailable; 1 - playing  2 - available for substitution

	// final stats
	Minutes     int
	Shots       int
	Goals       int
	Saves       int
	Tackles     int
	KeyPasses   int
	Assists     int
	Fouls       int
	YellowCards int
	RedCards    int

	// Auxiliary, used for AB calculation
	//
	ShotsOn  int
	ShotsOff int
	Conceded int

	// The ability change of the player in the game
	//
	StAb int
	TkAb int
	PsAb int
	ShAb int
}

type RosterPlayer struct {
	Name        string
	Age         int
	Nationality string
	PrefSide    string
	St          int
	Tk          int
	Ps          int
	Sh          int
	Stamina     int
	Ag          int
	ShContrib   int
	TkContrib   int
	PsContrib   int
	StAb        int
	TkAb        int
	PsAb        int
	ShAb        int
	Games       int
	Saves       int
	Tackles     int
	KeyPasses   int
	Shots       int
	Goals       int
	Assists     int
	DP          int
	Injury      int
	Suspension  int
	Fitness     int
}

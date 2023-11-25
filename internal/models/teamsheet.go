package models

import "github.com/google/uuid"

type TeamsheetPlayer struct {
	Player RosterPlayer
	Pos    string
}

type Teamsheet struct {
	ID      uuid.UUID
	Players []TeamsheetPlayer
}

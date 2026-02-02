package model

import "time"

type Match struct {
	Season string
	Date   time.Time

	HomeTeam string
	AwayTeam string

	FTHG int    // Full Time Home Goals
	FTAG int    // Full Time Away Goals
	FTR  string // H, A, D (Full time result)

	HTHG int    // HalfTimeHomeGoals
	HTAG int    // HalfTimeAwayGoals
	HTR  string // HalfTimeResult

	HomeShots         int
	AwayShots         int
	HomeShotsOnTarget int
	AwayShotsOnTarget int

	HomeFouls int
	AwayFouls int

	HomeYellow int
	AwayYellow int

	HomeRed int
	AwayRed int
}

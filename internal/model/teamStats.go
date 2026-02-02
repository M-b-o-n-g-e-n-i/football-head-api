package model

type TeamStats struct {
	Team string

	Played int
	Wins   int
	Draws  int
	Losses int

	GoalsFor     int
	GoalsAgainst int
	GoalDiff     int

	Points int
}

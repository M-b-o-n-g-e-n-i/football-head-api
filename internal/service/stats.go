package service

import (
	"sort"

	"football-head-api/internal/model"
)

// StatsService holds our match data and provides methods to query it
// this is basically our "database layer" - it stores data and lets us query it
type StatsService struct {
	matches []model.Match
}

// NewStatsService creates a new StatsService with the provided matches
func NewStatsService(matches []model.Match) *StatsService {
	return &StatsService{
		matches: matches,
	}
}

// GetAllTeams returns a sorted list of unique team names
func (s *StatsService) GetAllTeams() []string {
	uniqueTeams := make(map[string]bool)

	for _, match := range s.matches {
		uniqueTeams[match.HomeTeam] = true
		uniqueTeams[match.AwayTeam] = true
	}

	teams := make([]string, 0, len(uniqueTeams))
	for team := range uniqueTeams {
		teams = append(teams, team)
	}

	sort.Strings(teams)

	return teams
}

// GetMatchCount returns the total number of matches
// simple helper function - we'll use this to verify our data is loaded
func (s *StatsService) GetMatchCount() int {
	return len(s.matches)
}

package service

import (
	"fmt"
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

// ComputeLeagueTable calculates the full league standings for all teams
// This is the core business logic for GET /league/table
func (s *StatsService) ComputeLeagueTable() []model.TeamStats {
	// We use a map to accumulate stats for each team
	// Key = team name, Value = pointer to TeamStats
	// Why pointer? So we can modify the same TeamStats as we iterate
	statsMap := make(map[string]*model.TeamStats)

	// Get all teams to initialize their stats to zero
	teams := s.GetAllTeams()
	for _, team := range teams {
		statsMap[team] = &model.TeamStats{
			Team: team,
		}
	}

	// Process each match to update team stats
	for _, match := range s.matches {
		homeStats := statsMap[match.HomeTeam]
		awayStats := statsMap[match.AwayTeam]

		// Update games played
		homeStats.Played++
		awayStats.Played++

		// Update goals
		homeStats.GoalsFor += match.FTHG
		homeStats.GoalsAgainst += match.FTAG
		awayStats.GoalsFor += match.FTAG
		awayStats.GoalsAgainst += match.FTHG

		switch match.FTR {
		case "H": // Home team won
			homeStats.Wins++
			homeStats.Points += 3
			awayStats.Losses++
			// Away team gets 0 points

		case "A": // Away team won
			awayStats.Wins++
			awayStats.Points += 3
			homeStats.Losses++
			// Home team gets 0 points

		case "D": // Draw
			homeStats.Draws++
			homeStats.Points += 1
			awayStats.Draws++
			awayStats.Points += 1
		}

		// Update goal difference
		// Goal Difference = Goals For - Goals Against
		homeStats.GoalDiff = homeStats.GoalsFor - homeStats.GoalsAgainst
		awayStats.GoalDiff = awayStats.GoalsFor - awayStats.GoalsAgainst
	}

	// Convert map to slice for sorting
	table := make([]model.TeamStats, 0, len(statsMap))
	for _, stats := range statsMap {
		// Dereference the pointer to get the actual TeamStats struct
		table = append(table, *stats)
	}

	// Sort the table by Points, then Goal Difference, then Goals For
	sort.Slice(table, func(i, j int) bool {
		// If points are different, higher points wins
		if table[i].Points != table[j].Points {
			return table[i].Points > table[j].Points
		}

		// Points are equal, check goal difference
		if table[i].GoalDiff != table[j].GoalDiff {
			return table[i].GoalDiff > table[j].GoalDiff
		}

		// Goal difference is equal, check goals scored
		if table[i].GoalsFor != table[j].GoalsFor {
			return table[i].GoalsFor > table[j].GoalsFor
		}

		// Everything equal, sort alphabetically
		return table[i].Team < table[j].Team
	})

	return table
}

// Team specific stats
func (s *StatsService) GetTeamStats(teamName string) (*model.TeamStats, error) {
	teams := s.GetAllTeams()
	teamExists := false
	for _, team := range teams {
		if team == teamName {
			teamExists = true
			break
		}
	}

	if !teamExists {
		return nil, fmt.Errorf("team '%s' not found", teamName)
	}

	stats := &model.TeamStats{
		Team: teamName,
	}

	for _, match := range s.matches {
		// Check if this team played in this match
		if match.HomeTeam == teamName {
			// This team played at home
			s.updateStatsForMatch(stats, match.FTHG, match.FTAG, match.FTR == "H", match.FTR == "D")

		} else if match.AwayTeam == teamName {
			// This team played away
			s.updateStatsForMatch(stats, match.FTAG, match.FTHG, match.FTR == "A", match.FTR == "D")
		}
	}

	stats.GoalDiff = stats.GoalsFor - stats.GoalsAgainst

	return stats, nil
}

func (s *StatsService) updateStatsForMatch(stats *model.TeamStats, goalsFor, goalsAgainst int, won, drew bool) {
	// Update basic counters
	stats.Played++
	stats.GoalsFor += goalsFor
	stats.GoalsAgainst += goalsAgainst

	// Update result counters and points
	if won {
		stats.Wins++
		stats.Points += 3
	} else if drew {
		stats.Draws++
		stats.Points += 1
	} else {
		stats.Losses++
	}
}

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"football-head-api/internal/model"
	"football-head-api/internal/service"
)

// APIHandler holds our service and provides HTTP handler methods
// Think of this as the "HTTP layer" - it knows about HTTP but delegates business logic to the service
type APIHandler struct {
	statsService *service.StatsService
}

// NewAPIHandler creates a new APIHandler with the provided service
func NewAPIHandler(statsService *service.StatsService) *APIHandler {
	return &APIHandler{
		statsService: statsService,
	}
}

// Health handles GET /health
// This is a simple "ping" endpoint to check if the server is running
func (h *APIHandler) Health(w http.ResponseWriter, r *http.Request) {
	// Headers are metadata about the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// We create a struct to hold our response data
	// json tags (like `json:"status"`) tell the JSON encoder what field names to use
	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Matches int    `json:"total_matches"`
	}{
		Status:  "healthy",
		Message: "Football Head API is running",
		Matches: h.statsService.GetMatchCount(),
	}

	// json.NewEncoder(w) creates an encoder that writes to w (our response writer)
	// .Encode(response) converts our struct to JSON and writes it
	json.NewEncoder(w).Encode(response)
}

// Teams handles GET /teams
// Returns a list of all teams in the league
func (h *APIHandler) Teams(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teams := h.statsService.GetAllTeams()

	response := struct {
		Teams []string `json:"teams"`
		Count int      `json:"count"`
	}{
		Teams: teams,
		Count: len(teams),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// RegisterRoutes sets up all our HTTP routes
// This function tells Go's HTTP server which URLs map to which handlers
func (h *APIHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/teams", h.Teams)
	mux.HandleFunc("/league/table", h.LeagueTable)
	mux.HandleFunc("/teams/", h.TeamStats)
}

// LeagueTable handles GET /league/table
// Returns the full Premier League standings sorted by points
func (h *APIHandler) LeagueTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// The service processes all 380 matches and returns sorted standings
	table := h.statsService.ComputeLeagueTable()

	// We wrap the table in a response object to add metadata
	response := struct {
		Season string            `json:"season"`
		Table  []model.TeamStats `json:"table"`
		Count  int               `json:"team_count"`
	}{
		Season: "2023/24",
		Table:  table,
		Count:  len(table),
	}

	//Send the JSON response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// TeamStats handles GET /teams/{team}/stats
// Returns detailed statistics for a specific team
func (h *APIHandler) TeamStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path

	// Simple string manipulation approach
	// Remove the prefix "/teams/" to get "Arsenal/stats"
	// Then remove the suffix "/stats" to get "Arsenal"

	// First, check if path has the right format
	if len(path) < len("/teams/") || path[:7] != "/teams/" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid URL format. Use /teams/{team}/stats",
		})
		return
	}

	// Remove "/teams/" prefix (7 characters)
	remaining := path[7:]

	// Find where "/stats" starts
	statsIndex := len(remaining) - 6 // "/stats" is 6 characters
	if statsIndex < 0 || remaining[statsIndex:] != "/stats" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid URL format. Use /teams/{team}/stats",
		})
		return
	}

	// Extract team name (everything before "/stats")
	teamName := remaining[:statsIndex]

	if teamName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Team name cannot be empty",
		})
		return
	}

	// Get team stats from the service
	stats, err := h.statsService.GetTeamStats(teamName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Team '%s' not found", teamName),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

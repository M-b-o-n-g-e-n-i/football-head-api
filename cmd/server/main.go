package main

import (
	"fmt"
	"log"
	"net/http"

	"football-head-api/internal/data"
	"football-head-api/internal/handler"
	"football-head-api/internal/service"
)

func main() {
	//Print a startup message
	fmt.Println("Football Head API - Loading EPL 2023/2024 Data...")
	fmt.Println()

	// Load match data from the CSV file
	fmt.Println("Loading match data from CSV...")
	matches, err := data.LoadMatches("data/epl_2023_2024.csv")
	if err != nil {
		log.Fatalf("Error loading matches: %v", err)
	}
	fmt.Printf("Loaded %d matches\n", len(matches))

	// Create Service layer
	fmt.Println("Initializing stats service...")
	statsService := service.NewStatsService(matches)
	fmt.Printf("Service initialized with %d teams\n", len(statsService.GetAllTeams()))

	// Create HTTP Handler layer
	fmt.Println("Setting up HTTP handlers...")
	apiHandler := handler.NewAPIHandler(statsService)

	// Set up HTTP routes
	// ServeMux = "Server Multiplexer" = URL Router
	mux := http.NewServeMux()

	apiHandler.RegisterRoutes(mux)

	// Configure the HTTP server
	port := "8080"

	// http.Server is a struct that configures our server
	server := &http.Server{
		Addr:    ":" + port, // The address to listen on (":8080" = localhost:8080)
		Handler: mux,
	}

	fmt.Println()
	fmt.Println("Server ready!")
	fmt.Println()
	fmt.Println("Available endpoints:")
	fmt.Printf("GET http://localhost:%s/health\n", port)
	fmt.Printf("GET http://localhost:%s/teams\n", port)
	fmt.Println()
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Start the server
	// ListenAndServe starts the server and blocks (waits forever)
	// It only returns if there's an error
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

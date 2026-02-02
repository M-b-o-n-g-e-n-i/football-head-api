package main

import (
	"fmt"
	"log"
	"sort"

	"football-head-api/internal/data"
)

func main() {
	//Print a startup message
	fmt.Println("Football Head API - Loading EPL 2023/2024 Data...")
	fmt.Println()

	//Load the matches from CSV
	matches, err := data.LoadMatches("data/epl_2023_2024.csv")
	if err != nil {
		//log.Fatal prints the error and exits the program with status code 1
		log.Fatalf("Error loading matches: %v", err)
	}

	fmt.Printf("Total matches loaded: %d\n", len(matches))

	//Count unique teams
	//map[string]bool means: keys are strings (team names), values are booleans
	uniqueTeams := make(map[string]bool)

	for _, match := range matches {
		//Add both home and away teams to our set
		//Setting the value to true marks that we've seen this team
		uniqueTeams[match.HomeTeam] = true
		uniqueTeams[match.AwayTeam] = true
	}

	fmt.Printf("Unique teams: %d\n", len(uniqueTeams))

	fmt.Println("\nTeams in 2023/2024 Season:")

	teamList := make([]string, 0, len(uniqueTeams))
	for team := range uniqueTeams {
		teamList = append(teamList, team)
	}

	sort.Strings(teamList)

	for i, team := range teamList {
		fmt.Printf("%d. %s\n", i+1, team)
	}

	fmt.Printf("\nData loaded successfully!")
}

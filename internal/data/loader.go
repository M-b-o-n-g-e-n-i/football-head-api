package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"football-head-api/internal/model"
)

// LoadMatches reads the CSV file and returns a slice of Match structs
func LoadMatches(filepath string) ([]model.Match, error) {
	//Open the CSV file
	file, err := os.Open(filepath)
	if err != nil {
		//If we can't open the file, return nil slice and the error
		return nil, fmt.Errorf("Failed to open file: %w", err)
	}
	//Ensure the file is closed when we're done
	//defer means "execute this when the function returns"
	defer file.Close()

	//Create a CSV reader
	reader := csv.NewReader(file)

	//Read all records at once
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read CSV: %w", err)
	}

	//Check if we have data
	if len(records) < 2 {
		//We need at least 2 rows: header + 1 data row
		return nil, fmt.Errorf("CSV file is empty or contains only headers")
	}

	//Create a slice to hold our Match structs
	matches := make([]model.Match, 0, len(records)-1)

	//Skip the header row and process each data row
	//records[0] is the header, so we start from index 1
	for i, record := range records[1:] {
		//Parse each row into a Match struct
		match, err := parseMatchRecord(record)
		if err != nil {
			//i+2 because: i starts at 0 (records[1:]), +1 for header, +1 for human-readable line numbers
			return nil, fmt.Errorf("Failed to parse row %d: %w", i+2, err)
		}
		matches = append(matches, match)
	}
	return matches, nil
}

// parseMatchRecord converts a CSV record (slice of strings) into a Match struct
func parseMatchRecord(record []string) (model.Match, error) {
	//Validate we have enough columns
	if len(record) < 22 {
		return model.Match{}, fmt.Errorf("Expected 22 columns, got %d", len(record))
	}

	//Parse the date
	//time.Parse takes a layout string and the string to parse
	date, err := time.Parse("2006-01-02", record[1])
	if err != nil {
		return model.Match{}, fmt.Errorf("Invalid date format: %w", err)
	}

	//Helper function to parse integers
	//this avoids repeating code for each integer field
	parseInt := func(s string, fieldName string) (int, error) {
		//converts string to int
		val, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("Invalid %s: %w", fieldName, err)
		}
		return val, nil
	}

	// Parse all integer fields
	fthg, err := parseInt(record[4], "FullTimeHomeGoals")
	if err != nil {
		return model.Match{}, err
	}

	ftag, err := parseInt(record[5], "FullTimeAwayGoals")
	if err != nil {
		return model.Match{}, err
	}

	hthg, err := parseInt(record[7], "HalfTimeHomeGoals")
	if err != nil {
		return model.Match{}, err
	}

	htag, err := parseInt(record[8], "HalfTimeAwayGoals")
	if err != nil {
		return model.Match{}, err
	}

	homeShots, err := parseInt(record[10], "HomeShots")
	if err != nil {
		return model.Match{}, err
	}

	awayShots, err := parseInt(record[11], "AwayShots")
	if err != nil {
		return model.Match{}, err
	}

	homeShotsOnTarget, err := parseInt(record[12], "HomeShotsOnTarget")
	if err != nil {
		return model.Match{}, err
	}

	awayShotsOnTarget, err := parseInt(record[13], "AwayShotsOnTarget")
	if err != nil {
		return model.Match{}, err
	}

	homeFouls, err := parseInt(record[16], "HomeFouls")
	if err != nil {
		return model.Match{}, err
	}

	awayFouls, err := parseInt(record[17], "AwayFouls")
	if err != nil {
		return model.Match{}, err
	}

	homeYellow, err := parseInt(record[18], "HomeYellowCards")
	if err != nil {
		return model.Match{}, err
	}

	awayYellow, err := parseInt(record[19], "AwayYellowCards")
	if err != nil {
		return model.Match{}, err
	}

	homeRed, err := parseInt(record[20], "HomeRedCards")
	if err != nil {
		return model.Match{}, err
	}

	awayRed, err := parseInt(record[21], "AwayRedCards")
	if err != nil {
		return model.Match{}, err
	}

	//Create and return the Match struct
	//We use composite literal syntax to create the struct
	return model.Match{
		Season:            record[0],
		Date:              date,
		HomeTeam:          record[2],
		AwayTeam:          record[3],
		FTHG:              fthg,
		FTAG:              ftag,
		FTR:               record[6],
		HTHG:              hthg,
		HTAG:              htag,
		HTR:               record[9],
		HomeShots:         homeShots,
		AwayShots:         awayShots,
		HomeShotsOnTarget: homeShotsOnTarget,
		AwayShotsOnTarget: awayShotsOnTarget,
		HomeFouls:         homeFouls,
		AwayFouls:         awayFouls,
		HomeYellow:        homeYellow,
		AwayYellow:        awayYellow,
		HomeRed:           homeRed,
		AwayRed:           awayRed,
	}, nil
}

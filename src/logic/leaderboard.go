// Programmed by: HOUBLOUP Alexy

package hangman

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Struct for the scores
type Score struct {
	Player   string `json:"player"`
	Points   int    `json:"points"`
	DateTime string `json:"dateTime"`
}

var scores []Score

// Initializes the scores
func ParsePoints(pointsStr int) int {
	points, err := strconv.Atoi(strconv.Itoa(pointsStr))
	if err != nil {
		fmt.Print("Error parsing points: ", err)
		return 0
	}
	return points
}

// Save the scores into a file
func SaveScoresToFile(scores []Score) {
	file, err := os.Create("scores.json")
	if err != nil {
		fmt.Print("Error creating scores.json file: ", err)
		return
	}
	defer file.Close()

	encodedScores, err := json.MarshalIndent(scores, "", "    ")
	if err != nil {
		fmt.Print("Error encoding scores: ", err)
		return
	}

	_, err = file.Write(encodedScores)
	if err != nil {
		fmt.Print("Error writing scores to file: ", err)
	}
}

// Load the scores from a file
func LoadScoresFromFile(scores *[]Score) error {
	file, err := os.Open("scores.json")
	if err != nil {
		fmt.Print("Error opening scores.json file: ", err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(scores)
	if err != nil {
		fmt.Print("Error decoding scores from JSON: ", err)
		return err
	}

	return nil
}

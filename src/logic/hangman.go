// Adapted from "Hangman"

// Programmed by:
// - AMSELLEM--BOUSIGNAC Ryan
// - LIENARD Mathieu
// - DORGES Guillaume

package hangman

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Initialization of the struct with all of important informations of the game
type HangmanWeb struct {
	Username     string
	Difficulty   int
	Secretword   string
	Hideword     string
	Incorrectes  string
	Health       int
	UsedLetters  map[string]bool
	RevealedWord []string
	RevealedClue []int
	GameOver     bool
	Victory      bool
	Score        int
	Theme        int
	ImagePath    string
}

// We initialize the map to track the info per ID ( or per User)
var GameStartedPerID = make(map[int]bool)

// Check if the user is authentificated, the difficulty, the theme and play the game
func PlayHandler(w http.ResponseWriter, r *http.Request) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	sessionID := r.URL.Query().Get("id")

	difficultyStr := r.URL.Query().Get("difficulty")
	themeStr := r.URL.Query().Get("theme")

	difficulty, _ := strconv.Atoi(difficultyStr)
	theme, _ := strconv.Atoi(themeStr)

	var userSession Session
	for _, session := range sessions {
		if strconv.Itoa(session.ID) == sessionID {
			userSession = session
			break
		}
	}

	if userSession.Game.Victory || userSession.Game.GameOver {
		GameStartedPerID[userSession.ID] = false
		fmt.Println("Game reset for ID:", userSession.ID)
	}

	if !GameStartedPerID[userSession.ID] {
		userSession.Game.Init(difficulty, theme)
		GameStartedPerID[userSession.ID] = true
		fmt.Println("Game started for ID:", userSession.ID, "with difficulty:", difficulty, "and theme:", theme)
		fmt.Println("Secret word: " + userSession.Game.Secretword)
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		input := r.FormValue("userInput")
		userSession.Game.Guess(input)
	}

	RenderTemplateInGame(w, "templates/hangman.html", userSession.Game, &userSession)
}

// Check if user wants to leave the game
func WantToLeaveHandler(w http.ResponseWriter, r *http.Request) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	sessionID := r.URL.Query().Get("id")
	var userSession Session
	for _, session := range sessions {
		if strconv.Itoa(session.ID) == sessionID {
			userSession = session
			break
		}
	}

	GameStartedPerID[userSession.ID] = false

	fmt.Println("User wants to leave the game for ID:", userSession.ID)
	http.Redirect(w, r, "/index?id="+strconv.Itoa(userSession.ID), http.StatusSeeOther)
}

// Add the score to the leaderboard
func AddScoreHandler(w http.ResponseWriter, r *http.Request) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	sessionID := r.URL.Query().Get("id")
	var userSession Session
	for _, session := range sessions {
		if strconv.Itoa(session.ID) == sessionID {
			userSession = session
			break
		}
	}

	if userSession.Game.Victory {
		userSession.Game.Score = userSession.Game.Health * userSession.Game.Difficulty

		existingScoreIndex := -1
		for i, score := range scores {
			if score.Player == userSession.Username {
				existingScoreIndex = i
				break
			}
		}

		if existingScoreIndex != -1 {
			scores[existingScoreIndex].Points += userSession.Game.Score
		} else {
			scores = append(scores, Score{Player: userSession.Username, Points: userSession.Game.Score, DateTime: time.Now().Format("2006-01-02 15:04:05")})
		}

		SaveScoresToFile(scores)
	}

	fmt.Println("Score added for ID:", userSession.ID)
	http.Redirect(w, r, "/index?id="+strconv.Itoa(userSession.ID), http.StatusSeeOther)
}

// Func to Initialize the game
func (h *HangmanWeb) Init(difficulty int, theme int) {
	fmt.Println("Initializing the game...")

	// We fill all of the variables needed to begin a game
	h.Health = 10
	h.Difficulty = difficulty
	h.Secretword = WordSelection(h.Difficulty)
	h.Hideword = strings.Repeat("_ ", len(h.Secretword))
	h.Incorrectes = ""
	h.UsedLetters = make(map[string]bool)
	h.Theme = theme

	// We find the image per theme and health points
	if h.Theme != 1 {
		h.ImagePath = GetImagePath(h.Health, h.Theme)
	}

	rand.Seed(time.Now().UnixNano())
	n := len(h.Secretword)/2 - 1
	h.RevealedClue = make([]int, 0)
	Clue := make(map[int]bool)

	// We print the clues
	for i := 0; i < n; i++ {
		randomIndex := rand.Intn(len(h.Secretword))
		if !Clue[randomIndex] {
			h.RevealedClue = append(h.RevealedClue, randomIndex)
			Clue[randomIndex] = true
			h.UsedLetters[string(h.Secretword[randomIndex])] = true

			for j := 0; j < len(h.Secretword); j++ {
				if j != randomIndex && h.Secretword[j] == h.Secretword[randomIndex] {
					h.RevealedClue = append(h.RevealedClue, j)
					h.UsedLetters[string(h.Secretword[j])] = true
				}
			}
		}
	}

	// We add the clue Letter
	for i := 0; i < len(h.Secretword); i++ {
		revealed := false
		for _, index := range h.RevealedClue {
			if i == index {
				revealed = true
				break
			}
		}
		if revealed {
			revealedLetter := string(h.Secretword[i])
			h.Hideword = replaceAtIndex(h.Hideword, i*2, revealedLetter)
		}
	}

	h.RevealedWord = strings.Split(h.Hideword, " ")
	h.Victory, h.GameOver = false, false
}

// Func of Guess (So when the player play)
func (h *HangmanWeb) Guess(letter string) {
	if h.GameOver {
		return
	}

	// We adapt the input of the player
	letter = strings.TrimSpace(letter)
	letter = strings.ToLower(letter)

	// If the player try a complete word
	if len(letter) > 1 {
		if letter == h.Secretword {
			h.Victory = true
		} else {
			h.Health -= 2
			h.Incorrectes += letter + " " + " / "
		}

		// Or Only a letter
	} else if len(letter) == 1 {
		if h.UsedLetters[letter] == false {
			h.UsedLetters[letter] = true
			if strings.Contains(h.Secretword, letter) {
				for i, l := range h.Secretword {
					if letter == string(l) {
						h.RevealedClue = append(h.RevealedClue, i)
					}
				}
			} else {
				h.Incorrectes += letter + " " + " / "
				h.Health--
			}
		}
	} else {
		fmt.Println("Erreur lors de l'input")
	}

	// If the player loose it's a Gameover
	if h.Health <= 0 {
		h.GameOver = true
	}

	// Load the images
	if h.Theme != 1 {
		h.ImagePath = GetImagePath(h.Health, h.Theme)
	}

	for i := 0; i < len(h.Secretword); i++ {
		revealed := false
		for _, index := range h.RevealedClue {
			if i == index {
				revealed = true
				break
			}
		}

		// If the user finds a good letter, it swaps the underscore by the right letter
		if revealed {
			revealedLetter := string(h.Secretword[i])
			h.Hideword = replaceAtIndex(h.Hideword, i*2, revealedLetter)
		}
	}

	// We remove the spaces
	withoutspaces := ""
	for _, i := range h.Hideword {
		if i != ' ' {
			withoutspaces += string(i)
		}
	}

	// If it's the right word, it's a Victory !
	if h.Secretword == withoutspaces {
		h.Victory = true
		h.Secretword = withoutspaces
	}
}

// Func to switch in a word the underscore by the right letter
func replaceAtIndex(str string, index int, replacement string) string {
	return str[:index] + replacement + str[index+len(replacement):]
}

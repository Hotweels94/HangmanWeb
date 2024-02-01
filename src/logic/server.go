// Programmed by:
// - HOUBLOUP Alexy
// - LIENARD Mathieu
// - AMSELLEM--BOUSIGNAC Ryan

package hangman

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

// Struct for the actual session
type Session struct {
	Username string      `json:"username"`
	Token    string      `json:"token"`
	ID       int         `json:"id"`
	Game     *HangmanWeb `json:"game"`
}

var (
	users         []User
	sessions      []Session
	sessionsMutex sync.Mutex
)

const (
	port = ":3000"
)

// Start the server, load the scores and users plus it will handle all the pages
func Start() {

	fmt.Println("--------------------")
	fmt.Println("WELCOME TO HANGMAN")
	fmt.Println("--------------------")
	err := LoadScoresFromFile(&scores)
	if err != nil {
		scores = make([]Score, 0)
	}

	LoadUsersFromFile(&users)

	fmt.Println("Loaded", len(users), "users")

	InitPasswordAdmin()

	HandleAll()
}

// Func to load all of the pages
func HandleAll() {

	// Load the CSS and images
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	//Pages load
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/credits", CreditsHandler)
	http.HandleFunc("/leaderboard", LeaderboardHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/hangman", PlayHandler)
	http.HandleFunc("/login", LoginHandler)

	//Login load
	http.HandleFunc("/createuser", CreateUserHandler)
	http.HandleFunc("/loginuser", LoginUserHandler)

	//Leaderboard load
	http.HandleFunc("/addscore", AddScoreHandler)

	//Hangman load
	http.HandleFunc("/mainmenu", WantToLeaveHandler)

	//Admin load
	http.HandleFunc("/admin", AdminHandler)
	http.HandleFunc("/delete/", DeleteHandler)

	// Start the server in port 3000
	fmt.Println("Server started on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// Handles the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	idStr := r.URL.Query().Get("id")

	id, _ := strconv.Atoi(idStr)

	var userSession Session
	for _, session := range sessions {
		if session.ID == id {
			userSession = session
			break
		}
	}

	if userSession.ID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	RenderTemplateGlobal(w, "templates/index.html", userSession)
}

// Handles the credits page
func CreditsHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplateWithoutData(w, "templates/credits.html")
}

// Handles the leaderboard page
func LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].Points > scores[j].Points
	})

	RenderTemplateGlobal(w, "templates/leaderboard.html", scores)
}

// Handles the register page
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplateWithoutData(w, "templates/register.html")
}

// Handles the play page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplateWithoutData(w, "templates/login.html")
}

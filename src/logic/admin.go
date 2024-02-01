// Programmed by: HOUBLOUP Alexy

package hangman

import (
	"fmt"
	"math/rand"
	"net/http"
)

// Admin username and password
const (
	adminUsername  = "admin"
	passwordLength = 16
	passwordChars  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

var adminPassword string

// Randomly generates a password
func generateRandomPassword(length int) string {
	password := make([]byte, length)
	for i := 0; i < length; i++ {
		password[i] = passwordChars[rand.Intn(len(passwordChars))]
	}
	return string(password)
}

// Initializes the admin password
func InitPasswordAdmin() {
	adminPassword = generateRandomPassword(passwordLength)
	fmt.Println("Admin password:", adminPassword)
}

// Authenticates the admin user
func authenticateAdmin(username, password string) bool {
	return username == adminUsername && password == adminPassword
}

// Checks if the request is authenticated
func isAdminAuthenticated(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	return ok && authenticateAdmin(username, password)
}

// Handles the admin page
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdminAuthenticated(r) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	RenderTemplateGlobal(w, "templates/admin.html", scores)
}

// Allow to delete a score
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/delete/"):]

	if _, index := findScoreByPlayer(player); index != -1 {
		scores = append(scores[:index], scores[index+1:]...)
		SaveScoresToFile(scores)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)

	fmt.Println("Score deleted for player:", player)
}

// Finds a score by player name
func findScoreByPlayer(player string) (*Score, int) {
	for i, s := range scores {
		if s.Player == player {
			return &s, i
		}
	}
	return nil, -1
}

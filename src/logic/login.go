// Programmed by: HOUBLOUP Alexy

package hangman

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const SessionDuration = 72 * time.Hour

var (
	usersCounter = 0
	maxID        int
)

// Error message struct
type ErrorMessage struct {
	Message string
}

// Add a user to the json file
type User struct {
	Username     string `json:"username"`
	PasswordHash []byte `json:"password"`
	ID           int    `json:"id"`
	Game         *HangmanWeb
}

// Create a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")

		data := ErrorMessage{
			Message: "",
		}

		// Check if the username or password is empty
		if username == "" || password == "" || confirmPassword == "" {
			data.Message = "Please fill in all fields."
			RenderTemplateGlobal(w, "templates/register.html", data)
			return
		}

		// Check if the password and confirm password are the same
		if password != confirmPassword {
			data.Message = "Password and confirm password do not match."
			RenderTemplateGlobal(w, "templates/register.html", data)
			return
		}

		// Check if the username is banned
		bannedUsernames, err := readBannedUsernames("./word/banned.txt")
		if err != nil {
			http.Error(w, "Error reading banned usernames file", http.StatusInternalServerError)
			return
		}

		for _, bannedUsername := range bannedUsernames {
			if username == bannedUsername {
				data.Message = "This username is banned"
				RenderTemplateGlobal(w, "templates/register.html", data)
				return
			}
		}

		// Check if the username already exists
		if _, found := findUserByUsername(username); found {
			data.Message = "Username already exists."
			RenderTemplateGlobal(w, "templates/register.html", data)
			return
		}

		for _, user := range users {
			if user.ID > maxID {
				maxID = user.ID
			}
		}

		// Create the user
		userID := maxID + 1
		user := User{
			Username:     username,
			PasswordHash: hashPassword(password),
			ID:           userID,
			Game:         nil,
		}

		// Add the user to the json file
		users = append(users, user)

		SaveUsersToFile(users)

		// Redirect to the login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		fmt.Print("User created: ", user.Username, "\n")
	}
}

// Login of a user
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		data := ErrorMessage{
			Message: "",
		}

		// Check if the username or password is empty
		if username == "" || password == "" {
			http.Error(w, "Username or password is empty", http.StatusBadRequest)
			return
		}

		// Check if the username exists
		user, found := findUserByUsername(username)
		if !found {
			data.Message = "User not found"
			RenderTemplateGlobal(w, "templates/login.html", data)
			return
		}

		// Check if the password is correct
		if !verifyPassword(password, user.PasswordHash) {
			data.Message = "Incorrect password"
			RenderTemplateGlobal(w, "templates/login.html", data)
			return
		}

		loginUser(w, r, user)
	}
}

// ReadBannedUsernames reads the banned usernames from a file
func readBannedUsernames(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Print("Error opening banned usernames file: ", err)
		return nil, err
	}
	defer file.Close()

	var bannedUsernames []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bannedUsernames = append(bannedUsernames, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Print("Error reading banned usernames file: ", err)
		return nil, err
	}

	return bannedUsernames, nil
}

// Redirect to the main menu page
func loginUser(w http.ResponseWriter, r *http.Request, user User) {
	sessionToken, err := generateRandomToken()
	if err != nil {
		fmt.Print("Error generating random token: ", err)
		return
	}

	// Create a new session
	session := Session{
		Username: user.Username,
		Token:    sessionToken,
		ID:       user.ID,
		Game:     &HangmanWeb{},
	}

	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	sessions = append(sessions, session)

	http.Redirect(w, r, fmt.Sprintf("/index?id=%d", user.ID), http.StatusSeeOther)

	fmt.Print("User logged in: ", user.Username, "\n")
}

// Save the users into a file
func SaveUsersToFile(users []User) {
	file, err := os.Create("users.json")
	if err != nil {
		fmt.Print("Error creating users.json file: ", err)
		return
	}
	defer file.Close()

	encodedUsers, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		fmt.Print("Error encoding users: ", err)
		return
	}

	_, err = file.Write(encodedUsers)
	if err != nil {
		fmt.Print("Error writing users to file: ", err)
		return
	}
}

// Load the users from a file
func LoadUsersFromFile(users *[]User) error {
	file, err := os.Open("users.json")
	if err != nil {
		fmt.Print("Error opening users.json file: ", err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(users); err != nil {
		fmt.Print("Error decoding users: ", err)
		return err
	}

	return nil
}

// Find a user by username
func findUserByUsername(username string) (User, bool) {
	for _, u := range users {
		if u.Username == username {
			return u, true
		}
	}
	return User{}, false
}

// Generate a random token
func generateRandomToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}

// Verify the password
func verifyPassword(inputPassword string, hashedPassword []byte) bool {
	inputHash := hashPassword(inputPassword)
	return compareHashes(inputHash, hashedPassword)
}

// Hash a password
func hashPassword(password string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hasher.Sum(nil)
	return hashedPassword
}

// Compare the hashes
func compareHashes(hash1, hash2 []byte) bool {
	return hmac.Equal(hash1, hash2)
}

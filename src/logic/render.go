// Programmed by:
// - HOUBLOUP Alexy

package hangman

import (
	"fmt"
	"net/http"
	"text/template"
)

// Render the template for the game
func RenderTemplateInGame(w http.ResponseWriter, tmpl string, data interface{}, session *Session) {
	tmpt, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Print("Error parsing template: ", err)
		return
	}

	// Create a map with the data
	dataMap := map[string]interface{}{
		"ID":              session.ID,
		"Secretword":      session.Game.Secretword,
		"Hideword":        session.Game.Hideword,
		"Incorrectes":     session.Game.Incorrectes,
		"Health":          session.Game.Health,
		"UsedLetters":     session.Game.UsedLetters,
		"RevealedWord":    session.Game.RevealedWord,
		"RevealedIndices": session.Game.RevealedClue,
		"GameOver":        session.Game.GameOver,
		"Victory":         session.Game.Victory,
		"Score":           session.Game.Score,
		"ImagePath":       session.Game.ImagePath,
		"Theme":           session.Game.Theme,
	}

	err = tmpt.Execute(w, dataMap)
	if err != nil {
		fmt.Print("Error executing template: ", err)
		return
	}
}

// Render a template with the data
func RenderTemplateGlobal(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpt, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Print("Error parsing template: ", err)
		return
	}

	err = tmpt.Execute(w, data)
	if err != nil {
		fmt.Print("Error executing template: ", err)
		return
	}
}

// Render a template without data
func RenderTemplateWithoutData(w http.ResponseWriter, tmpl string) {
	tmpt, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Print("Error parsing template: ", err)
		return
	}

	err = tmpt.Execute(w, nil)
	if err != nil {
		fmt.Print("Error executing template: ", err)
		return
	}
}

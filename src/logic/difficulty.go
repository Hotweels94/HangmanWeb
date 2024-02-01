// Programmed by:
// - DORGES Guillaume

package hangman

import (
	"bufio"
	"math/rand"
	"os"
)

// Func to have a random word
func WordSelection(difficulty int) string {
	var file *os.File
	var err error
	if difficulty < 1 && difficulty > 3 {
		return "Problème."
	}

	// We open the right file based on the difficulty chosen
	switch difficulty {
	case 1:
		file, err = os.Open("./word/words.txt")
	case 2:
		file, err = os.Open("./word/words2.txt")
	case 3:
		file, err = os.Open("./word/words3.txt")
	case 4:
		file, err = os.Open("./word/words4.txt")
	}
	if err != nil {
		return "Problème."
	}

	// We scan and choose randomly a word in the .txt folder
	defer file.Close()
	scanner := bufio.NewScanner(file)
	words := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		words = append(words, line)
	}
	if len(words) > 0 {
		randomIndex := rand.Intn(len(words))
		randomWord := words[randomIndex]
		return randomWord
	}
	return words[0]
}

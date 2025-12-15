package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	cleanText := strings.Fields(lowerText)
	return cleanText
}

func repl() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		cleanUserInput := cleanInput(userInput)

		if len(cleanUserInput) == 0 {
			continue
		}

		fmt.Printf("Your Command Was word: %s \n", cleanUserInput[0])
	}

}

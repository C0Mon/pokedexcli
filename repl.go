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
		// Take user input
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		cleanUserInput := cleanInput(userInput)
		if len(cleanUserInput) == 0 {
			continue
		}

		// Find and run user command
		found := false
		for _, command := range commands {
			if command.name == cleanUserInput[0] {
				found = true
				command.callback()
			}
		}
		if !found {
			fmt.Println("Unknown command")

		}
	}
}

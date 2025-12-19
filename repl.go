package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	cleanText := strings.Fields(lowerText)
	return cleanText
}

func repl() {
	configs := config{
		Next:     "https://pokeapi.co/api/v2/location-area/1",
		Previous: "https://pokeapi.co/api/v2/location-area/1",
	}
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
				err := command.callback(&configs)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		if !found {
			fmt.Println("Unknown command")

		}
	}
}

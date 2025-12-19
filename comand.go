package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Lists all commands and their usage",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "List the next 20 names of location areas in the pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "List the last 20 names of location areas in the pokemon world",
			callback:    commandMapb,
		},
	}
}

func commandHelp(*config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage: \n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandExit(*config) error {
	fmt.Println("Closing the pokedex ... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) error {
	splitUrl, err := splitAtXAfterN(cfg.Next, '/', 6)
	if err != nil {
		return err
	}
	url := splitUrl[0]
	fmt.Printf("%s :  %s\n", splitUrl[0], splitUrl[1])
	id, err := strconv.Atoi(splitUrl[1])
	if err != nil {
		return err
	}

	err = mapArea(url, id)
	if err != nil {
		return err
	}
	cfg.Previous = cfg.Next
	cfg.Next = url + strconv.Itoa(id+20)
	return nil
}

func commandMapb(cfg *config) error {
	splitUrl, err := splitAtXAfterN(cfg.Previous, '/', 6)
	if err != nil {
		return err
	}
	url := splitUrl[0]
	fmt.Printf("%s :  %s\n", splitUrl[0], splitUrl[1])
	id, err := strconv.Atoi(splitUrl[1])
	if err != nil {
		return err
	}
	if id < 1 {
		fmt.Println("you're on the first page")
		return nil
	}
	err = mapArea(url, id)
	if err != nil {
		return err
	}
	cfg.Next = cfg.Previous
	cfg.Previous = url + strconv.Itoa(id-20)
	return nil
}

func mapArea(url string, id int) error {

	for i := 0; i < 20; i++ {
		newUrl := url + strconv.Itoa(i+id)
		data, err := getData(newUrl)
		if err != nil {
			return err
		}
		locationArea := LocationArea{}
		err = json.Unmarshal(data, &locationArea)
		if err != nil {
			return err
		}
		fmt.Println(locationArea.Name)
	}
	return nil
}
func splitAtXAfterN(word string, x byte, n int) ([]string, error) {
	var splitString []string
	count := 0
	for i := 0; i < len(word); i++ {
		if word[i] == x {
			count += 1
			if count == n {
				splitString = append(splitString, word[:i+1])
				splitString = append(splitString, word[i+1:])
			}
		}
	}
	if len(splitString) < 2 {
		return nil, errors.New("Config unable to split correctly")
	}
	return splitString, nil
}

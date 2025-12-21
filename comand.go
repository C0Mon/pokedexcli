package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/C0Mon/pokedexcli/internal/pokeerrors"
)

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
			description: "	Lists all commands and their usage",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "	Exits the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "	List the next 20 names of location areas in the pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "	List the last 20 names of location areas in the pokemon world",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "List all pokemon in an area. Enter as \"explore {area}\"",
			callback:    commandExplore,
		},
	}
}

func commandCatch(cfg *config) error {
	if len(cfg.Arguments) < 2 {
		fmt.Printf("Invalid command. Please type as %s {pokemon}\n", cfg.Arguments[0])
		return nil
	}
	name := cfg.Arguments[1]

	// Get url
	splitUrl, err := splitAtXAfterN(cfg.Next, '/', 5)
	if err != nil {
		return err
	}
	url := splitUrl[0] + "pokemon" + name

	// Get data from api
	data, err := cfg.pokeapiClient.GetData(url)
	if err != nil {
		fmt.Println("Pokemon not found")
		return nil
	}

	pokemon := Pokemon{}

	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		return fmt.Errorf("unexpected response type")
	}

	cfg.Pokedex[name] = pokemon
	return nil
}

func commandHelp(*config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage: \n\n")
	for _, command := range commands {
		fmt.Printf("%s:	%s\n", command.name, command.description)
	}
	return nil
}

func commandExit(*config) error {
	fmt.Println("Closing the pokedex ... Goodbye!")
	os.Exit(0)
	return nil
}

func commandExplore(cfg *config) error {
	if len(cfg.Arguments) < 2 {
		fmt.Println("Invalid command. Please type as explore {area}")
		return nil
	}

	// Get url
	splitUrl, err := splitAtXAfterN(cfg.Next, '/', 6)
	if err != nil {
		return err
	}
	url := splitUrl[0] + cfg.Arguments[1]

	// Get data from api
	data, err := cfg.pokeapiClient.GetData(url)
	if err != nil {
		fmt.Println("Area not found")
		return nil
	}

	// Unmarshal response into struct
	locationArea := LocationArea{}
	err = json.Unmarshal(data, &locationArea)
	if err != nil {
		return fmt.Errorf("unexpected response type")
	}

	// Output
	fmt.Println("Found Pokemon: ")
	for i, v := range locationArea.PokemonEncounters {
		fmt.Printf("%d.	%s\n", i, v.Pokemon.Name)
	}
	return nil
}

// Call pokeapi for the next 20 areas
func commandMap(cfg *config) error {
	// Get url and id
	splitUrl, err := splitAtXAfterN(cfg.Next, '/', 6)
	if err != nil {
		return err
	}
	url := splitUrl[0]
	id, err := strconv.Atoi(splitUrl[1])
	if err != nil {
		return err
	}

	// map the area
	err = mapArea(url, id, cfg)
	if err != nil {
		var nfe *pokeerrors.NotFoundError
		if errors.As(err, &nfe) {
			fmt.Println(err.Error())
			return nil
		}
		fmt.Println("Data not found")
		return nil
	}
	cfg.Previous = cfg.Next
	cfg.Next = url + strconv.Itoa(id+20)
	return nil
}

// Call pokeapi for the previous 20 areas
func commandMapb(cfg *config) error {
	splitUrl, err := splitAtXAfterN(cfg.Previous, '/', 6)
	if err != nil {
		return err
	}
	url := splitUrl[0]
	id, err := strconv.Atoi(splitUrl[1])
	if err != nil {
		return err
	}

	if id < 1 {
		fmt.Println("you're on the first page")
		return nil
	}
	err = mapArea(url, id, cfg)
	if err != nil {
		var nfe *pokeerrors.NotFoundError
		if errors.As(err, &nfe) {
			fmt.Println(err.Error())
			return nil
		}
		fmt.Println("Data not found")
		return nil
	}
	cfg.Next = cfg.Previous
	cfg.Previous = url + strconv.Itoa(id-20)
	return nil
}

func mapArea(url string, id int, cfg *config) error {

	for i := 0; i < 20; i++ {
		newUrl := url + strconv.Itoa(i+id)
		data, err := cfg.pokeapiClient.GetData(newUrl)
		if err != nil {
			return err
		}
		locationArea := LocationArea{}
		err = json.Unmarshal(data, &locationArea)
		if err != nil {
			if string(data[:]) == "Not Found" {
				return &pokeerrors.NotFoundError{
					Entity: "Areas",
				}
			}
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

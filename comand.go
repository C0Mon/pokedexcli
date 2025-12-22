package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
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
			description: "List all pokemon in an area\nCommand: explore {area}",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon\nCommand: catch {pokemon}",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect the stats of a pokemon\nCommand: inspect {pokemon}",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught pokemon in the pokedex",
			callback:    commandPokedex,
		},
	}
}

func commandPokedex(cfg *config) error {
	if len(cfg.Pokedex) == 0 {
		return fmt.Errorf("Your Pokedex is empty :(\nUse the catch command to catch some pokemon and fill it!\n")
	}
	fmt.Println("Your Pokedex:")
	for _, v := range cfg.Pokedex {
		fmt.Printf(" - %s\n", v.Name)
	}
	return nil
}

func commandInspect(cfg *config) error {
	if len(cfg.Arguments) < 2 {
		return fmt.Errorf("Invalid command. Please type as %s {pokemon}\n", cfg.Arguments[0])
	}
	name := cfg.Arguments[1]

	data, exists := cfg.Pokedex[name]
	if !exists {

		return fmt.Errorf("Unknown Pokemon\n")
	}

	fmt.Printf("Name: %s\n", data.Name)
	fmt.Printf("Height: %d\n", data.Height)
	fmt.Printf("Weigth: %d\n", data.Weight)
	fmt.Println("Stats:")
	for _, v := range data.Stats {
		fmt.Printf(" - %s:	%d\n", v.Stat.Name, v.BaseStat)
	}

	fmt.Println("Types:")
	for _, v := range data.Types {
		fmt.Printf(" - %s\n", v.Type.Name)
	}
	return nil
}

func commandCatch(cfg *config) error {
	if len(cfg.Arguments) < 2 {
		return fmt.Errorf("Invalid command. Please type as %s {pokemon}\n", cfg.Arguments[0])
	}
	name := cfg.Arguments[1]

	// Get url
	splitUrl, err := splitAtXAfterN(cfg.Next, '/', 5)
	if err != nil {
		return err
	}
	url := splitUrl[0] + "pokemon/" + name

	// Get data from api
	data, err := cfg.pokeapiClient.GetData(url)
	if err != nil {
		fmt.Errorf("Pokemon not found\n")
		return nil
	}

	pokemon := Pokemon{}

	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		return fmt.Errorf("Pokemon not found\n")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	chance := 100 * math.Log10(float64(pokemon.BaseExperience)*0.01)
	if chance < float64(rand.IntN(100)) {
		cfg.Pokedex[name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	fmt.Printf("You may now inspect it with the inspect command.")
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
		return fmt.Errorf("Invalid command. Please type as explore {area}\n")
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
		return fmt.Errorf("API request could not be made\n")
	}
	// Unmarshal response into struct
	locationArea := LocationArea{}
	err = json.Unmarshal(data, &locationArea)
	if err != nil {
		if string(data[:]) == "Not Found" {
			return &pokeerrors.NotFoundError{
				Entity: "Area",
			}
		}
		return fmt.Errorf("unexpected response type\n")
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
		return fmt.Errorf("Data not found\n")
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
		return err
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
					Entity: "Area",
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

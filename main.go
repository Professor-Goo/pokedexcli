package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type config struct {
	nextLocationURL     *string
	previousLocationURL *string
}

type RespShallowLocations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func main() {
	cfg := &config{}
	commands := getCommands()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := scanner.Text()

		cleanedInput := cleanInput(input)

		if len(cleanedInput) == 0 {
			continue
		}

		commandName := cleanedInput[0]

		command, exists := commands[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(cfg)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world. Each subsequent call displays the next 20 locations.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 location areas in the Pokemon world. It's a way to go back.",
			callback:    commandMapb,
		},
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the names of 20 location areas in the Pokemon world. Each subsequent call displays the next 20 locations.")
	fmt.Println("mapb: Displays the names of the previous 20 location areas in the Pokemon world. It's a way to go back.")
	fmt.Println()
	return nil
}

func commandMap(cfg *config) error {
	locationAreasURL := "https://pokeapi.co/api/v2/location-area"
	if cfg.nextLocationURL != nil {
		locationAreasURL = *cfg.nextLocationURL
	}

	res, err := http.Get(locationAreasURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var locationAreasResp RespShallowLocations
	err = json.Unmarshal(body, &locationAreasResp)
	if err != nil {
		return err
	}

	cfg.nextLocationURL = locationAreasResp.Next
	cfg.previousLocationURL = locationAreasResp.Previous

	for _, loc := range locationAreasResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.previousLocationURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	res, err := http.Get(*cfg.previousLocationURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var locationAreasResp RespShallowLocations
	err = json.Unmarshal(body, &locationAreasResp)
	if err != nil {
		return err
	}

	cfg.nextLocationURL = locationAreasResp.Next
	cfg.previousLocationURL = locationAreasResp.Previous

	for _, loc := range locationAreasResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func cleanInput(text string) []string {
	cleaned := strings.ToLower(strings.TrimSpace(text))

	if cleaned == "" {
		return []string{}
	}

	return strings.Fields(cleaned)
}

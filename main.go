package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Professor-Goo/pokedexcli/internal/pokecache"
)

type config struct {
	pokeapiClient       pokecache.Cache
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

type RespLocationArea struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func main() {
	pokeClient := pokecache.NewCache(5 * time.Minute)
	cfg := &config{
		pokeapiClient: pokeClient,
	}

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
		args := []string{}
		if len(cleanedInput) > 1 {
			args = cleanedInput[1:]
		}

		command, exists := commands[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(cfg, args...)
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
		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
		},
	}
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the names of 20 location areas in the Pokemon world. Each subsequent call displays the next 20 locations.")
	fmt.Println("mapb: Displays the names of the previous 20 location areas in the Pokemon world. It's a way to go back.")
	fmt.Println("explore <area_name>: Explore a location area")
	fmt.Println()
	return nil
}

func commandMap(cfg *config, args ...string) error {
	locationAreasURL := "https://pokeapi.co/api/v2/location-area"
	if cfg.nextLocationURL != nil {
		locationAreasURL = *cfg.nextLocationURL
	}

	// Check cache first
	if val, ok := cfg.pokeapiClient.Get(locationAreasURL); ok {
		var locationAreasResp RespShallowLocations
		err := json.Unmarshal(val, &locationAreasResp)
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

	// Make API request
	res, err := http.Get(locationAreasURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Add to cache
	cfg.pokeapiClient.Add(locationAreasURL, body)

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

func commandMapb(cfg *config, args ...string) error {
	if cfg.previousLocationURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	// Check cache first
	if val, ok := cfg.pokeapiClient.Get(*cfg.previousLocationURL); ok {
		var locationAreasResp RespShallowLocations
		err := json.Unmarshal(val, &locationAreasResp)
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

	// Make API request
	res, err := http.Get(*cfg.previousLocationURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Add to cache
	cfg.pokeapiClient.Add(*cfg.previousLocationURL, body)

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

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: explore <area_name>")
	}

	areaName := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	fmt.Printf("Exploring %s...\n", areaName)

	// Check cache first
	if val, ok := cfg.pokeapiClient.Get(url); ok {
		var locationAreaResp RespLocationArea
		err := json.Unmarshal(val, &locationAreaResp)
		if err != nil {
			return err
		}

		fmt.Println("Found Pokemon:")
		for _, encounter := range locationAreaResp.PokemonEncounters {
			fmt.Printf(" - %s\n", encounter.Pokemon.Name)
		}
		return nil
	}

	// Make API request
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Add to cache
	cfg.pokeapiClient.Add(url, body)

	var locationAreaResp RespLocationArea
	err = json.Unmarshal(body, &locationAreaResp)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationAreaResp.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
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

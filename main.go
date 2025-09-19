package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	// Create command registry
	commands := getCommands()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		// Scan for user input
		scanner.Scan()
		input := scanner.Text()

		// Clean the input
		cleanedInput := cleanInput(input)

		// Get the first word (command)
		if len(cleanedInput) == 0 {
			continue
		}

		commandName := cleanedInput[0]

		// Look up command in registry
		command, exists := commands[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		// Execute the command
		err := command.callback()
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
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println()
	return nil
}

func cleanInput(text string) []string {
	// Trim leading and trailing whitespace, then convert to lowercase
	cleaned := strings.ToLower(strings.TrimSpace(text))

	// If the cleaned string is empty, return empty slice
	if cleaned == "" {
		return []string{}
	}

	// Split by whitespace and return
	return strings.Fields(cleaned)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		// Scan for user input
		scanner.Scan()
		input := scanner.Text()

		// Clean the input
		cleanedInput := cleanInput(input)

		// Get the first word (command)
		if len(cleanedInput) > 0 {
			firstWord := cleanedInput[0]
			fmt.Printf("Your command was: %s\n", firstWord)
		}
	}
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

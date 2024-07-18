package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

func displaySpinner(description string, s *spinner.Spinner, duration time.Duration) {
	fmt.Println(description)
	s.Start()
	time.Sleep(duration)
	s.Stop()
	fmt.Println() // Add a newline for better separation between spinners
}

func main() {
	duration := 5 * time.Second
	exampleCharacterSets := []struct {
		description string
		charSet     []string
	}{
		{"Simple Rotation:", spinner.CharSets[9]},
		{"Block Building:", spinner.CharSets[1]},
		{"Dot Dancing:", []string{".", "o", "O", "@", "*"}},
		{"Arrow Rotation:", spinner.CharSets[18]},
	}

	// Advanced Character Sets Descriptions and Indices
	advancedCharacterSets := []struct {
		description string
		charSet     []string
	}{
		{"Fish Swimming:", spinner.CharSets[12]},
		{"Clock:", spinner.CharSets[37]},
		{"Globe Rotation:", spinner.CharSets[39]},
		// Custom character set
		{"Custom Spinner:", []string{"⚈", "⚆", "●", "⚆"}},
	}

	fmt.Println("Examples of Character Sets")
	fmt.Println("--------------------------")

	for _, set := range exampleCharacterSets {
		s := spinner.New(set.charSet, 100*time.Millisecond)
		displaySpinner(set.description, s, duration)
	}

	fmt.Println("Advanced Character Sets")
	fmt.Println("-----------------------")

	for _, set := range advancedCharacterSets {
		s := spinner.New(set.charSet, 100*time.Millisecond)
		displaySpinner(set.description, s, duration)
	}

	// Ensure all spinners stop before the program exits
	defer func() {
		fmt.Println("All spinners showcased.")
	}()
}

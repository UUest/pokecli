package commands

import (
	"fmt"
	"os"
)

func commandExit() error {
	err := fmt.Errorf("Closing the Pokedex...Goodbye!")
	fmt.Printf("%v\n", err)
	os.Exit(0)
	return err
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for k, v := range commands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type locationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type config struct {
	NextURL     string
	PreviousURL string
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Display help information",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Get next 20 locations from the PokeAPI",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Get previous 20 locations from the PokeAPI",
		callback:    commandMapb,
	},
}

func commandExit(cfg *config) error {
	err := fmt.Errorf("Closing the Pokedex... Goodbye!")
	fmt.Printf("%v\n", err)
	defer os.Exit(0)
	return err
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

var offset int = 0

func commandMap(cfg *config) error {
	url := cfg.NextURL
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20&offset=0"
	}
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Failed to fetch data from PokeAPI: %v", err)
	}
	defer res.Body.Close()
	locationAreaData := locationArea{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locationAreaData)
	if err != nil {
		return fmt.Errorf("Failed to parse data from PokeAPI: %v", err)
	}
	for _, location := range locationAreaData.Results {
		fmt.Printf("%v\n", location.Name)
	}
	cfg.NextURL = locationAreaData.Next
	cfg.PreviousURL = locationAreaData.Previous
	return nil
}

func commandMapb(cfg *config) error {
	url := cfg.PreviousURL
	if url == "" {
		return fmt.Errorf("you're on the first page")
	}
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Failed to fetch data from PokeAPI: %v", err)
	}
	defer res.Body.Close()
	locationAreaData := locationArea{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locationAreaData)
	if err != nil {
		return fmt.Errorf("Failed to parse data from PokeAPI: %v", err)
	}
	for _, location := range locationAreaData.Results {
		fmt.Printf("%v\n", location.Name)
	}
	cfg.NextURL = locationAreaData.Next
	cfg.PreviousURL = locationAreaData.Previous
	return nil
}

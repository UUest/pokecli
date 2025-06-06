package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/UUest/pokecli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
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

type Config struct {
	NextURL     string
	PreviousURL string
	Cache       *pokecache.Cache
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

func commandExit(cfg *Config) error {
	err := fmt.Errorf("Closing the Pokedex... Goodbye!")
	fmt.Printf("%v\n", err)
	defer os.Exit(0)
	return err
}

func commandHelp(cfg *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

var offset int = 0

func commandMap(cfg *Config) error {
	url := cfg.NextURL
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20&offset=0"
	}
	var raw []byte
	if v, ok := cfg.Cache.Get(url); ok {
		raw = v
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Failed to fetch data from PokeAPI: %v", err)
		}
		defer res.Body.Close()
		raw, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Failed to read response body: %v", err)
		}
		cfg.Cache.Add(url, raw)
	}
	var locationAreaData locationArea
	err := json.Unmarshal(raw, &locationAreaData)
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

func commandMapb(cfg *Config) error {
	url := cfg.PreviousURL
	if url == "" {
		return fmt.Errorf("you're on the first page")
	}
	var raw []byte
	if v, ok := cfg.Cache.Get(url); ok {
		raw = v
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Failed to fetch data from PokeAPI: %v", err)
		}
		defer res.Body.Close()
		raw, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Failed to read response body: %v", err)
		}
		cfg.Cache.Add(url, raw)
	}
	locationAreaData := locationArea{}
	err := json.Unmarshal(raw, &locationAreaData)
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

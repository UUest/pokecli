package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/UUest/pokecli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, ...string) error
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

type exploreLocation struct {
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
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
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
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
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

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height                 int           `json:"height"`
	HeldItems              []interface{} `json:"held_items"`
	ID                     int           `json:"id"`
	IsDefault              bool          `json:"is_default"`
	LocationAreaEncounters string        `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        interface{} `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  interface{} `json:"ability"`
			IsHidden bool        `json:"is_hidden"`
			Slot     int         `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastTypes []interface{} `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string      `json:"back_default"`
		BackFemale       interface{} `json:"back_female"`
		BackShiny        string      `json:"back_shiny"`
		BackShinyFemale  interface{} `json:"back_shiny_female"`
		FrontDefault     string      `json:"front_default"`
		FrontFemale      interface{} `json:"front_female"`
		FrontShiny       string      `json:"front_shiny"`
		FrontShinyFemale interface{} `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string      `json:"front_default"`
				FrontFemale  interface{} `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string      `json:"front_default"`
				FrontFemale      interface{} `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale interface{} `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string      `json:"back_default"`
				BackFemale       interface{} `json:"back_female"`
				BackShiny        string      `json:"back_shiny"`
				BackShinyFemale  interface{} `json:"back_shiny_female"`
				FrontDefault     string      `json:"front_default"`
				FrontFemale      interface{} `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale interface{} `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string      `json:"back_default"`
						BackFemale       interface{} `json:"back_female"`
						BackShiny        string      `json:"back_shiny"`
						BackShinyFemale  interface{} `json:"back_shiny_female"`
						FrontDefault     string      `json:"front_default"`
						FrontFemale      interface{} `json:"front_female"`
						FrontShiny       string      `json:"front_shiny"`
						FrontShinyFemale interface{} `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type Config struct {
	NextURL     string
	PreviousURL string
	Cache       *pokecache.Cache
	Pokedex     map[string]Pokemon
	mu          sync.RWMutex
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
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
		"explore": {
			name:        "explore",
			description: "Explore a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Display the Pokedex",
			callback:    commandPokedex,
		},
	}
}

func commandExit(cfg *Config, param ...string) error {
	err := fmt.Errorf("Closing the Pokedex... Goodbye!")
	fmt.Printf("%v\n", err)
	defer os.Exit(0)
	return err
}

func commandHelp(cfg *Config, param ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

var offset int = 0

func commandMap(cfg *Config, param ...string) error {
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

func commandMapb(cfg *Config, param ...string) error {
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

func commandExplore(cfg *Config, location ...string) error {
	if location == nil || len(location) == 0 {
		return fmt.Errorf("Please provide a location")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location[0])
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
	exploreLocationData := exploreLocation{}
	err := json.Unmarshal(raw, &exploreLocationData)
	if err != nil {
		return fmt.Errorf("Failed to parse data from PokeAPI: %v", err)
	}
	fmt.Printf("Exploring %s...\n", location[0])
	if exploreLocationData.PokemonEncounters == nil {
		return fmt.Errorf("No Pokemon encounters found")
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range exploreLocationData.PokemonEncounters {
		fmt.Printf("- %v\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *Config, pokemon ...string) error {
	if pokemon == nil || len(pokemon) == 0 {
		return fmt.Errorf("Please provide a Pokemon")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon[0])
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
	pkm := Pokemon{}
	err := json.Unmarshal(raw, &pkm)
	if err != nil {
		return fmt.Errorf("Failed to parse data from PokeAPI: %v", err)
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon[0])
	catchChance := rand.Intn(pkm.BaseExperience * rand.Intn(pkm.BaseExperience))
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if _, dup := cfg.Pokedex[pkm.Name]; dup {
		fmt.Printf("%v is already caught!\n", pokemon[0])
		return nil
	}
	if catchChance >= (pkm.BaseExperience * rand.Intn(pkm.BaseExperience)) {
		cfg.Pokedex[pkm.Name] = pkm
		fmt.Printf("%v was caught!\n", pokemon[0])
	} else {
		fmt.Printf("%v escaped!\n", pokemon[0])
	}
	return nil
}

func commandInspect(cfg *Config, pokemon ...string) error {
	if pokemon == nil || len(pokemon) == 0 {
		return fmt.Errorf("Please provide a Pokemon")
	}
	if pkm, ok := cfg.Pokedex[pokemon[0]]; ok {
		fmt.Printf("Inspecting %v...\n", pokemon[0])
		fmt.Printf("Name: %v\n", pkm.Name)
		fmt.Printf("Height: %v\n", pkm.Height)
		fmt.Printf("Weight: %v\n", pkm.Weight)
		fmt.Println("Stats:")
		fmt.Printf("   - HP: %v\n", pkm.Stats[0].BaseStat)
		fmt.Printf("   - Attack: %v\n", pkm.Stats[1].BaseStat)
		fmt.Printf("   - Defense: %v\n", pkm.Stats[2].BaseStat)
		fmt.Printf("   - Special Attack: %v\n", pkm.Stats[3].BaseStat)
		fmt.Printf("   - Special Defense: %v\n", pkm.Stats[4].BaseStat)
		fmt.Printf("   - Speed: %v\n", pkm.Stats[5].BaseStat)
		fmt.Println("Type:")
		for _, t := range pkm.Types {
			fmt.Printf("   - %v\n", t.Type.Name)
		}
		return nil
	}
	return fmt.Errorf("Pokemon not found in Pokedex")
}

func commandPokedex(cfg *Config, pokemon ...string) error {
	fmt.Println("Pokedex:")
	if len(cfg.Pokedex) == 0 {
		fmt.Println("No Pokemon caught yet!")
		return nil
	}
	for _, pkm := range cfg.Pokedex {
		fmt.Printf("   - %v\n", pkm.Name)
	}
	return nil
}

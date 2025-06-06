package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/UUest/pokecli/internal/pokecache"
)

func main() {
	cfg := &Config{
		Cache: pokecache.NewCache(5 * time.Second),
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		text := scanner.Text()
		cleanText := CleanInput(text)
		if len(cleanText) != 0 {
			if command, ok := commands[cleanText[0]]; ok {
				var err error
				if cleanText[0] == "explore" && len(cleanText) > 1 {
					err = command.callback(cfg, cleanText[1])
				} else {
					err = command.callback(cfg)
				}
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
			} else {
				fmt.Printf("Unknown command: %s\n", cleanText[0])
			}
		} else {
			fmt.Println("No command entered")
			continue
		}
	}
}

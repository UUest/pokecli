package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		text := scanner.Text()
		cleanText := CleanInput(text)
		if len(cleanText) != 0 {
			if command, ok := commands[cleanText[0]]; ok {
				err := command.callback()
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
		if cleanText[0] == "exit" {
			fmt.Println("Exiting...")
			break
		}
	}
}

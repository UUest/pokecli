package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/UUest/pokecli/cleaninput"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		text := scanner.Text()
		cleanText := cleaninput.CleanInput(text)
		if len(cleanText) != 0 {
			fmt.Printf("Your command was: %s\n", strings.ToLower(cleanText[0]))
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

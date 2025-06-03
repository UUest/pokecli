package cleaninput

import (
	"strings"
)

func CleanInput(text string) []string {
	return strings.Fields(text)
}

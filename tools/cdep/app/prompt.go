package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// boolPrompt asks for a bool value using the label
func boolPrompt(label string) bool {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" (y/n) ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s) == "y"
}

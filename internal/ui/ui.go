package ui

import (
	"os"
)

// Codes ANSI partagés par les messages du binaire. InitColors() les
// remplit au démarrage : codes réels si stdout est un terminal et que
// NO_COLOR est absent, chaînes vides sinon (Makefile / CI).
var (
	Red, Green, Yellow, Cyan, Reset string
)

// InitColors règle les codes ANSI selon le support du terminal.
func InitColors() {
	Reset = "\033[0m"

	if !isTerminal() || os.Getenv("NO_COLOR") != "" {
		Red, Green, Yellow, Cyan, Reset = "", "", "", "", ""
		return
	}

	Red = "\033[0;31m"
	Green = "\033[0;32m"
	Yellow = "\033[1;33m"
	Cyan = "\033[0;36m"
}

func isTerminal() bool {
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}

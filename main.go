// Command commit enchaîne lint, commit et push. Le point d'entrée traite
// --help et --version, initialise les couleurs puis délègue
// l'orchestration à internal.
package main

import (
	"fmt"
	"mzeahmed/commit/internal"
	"mzeahmed/commit/internal/ui"
	"os"
	"runtime/debug"
)

const usage = `Usage : commit <message> [branche]

  message     message du commit
  branche     garde-fou : doit être la branche courante (défaut : branche courante)

  -h, --help     affiche cette aide
      --version  affiche la version`

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-h", "--help":
			fmt.Println(usage)
			return
		case "--version":
			fmt.Println("commit", version())
			return
		}
	}

	ui.InitColors()

	if err := internal.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// version lit la version inscrite par Go dans le binaire (étiquette Git).
func version() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}

	return "inconnue"
}

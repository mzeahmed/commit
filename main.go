// Command commit est l'équivalent Go de tools/scripts/shell/commit.sh —
// le point d'entrée se contente d'initialiser les couleurs puis de
// déléguer l'orchestration à internal/commit (voir ce package).
package main

import (
	"fmt"
	"mzeahmed/commit/internal"
	"mzeahmed/commit/internal/ui"
	"os"
)

func main() {
	ui.InitColors()

	if err := internal.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

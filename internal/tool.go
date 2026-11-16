package internal

import (
	"fmt"
	"mzeahmed/commit/internal/ui"
	"os"
)

// requireBinary vérifie qu'un exécutable requis est bien présent.
func requireBinary(path string) error {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() || info.Mode()&0o111 == 0 {
		fmt.Printf("%s❌ Outil manquant ou non exécutable : %s%s\n", ui.Red, path, ui.Reset)
		fmt.Printf("%sAs-tu lancé 'composer install' ?%s\n", ui.Yellow, ui.Reset)

		return fmt.Errorf("outil requis manquant : %s", path)
	}

	return nil
}

package internal

import (
	"os"
	"os/exec"
	"strings"
)

// gitOutput lance git dans le répertoire courant et renvoie sa sortie
// standard trimmée (une ligne).
func gitOutput(args ...string) (string, error) {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// runGit lance git en laissant sa sortie filer vers le terminal
// (commit, push, add — comme le script bash qui ne les capture pas).
func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// gitExitCode renvoie le code de sortie de git sans le transformer en
// erreur d'exécution quand il est non nul (nécessaire pour les commandes
// dont le code de sortie fait partie du contrat, comme diff --quiet).
func gitExitCode(args ...string) (int, error) {
	err := exec.Command("git", args...).Run()
	if err == nil {
		return 0, nil
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return 0, err
	}

	return exitErr.ExitCode(), nil
}

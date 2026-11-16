package internal

import (
	"fmt"
	"mzeahmed/commit/internal/ui"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// runStep lance une commande en arrière-plan avec un spinner, capture le
// log dans un fichier temporaire (montré uniquement en cas d'échec) et
// renvoie une erreur d'échec explicite — l'équivalent exact de run_step()
// du script bash, spinner inclus (cadre |/-\, tick 100 ms).
func runStep(label string, name string, args ...string) error {
	fmt.Printf("%s%s...%s\n", ui.Cyan, label, ui.Reset)

	logFile, err := os.CreateTemp("", "commit-*.log")
	if err != nil {
		return fmt.Errorf("création du fichier de log : %w", err)
	}
	logPath := logFile.Name()
	defer os.Remove(logPath)

	cmd := exec.Command(name, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		logFile.Close()

		fmt.Printf("%s❌ %s échoué :%s\n", ui.Red, label, ui.Reset)
		fmt.Printf("%s❌ Commit/push annulés.%s\n", ui.Red, ui.Reset)

		return fmt.Errorf("lancement de %s : %w", filepath.Base(name), err)
	}

	done := make(chan struct{})
	finished := make(chan struct{})
	go spinner(done, finished)

	waitErr := cmd.Wait()
	close(done)
	<-finished // attend que le spinner ait nettoyé sa ligne
	logFile.Close()

	if waitErr != nil {
		logContent, readErr := os.ReadFile(logPath)
		if readErr == nil && len(logContent) > 0 {
			fmt.Printf("%s❌ %s échoué :%s\n", ui.Red, label, ui.Reset)
			os.Stdout.Write(logContent)
		} else {
			fmt.Printf("%s❌ %s échoué :%s\n", ui.Red, label, ui.Reset)
		}
		fmt.Printf("%s❌ Commit/push annulés.%s\n", ui.Red, ui.Reset)

		return fmt.Errorf("%s échoué : %w", label, waitErr)
	}

	fmt.Printf("%s✅ %s OK.%s\n", ui.Green, label, ui.Reset)

	return nil
}

// spinner anime la ligne courante jusqu'à la fermeture de done, puis
// efface son empreinte et signale la fin via finished.
func spinner(done <-chan struct{}, finished chan<- struct{}) {
	defer close(finished)

	frames := []rune{'|', '/', '-', '\\'}
	for i := 0; ; i++ {
		select {
		case <-done:
			fmt.Print("    \b\b\b\b")
			return
		default:
		}

		fmt.Printf(" [%c]  ", frames[i%len(frames)])
		time.Sleep(100 * time.Millisecond)
		fmt.Print("\b\b\b\b\b\b")
	}
}

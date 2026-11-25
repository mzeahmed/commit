// Package internal orchestre la commande commit : aperçu des changements,
// lint PHP préalable (Rector et PHP-CS-Fixer en dry-run, log montré
// seulement en cas d'échec, avec spinner), puis git add/commit/push.
//
// Usage :
//
//	commit <message> [branche]
//
// Les couleurs ANSI sont désactivées si la sortie n'est pas un terminal
// ou si NO_COLOR est défini (Makefile / CI).
package internal

import (
	"fmt"
	"mzeahmed/commit/internal/ui"
	"os"
	"strconv"
	"strings"
)

// Run orchestre le workflow complet : validation des arguments,
// positionnement à la racine du dépôt, contrôle de branche, aperçu, lint,
// commit et push.
func Run(args []string) error {
	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf(
			"%s❌ Merci de préciser le message du commit en argument.%s\n%sExemple : commit \"Mise à jour de la version de Yoostart\"%s",
			ui.Red, ui.Reset, ui.Yellow, ui.Reset,
		)
	}

	commitMessage := args[0]

	// Se placer à la racine du repo (au cas où), comme le script bash.
	repoRoot, err := gitOutput("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("résolution de la racine du dépôt : %w", err)
	}
	if err := os.Chdir(repoRoot); err != nil {
		return fmt.Errorf("changement de répertoire vers %s : %w", repoRoot, err)
	}

	// Branche courante. symbolic-ref fonctionne aussi sur un dépôt sans commit.
	currentBranch, err := gitOutput("symbolic-ref", "--short", "HEAD")
	if err != nil {
		return fmt.Errorf("aucune branche courante (HEAD détachée ?) : %w", err)
	}

	// L'argument « branche » est un garde-fou facultatif.
	branchName := currentBranch
	if len(args) > 1 && strings.TrimSpace(args[1]) != "" {
		branchName = args[1]
	}
	if branchName != currentBranch {
		// Garde anti-erreur : pas une panne de l'outil mais un mauvais
		// usage — sortie propre (exit 0) pour ne pas faire échouer make.
		fmt.Printf(
			"%s⚠️  Attention : la branche courante est '%s', mais tu as demandé '%s'.%s\n",
			ui.Yellow,
			currentBranch,
			branchName,
			ui.Reset,
		)
		return nil
	}

	// Montrer ce que `git add .` va embarquer avant de demander confirmation.
	changes, err := gitOutput("status", "--porcelain")
	if err != nil {
		return fmt.Errorf("lecture de l'état du dépôt : %w", err)
	}

	if changes == "" {
		fmt.Printf("%sℹ️  Aucun changement à committer.%s\n", ui.Yellow, ui.Reset)

		// Commits restés en local (push refusé ou échoué) : proposer l'envoi.
		pending, err := unpushedCommits(branchName)
		if err != nil {
			return err
		}
		if pending == 0 {
			return nil
		}
		fmt.Printf("%sℹ️  %d commit(s) non poussé(s) sur %s.%s\n", ui.Yellow, pending, branchName, ui.Reset)

		return askAndPush(branchName)
	}

	fmt.Printf("%s📋 Fichiers qui seront commités sur %s :%s\n", ui.Cyan, branchName, ui.Reset)
	if err := runGit("status", "--short", "--untracked-files=all"); err != nil {
		return err
	}

	confirmed, err := ui.Confirm(fmt.Sprintf("Valider le commit sur %s ?", branchName), true)
	if err != nil || !confirmed {
		// Annulation utilisateur (« No », Ctrl+C, EOF) : sortie propre
		// sans code d'erreur — l'outil a fait son travail, make reste
		// silencieux (`make commit` ne s'arrête donc pas sur ce cas).
		fmt.Printf("%sCommit annulé.%s\n", ui.Red, ui.Reset)
		return nil
	}

	if err := lint(); err != nil {
		return err
	}

	// Commit.
	fmt.Printf("%s📦 Ajout des modifications et commit...%s\n", ui.Cyan, ui.Reset)
	if err := runGit("add", "."); err != nil {
		return err
	}

	// `git diff --cached --quiet` renvoie 0 quand l'index est identique à
	// HEAD, 1 s'il y a des changements (tout autre code est une vraie erreur).
	switch exitCode, diffErr := gitExitCode("diff", "--cached", "--quiet"); {
	case diffErr != nil:
		return diffErr
	case exitCode == 0:
		fmt.Printf("%sℹ️  Aucun changement à committer. Rien poussé.%s\n", ui.Yellow, ui.Reset)
		return nil
	}

	if err := runGit("commit", "-m", commitMessage); err != nil {
		return err
	}

	fmt.Printf("%s✅ Commit effectué.%s\n", ui.Green, ui.Reset)

	return askAndPush(branchName)
}

// lint exécute en dry-run chaque outil PHP présent dans vendor/bin. Le
// binaire étant installé globalement, il doit rester utilisable partout :
// sans aucun outil, le lint est sauté avec un message.
func lint() error {
	steps := []struct {
		label string
		path  string
		args  []string
	}{
		{"🔧 Rector", "vendor/bin/rector", []string{"process", "--dry-run"}},
		{"🎨 PHP-CS-Fixer", "vendor/bin/php-cs-fixer", []string{"fix", "--dry-run", "--diff"}},
	}

	found := false
	for _, step := range steps {
		if _, err := os.Stat(step.path); err != nil {
			continue
		}
		found = true

		if err := requireBinary(step.path); err != nil {
			return err
		}
		if err := runStep(step.label, step.path, step.args...); err != nil {
			return err
		}
	}

	if !found {
		fmt.Printf("%sℹ️  Pas d'outils PHP (vendor/) dans ce dépôt — lint ignoré.%s\n", ui.Yellow, ui.Reset)
	}

	return nil
}

// unpushedCommits compte les commits locaux absents d'origin. Sans dépôt
// distant origin ou sans commit, il n'y a rien à pousser.
func unpushedCommits(branch string) (int, error) {
	if _, err := gitOutput("remote", "get-url", "origin"); err != nil {
		return 0, nil
	}
	if _, err := gitOutput("rev-parse", "--verify", "--quiet", "HEAD"); err != nil {
		return 0, nil
	}

	args := []string{"rev-list", "--count", "HEAD", "--not", "--remotes=origin"}
	remoteRef := "refs/remotes/origin/" + branch
	if _, err := gitOutput("rev-parse", "--verify", "--quiet", remoteRef); err == nil {
		args = []string{"rev-list", "--count", remoteRef + "..HEAD"}
	}

	out, err := gitOutput(args...)
	if err != nil {
		return 0, fmt.Errorf("comptage des commits non poussés : %w", err)
	}

	return strconv.Atoi(out)
}

// askAndPush ne pousse que sur confirmation. Un refus (« Non », Ctrl+C,
// EOF) n'est pas une erreur : les commits sont conservés en local.
func askAndPush(branch string) error {
	push, err := ui.Confirm(fmt.Sprintf("Pousser vers origin/%s ?", branch), true)
	if err != nil || !push {
		fmt.Printf(
			"%sℹ️  Rien poussé, le commit reste en local. Pour pousser : git push -u origin %s%s\n",
			ui.Yellow,
			branch,
			ui.Reset,
		)
		return nil
	}

	fmt.Printf("%s🚀 Push vers origin/%s...%s\n", ui.Cyan, branch, ui.Reset)
	if err := runGit("push", "-u", "origin", branch); err != nil {
		return err
	}
	fmt.Printf("%s✅ Push effectué.%s\n", ui.Green, ui.Reset)

	return nil
}

// Package internal Package commit est l'équivalent Go de tools/scripts/shell/commit.sh
// (port fidèle, le script bash reste en place) : outil de commit/push du
// dépôt racine yoostart — lint préalable (Rector + PHP-CS-Fixer en
// dry-run, log montré seulement en cas d'échec, avec spinner), puis
// git add/commit/push.
//
// Usage :
//
//	commit <message> <branche>
//
// Les couleurs ANSI sont désactivées si la sortie n'est pas un terminal
// ou si NO_COLOR est défini (Makefile / CI).
package internal

import (
	"fmt"
	"mzeahmed/commit/internal/ui"
	"os"
	"strings"
)

// Run orchestre le workflow complet : validation des arguments,
// positionnement à la racine du dépôt, contrôle de branche, vérification
// des outils, lint (Rector puis PHP-CS-Fixer) et commit/push.
func Run(args []string) error {
	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		fmt.Printf("%s❌ Merci de préciser le message du commit en argument.%s\n", ui.Red, ui.Reset)
		fmt.Printf("%sExemple : commit \"Mise à jour de la version de Yoostart\"%s\n", ui.Yellow, ui.Reset)
		os.Exit(1)
	}

	if len(args) < 2 || strings.TrimSpace(args[1]) == "" {
		fmt.Printf("%s❌ Merci de préciser le nom de la branche en argument.%s\n", ui.Red, ui.Reset)
		fmt.Printf("%sExemple : commit \"Mise à jour de la version de Yoostart\" main%s\n", ui.Yellow, ui.Reset)
		os.Exit(1)
	}

	commitMessage, branchName := args[0], args[1]

	// Se placer à la racine du repo (au cas où), comme le script bash.
	repoRoot, err := gitOutput("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("résolution de la racine du dépôt : %w", err)
	}
	if err := os.Chdir(repoRoot); err != nil {
		return fmt.Errorf("changement de répertoire vers %s : %w", repoRoot, err)
	}

	// Vérification de la branche courante.
	currentBranch, err := gitOutput("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return fmt.Errorf("lecture de la branche courante : %w", err)
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
		fmt.Printf("%sℹ️  Aucun changement à committer. Rien poussé.%s\n", ui.Yellow, ui.Reset)
		return nil
	}

	fmt.Printf("%s📋 Fichiers qui seront commités sur %s :%s\n", ui.Cyan, branchName, ui.Reset)
	if err := runGit("status", "--short", "--untracked-files=all"); err != nil {
		return err
	}

	confirmed, err := ui.Confirm("Valider le commit ?", true)
	if err != nil || !confirmed {
		// Annulation utilisateur (« No », Ctrl+C, EOF) : sortie propre
		// sans code d'erreur — l'outil a fait son travail, make reste
		// silencieux (`make commit` ne s'arrête donc pas sur ce cas).
		fmt.Printf("%sCommit annulé.%s\n", ui.Red, ui.Reset)
		return nil
	}

	// Lint PHP (Rector → PHP-CS-Fixer) : uniquement si le repo courant
	// porte des outils PHP (vendor/bin) — c'est le cas du monorepo
	// yoostart, pas des autres repos (mobile, v4/front, v4/api…). Le
	// binaire étant installé globalement (`/usr/local/bin/commit`), il
	// doit rester utilisable partout : sans vendor/, le lint est sauté
	// avec un message et le workflow continue directement au commit.
	_, rectorErr := os.Stat("vendor/bin/rector")
	if rectorErr == nil {
		// Vérification des outils requis (fail-fast avant toute étape).
		if err := requireBinary("vendor/bin/rector"); err != nil {
			return err
		}

		if err := requireBinary("vendor/bin/php-cs-fixer"); err != nil {
			return err
		}

		// Rector → PHP-CS-Fixer.
		if err := runStep("🔧 Rector", "vendor/bin/rector", "process", "--dry-run"); err != nil {
			return err
		}

		if err := runStep("🎨 PHP-CS-Fixer", "vendor/bin/php-cs-fixer", "fix", "--dry-run", "--diff"); err != nil {
			return err
		}
	} else {
		fmt.Printf("%sℹ️  Pas d'outils PHP (vendor/) dans ce dépôt — lint ignoré.%s\n", ui.Yellow, ui.Reset)
	}

	// Commit.
	fmt.Printf("%s📦 Ajout des modifications et commit...%s\n", ui.Cyan, ui.Reset)
	if err := runGit("add", "."); err != nil {
		return err
	}

	// Évite l'erreur si aucun changement : `git diff --cached --quiet`
	// renvoie 0 quand l'index est identique à HEAD, 1 s'il y a des
	// changements (tout autre code est une vraie erreur).
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

	// Push, uniquement sur confirmation. Un refus (« No », Ctrl+C, EOF)
	// n'est pas une erreur : le commit est conservé en local.
	push, err := ui.Confirm(fmt.Sprintf("Pousser vers origin/%s ?", branchName), true)
	if err != nil || !push {
		fmt.Printf(
			"%sℹ️  Rien poussé, le commit reste en local. Pour pousser : git push -u origin %s%s\n",
			ui.Yellow,
			branchName,
			ui.Reset,
		)
		return nil
	}

	fmt.Printf("%s🚀 Push vers origin/%s...%s\n", ui.Cyan, branchName, ui.Reset)
	if err := runGit("push", "-u", "origin", branchName); err != nil {
		return err
	}
	fmt.Printf("%s✅ Push effectué.%s\n", ui.Green, ui.Reset)

	return nil
}

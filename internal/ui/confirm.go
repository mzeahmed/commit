package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// Confirm pose une confirmation interactive Oui/Non à l'utilisateur via
// un `promptui.Select` (flèches haut/bas, Entrée pour valider).
//
// Le choix par défaut — l'élément pré-sélectionné — est piloté par
// `defaultYes` : les items sont ordonnés pour placer « Yes » en premier
// (sélection par défaut) s'il vaut true, « No » sinon. L'utilisateur
// voit toujours les deux options, seul le point de départ du curseur
// change ; appuyer sur Entrée sans bouger accepte donc le défaut.
//
// Renvoie true uniquement si l'utilisateur a sélectionné « Yes ».
// Renvoie une erreur si l'interaction échoue ou est interrompue —
// Ctrl+C (`promptui.ErrInterrupt`), EOF, terminal absent : dans ce cas
// le booléen n'est pas exploitable et l'appelant doit abandonner le
// workflow (c'est le contrat attendu par `internal.Run`, qui annule le
// commit sur toute erreur). Nécessite un terminal interactif : jamais
// utilisable en CI ou script non TTY.
func Confirm(label string, defaultYes bool) (bool, error) {
	items := []string{"Yes", "No"}
	if !defaultYes {
		items = []string{"No", "Yes"}
	}

	prompt := promptui.Select{
		Label: label,
		Items: items,
		Templates: &promptui.SelectTemplates{
			Selected: fmt.Sprintf(`{{ "%s" | green }} {{ "%s:" | bold }} {{ . | cyan }}`, promptui.IconGood, label),
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return result == "Yes", nil
}

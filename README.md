# Commit

> Une seule commande pour linter, commiter et pousser — depuis n'importe quel dépôt de ta machine.

[![CI](https://github.com/mzeahmed/commit/actions/workflows/ci.yml/badge.svg)](https://github.com/mzeahmed/commit/actions/workflows/ci.yml)
[![Release](https://github.com/mzeahmed/commit/actions/workflows/release.yml/badge.svg)](https://github.com/mzeahmed/commit/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mzeahmed/commit)](https://github.com/mzeahmed/commit/blob/main/go.mod)
[![Latest Release](https://img.shields.io/github/v/release/mzeahmed/commit)](https://github.com/mzeahmed/commit/releases/latest)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)

`commit` est une petite CLI en Go qui encapsule le cycle quotidien `git add / commit / push` derrière une commande interactive unique. Lance-la depuis n'importe quel répertoire : elle résout elle-même la racine du dépôt, contrôle la branche cible, exécute tes linters quand le projet en a, et pousse quand tout passe.

---

## Démarrage rapide

### 1. Installer

Une seule ligne (Linux, macOS) :

```bash
curl -fsSL https://raw.githubusercontent.com/mzeahmed/commit/main/install.sh | sh
```

Télécharge la dernière release pour ton système et ton architecture, vérifie son checksum et installe le binaire dans `/usr/local/bin` (sudo si nécessaire). Pour une version précise ou un autre dossier :

```bash
curl -fsSL https://raw.githubusercontent.com/mzeahmed/commit/main/install.sh | VERSION=v1.2.0 BINDIR=~/.local/bin sh
```

Voir [Installation](#installation) pour Windows, l'installation manuelle et la compilation depuis les sources.

### 2. Commiter

```bash
commit "Corriger le lien de l'avatar dans la messagerie" feature-branch
```

### 3. C'est tout

L'outil demande confirmation, lint si applicable, committe, puis demande s'il faut pousser.

---

## Fonctionnalités

- **Fonctionne partout** — lance `commit` depuis n'importe quel sous-dossier de n'importe quel dépôt Git ; la racine est résolue automatiquement
- **Confirmation interactive** — prompt Oui/Non avant que quoi que ce soit ne touche ton arbre de travail
- **Garde de branche** — refuse de pousser si la branche demandée n'est pas la branche courante (avertit et sort proprement)
- **Lint PHP intégré** — exécute Rector et PHP-CS-Fixer en dry-run quand `vendor/bin/` existe, avec un spinner ; la sortie n'est affichée qu'en cas d'échec
- **Non-PHP friendly** — dans les dépôts sans outillage PHP (front, mobile, …), le lint est sauté avec un message
- **No-op sûr** — quand il n'y a rien à committer, rien n'est poussé et le code de sortie reste propre
- **Pas de couleurs en CI** — la sortie ANSI est désactivée si stdout n'est pas un TTY ou si `NO_COLOR` est défini

---

## Installation

Des binaires pré-construits pour Linux, macOS et Windows sont disponibles sur la page des [GitHub Releases](https://github.com/mzeahmed/commit/releases/latest) (`amd64` ou `arm64`).

### Linux et macOS (script, recommandé)

```bash
curl -fsSL https://raw.githubusercontent.com/mzeahmed/commit/main/install.sh | sh
```

Le script détecte le système (`Linux`, `Darwin`) et l'architecture (`x86_64`, `arm64`), télécharge l'archive de la release, vérifie son checksum (`checksums.txt`) et installe `commit` dans `/usr/local/bin`, avec `sudo` si ce dossier n'est pas accessible en écriture. Variables :

- `VERSION` : release à installer (ex. `v1.2.0`), par défaut la dernière ;
- `BINDIR` : dossier d'installation, par défaut `/usr/local/bin`.

Relancer la même commande met à jour `commit`.

### macOS (manuel)

```bash
curl -LO https://github.com/mzeahmed/commit/releases/latest/download/commit_Darwin_arm64.tar.gz
tar -xzf commit_Darwin_arm64.tar.gz
sudo mv commit /usr/local/bin/
```

Utilise `commit_Darwin_x86_64.tar.gz` sur Mac Intel. Si macOS refuse d'exécuter le binaire (« cannot be opened because the developer cannot be verified »), retire le drapeau de quarantaine :

```bash
xattr -d com.apple.quarantine /usr/local/bin/commit
```

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://github.com/mzeahmed/commit/releases/latest/download/commit_Windows_x86_64.zip -OutFile commit.zip
Expand-Archive commit.zip -DestinationPath .
```

Utilise `commit_Windows_arm64.zip` sur les appareils ARM64. Déplace `commit.exe` vers un dossier de ton choix et ajoute-le au `PATH` pour lancer `commit` depuis n'importe quel terminal.

### Compiler depuis les sources

```bash
git clone https://github.com/mzeahmed/commit.git
cd commit
go install
```

---

## Utilisation

```bash
commit <message> <branche>
```

| Argument | Description |
|---|---|
| `message` | Le message de commit (mets-le entre guillemets s'il contient des espaces) |
| `branche` | La branche vers laquelle pousser — doit correspondre à la branche courante |

### Exemple de session

```text
📋 Fichiers qui seront commités sur develop :
 M internal/run.go
? Valider le commit ?  (← confirmation Oui / Non)
🔧 Rector...
✅ Rector OK.
🎨 PHP-CS-Fixer...
✅ PHP-CS-Fixer OK.
📦 Ajout des modifications et commit...
✅ Commit effectué.
? Pousser vers origin/develop ?  (← confirmation Oui / Non)
🚀 Push vers origin/develop...
✅ Push effectué.
```

### Comportement

| Situation | Résultat |
|---|---|
| Confirmation « No », Ctrl+C ou EOF | `Commit annulé.` — rien ne se passe, sortie propre (0) |
| Push refusé (« No », Ctrl+C ou EOF) | Commit conservé en local, rien n'est poussé, sortie propre (0) |
| Branche demandée ≠ branche courante | Avertissement, rien ne se passe, sortie propre (0) |
| Aucun changement à committer | Message (`Aucun changement à committer`), sortie propre (0) |
| Échec du lint | Log affiché en rouge, commit/push annulés, **exit 1** |
| Échec de git | Erreur, **exit 1** |
| Arguments manquants | Message d'usage, **exit 1** |

Quand tu chaînes l'outil dans un Makefile (`make commit`), changer d'avis ne doit pas faire échouer la chaîne de build avec une ligne
`make: *** Error 1`. Les vrais échecs gardent exit 1 pour rester chaînables (`&&`).

---

## Fonctionnement

1. **Racine du dépôt** — `git rev-parse --show-toplevel`, puis `chdir` vers elle.
2. **Garde de branche** — `git rev-parse --abbrev-ref HEAD` comparé à la branche demandée.
3. **Aperçu et confirmation** — liste des fichiers modifiés ou nouveaux, puis prompt Oui/Non (flèches, Entrée ; **Oui** par défaut). Sans changement, l'outil s'arrête ici.
4. **Lint** *(dépôts PHP uniquement)* — Rector puis PHP-CS-Fixer, tous deux en dry-run, exécutés en arrière-plan avec un spinner (`| / - \`) ; le log capturé n'est affiché que si une étape échoue.
5. **Commit** — `git add .`, sauté si l'index correspond à HEAD (`git diff --cached --quiet`), sinon `git commit`.
6. **Push** *(sur confirmation, **Oui** par défaut)* — `git push -u origin <branche>`.

---

## Prérequis

- Git
- Outillage PHP (`vendor/bin/rector`, `vendor/bin/php-cs-fixer` via Composer) — **seulement** si tu veux le lint ; l'absence d'outillage est détectée et le lint sauté
- Go (seulement pour compiler depuis les sources)

---

## Processus de release

Une GitHub Release est publiée automatiquement : pousse un tag `vX.Y.Z` et le [workflow Release](.github/workflows/release.yml) construit toutes les plateformes avec [GoReleaser](https://goreleaser.com), génère le changelog depuis les commits (en excluant `ci:`, `docs:`, `test:`, `chore:`) et attache les archives plus `checksums.txt`.

```bash
git tag v0.1.0
git push origin v0.1.0
```

---

## Structure du projet

```
main.go            point d'entrée : initialise les couleurs, délègue à internal/
internal/run.go    orchestration : arguments, racine du dépôt, garde de branche, lint, commit, push
internal/git.go    fins wrappers git (capture de sortie, streaming, codes de sortie)
internal/step.go   exécution d'étapes en arrière-plan + spinner
internal/tool.go   vérification des binaires requis
internal/ui/       couleurs ANSI + confirmation Oui/Non (promptui)
```

---

## Philosophie

- **Ne faire qu'une chose** — lint, commit, push. Rien d'autre.
- **Sûr par défaut** — rien ne touche le dépôt avant ta confirmation.
- **Silencieux en cas de succès** — les logs de lint ne sont montrés que quand ils comptent (échec).
- **Un no-op n'est pas une erreur** — changer d'avis ne doit jamais casser ta chaîne de build.

---

## Licence

[MIT](LICENSE). Historique des versions : [CHANGELOG](CHANGELOG.md).

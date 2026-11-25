# Changelog

Toutes les modifications notables de ce projet sont documentées ici.
Le format suit [Keep a Changelog](https://keepachangelog.com/fr/1.1.0/) et le projet adhère au [versionnement sémantique](https://semver.org/lang/fr/).

## [Unreleased]

### Ajouté
- Liste des fichiers qui seront commités, affichée avant la confirmation.

### Modifié
- La branche est vérifiée avant toute question, et l'absence de changement arrête l'outil avant le lint.
- Push avec `-u` : le suivi de la branche distante est configuré.

## [0.1.5] - 2026-11-25

### Modifié
- Après le commit, une confirmation demande s'il faut pousser (Oui par défaut). En cas de refus, le commit reste en local.

### Supprimé
- Étape Composer (`composer update --no-dev` puis `composer install`), annoncée comme retirée en 0.1.4 mais toujours présente dans le code de cette version.

## [0.1.4] - 2026-11-24

### Ajouté
- Script `install.sh` pour une installation en une commande `curl` (Linux, macOS) : détection du système et de l'architecture, vérification du checksum via `checksums.txt`, variables `VERSION` et `BINDIR`.
- Fichier `LICENSE` (MIT) et ce `CHANGELOG.md`.

### Modifié
- README : installation par `curl` via `install.sh`.
- `make install` retire le binaire local (`bin/`) après l'installation dans `/usr/local/bin`.

### Supprimé
- Étape Composer avant le commit (`composer update --no-dev`, puis `composer install`), restée active en 0.1.2 et 0.1.3 malgré le retrait annoncé.

## [0.1.3] - 2026-11-23

### Modifié
- Retrait annoncé de l'étape Composer (#2). La fusion n'a apporté aucune modification de code : l'étape Composer est toujours exécutée dans cette version.

## [0.1.2] - 2026-11-23

### Ajouté
- Étape Composer, quand `composer` est disponible : `composer update --no-dev` juste avant le commit (le `composer.lock` mis à jour est inclus), puis `composer install` à la fin pour réinstaller les dépendances de dev (#1).

## [0.1.1] - 2026-11-16

### Ajouté
- `Makefile` de développement (build, install, fmt, vet, test, check, tidy, release).

### Corrigé
- Les arguments sont validés avant d'afficher la confirmation (message d'erreur plus explicite).

## [0.1.0] - 2026-11-16

### Ajouté
- Première version de la CLI `commit` : `git add / commit / push` derrière une commande interactive unique, résolution de la racine du dépôt, garde de branche, lint Rector et PHP-CS-Fixer quand `vendor/bin/` existe.
- README.

[Unreleased]: https://github.com/mzeahmed/commit/compare/v0.1.5...HEAD
[0.1.5]: https://github.com/mzeahmed/commit/compare/v0.1.4...v0.1.5
[0.1.4]: https://github.com/mzeahmed/commit/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/mzeahmed/commit/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/mzeahmed/commit/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/mzeahmed/commit/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/mzeahmed/commit/releases/tag/v0.1.0

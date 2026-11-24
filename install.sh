#!/bin/sh
# Installe la dernière release de commit dans /usr/local/bin (Linux, macOS).
#   curl -fsSL https://raw.githubusercontent.com/mzeahmed/commit/main/install.sh | sh
# Variables : BINDIR (défaut /usr/local/bin), VERSION (défaut : dernière release, ex. v1.2.0)
set -eu

REPO="mzeahmed/commit"
BINDIR="${BINDIR:-/usr/local/bin}"

case "$(uname -s)" in
  Linux) OS="Linux" ;;
  Darwin) OS="Darwin" ;;
  *) echo "Erreur : système non supporté : $(uname -s) (Windows : voir le README)." >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64 | amd64) ARCH="x86_64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *) echo "Erreur : architecture non supportée : $(uname -m)" >&2; exit 1 ;;
esac

ARCHIVE="commit_${OS}_${ARCH}.tar.gz"
if [ -n "${VERSION:-}" ]; then
  BASE="https://github.com/${REPO}/releases/download/${VERSION}"
else
  BASE="https://github.com/${REPO}/releases/latest/download"
fi

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "Téléchargement de ${ARCHIVE}..."
curl -fsSL "${BASE}/${ARCHIVE}" -o "${TMP}/${ARCHIVE}"
curl -fsSL "${BASE}/checksums.txt" -o "${TMP}/checksums.txt"

echo "Vérification du checksum..."
if command -v sha256sum >/dev/null 2>&1; then
  (cd "$TMP" && grep " ${ARCHIVE}\$" checksums.txt | sha256sum -c -)
else
  (cd "$TMP" && grep " ${ARCHIVE}\$" checksums.txt | shasum -a 256 -c -)
fi

tar -xzf "${TMP}/${ARCHIVE}" -C "$TMP" commit

SUDO=""
[ -w "$BINDIR" ] || SUDO="sudo"
$SUDO install -m 755 "${TMP}/commit" "${BINDIR}/commit"

echo "commit installé dans ${BINDIR}/commit"

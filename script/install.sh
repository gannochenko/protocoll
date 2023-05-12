#!/usr/bin/env bash

set -eu

FOLDER_PATH="/.protocoll/bin"
FILE_PATH="/protocoll"

DST="${HOME}${FOLDER_PATH}"
DST_FILE="${DST}${FILE_PATH}"
VERSION="v0.0.5"

check_command() {
  command -v "$1" > /dev/null 2>&1
  return $?
}

ensure_command() {
  if ! check_command "$1"; then
    echo "Command not found: '$1'"
    exit 1
  fi
}

ensure_success() {
  $* > /dev/null

  if [[ $? -ne 0 ]]; then
    echo "Command execution failed: $*"
    exit 1
  fi
}

download_and_run() {
  ensure_success mkdir -p "$DST"

  local url="https://github.com/gannochenko/protocoll/releases/download/v${VERSION}/protocoll-darwin-amd64"

  rm -f "${DST_FILE}"

  ensure_success curl -LJ -o "${DST_FILE}" ${url}
  ensure_success chmod +x "$DST_FILE"
}

echo "Installing protocoll, version v${VERSION}"
echo ""

ensure_command mkdir
ensure_command curl

download_and_run || exit 1

echo ""
echo "Success. Make sure to add this line to your rc file: "
echo "export PATH=\$HOME${FOLDER_PATH}:\$PATH;"

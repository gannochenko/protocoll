#!/usr/bin/env bash

set -eu

FOLDER_PATH="/.protocoll/bin"
FILE_PATH="/protocoll"

DST="${HOME}${FOLDER_PATH}"
DST_FILE="${DST}${FILE_PATH}"
VERSION="#VERSION#"

error_and_exit() {
  echo "$1"
  exit 1
}

check_command() {
  command -v "$1" > /dev/null 2>&1
  return $?
}

ensure_command() {
  if ! check_command "$1"; then
    error_and_exit "Command not found: '$1'"
  fi
}

ensure_success() {
  $* > /dev/null

  if [[ $? -ne 0 ]]; then
    error_and_exit "Command execution failed: $*"
  fi
}

get_architecture() {
  local os="$(uname -s)"
  local cpu="$(uname -m)"

  case "$os" in
      Linux)
          local os="linux"
          ;;
      Darwin)
          local os="darwin"
          ;;
      MINGW* | MSYS* | CYGWIN*)
          local os="windows"
          ;;
      *)
          error_and_exit "OS type is not supported: $os"
          ;;
  esac

  case "$cpu" in
      x86_64)
          local cpu="amd64"
          ;;
      *)
          error_and_exit "CPU architecture is not supported: $cpu"
          ;;
  esac

  RETVAL="${os}-${cpu}"
}

download_and_run() {
  ensure_success mkdir -p "$DST"

  get_architecture
  local architecture=$RETVAL

  echo $architecture

  local url="https://github.com/gannochenko/protocoll/releases/download/v${VERSION}/protocoll-${architecture}"

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

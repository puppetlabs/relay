#!/usr/bin/env bash

#
# Commands
#

FIND="${FIND:-find}"
GIT="${GIT:-git}"
GSUTIL="${GSUTIL:-gsutil}"
SHA256SUM="${SHA256SUM:-shasum -a 256}"

#
#
#

relay::cli::default_programs() {
  local DEFAULT_PROGRAMS
  DEFAULT_PROGRAMS=( relay )

  for DEFAULT_PROGRAM in ${DEFAULT_PROGRAMS[@]}; do
    printf "%s\n" "${DEFAULT_PROGRAM}"
  done
}

relay::cli::git_tag() {
  printf "%s\n" "${GIT_TAG_OVERRIDE:-$( $GIT tag --points-at HEAD 'v*.*.*' )}"
}

relay::cli::sha256sum() {
  $SHA256SUM | cut -d' ' -f1
}

relay::cli::escape_shell() {
  printf '%s\n' "'${*//\'/\'\"\'\"\'}'"
}

relay::cli::release_version() {
  local GIT_TAG GIT_CHANGED_FILES
  GIT_TAG="$( relay::cli::git_tag )"
  GIT_CHANGED_FILES="$( $GIT status --short )"

  # Check for releasable version: if we have no tags or any changed files, we
  # can't release.
  if [ -z "${GIT_TAG}" ] || [ -n "${GIT_CHANGED_FILES}" ]; then
    return 1
  fi

  # Arbitrarily pick the first line.
  read GIT_TAG_A <<<"${GIT_TAG}"

  printf "%s\n" "${GIT_TAG_A#v}"
}

relay::cli::release_check() {
  if ! relay::cli::release_version >/dev/null; then
    echo "$0: no release tag (this commit must be tagged with the format vX.Y.Z)" >&2
    return 2
  fi
}

relay::cli::release_vars() {
  RELEASE_VERSION="$( relay::cli::release_version || true )"
  if [ -z "${RELEASE_VERSION}" ]; then
    printf 'RELEASE_VERSION=\n'
    return
  fi

  # Parse the version information.
  IFS='.' read RELEASE_VERSION_MAJOR RELEASE_VERSION_MINOR RELEASE_VERSION_PATCH <<<"${RELEASE_VERSION}"

  printf 'RELEASE_VERSION=%s\n' "$( relay::cli::escape_shell "${RELEASE_VERSION}" )"
  printf 'RELEASE_VERSION_MAJOR=%s\n' "$( relay::cli::escape_shell "${RELEASE_VERSION_MAJOR}" )"
  printf 'RELEASE_VERSION_MINOR=%s\n' "$( relay::cli::escape_shell "${RELEASE_VERSION_MINOR}" )"
  printf 'RELEASE_VERSION_PATCH=%s\n' "$( relay::cli::escape_shell "${RELEASE_VERSION_PATCH}" )"
}

relay::cli::release_vars_local() {
  printf 'local RELEASE_VERSION RELEASE_VERSION_MAJOR RELEASE_VERSION_MINOR RELEASE_VERSION_PATCH\n'
  relay::cli::release_vars "$@"
}

relay::cli::release() {
  if [[ "$#" -lt 2 ]]; then
    echo "usage: ${FUNCNAME[0]} <release-name> <filename> [dist-ext [dist-prefix]]" >&2
    return 1
  fi

  relay::cli::release_check
  eval "$( relay::cli::release_vars )"

}

relay::cli::version() {
  eval "$( relay::cli::release_vars )"

  if [ -n "${RELEASE_VERSION}" ]; then
    printf "%s\n" "v${RELEASE_VERSION}"
  else
    $GIT describe --tags --always --dirty
  fi
}

relay::cli::cli_vars() {
  local GO GOOS GOARCH
  GO="${GO:-go}"
  GOOS="$( $GO env GOOS )"
  GOARCH="$( $GO env GOARCH )"

  local EXT=
  [[ "${GOOS}" == "windows" ]] && EXT=.exe

  printf 'CLI_NAME=%s\n' "$( echo relay )"
  printf 'CLI_VERSION=%s\n' "$( relay::cli::version )"
  printf 'CLI_FILE_PREFIX="${CLI_NAME}-${CLI_VERSION}"-%s-%s\n' \
    "$( relay::cli::escape_shell "${GOOS}" )" \
    "$( relay::cli::escape_shell "${GOARCH}" )"
  printf 'CLI_FILE_BIN="${CLI_FILE_PREFIX}%s"\n' "${EXT}"
}

relay::cli::cli_vars_local() {
  printf 'local CLI_NAME CLI_FILE_PREFIX CLI_FILE_BIN\n'
  relay::cli::cli_vars "$@"
}

relay::cli::cli_artifacts() {
  if [[ "$#" -ne 2 ]]; then
    echo "usage: ${FUNCNAME[0]} <program> <directory>" >&2
    return 1
  fi

  eval "$( relay::cli::cli_vars_local "$1" )"

  local CLI_MATCH
  CLI_MATCH="${CLI_NAME}-${CLI_VERSION}-"

  $FIND "$2" -type f -name "${CLI_MATCH}"'*'
}

relay::cli::cli_platform_ext() {
  if [[ "$#" -ne 2 ]]; then
    echo "usage: ${FUNCNAME[0]} <program> <package-file>" >&2
    return 1
  fi

  eval "$( relay::cli::cli_vars_local "$1" )"

  local CLI_FILE
  CLI_FILE="$( basename "$2" )"

  printf "%s\n" "${CLI_FILE##${CLI_NAME}-${CLI_VERSION}-}"
}

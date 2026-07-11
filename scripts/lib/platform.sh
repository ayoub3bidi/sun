#!/usr/bin/env sh
# Shared platform detection for sun install/run scripts.
set -e

sun_cache_dir() {
  if [ -n "${XDG_CACHE_HOME:-}" ]; then
    printf '%s/sun' "$XDG_CACHE_HOME"
  else
    printf '%s/.cache/sun' "${HOME}"
  fi
}

sun_detect_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$os" in
    linux) printf 'linux' ;;
    darwin) printf 'darwin' ;;
    *)
      echo "unsupported OS: $os (sun supports Linux and macOS)" >&2
      exit 1
      ;;
  esac
}

sun_detect_arch() {
  arch=$(uname -m)
  case "$arch" in
    x86_64|amd64) printf 'amd64' ;;
    arm64|aarch64) printf 'arm64' ;;
    *)
      echo "unsupported architecture: $arch" >&2
      exit 1
      ;;
  esac
}

sun_repo() {
  printf '%s' "${SUN_REPO:-ayoub3bidi/sun}"
}

sun_latest_version() {
  latest=$(curl -fsSL https://api.github.com/repos/ayoub3bidi/sun/releases/latest 2>/dev/null \
    | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')
  if [ -n "$latest" ]; then
    printf '%s' "$latest"
  else
    printf '%s' "v1.0.0"
  fi
}

sun_version() {
  if [ -n "${SUN_VERSION:-}" ]; then
    printf '%s' "$SUN_VERSION"
  else
    sun_latest_version
  fi
}

sun_version_raw() {
  version=$(sun_version)
  printf '%s' "${version#v}"
}

sun_asset_name() {
  os=$1
  arch=$2
  printf 'sun_%s_%s_%s.tar.gz' "$(sun_version_raw)" "$os" "$arch"
}

sun_download_url() {
  os=$1
  arch=$2
  repo=$(sun_repo)
  version=$(sun_version)
  asset=$(sun_asset_name "$os" "$arch")
  printf 'https://github.com/%s/releases/download/%s/%s' "$repo" "$version" "$asset"
}

#!/usr/bin/env sh
# Optional persistent install of sun to ~/.local/bin
set -e

# When piped to sh, $0 is the shell and there is no file to source from.
# Re-exec from a downloaded copy so lib/platform.sh is available.
INSTALL_URL="https://raw.githubusercontent.com/ayoub3bidi/sun/main/scripts"
if [ ! -f "$0" ] || [ ! -r "$0" ]; then
  _tmpdir=$(mktemp -d)
  trap 'rm -rf "$_tmpdir"' EXIT INT HUP
  mkdir -p "${_tmpdir}/lib"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "${INSTALL_URL}/lib/platform.sh" -o "${_tmpdir}/lib/platform.sh"
    curl -fsSL "${INSTALL_URL}/install.sh" -o "${_tmpdir}/install.sh"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "${_tmpdir}/lib/platform.sh" "${INSTALL_URL}/lib/platform.sh"
    wget -qO "${_tmpdir}/install.sh" "${INSTALL_URL}/install.sh"
  else
    echo "curl or wget is required" >&2
    exit 1
  fi
  exec sh "${_tmpdir}/install.sh" "$@"
fi

. "$(dirname "$0")/lib/platform.sh"

SUN_VERSION="${SUN_VERSION:-$(sun_latest_version)}"
SUN_VERSION_RAW="${SUN_VERSION#v}"

sun_install_to_cache() {
  os=$(sun_detect_os)
  arch=$(sun_detect_arch)
  install_dir="$(sun_cache_dir)/${SUN_VERSION}"
  bin_path="${install_dir}/sun"

  if [ -x "$bin_path" ]; then
    printf '%s' "$bin_path"
    return 0
  fi

  mkdir -p "$install_dir"
  url=$(sun_download_url "$os" "$arch")

  tmp=$(mktemp -d)
  trap 'rm -rf "$tmp"' EXIT INT HUP

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "${tmp}/sun.tar.gz"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "${tmp}/sun.tar.gz" "$url"
  else
    echo "curl or wget is required" >&2
    exit 1
  fi

  tar -xzf "${tmp}/sun.tar.gz" -C "$tmp"
  if [ -f "${tmp}/sun" ]; then
    mv "${tmp}/sun" "$bin_path"
  elif [ -f "${tmp}/sun_${SUN_VERSION_RAW}_${os}_${arch}/sun" ]; then
    mv "${tmp}/sun_${SUN_VERSION_RAW}_${os}_${arch}/sun" "$bin_path"
  else
    echo "sun binary not found in release archive" >&2
    exit 1
  fi
  chmod +x "$bin_path"
  printf '%s' "$bin_path"
}

bin=$(sun_install_to_cache)
target="${SUN_INSTALL_DIR:-${HOME}/.local/bin}"
mkdir -p "$target"
cp "$bin" "${target}/sun"
chmod +x "${target}/sun"

echo "sun installed to ${target}/sun"
case ":${PATH}:" in
  *":${target}:"*) ;;
  *) echo "Add to PATH: export PATH=\"${target}:\$PATH\"" ;;
esac
echo "Run: sun create"

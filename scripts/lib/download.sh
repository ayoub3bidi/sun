#!/usr/bin/env sh
# Download and verify sun release binary into cache.
set -e

. "$(dirname "$0")/platform.sh"

sun_install_to_cache() {
  os=$(sun_detect_os)
  arch=$(sun_detect_arch)
  version=$(sun_version)
  cache_root=$(sun_cache_dir)
  install_dir="${cache_root}/${version}"
  bin_path="${install_dir}/sun"
  mkdir -p "$install_dir"

  if [ -x "$bin_path" ]; then
    printf '%s' "$bin_path"
    return 0
  fi

  url=$(sun_download_url "$os" "$arch")
  checksum_url="${url%.tar.gz}.sha256"

  tmp=$(mktemp -d)
  trap 'rm -rf "$tmp"' EXIT INT HUP

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "${tmp}/sun.tar.gz"
    if curl -fsSL "$checksum_url" -o "${tmp}/sun.tar.gz.sha256" 2>/dev/null; then
      (
        cd "$tmp" && sha256sum -c sun.tar.gz.sha256
      )
    fi
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "${tmp}/sun.tar.gz" "$url"
  else
    echo "curl or wget is required" >&2
    exit 1
  fi

  tar -xzf "${tmp}/sun.tar.gz" -C "$tmp"
  if [ -f "${tmp}/sun" ]; then
    mv "${tmp}/sun" "$bin_path"
  elif [ -f "${tmp}/sun_${version#v}_${os}_${arch}/sun" ]; then
    mv "${tmp}/sun_${version#v}_${os}_${arch}/sun" "$bin_path"
  else
    echo "could not find sun binary in release archive" >&2
    exit 1
  fi
  chmod +x "$bin_path"
  printf '%s' "$bin_path"
}

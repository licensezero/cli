#!/bin/sh
set -e

fail() {
  echo "$1" > /dev/stderr
  exit 1
}

# Determine Machine Architecture
architecture=$(uname -m)
case $architecture in
  amd64) architecture="amd64";;
  x86_64) architecture="amd64";;
  i386) architecture="i386";;
  *) fail "Could not identify architecture.";;
esac

# Determine Operating System
os="$(uname | tr '[:upper:]' '[:lower:]')"
case "$os" in
  darwin) os="darwin";;
  linux) os="linux";;
  freebsd) os="freebsd";;
  mingw*) os="windows";;
  msys*) os="windows";;
  *) fail "Could not identify operating system.";;
esac

[ "$os" = "windows" ] && executable_suffix=".exe"

GITHUB_REPOSITORY="https://github.com/licensezero/cli"

# Find Latest Release
ACCEPT_JSON="Accept: application/json"
release_url="$GITHUB_REPOSITORY/releases/latest"
if [ -x "$(command -v curl)" ]; then
  response=$(curl --silent --location --write-out 'HTTPSTATUS:%{http_code}' --header "$ACCEPT_JSON" "$release_url")
  releases=$(echo "$response" | sed -e 's/HTTPSTATUS\:.*//g')
  code=$(echo "$response" | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
elif [ -x "$(command -v wget)" ]; then
  error_output=$(mktemp)
  trap "rm -f $error_output" EXIT
  releases=$(wget --quiet --header="$ACCEPT_JSON" -O - --server-response "$release_url" 2> "$error_output")
  code=$(awk '/^  HTTP/{print $2}' < "$error_output" | tail -1)
else
  fail "Could not find curl or wget to download release information."
fi
[ "$code" != 200 ] && fail "Release data request failed with status $code."
tag="$(echo "$releases" | tr -s '\n' ' ' | sed 's/.*"tag_name":"//' | sed 's/".*//')"

echo "Latest Release: $tag"

# Download Executable
executable="licensezero-${os}-${architecture}${executable_suffix}"
tmp=$(mktemp)
trap "rm -f $tmp" EXIT
executable_url="$GITHUB_REPOSITORY/releases/download/$tag/$executable"
if [ -x "$(command -v curl)" ]; then
  statusCode=$(curl -s -w '%{http_code}' -L "$executable_url" -o "$tmp")
elif [ -x "$(command -v wget)" ]; then
  statusCode=$(wget -q -O "$tmp" --server-response "$executable_url" 2>&1 | awk '/^  HTTP/{print $2}' | tail -1)
else
  fail "Could not find curl or wget to download files."
fi
[ "$statusCode" != 200 ] && fail "Error: github.com responed $statusCode."

# Install Executable
chmod +x "$tmp"
if [ ! -z "$BINDIR" ]; then
  install_path="$BINDIR"
elif [ ! -z "$PREFIX" ]; then
  install_path="$PREFIX/bin"
elif [ -d "$HOME/.local/bin" ]; then
  install_path="$HOME/.local/bin"
elif [ -d "$HOME/local/bin" ]; then
  install_path="$HOME/local/bin"
elif [ -d "$HOME/bin" ]; then
  install_path="$HOME/bin"
elif [ -d "/usr/local/bin" ]; then
  install_path="/usr/local/bin"
else
  echo "Could not determine where to install the command line interface." >/dev/stderr
  echo "Downloading to the current directory." >/dev/stderr
  install_path="$PWD"
fi
install_path="$install_path/licensezero$executable_suffix"
if [ -w "$(dirname "$install_path")" ]; then
  if [ -f "$install_path" ] && ! [ -w "$install_path" ]; then
    sudo mv "$tmp" "$install_path"
  else
    mv "$tmp" "$install_path"
  fi
else
  sudo mv "$tmp" "$install_path"
fi
echo "Installed To: $install_path"

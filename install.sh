#!/usr/bin/env bash
set -e

REPO_OWNER="adamdlear"
REPO_NAME="nvmgr"
INSTALL_DIR="$HOME/.local/bin"
BIN_NAME="nvmgr"

# Detect platform
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
  Darwin)  OS="Darwin" ;;
  Linux)   OS="Linux" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH="x86_64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Build download URL
TARBALL="${REPO_NAME}_${OS}_${ARCH}.tar.gz"
LATEST_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/latest/download/${TARBALL}"

echo "Downloading $TARBALL..."
curl -fsSL "$LATEST_URL" -o "/tmp/$TARBALL"

echo "Extracting archive..."
tar -xzf "/tmp/$TARBALL" -C /tmp

# Ensure install directory exists
mkdir -p "$INSTALL_DIR"

echo "Installing to $INSTALL_DIR..."
install -m 755 "/tmp/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"

# Check PATH
if ! echo "$PATH" | grep -q "$HOME/.local/bin"; then
  echo "Adding ~/.local/bin to PATH..."

  # Bash
  if [ -f "$HOME/.bashrc" ]; then
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
  fi

  # Zsh
  if [ -f "$HOME/.zshrc" ]; then
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.zshrc"
  fi

  # Fish
  if [ -d "$HOME/.config/fish" ]; then
    echo 'set -gx PATH $HOME/.local/bin $PATH' >> "$HOME/.config/fish/config.fish"
  fi

  ADDED_PATH=true
else
  ADDED_PATH=false
fi

echo ""
echo "Installed $BIN_NAME successfully."
echo ""
echo "Location: $INSTALL_DIR/$BIN_NAME"
echo ""

if [ "$ADDED_PATH" = true ]; then
  echo "PATH updated. Restart your shell or run:"
  echo ""
  echo "    source ~/.bashrc    # bash"
  echo "    source ~/.zshrc     # zsh"
  echo "    exec fish           # fish"
  echo ""
fi

echo "Run the program with:"
echo ""
echo "    $BIN_NAME"
echo ""


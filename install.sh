#!/bin/bash

# Comment Remover Installation Script for Unix/Linux/macOS
# This script builds and installs 'commenter' to /usr/local/bin

set -e # Exit on any error

BINARY_NAME="commenter"
INSTALL_DIR="/usr/local/bin"
BUILD_DIR="build"

echo "🚀 Installing Comment Remover (${BINARY_NAME})..."

# Check if Go is installed
if ! command -v go &>/dev/null; then
  echo "❌ Go is not installed. Please install Go from https://golang.org/"
  exit 1
fi

# Create build directory
mkdir -p "$BUILD_DIR"

echo "🔨 Building ${BINARY_NAME}..."
go build -ldflags "-s -w" -o "$BUILD_DIR/$BINARY_NAME" .

# Check if the build was successful
if [ ! -f "$BUILD_DIR/$BINARY_NAME" ]; then
  echo "❌ Build failed!"
  exit 1
fi

echo "📦 Installing to ${INSTALL_DIR}..."

# Check if we need sudo
if [ -w "$INSTALL_DIR" ]; then
  cp "$BUILD_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
  chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
  sudo cp "$BUILD_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
  sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

echo "✅ Successfully installed ${BINARY_NAME} to ${INSTALL_DIR}!"
echo ""
echo "Usage:"
echo "  ${BINARY_NAME} <file>           # Preview comment removal"
echo "  ${BINARY_NAME} -r src/          # Process directory recursively"
echo "  ${BINARY_NAME} -w <file>        # Remove comments and save (short)"
echo "  ${BINARY_NAME} -w -r project/   # Process and save recursively"
echo "  ${BINARY_NAME} -h               # Show help (short)"
echo ""
echo "Supported file types: .ts, .tsx, .js, .jsx, .go, .sql, .json"
echo ""
echo "To uninstall, run: sudo rm ${INSTALL_DIR}/${BINARY_NAME}"

#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2025-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

#
# Start dogfooding Context CLI in a fresh project
#
# Usage: ./hack/start-dogfood.sh [--no-git] <target-folder>
#   target-folder: The directory to initialize (e.g., ~/WORKSPACE/ctx-dogfood)
#   --no-git:      Skip git repository initialization
#
# Prerequisites:
#   - ctx must be in your PATH (run: make install)
#
# This script:
#   1. Verifies ctx is in PATH
#   2. Creates the target folder
#   3. Initializes git repo (unless --no-git)
#   4. Runs ctx init
#   5. Copies specs/ for reference
#   6. Copies TASKS_DOGFOOD.md with rebuild goals
#   7. Copies PROMPT.md for Ralph Loop usage

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get script directory (where ctx source lives)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE_DIR="$(dirname "$SCRIPT_DIR")"

# Parse arguments
INIT_GIT=true
TARGET_DIR=""

while [[ $# -gt 0 ]]; do
  case $1 in
    --no-git)
      INIT_GIT=false
      shift
      ;;
    -*)
      echo -e "${RED}Error: Unknown option $1${NC}"
      exit 1
      ;;
    *)
      TARGET_DIR="$1"
      shift
      ;;
  esac
done

if [ -z "$TARGET_DIR" ]; then
  echo -e "${RED}Error: Target folder required${NC}"
  echo ""
  echo "Usage: $0 [--no-git] <target-folder>"
  echo "  Example: $0 ~/WORKSPACE/ctx-dogfood"
  echo "  Example: $0 --no-git ~/WORKSPACE/ctx-dogfood"
  exit 1
fi

# Expand ~ to home directory
TARGET_DIR="${TARGET_DIR/#\~/$HOME}"

echo "========================================="
echo "Context CLI Dogfooding Setup"
echo "========================================="
echo ""

# Step 1: Verify ctx is in PATH
if ! command -v ctx &> /dev/null; then
  echo -e "${RED}Error: ctx is not in your PATH${NC}"
  echo ""
  echo "Dogfooding requires ctx to be installed globally, just like a real user would have it."
  echo ""
  echo "To install ctx:"
  echo "  1. Build:   make build"
  echo "  2. Install: sudo make install"
  echo ""
  echo "Or manually:"
  echo "  sudo cp ./ctx /usr/local/bin/"
  echo ""
  echo "Then try again."
  exit 1
fi

CTX_VERSION=$(ctx --version 2>/dev/null || echo "unknown")
echo -e "${GREEN}Found ctx in PATH:${NC} $(command -v ctx)"
echo -e "${GREEN}Version:${NC} ${CTX_VERSION}"
echo ""
echo "Source:  ${SOURCE_DIR}"
echo "Target:  ${TARGET_DIR}"
echo ""

# Step 2: Create target folder
if [ -d "$TARGET_DIR" ]; then
  echo -e "${YELLOW}Warning: Target folder already exists${NC}"
  read -p "Continue and reinitialize? (y/N) " -n 1 -r
  echo ""
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 0
  fi
else
  echo "Creating target folder..."
  mkdir -p "$TARGET_DIR"
fi

# Step 3: Change to target folder
cd "$TARGET_DIR"
echo -e "${GREEN}Working directory:${NC} $(pwd)"
echo ""

# Step 4: Initialize git repo if not already (unless --no-git)
if [ "$INIT_GIT" = true ] && [ ! -d ".git" ]; then
  echo "Initializing git repository..."
  git init -b main
  echo ""
elif [ "$INIT_GIT" = false ]; then
  echo -e "${YELLOW}Skipping git init (--no-git)${NC}"
  echo ""
fi

# Step 5: Run ctx init
echo "Running ctx init..."
ctx init --merge
echo ""

# Step 6: Copy specs/ for reference
if [ -d "${SOURCE_DIR}/specs" ]; then
  echo "Copying specs/ for reference..."
  cp -r "${SOURCE_DIR}/specs" "./specs"
  echo -e "${GREEN}Installed:${NC} specs/ ($(find specs -maxdepth 1 -name "*.md" 2>/dev/null | wc -l) spec files)"
  echo ""
fi

# Step 7: Copy TASKS_DOGFOOD.md with rebuild goals
if [ -f "${SCRIPT_DIR}/TASKS_DOGFOOD.md" ]; then
  echo "Setting up dogfood tasks..."
  cp "${SCRIPT_DIR}/TASKS_DOGFOOD.md" ".context/TASKS.md"
  echo -e "${GREEN}Installed:${NC} .context/TASKS.md (ctx rebuild goals)"
  echo ""
fi

# Step 8: Copy PROMPT.md for Ralph Loop
if [ -f "${SOURCE_DIR}/PROMPT.md" ]; then
  echo "Copying PROMPT.md for Ralph Loop..."
  cp "${SOURCE_DIR}/PROMPT.md" "./PROMPT.md"
  echo -e "${YELLOW}Note: You may want to customize PROMPT.md for your project${NC}"
  echo ""
fi

# Done
echo "========================================="
echo -e "${GREEN}Dogfooding setup complete!${NC}"
echo "========================================="
echo ""
echo "Next steps:"
echo "  1. cd ${TARGET_DIR}"
echo "  2. Edit PROMPT.md to describe your project"
echo "  3. Run: claude --dangerously-skip-permissions"
echo "  4. Start Ralph Loop: /ralph-loop --file PROMPT.md"
echo ""
echo "Or run ctx commands directly:"
echo "  ctx status"
echo "  ctx agent --budget 4000"
echo ""

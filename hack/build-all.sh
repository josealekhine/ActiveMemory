#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

#
# Build Context CLI for multiple platforms
#
# Usage: ./hack/build-all.sh [version]
#   version: The version string to embed (default: dev)
#
# Output: Binaries are placed in ./dist/

set -e

VERSION="${1:-dev}"
OUTPUT_DIR="dist"
BINARY_NAME="ctx"
MODULE_PATH="./cmd/ctx"

# Build targets: OS/ARCH pairs
TARGETS=(
  "darwin/amd64"
  "darwin/arm64"
  "linux/amd64"
  "linux/arm64"
  "windows/amd64"
  "windows/arm64"
)

echo "Building Context CLI v${VERSION}"
echo "========================================="

# Clean and create output directory (preserve RELEASE_NOTES.md if it exists)
if [ -f "${OUTPUT_DIR}/RELEASE_NOTES.md" ]; then
  mv "${OUTPUT_DIR}/RELEASE_NOTES.md" /tmp/RELEASE_NOTES.md.bak
fi
rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}"
if [ -f /tmp/RELEASE_NOTES.md.bak ]; then
  mv /tmp/RELEASE_NOTES.md.bak "${OUTPUT_DIR}/RELEASE_NOTES.md"
fi

# Build for each target
for target in "${TARGETS[@]}"; do
  GOOS="${target%/*}"
  GOARCH="${target#*/}"

  output_name="${BINARY_NAME}-${VERSION}-${GOOS}-${GOARCH}"
  if [ "${GOOS}" = "windows" ]; then
    output_name="${output_name}.exe"
  fi

  echo "Building ${GOOS}/${GOARCH}..."

  CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" go build \
    -ldflags="-s -w -X main.Version=${VERSION}" \
    -o "${OUTPUT_DIR}/${output_name}" \
    "${MODULE_PATH}"
done

echo ""
echo "Build complete. Binaries:"
ls -lh "${OUTPUT_DIR}/"

# Create individual checksum files
echo ""
echo "Creating checksums..."
cd "${OUTPUT_DIR}"
for binary in ctx-*; do
  # Skip if it's already a checksum file
  [[ "${binary}" == *.sha256 ]] && continue
  if command -v sha256sum &> /dev/null; then
    sha256sum "${binary}" > "${binary}.sha256"
  elif command -v shasum &> /dev/null; then
    shasum -a 256 "${binary}" > "${binary}.sha256"
  fi
done
cd ..

echo ""
echo "Done! Binaries are in ${OUTPUT_DIR}/"

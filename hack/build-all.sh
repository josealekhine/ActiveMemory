#!/bin/bash
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

# Clean and create output directory
rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}"

# Build for each target
for target in "${TARGETS[@]}"; do
    GOOS="${target%/*}"
    GOARCH="${target#*/}"

    output_name="${BINARY_NAME}-${GOOS}-${GOARCH}"
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

# Create checksums
echo ""
echo "Creating checksums..."
cd "${OUTPUT_DIR}"
if command -v sha256sum &> /dev/null; then
    sha256sum ctx-* > checksums.txt
elif command -v shasum &> /dev/null; then
    shasum -a 256 ctx-* > checksums.txt
fi
cd ..

echo ""
echo "Done! Binaries are in ${OUTPUT_DIR}/"

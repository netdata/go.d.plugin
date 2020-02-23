#!/usr/bin/env bash

set -e

PLATFORMS=(
  darwin/386
  darwin/amd64
  freebsd/386
  freebsd/amd64
  freebsd/arm
  linux/386
  linux/amd64
  linux/arm
  linux/arm64
  linux/ppc64
  linux/ppc64le
  linux/mips
  linux/mipsle
  linux/mips64
  linux/mips64le
)

getos() {
  local IFS=/ && read -ra array <<< "$1" && echo "${array[0]}"
}

getarch() {
  local IFS=/ && read -ra array <<< "$1" && echo "${array[1]}"
}

WHICH="$1"

VERSION="${TRAVIS_TAG:-$(git describe --tags --always --dirty)}"

GOLDFLAGS=${GLDFLAGS:-}
GOLDFLAGS="$GOLDFLAGS -w -s -X main.version=$VERSION"

build_all_platforms() {
  for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=$(getos "$PLATFORM")
    GOARCH=$(getarch "$PLATFORM")
    FILE="bin/go.d.plugin-${VERSION}.${GOOS}-${GOARCH}"

    echo "Building for os/arch ${GOOS}/${GOARCH}"
    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "${GOLDFLAGS}" -o "${FILE}" github.com/netdata/go.d.plugin/cmd/godplugin

    ARCHIVE="${FILE}.tar.gz"
    tar -C bin -cvzf "${ARCHIVE}" "${FILE/bin\//}"
    rm "${FILE}"
  done
}

build_specific_platform() {
  GOOS=$(getos "$1")
  GOARCH=$(getarch "$1")
  : "${GOARCH:=amd64}"

  echo "Building for os/arch ${GOOS}/${GOARCH}"
  CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/godplugin -ldflags "${GOLDFLAGS}" github.com/netdata/go.d.plugin/cmd/godplugin
}

build_current_platform() {
  eval "$(go env | grep -e "GOHOSTOS" -e "GOHOSTARCH")"
  GOOS=${GOOS:-$GOHOSTOS}
  GOARCH=${GOARCH:-$GOHOSTARCH}

  echo "Building for os/arch ${GOOS}/${GOARCH}"
  CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/godplugin -ldflags "${GOLDFLAGS}" github.com/netdata/go.d.plugin/cmd/godplugin
}

echo "Building binaries for version: $VERSION"
if [[ "$WHICH" == "all" ]]; then
  build_all_platforms
elif [[ -n "$WHICH" ]]; then
  build_specific_platform "$WHICH"
else
  build_current_platform
fi

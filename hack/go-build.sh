#!/usr/bin/env bash

set -e

WHICH="$1"
PLATFORMS=(darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm linux/arm64 linux/ppc64 linux/ppc64le linux/mips linux/mipsle linux/mips64 linux/mips64le)

VERSION="${TRAVIS_TAG}"
: "${VERSION:=$(git describe --tags --always --dirty)}"
GOFLAGS=${GOFLAGS:-}
GLDFLAGS=${GLDFLAGS:-}
GLDFLAGS="$GLDFLAGS -X main.version=$VERSION"

echo "Building binaries for version: $VERSION"
if [[ "${WHICH}" != "all" ]]; then
	eval $(go env | grep -e "GOHOSTOS" -e "GOHOSTARCH")
	: "${GOOS:=${GOHOSTOS}}"
	: "${GOARCH:=${GOHOSTARCH}}"
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build ${GOFLAGS} -ldflags "${GLDFLAGS}" -o bin/godplugin github.com/netdata/go.d.plugin/cmd/godplugin
else
	for PLATFORM in "${PLATFORMS[@]}"; do
		PLTFRM_SPLT=(${PLATFORM//\// })
		GOOS=${PLTFRM_SPLT[0]}
		GOARCH=${PLTFRM_SPLT[1]}
		FILE="bin/go.d.plugin-${VERSION}.${GOOS}-${GOARCH}"
		CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build ${GOFLAGS} -ldflags "${GLDFLAGS}" -o "${FILE}" github.com/netdata/go.d.plugin/cmd/godplugin
		ARCHIVE="${FILE}.tar.gz"
		tar -C bin -cvzf ${ARCHIVE} ${FILE/bin\//}
		rm ${FILE}
	done
fi

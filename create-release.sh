#!/usr/bin/env bash

set -x -o

VERSION=${VERSION:-"0.0.0"}
DIST=("linux" "darwin" "windows")
DIR_VERSION=kubeconnect-${VERSION}
RELEASE_DIR=$PWD/release

mkdir -p build "${RELEASE_DIR}"

for dist in ${DIST[*]}; do
  BUILD_DIR=build/$dist
  DIR=${BUILD_DIR}/${DIR_VERSION}
  rm -rf ${DIR}
  mkdir -p ${DIR}

  GOOS=$dist GOARCH=amd64 go build \
    -tags release \
    -ldflags "-X kubeconnect/cmd.Version=${VERSION}" \
    -o ${DIR}/kubeconnect;

  if [[ "$dist" == "windows" ]]; then
    cd "$BUILD_DIR"; zip -r "${RELEASE_DIR}/kubeconnect-$dist-${VERSION}.zip" "${DIR_VERSION}"; cd -
  else
    tar -czf "${RELEASE_DIR}/kubeconnect-$dist-${VERSION}.tar.gz" -C "$BUILD_DIR" .
  fi
done

rm -rf build

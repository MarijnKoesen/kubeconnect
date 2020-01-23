#!/usr/bin/env bash

set -x -o

VERSION=$(grep Version cmd/root.go | awk -F'"' '{ print $2 }')
DIR=kubeconnect-${VERSION}

rm -rf ${DIR}
mkdir ${DIR}

rm ${DIR}/kubeconnect
GOOS=linux go build -o ${DIR}/kubeconnect
tar -czf kubeconnect-linux-${VERSION}.tar.gz ${DIR}/kubeconnect

rm ${DIR}/kubeconnect
GOOS=darwin go build -o ${DIR}/kubeconnect
tar -czf kubeconnect-darwin-${VERSION}.tar.gz ${DIR}/kubeconnect

rm ${DIR}/kubeconnect
GOOS=windows go build -o ${DIR}/kubeconnect
zip -r kubeconnect-windows-${VERSION}.zip ${DIR}/kubeconnect

rm -rf ${DIR}

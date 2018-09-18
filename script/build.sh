#!/bin/bash
set -e

targets="linux_386 linux_amd64 darwin_amd64 windows_386 windows_amd64"

mkdir -p build
for target in $targets; do
    echo Building ${target}...
    mkdir -p tmp/${target}
    bin=kouch
    if [ "${target%_*}" == "windows" ]; then
        bin=kouch.exe
    fi
    GOOS=${target%_*} GOARCH=${target##*_} go build -o tmp/${target}/${bin} ./cmd/kouch
    cp LICENSE.md README.md tmp/${target}

    tar -czvpf build/kouch-${TRAVIS_TAG/-/_}-${target}.tar.gz tmp/${target}
    zip -9r build/kouch-${TRAVIS_TAG/-/_}-${target}.zip tmp/${target}
done

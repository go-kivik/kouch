#!/bin/bash
set -e

dir=$(go version | cut -d' ' -f4 | sed -e 's|/|-|')
mkdir -p tmp/${dir}
go build -o tmp/${dir}/kouch ./cmd/kouch
cp LICENSE.md README.md tmp/${dir}

mkdir -p build
tar -czvpf build/kouch-${TRAVIS_TAG}-${dir}.tar.gz tmp/${dir}
zip -9r build/kouch-${TRAVIS_TAG}-${dir}.zip tmp/${dir}

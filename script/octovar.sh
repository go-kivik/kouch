#!/bin/bash
set -e

case "$1" in
    prerelease)
        if [[ "${TRAVIS_TAG}" == *"-"* ]]; then
            echo true
        else
            echo false
        fi
    ;;
esac

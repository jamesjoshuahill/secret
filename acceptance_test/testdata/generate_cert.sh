#!/usr/bin/env bash

dir="$(cd "$(dirname "${BASH_SOURCE[0]}" )" && pwd)"

pushd $dir
    go run "$(go env GOROOT)/src/crypto/tls/generate_cert.go" \
        --host localhost,::1 \
        --ca \
        --rsa-bits 1024 \
        --start-date "Jan 1 00:00:00 1970" \
        --duration=1000000h
popd

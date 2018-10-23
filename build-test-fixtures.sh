#!/bin/bash

if [ ! -f ./go-mockgen ]; then
    function finish {
        echo "Removing binary..."
        rm ./go-mockgen
    }

    echo "Binary not found, building..."
    go build
    trap finish EXIT
fi

echo "Clearing old mocks..."
rm -f ./internal/e2e-tests/mock/*.go

echo "Generating mocks..."
./go-mockgen github.com/efritz/go-mockgen/internal/e2e-tests/iface -d ./internal/e2e-tests/mock -f

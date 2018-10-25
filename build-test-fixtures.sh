#!/bin/bash

binname="go-mockgen"
srcpath="github.com/efritz/go-mockgen/internal/e2e-tests"
genpath="./internal/e2e-tests/mocks"

if [ ! -f "./${binname}" ]; then
    function finish {
        echo "Removing binary..."
        rm "./${binname}"
    }

    echo "Binary not found, building..."
    go build
    trap finish EXIT
fi

echo "Clearing old mocks..."
rm -f "${genpath}/*.go"

echo "Generating mocks..."
"./${binname}" -d "${genpath}" -f "${srcpath}"

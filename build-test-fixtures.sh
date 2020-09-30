#!/bin/bash

binname="go-mockgen"
srcpath="./internal/testdata"
genpath="./internal/testdata/mocks"

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

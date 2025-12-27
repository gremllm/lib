#!/usr/bin/env bash

set -e

mkdir -p build
go build -buildmode=c-shared -o build/libschema.so ./cmd/libschema/

cd ffi_tests/nodejs/
if [ ! -d "node_modules" ]; then
    npm install
fi
cd -

echo "Running Python FFI test..."
python3 ffi_tests/python/main.py

echo "Running NodeJS FFI test..."
node ffi_tests/nodejs/main.js
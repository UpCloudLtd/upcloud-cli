#!/usr/bin/env bash

BIN_DIR=bin
pushd ${BIN_DIR}
sha256sum * > SHA256SUMS
popd
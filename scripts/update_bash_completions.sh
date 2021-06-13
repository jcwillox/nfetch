#!/usr/bin/env bash

# ensure we are in the repo root
DIR=$(dirname "$0")
cd "$DIR/.." || exit

go run . completion bash > ./completions/nfetch.sh
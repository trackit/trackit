#!/bin/bash

set -e

build_number=$(git rev-parse --short=8 HEAD)
if [[ $(git status --porcelain 2>/dev/null | wc -l) -gt 0 ]]
then
	build_number="$build_number-devel"
fi

go build \
	-ldflags "-X main.buildNumber=$build_number" \
	-o main

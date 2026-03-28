#!/usr/bin/env bash

cat ./internal/cli/cli.go | grep "case" | sed "s/case //; s/://; s/\t\t//; s/^/option: /"
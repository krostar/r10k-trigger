#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o posix

go install github.com/vektra/mockery/cmd/mockery

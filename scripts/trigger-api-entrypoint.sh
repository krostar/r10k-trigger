#!/usr/bin/env sh

set -o errexit
set -o nounset
set -o pipefail

exec trigger-api "$@"

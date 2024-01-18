#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

make install

parallel --lb -j2 --halt-on-error 2 ::: './start.sh' './integration-test/test.sh'
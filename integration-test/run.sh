#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
set -x

make install

parallel --lb -j2 --halt now,done=1 ::: './start.sh >> node_logs.txt 2>&1' './integration-test/test.sh'
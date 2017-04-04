#!/bin/sh

WHOAMI=`python -c 'import os, sys; print os.path.realpath(sys.argv[1])' $0`
ROOT=`dirname $WHOAMI`

${ROOT}/setup-ubuntu.sh
${ROOT}/setup-libvips.sh
${ROOT}/setup-go.sh

#!/bin/sh

# This is copied in to the built Dockerfile by "docker build -t go-iiif ."
# It exists principally so we can configure ${IIIF_SERVER_CONFIG} dynamically
# (20180620/thisisaaronland)

/bin/iiif-server -host 0.0.0.0 -config ${IIIF_SERVER_CONFIG}
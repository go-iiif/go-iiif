#!/bin/sh

# This is copied in to the built Dockerfile by "docker build -t go-iiif ."
# It exists principally so we can configure ${IIIF_SERVER_CONFIG} dynamically
# (20180620/thisisaaronland)

CONFIG="/etc/go-iiif/config.json"
FORMAT="jpg"
QUALITY="default"
REGION="full"
ROTATION="0"
SIZE="full"

if [ "${IIIF_CONFIG}" != "" ]
then
    CONFIG=$IIIF_CONFIG
fi

if [ "${IIIF_FORMAT}" != "" ]
then
    FORMAT=$IIIF_FORMAT
fi

if [ "${IIIF_QUALITY}" != "" ]
then
    FORMAT=$IIIF_QUALITY
fi

if [ "${IIIF_REGION}" != "" ]
then
    FORMAT=$IIIF_REGION
fi

if [ "${IIIF_ROTATION}" != "" ]
then
    FORMAT=$IIIF_ROTATION
fi

if [ "${IIIF_SIZE}" != "" ]
then
    FORMAT=$IIIF_SIZE
fi

/bin/iiif-transform -host 0.0.0.0 -config ${CONFIG} -format ${FORMAT} -quality ${QUALITY} -region ${REGION} -rotation ${ROTATION} -size ${SIZE} ${INPUT} ${OUTPUT}

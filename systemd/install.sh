#!/bin/sh

export PATH="${PATH}:/usr/local/bin"

PYTHON=`which python`
GOLANG=`which go`

if [ ! -x ${PYTHON} ]
then
    echo "Missing python binary"
    exit 1
fi

if [ ! -x ${GOLANG} ]
then
    echo "Missing go binary"
    exit 1
fi


WHOAMI=`${PYTHON} -c 'import os, sys; print os.path.realpath(sys.argv[1])' $0`

SYSTEMD=`dirname ${WHOAMI}`
GO_IIIF=`dirname ${SYSTEMD}`

USER="iiif-server"
GROUP="iiif-server"

CONFIG="/etc/iiif-server/config.json"
CONFIG_ROOT=`dirname ${CONFIG}`

SERVICE="/lib/systemd/system/iiif-server.service"

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit 1
fi

if getent passwd ${USER} > /dev/null 2>&1; then
    echo "${USER} user account already exists"
else
    useradd ${USER} -s /sbin/nologin -M
fi

cd ${GO_IIIF}
export GOPATH="${GO_IIIF}"
${GOLANG} build -o /usr/local/bin/iiif-server cmd/iiif-server.go
cd -


if [ ! -d ${CONFIG_ROOT} ]
then
    mkdir -p ${CONFIG_ROOT}
fi

if [ -f ${CONFIG} ]
then 
    # MTIME=`stat -c %Y /etc/iiif-server/config.json`
    # mv ${CONFIG} ${CONFIG}.{$MTIME}

    echo "${CONFIG} already exists, so leaving it in place"
else

    cp ${GO_IIIF}/config.json.example ${CONFIG}
    chmod 0644 /etc/iiif-server/config.json
    chgrp ${GROUP} /etc/iiif-server/config.json
fi

if [ -f ${SERVICE} ]
then
    echo "${SERVICE} already exists, so leaving it in place"
else
    cp ${SYSTEMD}/iiif-server.service.example ${SERVICE}
    sudo chmod 644 ${SERVICE}
fi

echo ""
echo "system stuff installed - you will still need to run the following, manually:"
echo "	systemctl enable iiif-server.service"
echo "	systemctl start iiif-server.service"
echo "also remember that any path-based sources defined in your config (images.source, derivatives.cache, etc.) will need to be writable by the ${USER} user"

exit 0

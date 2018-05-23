#!/bin/sh

# install go

VERSION="1.10.2"
DIST="go${VERSION}.linux-amd64.tar.gz"
SOURCE="https://storage.googleapis.com/golang/${DIST}"

HASH="4b677d698c65370afa33757b6954ade60347aaca310ea92a63ed717d7cb0c2ff"

if [ ! -d /usr/local/go${VERSION} ]
then
    cd /tmp

    wget ${SOURCE}

    FNAME=`basename ${SOURCE}`
    SRC_HASH=`shasum -a 256 /tmp/${FNAME} | awk '{ print $1 }'`

    if [ "${SRC_HASH}" != "${HASH}" ]
    then
	echo "Weird hash (${SRC_HASH}), expected ${HASH}"
	exit 1
    fi

    tar -xvzf ${DIST}

    if [ -f /tmp/${DIST} ]
    then
	rm /tmp/${DIST}
    fi

    mv /tmp/go /usr/local/go${VERSION}

    if [ -L /usr/local/go ]
    then
	rm /usr/local/go
    fi

    ln -s /usr/local/go${VERSION} /usr/local/go

    for BIN in go gofmt godoc
    do
	
	if [ -L /usr/local/bin/${BIN} ]
	then
	    rm /usr/local/bin/${BIN}
	fi
	
	ln -s /usr/local/go/bin/${BIN} /usr/local/bin/${BIN}
    done

    cd -
fi

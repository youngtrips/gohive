#!/bin/sh

GOPATH=`pwd`
DST_PATH=.bin
PROJ_NAME=gohive

rm -rf .bin
mkdir -p $DST_PATH

cp $GOPATH/bin/ucd $DST_PATH/
cp $GOPATH/bin/gamed $DST_PATH
cp $GOPATH/bin/gated $DST_PATH

docker build --no-cache --rm=true -t $PROJ_NAME .

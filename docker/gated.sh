#!/bin/sh

NAME=gated1
PWD=`pwd`
IMAGES=lab.qeetap.com:5000/gohive:latest
CMD=./bin/gated

docker stop $NAME
docker rm -f $NAME

docker run -d -v $PWD/conf:/gohive/conf -v $PWD/logs:/gohive/logs --network="host" --restart=always --name $NAME $IMAGES $CMD

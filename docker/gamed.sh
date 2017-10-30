#!/bin/sh

NAME=gamed1
PWD=`pwd`
IMAGES=lab.qeetap.com:5000/gohive:latest
CMD=./bin/proxy

docker stop $NAME
docker rm -f $NAME

docker run -d -v $PWD/conf:/gohive/conf -v $PWD/logs:/gohive/logs -v $PWD/data:/gohive/data --network="host" --restart=always --name $NAME $IMAGES $CMD

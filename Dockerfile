FROM alpine
MAINTAINER Tuz <youngtrips@gmail.com>

RUN apk add tzdata --update && cp /usr/share/zoneinfo/Asia/Chongqing /etc/localtime 
RUN apk add ca-certificates
RUN rm -rf /var/cache/apk/*

WORKDIR /gohive
VOLUME /gohive/conf
VOLUME /gohive/data
VOLUME /gohive/logs

COPY .bin /gohive/bin

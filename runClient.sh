#!/bin/bash

if [ "$#" -ne 1 ] ; then
    echo "Usage: ./runClient.sh <container-name-suffix>"
    exit
fi

export CHAT_USERNAME=$1

docker stop chat-$CHAT_USERNAME && \
docker rm chat-$CHAT_USERNAME

docker-compose up -d chat-client && \
docker attach chat-$CHAT_USERNAME

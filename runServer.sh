#!/bin/bash

docker stop chat-server && \
docker rm chat-server

docker-compose up -d chat-server

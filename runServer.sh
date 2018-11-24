#!/bin/bash

docker stop chat-server && \
docker rm chat-server

docker-compose up chat-server

#!/usr/bin/bash
swag init --pd &&
go build -v -ldflags='-s -w' && 
rsync -azP $(pwd)/env.example $1:/opt/workerengine/env &&
rsync -azP $(pwd)/workerengine $1:/opt/workerengine/workerengine
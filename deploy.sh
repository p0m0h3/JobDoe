#!/usr/bin/bash
swag init --pd &&
go build -v -ldflags='-s -w' && 
rsync -azP $(pwd)/env $1:/opt/workerengine/env &&
rsync -azP --delete $(pwd)/tools/*.toml $1:/opt/workerengine/tools/ &&
rsync -azP $(pwd)/workerengine $1:/opt/workerengine/workerengine
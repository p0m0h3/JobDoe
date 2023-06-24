#!/usr/bin/bash
go build && 
swag init --pd &&
rsync -azP $(pwd)/env $1:/opt/workerengine/env &&
rsync -azP --delete $(pwd)/tools/*.toml $1:/opt/workerengine/tools/ &&
rsync -azP $(pwd)/workerengine $1:/opt/workerengine/workerengine
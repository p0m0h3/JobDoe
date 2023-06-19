#!/usr/bin/bash
go build && 
swag init &&
rsync -azP $(pwd)/env $1:~/env &&
rsync -azP --delete $(pwd)/tools/*.toml $1:~/tools/ &&
rsync -azP $(pwd)/workerengine $1:~/workerengine
#!/usr/bin/bash
go build
rsync -azP $(pwd)/env $1:~/env
rsync -azP $(pwd)/workerengine $1:~/workerengine
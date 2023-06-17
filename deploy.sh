#!/usr/bin/bash
go build
rsync -azP $(pwd)/workerengine $1:~/workerengine
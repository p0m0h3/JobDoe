#!/usr/bin/bash
go build
rsync -azP $(pwd)/env $1:~/env
rsync -azP $(pwd)/tools/*.toml $1:~/tools/
rsync -azP $(pwd)/workerengine $1:~/workerengine
PACKAGE_FILES=workerengine workerengine.service install.sh env.example regauth.json.example

swag:
	swag init --quiet --parseDependency

build: swag
	go build

run:
	./workerengine

pack: swag
	go build -v -ldflags='-s -w'
	tar -czf workerengine.tar.gz ${PACKAGE_FILES}
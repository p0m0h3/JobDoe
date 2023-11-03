
DEBIAN_FILES=control postinst
OPT_FILES=workerengine env.example regauth.json.example

deb: clean
	go build -v -ldflags='-s -w'
	mkdir -p ./workerengine_amd64/DEBIAN
	cp ${DEBIAN_FILES} ./workerengine_amd64/DEBIAN/
	mkdir -p ./workerengine_amd64/usr/lib/systemd/system
	cp workerengine.service ./workerengine_amd64/usr/lib/systemd/system/
	mkdir -p ./workerengine_amd64/opt/workerengine
	cp -r ${OPT_FILES} ./workerengine_amd64/opt/workerengine/
	dpkg-deb --build --root-owner-group workerengine_amd64

build: 
	go build

clean:
	rm -r -f workerengine_amd64 workerengine_amd64.deb

run:
	./workerengine

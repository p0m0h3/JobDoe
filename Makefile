VERSION=v0.2.0
DEBIAN_FILES=deb/control deb/postinst
OPT_FILES=workerengine config.json.example regauth.json.example
PACKAGE_DIR=workerengine_${VERSION}


build: clean
	go build -v -ldflags='-s -w'
	mkdir -p ./${PACKAGE_DIR}/DEBIAN
	cp ${DEBIAN_FILES} ./${PACKAGE_DIR}/DEBIAN/
	mkdir -p ./${PACKAGE_DIR}/usr/lib/systemd/system
	cp deb/workerengine.service ./${PACKAGE_DIR}/usr/lib/systemd/system/
	mkdir -p ./${PACKAGE_DIR}/opt/workerengine
	cp -r ${OPT_FILES} ./${PACKAGE_DIR}/opt/workerengine/
	dpkg-deb --build --root-owner-group ${PACKAGE_DIR}

install: build
	dpkg -i ${PACKAGE_DIR}.deb

clean:
	rm -r -f ${PACKAGE_DIR} ${PACKAGE_DIR}.deb

#!/usr/bin/sh

# Make sure only root can run install
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root." 1>&2
   exit 1
fi

mkdir -p /opt/workerengine
cp ./workerengine ./env.example /opt/workerengine
cp ./workerengine.service /etc/systemd/system
systemctl daemon-reload
chmod -R 700 /opt/workerengine
chmod +x /opt/workerengine/workerengine
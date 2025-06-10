#!/bin/bash

chmod 666 /var/run/docker.sock

exec "$@"
# exec su docker -c ./start.sh
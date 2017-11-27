#!/bin/bash

set -ex

if [[ ! -f "server/main" ]]
then
	pushd "server"
	./buildstatic.sh
	popd
fi

pushd docker
docker-compose build
../scripts/awsenv default docker-compose up
popd

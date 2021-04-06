#!/bin/bash
export IP_DB=`docker inspect docker_sql_1 | grep "IPAddress" | tail -n 1 | cut -d \" -f 4`
docker exec -it docker_api_1 /main --aws-region us-west-2 --http-address "127.0.0.1:8080" --sql-address "trackit:trackitpassword@tcp($IP_DB)/trackit?parseTime=true" --es-address "es:9200" --redis-address "redis:6379" --pretty-json-responses --task $argv

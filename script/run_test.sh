#!/bin/sh
DB_CONN_STR="host=db user=postgres password=example dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Taipei"
REDIS_HOST="redis:6379"
REDIS_DB=0
docker run --workdir="/go/app" -v "$(pwd):/go/app" --network=deployment_example -e REDIS_DB=${REDIS_DB} -e REDIS_HOST=${REDIS_HOST} -e DB_CONN_STR="${DB_CONN_STR}" -it golang:1.15 go test ./...
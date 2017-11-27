#!/bin/bash

sql_user="${1:-root}"
sql_password="${2:-rootpassword}"
sql_host="${3:-localhost}"
sql_database="${4:-trackit}"

xo \
	"mysql://${sql_user}:${sql_password}@${sql_host}/${sql_database}" \
	--out models \
	--ignore-fields \
		created \
		modified

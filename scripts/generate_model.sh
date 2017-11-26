#!/bin/bash

sql_user="${1:-trackit}"
sql_password="${2:-trackitpassword}"
sql_host="${3:-localhost}"
sql_database="${4:-trackit}"

xo \
	"mysql://${sql_user}:${sql_password}@${sql_host}/${sql_database}" \
	--verbose \
	--out models \
	--ignore-fields \
		created \
		modified

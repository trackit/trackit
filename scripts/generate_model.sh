#!/bin/bash

xo \
	'mysql://root:rootpassword@localhost/db' \
	--out models \
	--ignore-fields \
		created \
		modified

#!/bin/bash

xo \
	'mysql://root:ooquuphieMohnei8pei3mias7pee8Yae@localhost/trackit' \
	--out models \
	--ignore-fields \
		created \
		modified

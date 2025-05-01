#!/bin/sh

: "${DIR:=/}"

stat_it() {
	echo "hyphen $1 $(tree "$DIR" | grep - | grep -cE "\\.$1$")"
	echo "unscor $1 $(tree "$DIR" | grep _ | grep -cE "\\.$1$")"
	echo "neithe $1 $(tree "$DIR" | grep -vE "_|-" | grep -cE "\\.$1$")"
}

stat_it sh
stat_it bash
stat_it js
stat_it py
stat_it go
stat_it cmake
stat_it cpp
stat_it cc
stat_it yml
stat_it json
stat_it c

#!/usr/bin/env bash
set -e

declare -A BUFF_M
declare -A MAIN_M

rows=5
cols=10
iterations=20

# Init
for ((y = 1; y <= rows; y++)); do
	for ((x = 1; x <= cols; x++)); do
		MAIN_M["$x,$y"]=0
	done
done

# Glider
# NOTE: https://github.com/mvdan/sh/issues/226
MAIN_M["4,2"]=1
MAIN_M["2,3"]=1
MAIN_M["4,3"]=1
MAIN_M["3,4"]=1
MAIN_M["4,4"]=1

# Print
function print() {
	for ((y = 1; y <= rows; y++)); do
		for ((x = 1; x <= cols; x++)); do
			get="$x,$y"
			printf "%2s" "${MAIN_M[$get]}"
		done
		echo
	done
	echo
}

# Copy Matrix
function cpym() {
	for ((y = 1; y <= rows; y++)); do
		for ((x = 1; x <= cols; x++)); do
			get="$x,$y"
			BUFF_M[$get]=${MAIN_M[$get]}
		done
	done
}

# Cell Check
function cell() {
	x=$1
	y=$2
	get="$x,$y"
	R=${BUFF_M[$get]}
	C=-$R
	for ((i = -1; i <= 1; i++)); do
		for ((j = -1; j <= 1; j++)); do
			xj=$((x + j - 1 + cols))
			cx=$((xj % cols + 1))
			yi=$((y + i - 1 + rows))
			cy=$((yi % rows + 1))
			get="$cx,$cy"
			C=$((C + ${BUFF_M[$get]}))
		done
	done
	if ((C < 2)) || ((C > 3)); then
		R=0
	fi
	if ((C == 3)); then
		R=1
	fi
}

# Mainloop
print
for ((g = 0; g <= iterations; g++)); do
	cpym
	for ((y = 1; y <= rows; y++)); do
		for ((x = 1; x <= cols; x++)); do
			cell "$x" "$y"
			MAIN_M["$x,$y"]=$R
		done
	done
	# print
done
sleep 1
print

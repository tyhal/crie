#!/bin/sh
set -e

# TODO Put into build process of standards
PATCHOFFSET=13

STD_MAJOR="$(script/crie version --major | tr -d "[:space:]")"
STD_MINOR="$(script/crie version --minor | tr -d "[:space:]")"
STD_PATCH="$(script/crie version --patch | tr -d "[:space:]")"

STD_FULL="$STD_MAJOR.$STD_MINOR.$STD_PATCH"

TEST_FULL="$(script/crie version | tr -d "[:space:]")"
if [ "$STD_FULL" != "$TEST_FULL" ]; then
	echo "FAIL! THE WHOLE DOES NOT MATCH THE SUM OF ITS PARTS, WHOLE: $STD_FULL PARTS: $TEST_FULL"
	exit 1
else
	echo "SUB VERSIONS ARE CONSISTENT WITH FULL VERSION"
fi

# Calculated from Git Commits
RAW_PATCH="$(git log --no-merges --pretty=format:'' | wc -l)"
TEST_PATCH="$((RAW_PATCH - PATCHOFFSET))"
if [ "$TEST_PATCH" != "$STD_PATCH" ]; then
	echo "FAIL! INCORRECT PATCH, TOOL RETURNS $STD_PATCH BUT IT SHOULD RETURN $TEST_PATCH"
	exit 1
else
	echo "PATCH VERSION MATCH!"
fi

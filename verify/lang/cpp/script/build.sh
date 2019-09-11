#!/usr/bin/env bash

rm -rf build
mkdir build
cd build || exit 1
cmake ..
cmake --build .

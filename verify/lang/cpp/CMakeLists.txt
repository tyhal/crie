cmake_minimum_required(VERSION 3.0.2)

# really long line to test line length and make sure there is no regression or something like that because that would be inconventient for everyone

set(CMAKE_CXX_FLAGS "-std=c++11 ${CMAKE_CXX_FLAGS}")

set(SOURCE_FILES src/example.cc)
add_executable(example ${SOURCE_FILES})

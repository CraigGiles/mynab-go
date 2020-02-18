#!/bin/bash

MAIN_FILE=main.go
EXE_NAME=application

CTIME_EXEC="ctime/ctime"
CTIME_TIMING_FILE="ctime_info.ctm"

# Build ctime
if [ ! -f "$CTIME_EXEC" ]; then
  ./"ctime/build.sh"
fi
	
function clean_build_directory {
    echo "cleaning up old build folder"
    rm -rf build
    mkdir -p build
}

function build_app {
    echo "building application"
    go build -o build/$EXE_NAME ./src/$MAIN_FILE
}

function start_compile_timing {
    $CTIME_EXEC -begin "$CTIME_TIMING_FILE"
}

function stop_compile_timing {
    $CTIME_EXEC -end "$CTIME_TIMING_FILE" $LAST_ERROR
}

function count_lines_of_code {
    echo ""
    ./cloc.sh
}

function run_app {
    ./build/$EXE_NAME
}

if [ "$1" = "run" ]
then
    run_app
else
    clean_build_directory

    start_compile_timing
    build_app
    stop_compile_timing

    count_lines_of_code
fi


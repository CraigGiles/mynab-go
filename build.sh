#!/bin/bash

EXE_NAME=application

BUILD_DIR="build"
PROJECT_NAME="main"
MODULE_NAME="mynab"

CTIME_EXEC="ctime/ctime"
CTIME_TIMING_FILE="ctime_info.ctm"

# Build ctime
if [ ! -f "$CTIME_EXEC" ]; then
  ./"ctime/build.sh"
fi
	
function clean_build_directory {
    rm -rf /$BUILD_DIR
    mkdir -p ./$BUILD_DIR
}

function build_app {
    go build -o ./$BUILD_DIR/$MODULE_NAME "./$MODULE_NAME" 
}

function test_app {
    go test "./$MODULE_NAME" 
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
    ./$BUILD_DIR/$MODULE_NAME
}

function launch_development_environment {
    docker-compose up
}

function migrate_database {
    go run sql/migrate.go
}

if [ "$1" = "run" ]
then
    run_app
elif [ "$1" = "test" ]
then
    test_app
elif [ "$1" = "devenv" ]
then
    launch_development_environment
elif [ "$1" = "migrate" ]
then
    migrate_database
else
    clean_build_directory

    start_compile_timing
    build_app
    stop_compile_timing

    count_lines_of_code
fi


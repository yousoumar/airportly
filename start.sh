#!/bin/bash

set -e

cleanup() {
    echo "Process stopping..."
    for package in "${packages[@]}"; do
        pkill -TERM -f "$build_dir/$package"
    done
    wait
}
trap 'cleanup' EXIT

source_dir=./cmd
build_dir=./bin

if [ ! -d "$build_dir" ]; then
    mkdir "$build_dir"
fi

if [ "$(ls -A "$build_dir")" ]; then
    rm -r "${build_dir:?}/"*
fi

packages=("sensor" "file-recorder" "alert-manager" "database-recorder" "http-rest-server")

echo "Building all packages..."

for package in "${packages[@]}"; do
    echo "Building $package"
    go build -o "$build_dir/$package" "$source_dir/$package"
done

echo "Building the IHM..."

cd ./ihm && npm install && npm run build && cd ..

echo "Build completed successfully. Starting all services..."

for package in "${packages[@]}"; do
    echo "Running $package"
    if [ "$package" == "sensor" ]; then
        ."/$build_dir/$package" CDG pressure 1007 &
        ."/$build_dir/$package" RAK pressure 1006 &
        ."/$build_dir/$package" CDG temperature 35 &
        ."/$build_dir/$package" RAK temperature 16 &
        ."/$build_dir/$package" CDG wind-speed 35 &
        ."/$build_dir/$package" RAK wind-speed 34 &
        continue
    fi
    ."/$build_dir/$package" &
done

cd ./ihm && npm run preview

echo "All services started successfully. Press Ctrl+C to stop all services."
wait

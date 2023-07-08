#!/bin/bash

GOOS=linux
GOARCH=arm64
CGO_ENABLED=0

echo "Compiling for $GOOS/$GOARCH..."

output_name=your_output_name_$GOOS\_$GOARCH

go build -o $output_name

if [ $? -eq 0 ]; then
    echo "Compilation successful: $output_name"
else
    echo "Compilation failed for $GOOS/$GOARCH"
    exit 1
fi

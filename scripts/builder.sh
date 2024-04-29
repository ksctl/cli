#!/bin/bash

cd .. || echo -e "\033[1;31mUnable to cd into ksctl root\033[0m\n"

# Get the binary from the source code
cd cli || echo -e "\033[1;31mPath couldn't be found\033[0m\n"
# Check if sudo access
go get -d
go build -v -o ksctl .
chmod +x ksctl

sudo mv -v ksctl /usr/local/bin/ksctl

echo -e "\033[1;32mINSTALL COMPLETE\033[0m\n"

cd - || echo -e "\033[1;31mFailed to move to previous directory\033[0m\n"

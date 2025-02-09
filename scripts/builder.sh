#!/bin/bash

cd .. || echo -e "\033[1;31mUnable to cd into ksctl root\033[0m\n"

go get -d
go build -v -o ksctl .
chmod +x ksctl

sudo mv -v ksctl /usr/local/bin/ksctl

echo -e "\033[1;32mINSTALL COMPLETE\033[0m\n"

#!/bin/sh
go run bf.go $1 > /dev/null
make out > /dev/null
./out
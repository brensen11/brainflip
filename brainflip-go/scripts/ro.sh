#!/bin/sh
go run bf.go -O $1 > /dev/null
make out > /dev/null
./out
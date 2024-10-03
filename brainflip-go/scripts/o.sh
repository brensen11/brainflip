#!/bin/sh
go run bf.go -O $1
make out
./out < test-input.dat > debug/$1-O.dat
python scripts/transform.py debug/$1-O.dat > debug/$1-O.data
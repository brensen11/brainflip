#!/bin/sh
go run bf.go $1
make out
./out < test-input.dat > debug/$1.dat
python transform.py debug/$1.dat > debug/$1.data
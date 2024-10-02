#!/bin/sh
go run bf.go -O $1
make out
./out > debug/$1-O.dat
python transform.py debug/$1-O.dat > debug/$1-O.data
#!/bin/sh
rm -rf debug/*
go run bf.go ../bfcheck/prog-$1.b
make out
./out < ../bfcheck/input.dat > debug/output.dat
python transform.py debug/output.dat > debug/output.data

cp ../bfcheck/output-$1.dat debug/
python transform.py debug/output-$1.dat > debug/output-$1.data

diff debug/output-$1.dat debug/output.dat
diff debug/output-$1.data debug/output.data
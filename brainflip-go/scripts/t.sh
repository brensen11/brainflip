#!/bin/sh
rm -rf debug/*
go run bf.go ../bftest/prog-$1.b
make out
./out < ../bftest/input.dat > debug/output.dat
python scripts/transform.py debug/output.dat > debug/output.data

cp ../bftest/output-$1.dat debug/
python scripts/transform.py debug/output-$1.dat > debug/output-$1.data

diff debug/output-$1.dat debug/output.dat
diff debug/output-$1.data debug/output.data
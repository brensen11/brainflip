#!/bin/sh
./scripts/r.sh $1 > normal.dat
./scripts/ro.sh $1 > optimized.dat
diff normal.dat optimized.dat
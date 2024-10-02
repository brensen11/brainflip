#!/bin/sh
./c.sh $1
./o.sh $1
diff debug/$1.data debug/$1-O.data
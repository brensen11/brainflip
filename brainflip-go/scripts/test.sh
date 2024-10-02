#!/bin/bash

# Loop from 0 to 50
for N in $(seq -w 0 50); do
    echo testing $N
    ./t.sh 00"$N"
done

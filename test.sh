#!/bin/bash

rm ./incdb
rm ./data/incdb.data
go build -o incdb *.go

./incdb r
./incdb w 1 a
./incdb w 2 b
./incdb r
./incdb w 1 c
./incdb r
./incdb w 3 d
./incdb r

#!/bin/bash

rm -f ./incdb
rm -f ./data/incdb.data
go build -o incdb *.go

./incdb r
./incdb w 1 a
./incdb w 2 b
./incdb w abc def
./incdb r
./incdb w 1 c
./incdb r
./incdb w 3 d
./incdb r
./incdb r limit
./incdb r limit abc
./incdb r limit 1
./incdb r limit 10
./incdb r 4
./incdb r 1
./incdb r abc

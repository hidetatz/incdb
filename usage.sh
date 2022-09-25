#!/bin/bash

set -x

make clean
make
./incdb 'create table item (id string, name string)'
./incdb 'insert into item (id, name) values ("1", "laptop")'
./incdb 'insert into item (id, name) values ("2", "iPhone")'
./incdb 'insert into item values ("3", "radio")'
./incdb 'r item'
./incdb 'r item "2"'
./incdb 'r item order by desc limit 2 offset 1'

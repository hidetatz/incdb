#!/bin/bash

kill_incdbd_if_exists() {
	pid=$(ps aux | grep incdbd | grep -v grep | awk '{print $2}')
	if [ "$pid" != "" ]; then
		kill -9 "$pid"
	fi
}

make clean
make

kill_incdbd_if_exists
./incdbd &
echo "waiting for incdbd gets up and running..." && sleep 1

./incdb 'create table item (id string, name string)'
./incdb 'insert into item (id, name) values ("1", "laptop")'
./incdb 'insert into item (id, name) values ("2", "iPhone")'
./incdb 'insert into item values ("3", "radio")'
./incdb 'select * from item'
./incdb 'select name from item'
./incdb 'select * from item where id = "2"'
./incdb 'select * from item order by desc limit 2 offset 1'

kill_incdbd_if_exists

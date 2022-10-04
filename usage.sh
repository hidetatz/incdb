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

./incdb 'create table person (id string, name string, lang string)'
./incdb 'insert into person (id, name, lang) values ("1", "alice", "En")'
./incdb 'insert into person (id, name, lang) values ("2", "bob", "Ja")'
./incdb 'insert into person values ("3", "chris", "Ja")'
./incdb 'insert into person values ("4", "donald", "En")'
./incdb 'insert into person values ("5", "eddie", "Ch")'
./incdb 'insert into person values ("6", "fred", "Ch")'
./incdb 'select * from person'
./incdb 'select name, lang from person'
./incdb 'select * from person where id = "2"'
./incdb 'select * from person where lang != "Ja"'
./incdb 'select * from person order by name desc limit 5 offset 3'
./incdb 'select * from person order by lang'

kill_incdbd_if_exists

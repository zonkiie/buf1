#!/bin/sh

#buffalo db generate model Book Name:text -s -d
#buffalo db generate model User Username:text Email:text Passwordhash:text Book:[]Book -s -d

#soda generate fizz User -d
#soda create -d database.yml
#buffalo db create -a

sqlite3 db.sqlite "create table user(id text primary key, created_at datetime, updated_at datetime, username text, email text, passwordhash text);"
sqlite3 db.sqlite "create table book(id text primary key, created_at datetime, updated_at datetime, name text, user_id text, FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE RESTRICT ON UPDATE RESTRICT);"

#!/bin/bash

# in another terminal, run the following to start a dev mysql server
# docker run --name mysqldev -e MYSQL_ROOT_PASSWORD=password123 -v ~/mysql_data:/var/lib/mysql -p 3306:3306 -d mysql

# general updates for system
sudo apt update && \
sudo apt upgrade -y &&

# install netcat
sudo apt install netcat -y

# install mariadb-client to connect to mysql
sudo apt install mariadb-client -y

# connect to db cmd
# mysql -h 172.17.0.2 -u root -p
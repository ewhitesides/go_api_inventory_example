#!/bin/bash

# run this outside of the devcontainer
docker run --name mysqldev -e MYSQL_ROOT_PASSWORD=password123 -v ~/mysql_data:/var/lib/mysql -p 3306:3306 -d mysql

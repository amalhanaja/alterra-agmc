docker run --rm \
 --name=mysql-runtime \
 -e MYSQL_ROOT_PASSWORD=password \
 -e MYSQL_DATABASE=development \
 -p 3306:3306 \
 -d mysql:latest
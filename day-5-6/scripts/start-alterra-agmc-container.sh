docker run \
 --name=alterra-agmc \
 -e DB_DSN="root:password@tcp(host.docker.internal:3306)/development?charset=utf8mb4&parseTime=True&loc=Local" \
 -e MONGO_URI="mongodb://root:password@host.docker.internal:27017" \
 -p 8080:8080 \
 alterra-agmc:1.0.0
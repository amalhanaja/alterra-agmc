docker run --rm \
 --name=mongo-runtime \
 -e MONGO_INITDB_ROOT_USERNAME=root \
 -e MONGO_INITDB_ROOT_PASSWORD=password \
 -e MONGO_INITDB_DATABASE=development \
 -p 27017:27017 \
 -d mongo:latest
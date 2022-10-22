# go-rest-api
possible docker build command:
docker build . -t go-rest-api

possible docker run command:
docker run -p 8080:8080 -e MONGO_URI=mongodb://172.17.0.3:27017 -e MONGO_DB_NAME=pub -e PRIMARY_COLLECTION=pub_customers go-rest-api

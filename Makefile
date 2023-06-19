start-mongodb:
	docker run --name mongodb -d -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=pass mongodb/mongodb-community-server:latest

stop-mongodb:
	docker stop mongodb
	docker rm mongodb

start-server:
	npm run dev-server

test-server:
	npm run test-server
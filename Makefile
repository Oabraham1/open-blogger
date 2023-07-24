start-mongodb:
	docker run --name mongodb -d -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=pass mongodb/mongodb-community-server:latest

stop-mongodb:
	docker stop mongodb
	docker rm mongodb

start-server:
	cd server && go run cmd/main.go

test-server:
	cd server && go test -v ./...

start-client:
	cd client && npm run dev -- --open

test-client-unit:
	cd client && npm run test:unit

test-client-int:
	cd client && npm run test:integration
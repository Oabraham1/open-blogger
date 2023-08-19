DB_URL = postgresql://root:secret@localhost:5432/openBloggerDB?sslmode=disable

start-postgres-server:
	@echo "Starting postgres server"
	docker run --name openBloggerPostgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

stop-postgres-server:
	@echo "Stopping postgres server"
	docker stop openBloggerPostgres

start-mongo-server:
	@echo "Starting mongo server"
	docker run --name openBloggerMongo -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=secret -d mongo:7.0.0

stop-mongo-server:
	@echo "Stopping mongo server"
	docker stop openBloggerMongo

createsqldb:
	docker exec -it openBloggerPostgres createdb --username=root --owner=root openBloggerDB

dropsqldb:
	docker exec -it openBloggerPostgres dropdb openBloggerDB

migrateUp:
	cd server && migrate -path ./db/migration -database "$(DB_URL)" -verbose up

migrateDown:
	cd server && migrate -path ./db/migration -database "$(DB_URL)" -verbose down

sqlc-gen:
	cd server && sqlc generate

mock-gen:
	cd server && mockgen -source=db/sqlc/store.go -destination=db/mock/store.go -package=mockdb

start-server:
	cd server && go run cmd/main.go

test-server:
	cd server && go test -v -cover ./...

setup-server-test-env-for-ci:
	touch .env
	echo "DB_URL=$(DB_URL)" >> .env
	echo "DB_DRIVER=postgres" >> .env
	echo "ENVIRONMENT=development" >> .env
	echo "HTTP_SERVER_ADDRESS=0.0.0.0:8080" >> .env

format-go:
	go fmt ./...

start-client:
	cd client && npm run dev -- --open

test-client-unit:
	cd client && npm run test:unit

test-client-int:
	cd client && npm run test:integration
DB_URL = postgresql://root:secret@localhost:5432/openBloggerDB?sslmode=disable

start-postgres-server:
	@echo "Starting postgres server"
	docker run --name openBloggerPostgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it openBloggerPostgres createdb --username=root --owner=root openBloggerDB

dropdb:
	docker exec -it openBloggerPostgres dropdb openBloggerDB

migrateUp:
	cd server && migrate -path ./db/migration -database "$(DB_URL)" -verbose up

migrateDown:
	cd server && migrate -path ./db/migration -database "$(DB_URL)" -verbose down

sqlc-gen:
	cd server && sqlc generate

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
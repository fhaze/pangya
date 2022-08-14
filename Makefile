include .env
DB_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

deps:
	go mod download

build:
	mkdir -p ./bin
	go build -o ./bin/login_server ./src/cmd/login_server/main.go

test:
	go test ./...

dev-migrate-up:
	migrate -source file://migrations -database ${DB_DSN} up

dev-migrate-down:
	migrate -source file://migrations -database ${DB_DSN} down

dev-workspace:
	docker-compose up -d

dev: dev-workspace dev-migrate-up

run: dev
	go run src/cmd/login_server/main.go

clean:
	docker-compose down

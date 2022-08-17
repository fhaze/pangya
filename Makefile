include .env
DB_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

deps:
	go mod download

build:
	mkdir -p ./bin
	go build -o ./bin/sync_server ./src/cmd/sync_server/main.go
	go build -o ./bin/login_server ./src/cmd/login_server/main.go
	go build -o ./bin/game_server ./src/cmd/game_server/main.go

test:
	go test ./...

dev-migrate-up:
	migrate -source file://migrations -database ${DB_DSN} up

dev-migrate-down:
	migrate -source file://migrations -database ${DB_DSN} down

dev-compose-up:
	docker-compose up -d

dev: build dev-compose-up dev-migrate-up

run: dev
	tmux new-session -d -s PangyaServer -n Shell -d "bin/sync_server; sleep 100"
	tmux split-window -t PangyaServer "bin/login_server; sleep 100"
	tmux split-window -t PangyaServer "bin/game_server; sleep 100"
	tmux select-layout -t PangyaServer even-vertical
	tmux attach -t PangyaServer

run2: dev
	tmux new-session -d -s PangyaServer -n Shell -d "bin/sync_server; sleep 100"
	tmux split-window -t PangyaServer "bin/login_server; sleep 100"
	tmux split-window -t PangyaServer "cd bin; dotenv -f ../.env run ./game_server; sleep 100"
	tmux split-window -t PangyaServer "cd bin; dotenv -f ../.env2 run ./game_server; sleep 100"
	tmux select-layout -t PangyaServer tiled
	tmux attach -t PangyaServer

dev-compose-down:
	docker-compose down --volume

clean: dev-compose-down
	rm -r ./bin

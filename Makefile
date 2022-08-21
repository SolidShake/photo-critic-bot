build:
	go build -o ./out/photo-critic-bot

run:
	go run ./cmd/main.go

lint:
	golangci-lint run

start:
	docker-compose up -d

stop:
	docker-compose down

clear:
	docker volume rm photo-critic-bot_postgres_data

restart:
	make stop && make clear && make start

# connect to dev db
connect:
	docker exec -it photo-critic-bot_db_1 psql -U postgres

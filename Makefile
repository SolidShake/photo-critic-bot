build:
	go build -o ./out/photo-critic-bot

run:
	go run .

lint:
	golangci-lint run


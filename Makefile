.PHONY: up down run build clean logs

up:
	docker compose up -d

down:
	docker compose down

run: up
	go run main.go

build:
	go build -o bin/playground main.go

clean:
	docker compose down -v
	rm -f bin/playground

logs:
	docker compose logs -f postgres

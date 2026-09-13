.PHONY: run build docker-up clean

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

clean:
	rm -rf bin/

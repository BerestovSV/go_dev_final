.PHONY: build run test clean docker-build docker-run

build:
	go build -o todo-server .

run:
	go run main.go

test:
	go test ./... -v

clean:
	rm -f todo-server
	rm -f *.db

docker-build:
	docker build -t todo-server .

docker-run:
	docker run -p 7540:7540 -e TODO_PASSWORD=mysecretpassword todo-server

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down
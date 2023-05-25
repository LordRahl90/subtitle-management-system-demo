start:
	go run ./cmd

test:
	go test ./... --cover

test-with-race:
	go test -race ./... --cover

build:
	docker build -t lordrahl/translations:latest .

docker-start: build
	docker-compose up

twr: test-with-race
ds: docker-start
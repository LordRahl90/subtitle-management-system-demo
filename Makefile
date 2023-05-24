start:
	go run ./cmd

test:
	go test ./... --cover

test-with-race:
	go test -race ./... --cover

twr: test-race
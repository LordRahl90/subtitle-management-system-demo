start:
	go run ./cmd

test:
	go test ./... --cover -v

test-with-race:
	go test -race ./... --cover

twr: test-with-race
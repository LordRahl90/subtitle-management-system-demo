start:
	go run ./cmd

test:
	go test ./... --cover

test-race:
	go test -race ./... --cover

tr: test-race
.PHONY: build test clean run

BINARY_NAME=trackpulse

build:
	go build -o $(BINARY_NAME) ./cmd/trackpulse

build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe ./cmd/trackpulse

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) ./cmd/trackpulse

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	go clean

run:
	go run ./cmd/trackpulse

migrate:
	go run ./migrations/migrate.go

uuid:
	go run ./scripts/generate_uuid.go
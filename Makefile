build:
	CGO_ENABLED=0 go build -o etherecho ./cmd/main.go

test:
	go test ./...
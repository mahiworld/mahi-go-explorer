dev:
	@ go run cmd/app/main.go

test:
	@ go test ./...

build:
	@ mkdir -p ./dist
	@ go mod download
	@ GOOS=linux GOARCH=amd64 go build -o ./dist/mahi-go-explorer ./cmd/app/

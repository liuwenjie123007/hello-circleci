.PHONY: error
error:
	exit 1

.PHONY: dev
dev:
	go run cmd/main.go


.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/admin-back cmd/admin-back/main.go

.PHONY: test
test:
	go test ./...

.PHONY: docker
docker:
	apt update -y
	apt install -y clang-format=1:11.0-51+nmu5

.PHONY: fmt
fmt:
	go fmt ./...

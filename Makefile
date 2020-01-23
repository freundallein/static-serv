export BIN_DIR=bin
export STATIC_ROOT=/storage
export PREFIX=/static
export PORT=8000
export IMAGE_NAME=freundallein/staticserv:latest

init:
	git config core.hooksPath .githooks
run:
	go run main.go
test:
	go test -cover ./...
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o $$BIN_DIR/staticserv
dockerbuild:
	docker build -t $$IMAGE_NAME -f Dockerfile .
distribute:
	echo "$$DOCKER_PASSWORD" | docker login -u "$$DOCKER_USERNAME" --password-stdin
	docker build -t $$IMAGE_NAME .
	docker push $$IMAGE_NAME
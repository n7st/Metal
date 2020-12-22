build:
	go build -o metal cmd/metal/main.go

build_plugins:
	bash ./scripts/build-plugins.sh

test:
	go test -v ./...

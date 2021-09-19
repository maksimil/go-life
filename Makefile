run:
	go run ./cmd/life/main.go

build:
	go build -ldflags "-s -w" ./cmd/life/main.go

install:
	go install ./cmd/navgo
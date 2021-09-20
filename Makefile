run: gen
	go run ./cmd/life

build: gen
	go build -ldflags "-s -w" ./cmd/life

install: gen
	go install ./cmd/life

gen:
	go run ./cmd/gen
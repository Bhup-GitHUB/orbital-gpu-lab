.PHONY: all bitflip protect radiator test

all: bitflip protect radiator

bitflip:
	go run ./cmd/bitflip

protect:
	go run ./cmd/protect

radiator:
	go run ./cmd/radiator

test:
	go test ./...

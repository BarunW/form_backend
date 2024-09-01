all : build run test
.PHONY : all

build:
	go build .

run:
	go run .

test:
	go test .

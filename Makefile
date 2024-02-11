.PHONY: help test run build clean

help:
	@echo "test - run tests"
	@echo "run - run the application"
	@echo "build - build the application"
	@echo "clean - remove the bin directory"

test:
	go test -v tests/*.go
run:
	go run cmd/main.go
build:
	go build -o bin/main cmd/main.go
clean:
	rm -rf bin/*
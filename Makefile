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
lint:
	golangci-lint run

docker-build:
	docker build -t mq-service .
docker-run:
	docker run -p 8080:8080 mq-service

helm-install:
	helm install mq-service ./deployments/charts/app
helm-upgrade:
	helm upgrade mq-service ./deployments/charts/app
helm-uninstall:
	helm uninstall mq-service
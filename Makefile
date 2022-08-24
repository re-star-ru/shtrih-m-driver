
run:
	go run ./app

docker:
	docker build --tag shtrihdriver .

lint:
	golangci-lint run

.PHONY: build, docker, lint

build:

podman:
	podman build --tag shtrihdriver .

docker:
	docker build --tag shtrihdriver .

lint:
	golangci-lint.exe run

.PHONY: build, podman, docker, lint
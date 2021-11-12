
build:

podman:
	podman build --tag shtrihdriver .

docker:
	docker build --tag shtrihdriver .

.PHONY: build, podman, docker
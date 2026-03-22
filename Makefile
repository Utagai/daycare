.PHONY: build-docker test-e2e

build:
	go build .

build-docker:
	docker build -t daycare-test .

test-unit:
	go test ./... -count=1

test-e2e: build-docker
	docker run --rm daycare-test

test: test-unit test-e2e

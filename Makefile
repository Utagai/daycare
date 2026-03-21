.PHONY: build-docker build-fast test-e2e docker-shell docker-clean

build:
	go build .

build-docker:
	docker build --no-cache -t daycare-test .

build-fast:
	docker build -t daycare-test .

test-unit:
	go test ./... -count=1

test-e2e: build-docker
	docker run --rm -e CI=true daycare-test

test: test-unit test-e2e

docker-shell: build-docker
	docker run -it --rm -e CI=true daycare-test bash

docker-clean:
	docker rm -f daycare-test 2>/dev/null || true

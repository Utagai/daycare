.PHONY: build-docker build-fast test-e2e docker-shell docker-clean

build-docker:
	docker build --no-cache -t daycare-test .

build-fast:
	docker build -t daycare-test .

test-e2e: build-docker
	docker run --rm -e CI=true daycare-test

docker-shell: build-docker
	docker run -it --rm -e CI=true daycare-test bash

docker-clean:
	docker rm -f daycare-test 2>/dev/null || true

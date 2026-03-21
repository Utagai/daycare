.PHONY: build-docker build-fast test-e2e docker-shell docker-clean

build-docker:
	docker build --no-cache -t daycare-test .

build-fast:
	docker build -t daycare-test .

test-e2e: build-docker
	docker compose run --rm e2e-tests

docker-shell: build-docker
	docker compose run --rm e2e-tests bash

docker-clean:
	docker compose down -v --rmi local 2>/dev/null || true

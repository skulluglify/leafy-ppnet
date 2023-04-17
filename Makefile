.PHONY: update build

sync:
	git submodule sync --recursive

update:
	git submodule update --init --recursive
	git submodule update --recursive --remote

build:
	mkdir -p build/bin
	env DOCKER_BUILDKIT=1 docker-compose -f compose.yml -p leafy-ppnet build

clean:
	rm -rf build
	docker image prune

up:
	docker-compose up -d

down:
	docker-compose down

.PHONY: update build

update:
	git submodule sync --recursive
	git submodule update --init --recursive
	git submodule update --recursive --remote

build: update
	mkdir -p build/bin
	env CGO_ENABLED=0 go build -o build/bin/app.exe main.go

docker-compose-build:
	sudo env DOCKER_BUILDKIT=1 docker-compose -f compose.yml -p leafy build

docker-compose-up: docker-compose-build
	sudo docker-compose up

docker-build:
	sudo env env DOCKER_BUILDKIT=1 docker build -f Dockerfile -t leafy .

docker-down:
	docker-compose down

docker-up: docker-compose-up

clean:
	rm -rf build
	docker image prune

dump:
	sudo -u postgres pg_dump --dbname leafy --schema-only >postgresql.sql

go-test:
	go test -v test/transaction_test.go
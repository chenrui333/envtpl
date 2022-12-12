.PHONY: clean docker test

build:
	go build -ldflags "-s -w" -o bin/envtpl "./cmd/envtpl/."

clean:
	rm -rf bin

image:
	docker build -t subfuzion/envtpl .

test:
	docker-compose -f docker-compose.test.yml build
	docker-compose -f docker-compose.test.yml run --rm sut

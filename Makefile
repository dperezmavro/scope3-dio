.PHONY: build lint run test
TOKEN=$(shell cat token.secret.txt) # this is a one-line file containing an auth token

build:
	CGO_ENABLED=0 cd src && go build -o ../bin/main

docker:
	docker build \
		--platform linux/amd64 \
		--tag "dio-scope3" \
		-f Dockerfile .

lint:
	golangci-lint run ./...

run:
	make build
	SCOPE3_API_TOKEN="${TOKEN}" VERSION=1 PORT=3000 ENV=dev SERVICE=dio-scope3 ./bin/main

test:
	cd src && go test ./...
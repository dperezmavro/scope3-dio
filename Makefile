.PHONY: build lint run test
TOKEN=$(shell echo ${SCOPE3_API_TOKEN}) # expects the SCOPE3_API_TOKEN to be defined

build:
	CGO_ENABLED=0 cd src && go build -o ../bin/main

docker-build:
	docker build \
		--platform linux/amd64 \
		--tag "dio-scope3" \
		-f Dockerfile .

docker-run:
	make docker-build
	docker run -p \
		3000:3000 \
		-e SCOPE3_API_TOKEN=${TOKEN} \
		--rm -ti \
		--name dio-scope3 \
		dio-scope3:latest

lint:
	golangci-lint run ./...

run:
	make build
	SCOPE3_API_TOKEN=${TOKEN} VERSION=1 PORT=3000 ENV=dev SERVICE=dio-scope3 ./bin/main

test:
	cd src && go test ./...

bench:
	go test \
		-benchmem \
		-run='' \
		-bench 'BenchmarkChannels' \
		-benchtime=10000x \
		-memprofile memprofile.out \
		github.com/dperezmavro/scope3-dio/src/clients/scope3
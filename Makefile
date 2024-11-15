.PHONY: build run

build:
	CGO_ENABLED=0 cd src && go build -o ../bin/main

run:
	make build
	SCOPE3_API_TOKN=qwerty PORT=3000 ENV=local SERVICE=abc ./bin/main
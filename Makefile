.PHONY: build

build:
	CGO_ENABLED=0 cd src && go build -o ../bin/main 
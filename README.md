# About

This is my submissions for the API Lead Take Home Challenge from Scope3.

# Overview
* choise of language
* choice of libraries (ristretto, zerolog)
* current implementation
    * metrics/healthcheck
    * wait-limit for fetching missing records
    * stripped-out response (limitation)
    * versioned paths mirroring actual-api
    * weighted inputs

# Building
You can build either the standalone binary, or a docker container containing the binary. 

To build the binary, run `make build`. This will generate a file in `./bin/main`. For this to work you need to have installed Go on your machine.

To build the Docker container, run `make docker-build`. This will make a container tagged as `dio-scope3:latest`. For this to work you need to have installed and configured docker on your machine. 

# Running
* Environment variables
* configuration
* wait-limit
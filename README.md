# About

This is my submission for the API Lead Take Home Challenge from Scope3.

I wrote this service from scratch!

# Overview
The task prompt:
> We are building an api to return emission measurement values for a given property and need to optimize for speed of response over completeness. [...] The team has asked you to design and build a lightweight cache with api to sit between our internal measurements api and our customers 

## General tech Choices

The service itself is written in Go, as Go provides an excellent balance between speed of development, runtime safety, runtime performance and ease of distribution. 

The in-memory caching aspect itself (key eviction, TTL management etc) is offloaded to a third-party library, Ristretto (https://github.com/dgraph-io/ristretto) because it is a proven high-performance in-memory cache.

Because this is meant to be a high-throughput service that will scale horizontally, I implemented a coupled of features in order to make runtime management and debugging easier using these endpoints:
* a `/healthcheck` endpoint that would be used during deployment and runtime to determine if the service is healthy
* a `/metrics` endpoint that exposes some Ristretto metrics to help us benchmark performance

These two endpoints, whilst not directly related to the challenge, is what I would expect to see from any production-ready service. This service also offers structured logging with TraceIds in the logs, to make debugging in the presence of multiple instances or concurrent events easier. I used zerolog for this as it is a fast, 0-allocation logger with minimal overhead.

## API implementation 
For the actual API implementation I exposed a `/v2/measure` API that takes the same input contract as Scope3's API. Keeping the same API interface (versioned API endpoint and request types) would make it easier for clients to switch to this faster interface with minimal configuration changes (just repointing the URL) and also minimise the likelihood of breaking changes on their end. 

For the purposes of this interview, I slimmed down my API responses quite a bit so that you only get the `totalEmissions`, `utcDateTime`, `weight` and `propertyName` as part of the response. I also implemented two features that would help this become more useful:
* waiting for results: You can configure this service (using the `WAIT_FOR_MISSING=true` environment variable) to wait for results in case there is a cache miss. This works by configuring a 45ms timeout during which the service is waiting to see if the in-memory cache has these results available so that they can be included in the response.

* weighting of results: I added a new `weight` parameter (type `int`) to the incoming request to allow the client to specify a higher weight for the result to be cached (higher weight in Ristretto means less likelihood of eviction), making the cached result less likely to evict.

The service itself is making heavy use of channels and goroutines in order to fetch information asynchronously from the Scope3 API, always being mindful of keeping the critical path as fast as possible. I had to make a design choice around handling cache misses, and I chose to return an empty result if there's a cache miss and use a go-routine in the background to fetch the data and save it in the cache for subsequent requests. The diagram bellow shows the default flow:

![](flowchart.drawio.svg "Default flowchart when WAIT_FOR_MISSING=false")

When waiting is configured, the service behaves differently. The service defines a 45ms window during which it will wait and query the cache for the missing results. Missing records are requested in the background using goroutines, so there is a chance that they will have succeeded during that window. This would provide a more user-friendly experience and avoid returning empty results. A flow-chart of this is shown below:

![](flowchart-with-wait.drawio.svg "Default flowchart when WAIT_FOR_MISSING=true")

# Building
You can build either the standalone binary, or a docker container containing the binary. 

To build the binary, run `make build`. This will generate a file in `./bin/main`. For this to work you need to have installed Go on your machine.

To build the Docker container, run `make docker-build`. This will make a container tagged as `dio-scope3:latest`. For this to work you need to have installed and configured docker on your machine. 

# Running
To run this locally with default values, you can run `SCOPE3_API_TOKEN="TOKEN_HERE" make run` - this will compile the binary and run it on your local host.

There are a few variables you can tweak in this service, and you can access them all through environment config:
* `ENV`: the name of the environment (defaults to `local`)
* `PORT`: the port to listen on (defaults to `3000`)
* `SERVICE`: the name of the service (defaults to `localService`)
* `VERSION`: the version of the service (defaults to `0`)
* `WAIT_FOR_MISSING`: this configures the cache to wait for 45ms in case some of the missing results appear in the meantime from the async fetching, in order to avoid returning an empty response (defaults to false, no waiting)
* `SCOPE3_API_TOKEN`: token for authenticating to the scope3 api. This is mandatory and there are no defaults
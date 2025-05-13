## Go User Service


Go service skeleton that includes a lot of boilerplate code needed by most of the Go services:
- HTTP server and router exposing `/health` and `/users` endpoints
- metrics server exposing metrics on `/metrics` endpoint
- optional pprof server exposing data on separate port
- metrics and logging middleware for HTTP requests 
- integration test for `/health` endpoint
- optimised Dockerfile
- Makefile with optimise go imports, go linting, running tests, building Docker image and running Docker container for the service

### Build service locally

```shell
make
```

### Run service locally in Docker

```shell
make run
```
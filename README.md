# Overview 
Golang implementation of microservices that manage Shattered Realms Online.

## Components
* `api` Proto definitions used in the architecture
* `cmd` Microservices
  * `characters` User playable character management
  * `chat` Chat service for sending messages between players and in chat channels
  * `gamebackend` Management of game servers, connections and settings
* `pkg`: Global components shared between microservices
  * `config` Configuration parameters
  * `helpers` Utility functions used throughout the microservices
  * `pb` Auto-generated go protobuf and gRPC files
  * `repository` Abstraction layer to the underlying database and storage mechanisms
  * `service` Services that use the repositories to perform necessary actions
  * `srv` gRPC service implementations
* `test` Configuration files and data for testing
  * `db` Temporary testing database connection 
* `third_party` Third party proto definitions

# Development
The `Makefile` is located within the `build` folder within the project root directory. All make commands should be run from there. A docker compose file has been included to run all required services.

## Requirements
* Docker
* Docker compose
* Golang 19

## Environments
This project uses environment variables during the build process which should be stored within a `.env` file located within the project root directory. If one is not configured, rename `.env.template` to `.env` and configure the variables for deployment.

## Commands
### Building
**Binary:** To build all binarys run `make build`. The output result will be placed in the `bin` folder in the project root directory. To build a specific file run `make build-<app-name>`\
**Docker:** To build the docker image run `make buildi` and a image called `sro-<app-name>` will be generated. To build a specific docker image run `make build-<app-name>-image`

### Testing
To run all tests and see the coverage report use `make test`. To view HTML results, simply run `make report`.

### Deployment
Deployment is done using docker. If using an AWS docker repository, running `make aws-docker-login` will authenticate with the default aws credential context. To push the images, run `make push`. This will build the image and push them to the docker repository.

### Full backend orchestration
To test full integration using HTTP2 and envoy proxies, use docker compose. Then request from a frontend service can be made through the envoy proxy `https://localhost:9090` requesting a specific service such as chat `https://localhost:9090/chat`. 
```
docker-compose pull
docker-compose up --build -d
```

#### Envoy: localhost network
To run services locally for quicker build iterations, you can use a single envoy proxy that queries the full backend orchestarted docker-compose setup.
```
docker build ./build/envoy -t sro-local-envoy -f ./build/envoy/local.Dockerfile
docker run --rm --add-host host.docker.internal:host-gateway -p 9091:9091 --name="sro-local-envoy" sro-local-envoy
```

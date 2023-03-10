# Overview 

* `characters` - User playable characters management.
* `chat` - Chat service for sending messages between players and in channels
* `gamebackend` - Management of game servers and settings

# Development
The `Makefile` is located within the `build` folder within the project root directory. All make commands should be run from there. A docker compose file has been included to run all required services.

## Requirements
* Docker
* Docker compose
* Golang 19

## Environments
This project uses environment variables which should be stored within a `.env` file located within the project root directory. If one is not configured, rename `.env.template` to `.env` and configure the variables for deployment. 

## Commands
### Building
**Binary:** To build all binarys run `make build`. The output result will be placed in the `bin` folder in the project root directory. To build a specific file run `make build-<app-name>`\
**Docker:** To build the docker image run `make build-image` and a image called `sro-<app-name>` will be generated. To build a specific docker image run `make build-<app-name>-image`

### Testing
To run all tests and see the coverage report use `make test`. To view HTML results, simply run `make report`.

### Deployment
Deployment is done using docker. If using an AWS docker repository, running `make aws-docker-login` will authenticate with the default aws credential context. To push the images, run `make push`. This will build the image and push them to the docker repository. To push to a specific environment, run `make push-dev`, `make push-qa`, or `make push-prod`. By default, it will push to the development environment.
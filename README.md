# Overview 

* `accounts` - User accounts authentication and and authorization.
* `characters` - User playable characters management.

# Development
The `Makefile` is located within the `build` folder within the project root directory. All make commands should be run from there.

## Environments
This project uses environment variables which should be stored within a `.env` file located within the project root direcory. If one is not configured, rename `.env.template` to `.env` and configure the variables for deployment. These variables can be overwritten in the OS, in a docker env file, kubernetes env file, and at runtime by supplying them before the run command.

## Commands
### Building
**Binary:** To build all binarys run `make build`. The output result will be placed in the `bin` folder in the project root directory. To build a specific file run `make build-<app-name>`\
**Docker:** To build the docker image run `make build-image` and a image called `sro-<app-name>` will be generated. To build a specific docker image run `make build-<app-name>-image`

### Testing
To run all tests and see the coverage report use `make test`. To view a the HTML results, simply run `make report`.

### Deployment
Deployment is done using docker. If using an AWS docker repository, running `make aws-docker-login` will authenticate with the default aws credential context. To push the images, run `make push`. This will build the image and push them to the docker repository. To push to a specific environment, run `make push-dev`, `make push-qa`, or `make push-prod`. By default it will push to the development environment.
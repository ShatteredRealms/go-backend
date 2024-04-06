#)####################################################################################
#   _____ _           _   _                    _   _____            _               #
#  / ____| |         | | | |                  | | |  __ \          | |              #
# | (___ | |__   __ _| |_| |_ ___ _ __ ___  __| | | |__) |___  __ _| |_ __ ___  ___ #
#  \___ \| '_ \ / _` | __| __/ _ \ '__/ _ \/ _` | |  _  // _ \/ _` | | '_ ` _ \/ __|#
#  ____) | | | | (_| | |_| ||  __/ | |  __/ (_| | | | \ \  __/ (_| | | | | | | \__ \#
# |_____/|_| |_|\__,_|\__|\__\___|_|  \___|\__,_| |_|  \_\___|\__,_|_|_| |_| |_|___/#
#####################################################################################

#
# Makefile for building, running, and testing
#

# Import dotenv
ifneq (,$(wildcard ../.env))
	include ../.env
	export
endif

# Application versions
BASE_VERSION = $(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')
COMMIT_HASH = $(shell git rev-parse --short HEAD)


# Gets the directory containing the Makefile
ROOT_DIR = $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# Base container registry
SRO_BASE_REGISTRY ?= 779965382548.dkr.ecr.us-east-1.amazonaws.com
SRO_REGISTRY ?= $(SRO_BASE_REGISTRY)/sro

# The registry for this service
REGISTRY = $(SRO_REGISTRY)/accounts
time=$(shell date +%s)

PROTO_DIR=$(ROOT_DIR)/api

PROTO_FILES = $(shell find $(PROTO_DIR) -name '*.proto')

MOCK_INTERFACES = $(shell egrep -rl --include="*.go" "type (\w*) interface {" $(ROOT_DIR)/pkg | sed "s/.go$$//")

ALL_CMDS = $(shell find $(ROOT_DIR)/cmd -maxdepth 1 -mindepth 1 -type d | sed "s/^.*\///")
CMDS = $(filter-out stressy, $(ALL_CMDS))

# Output folder for TS files for the frontend
# If it's empty, nothing will be generated
TS_ROOT_DIR = $(shell realpath $(ROOT_DIR)/../)/frontend
TS_OUT_DIR = $(TS_ROOT_DIR)/src/app/generated

#   _____                    _
#  |_   _|                  | |
#    | | __ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    \_/\__,_|_|  \__, |\___|\__|___/
#                  __/ |
#                 |___/

.PHONY: test report mocks clean-mocks report-watch
test:
	ginkgo --race -p --cover -covermode atomic -coverprofile=coverage.out --output-dir $(ROOT_DIR)/ $(ROOT_DIR)/pkg/...

test-e2e:
	@echo "Starting e2e integration tests"
	ginkgo --race -p $(ROOT_DIR)/test/e2e/...

test-watch:
	ginkgo watch --race -p --cover -covermode atomic -output-dir=$(ROOT_DIR) $(ROOT_DIR)/...

report: test
	go tool cover -func=$(ROOT_DIR)/coverage.out
	go tool cover -html=$(ROOT_DIR)/coverage.out

report-watch:
	while inotifywait -e close_write $(ROOT_DIR)/coverage.out; do \
		go tool cover -func=$(ROOT_DIR)/coverage.out; \
		go tool cover -html=$(ROOT_DIR)/coverage.out; \
	done

dev-watch: test-watch report-watch

mocks: clean-mocks
	mkdir -p $(ROOT_DIR)/pkg/mocks
	@for file in $(MOCK_INTERFACES); do \
		mockgen -package=mocks -source=$${file}.go -destination="$(ROOT_DIR)/pkg/mocks/$${file##*/}_mock.go"; \
	done

clean-mocks:
	rm -rf $(ROOT_DIR)/pkg/mocks

build: $(addprefix build-, $(CMDS)) buildclis
build-%:
	go build -ldflags="-X 'github.com/ShatteredRealms/go-backend/pkg/config/default.Version=$(BASE_VERSION)'" -o $(ROOT_DIR)/bin/$* $(ROOT_DIR)/cmd/$*  
buildclis:
	go build -o $(ROOT_DIR)/bin/stressy $(ROOT_DIR)/cmd/stressy

run: $(addprefix run-, $(CMDS))
run-%:
	go run $(ROOT_DIR)/cmd/$*

run-watch: $(addprefix run-watch-, $(CMDS))
run-watch-%:
	gow run $(ROOT_DIR)/cmd/$*

deploy: aws-docker-login push

buildi: $(addprefix buildi-, $(CMDS))
buildi-%:
	docker build --build-arg APP_VERSION=$(BASE_VERSION) -t sro-$* -f build/$*.Dockerfile .

aws-docker-login:
	aws ecr get-login-password | docker login --username AWS --password-stdin $(SRO_BASE_REGISTRY)

pushf: $(addprefix pushf-, $(CMDS))
pushf-%:
	docker tag sro-$* $(SRO_REGISTRY)/$*:latest
	docker tag sro-$* $(SRO_REGISTRY)/$*:$(BASE_VERSION)
	docker tag sro-$* $(SRO_REGISTRY)/$*:$(BASE_VERSION)-$(COMMIT_HASH)
	docker push $(SRO_REGISTRY)/$*:latest
	docker push $(SRO_REGISTRY)/$*:$(BASE_VERSION)
	docker push $(SRO_REGISTRY)/$*:$(BASE_VERSION)-$(COMMIT_HASH)

push: $(addprefix push-, $(CMDS))
push-%: buildi-%
	docker tag sro-$* $(SRO_REGISTRY)/$*:latest
	docker tag sro-$* $(SRO_REGISTRY)/$*:$(BASE_VERSION)
	docker tag sro-$* $(SRO_REGISTRY)/$*:$(BASE_VERSION)-$(COMMIT_HASH)
	docker push $(SRO_REGISTRY)/$*:latest
	docker push $(SRO_REGISTRY)/$*:$(BASE_VERSION)
	docker push $(SRO_REGISTRY)/$*:$(BASE_VERSION)-$(COMMIT_HASH)

.PHONY: clean-protos protos $(PROTO_FILES)

clean-protos:
	rm -rf "$(ROOT_DIR)/pkg/pb"

protos: clean-protos $(PROTO_FILES) mocks

$(PROTO_FILES):
	protoc "$@" \
		-I "$(PROTO_DIR)" \
		--go_out="$(ROOT_DIR)" \
		--go-grpc_out="$(ROOT_DIR)" \
		--grpc-gateway_out="$(ROOT_DIR)" \
		--grpc-gateway_opt "logtostderr=true" \
		--plugin=protoc-gen-ts="$(TS_ROOT_DIR)/node_modules/.bin/protoc-gen-ts" \
		--ts_out=service=grpc-web:"$(TS_OUT_DIR)" \
		--js_out=import_style=commonjs:"$(TS_OUT_DIR)"

download:
	go mod download

install-tools:
	  cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

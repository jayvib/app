BINARY="engine.linux"
NOW=`date +'%y%m%d%H%M%S'`

################TARGETS##################
test: preparetest ## Runs unit and integration test
	go test -tags="unit integration"-v -cover -covermode=atomic ./...

###############TOOL VERSIONS############
PROTOBUF_FILENAME="protobuf-all-3.9.1.tar.gz"
PROTOBUF_VERSION="3.9.1"

unit-test: preparetest ## Runs unit test only
	@go test -tags=unit -covermode=atomic -short ./... | grep -v '^?'

integration-test: preparetest ## Runs integration test only
	@go test -tags="integration elasticsearch mysql dynamo"-short ./... | grep -v '^?'

start: ## Start the app application
	docker-compose start

stop: ## Stop the app application
	docker-compose stop

run: docker ## Build and run the application
	docker-compose -f docker-compose.yml up -d

teardown:  ## Stop and remove the created containers and images
	docker-compose -f docker-compose.yml down

run-development-services: ## Build and run services use for testing.
	docker-compose -f docker-compose-services.yml up -d

teardown-development-services: ## Stop and remove the created containers and images
	docker-compose -f docker-compose-services.yml down

start-development-services: ## Start the services use for testing
	docker-compose -f docker-compose-services.yml start

stop-development-service: ## Stop the services use for testing
	docker-compose -f docker-compose-services.yml stop

run-es-development: ## Build and run the Elasticsearch for testing or development
	docker-compose -f docker-compose-test-elasticsearch.yml up -d

teardown-es-development: ## Stop and remove the Elasticsearch containers and images
	docker-compose -f docker-compose-test-elasticsearch.yml down -d

start-es-development: ## Start the Elasticsearch for testing or development
	docker-compose -f docker-compose-test-elasticsearch.yml start

stop-es-development: ## Stop the Elasticsearch for testing or development
	docker-compose -f docker-compose-test-elasticsearch.yml stop

fmt: ## Format source files excluding the vendor directory
	@go fmt -x ./...

mod: ## Downloads the module use in the project.
	go mod download

engine: mod ## Build the binary file of the project
	go build -o ${BINARY} .

docker:  ## Build a docker image of the project
	sudo docker build -t clean-architecture .

dockerpush: ## Push the docker image to the repository
	@echo ${NOW}
	@echo "Setting image tag -> latest ${NOW}"
	@sudo docker tag clean-architecture jayvib/clean-architecture:latest
	@sudo docker tag clean-architecture jayvib/clean-architecture:${NOW}
	@echo "Pushing image 	 -> jayvib/clean-architecture"
	@sudo docker push jayvib/clean-architecture

list-targets: ## List the available targets that can be user
	@echo '#####Makefile Targets######'
	@grep '^[^#[:space:]].*:' Makefile | sed 's/:.*//g'

help: ## Display the available targets and its description
	# Learned this from:
	# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## Lints the project source code
	golint $(shell go list ./... | grep -v /vendor/)

preparetest: ## Pre-process before actual testing start
	@mkdir -p ~/.app
	@cp devel-config.json ~/.app/devel-config.json
	@cp devel-config.json ~/devel-config.json

preparestaging: ## Pre-process before actual staging start
	@mkdir -p ~/.app
	@cp staging-config.json ~/.app/staging-config.json
	@cp staging-config.json ~/staging-config.json

prepareprod: ## Pre-process before actual production start
	@mkdir -p ~/.app
	@cp config.json ~/config.json
	@cp config.json ~/.app/config.json

preparetools: ## Downloads the tools use in the project.
	@echo Installing mockery
	@go get github.com/vektra/mockery/cmd/mockery
	@echo Installing golint
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/golang/protobuf/protoc-gen-go

install-protobuf-compiler: ## Installing the protocol buffers compiler
	sudo apt-get install autoconf automake libtool curl make g++ unzip
	mkdir -p ~/Downloads/protobuf
	wget -x -O ~/Downloads/protobuf/protobuf.tar.gz https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/${PROTOBUF_FILENAME}
	cd ~/Downloads/protobuf && \
		tar xvzf protobuf.tar.gz && \
		./protobuf-${PROTOBUF_VERSION}/configure &&\
		make && \
		make check && \
		sudo make install && \
		sudo ldconfig

.PHONY: test unittest integrationtest \
		start stop run teardown runtestservice \
		teardowntestservice starttestservice  \
		stoptestservice fmt mod engine docker \
		dockerpush list lint-prepare lint \
		preparetest tools-prepare run-es-development \
		teardown-es-development start-es-development \
		stop-es-development list-targets help \
		preparestaging prepareprod preparetools \
		install-protobuf-compiler



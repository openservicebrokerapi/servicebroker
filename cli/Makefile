.PHONY: all cli test-unit test clean validate

EXE=sc
SRC_FILES=$(shell find . -name \*.go -print | grep -v _test.go | grep -v \\./tests)
TEST_UNIT_DIRS=$(shell find . -name \*_test.go -exec dirname {} \; | grep -v \\./tests | sort -u)

all: validate cli test-unit test

#####################

cli: $(EXE)

$(EXE): $(SRC_FILES)
ifeq (,$(shell which go))
	@$(MAKE) --no-print-directory cli-docker
else
	@$(MAKE) --no-print-directory cli-local
endif

cli-docker:
	docker run -ti \
	  -w /src/github.com/cncf/servicebroker/cli \
	  -v ${CURDIR}:/src/github.com/cncf/servicebroker/cli \
	  -e GOPATH=/ \
	  golang \
	sh -c "go get -d ./... && make cli-local"

cli-local:
	go get -d ./...
	go build -o $(EXE) .

#####################

test-unit: cli
	@echo "\n### Running the unit tests"
ifeq (,$(shell which go))
	@$(MAKE) --no-print-directory test-unit-docker
else
	@$(MAKE) --no-print-directory test-unit-local
endif

test-unit-docker:
	docker run -ti \
	  -w /src/github.com/cncf/servicebroker/cli \
	  -v ${CURDIR}:/src/github.com/cncf/servicebroker/cli \
	  -e GOPATH=/ \
	  golang \
	sh -c "go get -d ./... && make --no-print-directory test-unit-local"

test-unit-local:
	@go test $(TEST_UNIT_DIRS)

#####################

test: cli
	@echo "\n### Running the integration tests"
ifeq (,$(shell which go))
	@$(MAKE) --no-print-directory test-docker
else
	@$(MAKE) --no-print-directory test-local
endif

test-docker:
	docker run -ti \
	  -w /src/github.com/cncf/servicebroker/cli \
	  -v ${CURDIR}:/src/github.com/cncf/servicebroker/cli \
	  -e GOPATH=/ \
	  --net host \
	  golang \
	sh -c "go get -d ./... && make --no-print-directory test-local"

test-local: cli
	echo No tests yet
	# @go test -v ./tests

#####################

clean:
	-rm -f $(EXE)

#####################

validate: 
ifeq (,$(shell which go))
	@$(MAKE) --no-print-directory validate-docker
else
	@$(MAKE) --no-print-directory validate-local
endif

validate-docker:
	docker run -ti \
	  -w /src/github.com/cncf/servicebroker/cli \
	  -v ${CURDIR}:/src/github.com/cncf/servicebroker/cli \
	  -e GOPATH=/ \
	  golang \
	sh -c "go get -d ./... && make --no-print-directory validate-local"

validate-local:
	@if ! [ "`find . -name \*.go -exec gofmt -s -l {} \;`" = "" ]; then \
		echo Please gofmt the following files: \\n$(shell find . -name \*.go -exec gofmt -s -l {} \;) \\n ; \
		exit 1 ; \
	fi

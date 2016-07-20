# .PHONY: all binary test-unit test clean server

EXE:=service_controller
IMAGE:=service-controller
GO_FILES=$(shell find . -name \*.go -print | grep -v _test.go)
TEST_UNIT_DIRS=$(shell find . -name \*_test.go -exec dirname {} \; | grep -v ./tests | sort -u)

all: binary

binary: $(EXE)

$(EXE): $(GO_FILES)
	docker run \
	  -w /src/github.com/cncf/servicebroker/k8s/service_controller \
	  -v ${CURDIR}:/src/github.com/cncf/servicebroker/k8s/service_controller \
	  -e GOPATH=/ \
	  golang \
	  sh -c "go get -d ./... && go build ."
	docker build -t $(IMAGE):latest .

test-unit: binary
	docker run \
	  -w /src/github.com/cncf/servicebroker/k8s/service_controller \
	  -v ${CURDIR}:/src/github.com/cncf/servicebroker/k8s/service_controller \
	  -e GOPATH=/ \
	  golang \
	  sh -c "go get -d ./... && make --no-print-directory test-unit-local"
	
test-unit-local:
	@TEST_UNIT_DIRS="$(shell find . -name \*_test.go -exec dirname {} \; | sort -u)"
	@echo "\n### Running the unit tests"
	@go test $(TEST_UNIT_DIRS)

test: binary
	docker run \
	  -w /src/github.com/cncf/servicebroker/k8s/service_controller \
	  -v ${CURDIR}:/src/github.com/cncf/servicebroker/k8s/service_controller \
	  -e GOPATH=/ \
	  --net host \
	  golang \
	  sh -c "go get -d ./... && make --no-print-directory test-local"

test-local: binary
	@echo "\n### Running the integration tests"
	@go test ./tests

clean:
	rm -f $(EXE)
	docker rmi -f $(IMAGE) || echo gone

server: binary
	docker run -ti -p 10000:10000 $(IMAGE)

# .PHONY: all binary test-unit clean server

EXE:=service_controller
IMAGE:=service-controller
GO_FILES=$(shell find . -name \*.go -print | grep -v _test.go)
TEST_DIRS=$(shell find . -name \*_test.go -exec dirname {} \; | sort -u)

all: binary

binary: $(EXE)

$(EXE): $(GO_FILES)
	docker run \
	  -w /src/github.com/cncf/servicebroker/k8s/service_controller \
	  -v ${CURDIR}:/src/github.com/cncf/servicebroker/k8s/service_controller \
	  -e GOPATH=/ \
	  golang \
	  sh -c "go get ./... && go build ."
	docker build -t $(IMAGE):latest .

test-unit: binary
	docker run \
	  -w /src/github.com/cncf/servicebroker/k8s/service_controller \
	  -v ${CURDIR}:/src/github.com/cncf/servicebroker/k8s/service_controller \
	  -e GOPATH=/ \
	  golang \
	  sh -c "go get ./... && make --no-print-directory test-unit-local"
	
test-unit-local:
	@TEST_DIRS=$(shell find . -name \*_test.go -exec dirname {} \; | sort -u)
	@echo "\n### Running the unit tests"
	@go test $(TEST_DIRS)

clean:
	rm -f $(EXE)
	docker rmi -f $(IMAGE) || echo gone

server: binary
	docker run -ti -p 10000:10000 $(IMAGE)

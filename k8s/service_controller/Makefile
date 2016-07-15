# .PHONY: all binary test clean server

EXE:=service_controller
IMAGE:=service-controller
GO_FILES=$(shell find . -name \*.go -print)

all: binary test

binary: $(EXE)

$(EXE): $(GO_FILES)
	docker run \
	  -w /src/github.com/cncf/servicebroker/k8s/service_controller \
	  -v ${CURDIR}:/src/github.com/cncf/servicebroker/k8s/service_controller \
	  -e GOPATH=/ \
	  golang \
	  sh -c "go get ./... && go build ."
	docker build -t $(IMAGE):latest .

test: binary

clean:
	rm -f $(EXE)
	docker rmi -f $(IMAGE) || echo gone

server: binary
	docker run -ti -p 10000:10000 $(IMAGE)

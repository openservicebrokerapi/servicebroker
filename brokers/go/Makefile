EXE=gobroker

all: $(EXE)

$(EXE): gobroker.go
ifeq (,$(shell which go))
	@$(MAKE) --no-print-directory binary-docker
else
	@$(MAKE) --no-print-directory binary-local
endif
	docker build -t gobroker .

binary-docker:
	docker run -ti \
	  -w /src/github.com/servicebroker/servicebroker/k8s/service_controller/brokers/go \
	  -v $(shell cd ../.. && pwd):/src/github.com/servicebroker/servicebroker/k8s/service_controller \
	  -e GOPATH=/ \
	  golang \
	  sh -c "go get -d ./... && make binary-local"

binary-local:
	go get -d ./...
	go build -o $(EXE) -tags netgo -installsuffix netgo gobroker.go

broker: $(EXE)
	docker run --net=host -ti gobroker

clean:
	-rm -f $(EXE)
	-docker rmi -f gobroker 2> /dev/null

test:
	echo nothing yet

test-unit:
	echo nothing yet

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
	docker run \
	  -w /tmp \
	  -v ${CURDIR}:/tmp \
	  -e GOPATH=/ \
	  golang \
	  sh -c "go get -d ./... && make binary-local"

binary-local:
	go build -o $(EXE) -tags netgo -installsuffix netgo gobroker.go

clean:
	-rm -f $(EXE)
	-docker rmi -f gobroker 2> /dev/null

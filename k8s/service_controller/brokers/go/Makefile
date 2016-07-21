EXE=gobroker

all: $(EXE)

$(EXE): gobroker.go
	go build -o $(EXE) -tags netgo -installsuffix netgo gobroker.go
	docker build -t gobroker .

clean:
	rm -f $(EXE)
	docker rmi -f gobroker || echo gone

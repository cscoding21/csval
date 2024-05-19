cov:
	go vet ./... && \
	go test -coverprofile=c.out ./... && \
	go tool cover -html=c.out;
	
qual:
	go vet . && \
	golint .;

install:
	go install ./cmd/runner/csval.go
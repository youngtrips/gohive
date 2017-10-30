build: fmt
	go version > TESTED_GO_RELEASE
	go build -x

install:
	go install

fmt:
	gofmt -s=true -w */*.go *.go

test:
	go test -v

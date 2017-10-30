## env
CGO_ENABLED=0
GOARCH=amd64
GOOS=linux
GO=go

##
BASE_PATH=$(shell pwd)
GOPATH=$(BASE_PATH)/.build
PROJ_NAME = $(shell pwd |sed 's/^\(.*\)[/]//' )
APPS=$(shell ls cmd)

all: build

g:
	./gen_proto.sh

build:
	@if [ ! -L "$(GOPATH)/src/$(PROJ_NAME)" ]; then \
		mkdir -p $(GOPATH)/src ; \
	   	ln -s $(BASE_PATH) $(GOPATH)/src/$(PROJ_NAME) ; \
	fi
	@for APP in $(APPS) ; do \
		echo building $$APP ; \
		CGO_ENABLED=$(CGO_ENABLED) GOPATH=$(GOPATH) $(GO) install $(PROJ_NAME)/cmd/$$APP ; \
		mkdir -p bin ; \
		cp -f $(GOPATH)/bin/* bin/ ; \
	done

clean:
	rm -rf $(GOPATH)/pkg
	rm -rf libs/pb/*.pb.go

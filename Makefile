PLUGINS =$(shell cd plugins && find . -mindepth 1 -maxdepth 1 -type d | grep -v core | grep -v shared)

all: deps test

deps:
	go get -d -v ./...

test: deps
	go test ./...

lint: deps
	go get -u github.com/golang/lint/golint
	golint ./...

vet: deps
	go get golang.org/x/tools/cmd/vet
	go vet ./...

plugins: deps
	cd "plugins" && go install -v $(PLUGINS)

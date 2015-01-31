PLUGINS =$(shell cd plugins && find . -mindepth 1 -maxdepth 1 -type d | grep -v core | grep -v shared)

all: deps test

deps:
	go get -d -v ./...

test: deps
	go test ./...

lint: deps
	golint ./...

vet: deps
	go vet ./...

plugins: deps
	cd "plugins" && go install -v $(PLUGINS)

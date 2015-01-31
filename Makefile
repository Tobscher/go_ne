PLUGINS =$(shell cd plugins && find . -mindepth 1 -maxdepth 1 -type d | grep -v core | grep -v shared)

all: deps test

deps:
	go get -d -v ./...

test: deps
	go test ./...

lint: deps
	command -v golint >/dev/null 2>&1 || { go get -u github.com/golang/lint/golint; }
	
	golint ./...

vet: deps
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	go get golang.org/x/tools/cmd/vet
	go vet ./...

plugins: deps
	cd "plugins" && go install -v $(PLUGINS)

PLUGINS = $(shell cd plugins && find . -mindepth 1 -maxdepth 1 -type d | grep -v core | grep -v shared)
OS = "linux darwin"
ARCH = "386 amd64"

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

kiss-xc: deps
	goxc -os=$(OS) -arch=$(ARCH) -main-dirs-exclude=Godeps,agent,plugins,downloader

kiss-bintray: kiss-xc
	goxc bintray

kiss-bump-patch:
	goxc bump

kiss-bump-minor:
	goxc bump -dot=1

kiss-bump-major:
	&& goxc bump -dot=0

kiss-release-patch: kiss-bump-patch kiss-bintray

kiss-release-minor: kiss-bump-minor kiss-bintray

kiss-release-major: kiss-bump-major kiss-bintray

agent-xc: deps
	cd "agent" && goxc -os=$(OS) -arch=$(ARCH)

agent-bintray: agent-xc
	cd "agent" && goxc bintray

agent-bump-patch:
	cd "agent" && goxc bump

agent-bump-minor:
	cd "agent" && goxc bump -dot=1

agent-bump-major:
	cd "agent" && goxc bump -dot=0

agent-release-patch: agent-bump-patch agent-bintray

agent-release-minor: agent-bump-minor agent-bintray

agent-release-major: agent-bump-major agent-bintray

plugins-xc: deps
	cd "plugins" && $(foreach plugin,$(PLUGINS),(cd $(plugin) && goxc -os=$(OS) -arch=$(ARCH));)

plugins-bintray: plugins-xc
	cd "plugins" && $(foreach plugin,$(PLUGINS),(cd $(plugin) && goxc bintray);)

plugins-bump-patch:
	cd "plugins" && $(foreach plugin,$(PLUGINS),(cd $(plugin) && goxc bump);)

plugins-bump-minor:
	cd "plugins" && $(foreach plugin,$(PLUGINS),(cd $(plugin) && goxc bump -dot=1);)

plugins-bump-major:
	cd "plugins" && $(foreach plugin,$(PLUGINS),(cd $(plugin) && goxc bump -dot=0);)

plugins-release-patch: plugins-bump-patch plugins-bintray

plugins-release-minor: plugins-bump-minor plugins-bintray

plugins-release-major: plugins-bump-major plugins-bintray

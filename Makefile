PACKAGES ?= $(shell go list ./...|grep -v vendor)
VERSION ?= $(shell git describe)
GITSHA ?= $(shell git rev-parse HEAD)
LDFLAGS ?= -w -s -X main.Version=$(VERSION) -X main.GitSHA=$(GITSHA)
LINUX_PACKAGE ?= mesos-cli-linux-amd64
DARWIN_PACKAGE ?= mesos-cli-darwin-amd64
DOCKER_IMAGE := mesanine/mesos-cli

.PHONY: test release docker deploy

all: test

test:
	go $@ -v $(PACKAGES)
	go vet $(PACKAGES)
	cd filter && go test -test.bench Messages*

release:
	if ! [ -d ./release ]; then mkdir ./release ; fi
	@GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./release/$(LINUX_PACKAGE)-$(VERSION)
	if [ -h ./release/$(LINUX_PACKAGE) ]; then rm -v ./release/$(LINUX_PACKAGE); fi
	cd ./release && ln -sv $(LINUX_PACKAGE)-$(VERSION) $(LINUX_PACKAGE)
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./release/$(DARWIN_PACKAGE)-$(VERSION)
	if [ -h ./release/$(DARWIN_PACKAGE) ]; then rm -v ./release/$(DARWIN_PACKAGE); fi
	cd ./release && ln -sv $(DARWIN_PACKAGE)-$(VERSION) $(DARWIN_PACKAGE)

docker:
	$@ build -t $(DOCKER_IMAGE):$(VERSION) .
	$@ tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

deploy:
	docker login -u $$DOCKER_LOGIN -p $$DOCKER_PASSWORD
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

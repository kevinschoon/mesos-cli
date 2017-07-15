PACKAGES ?= $(shell go list ./...|grep -v vendor)
.PHONY: all
all: test

.PHONY: test
test:
	go $@ -v $(PACKAGES)
	go vet $(PACKAGES)


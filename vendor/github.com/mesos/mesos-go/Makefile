PROTO_PATH := ${GOPATH}/src/:./vendor/:.
PROTO_PATH := ${PROTO_PATH}:./vendor/github.com/gogo/protobuf/protobuf
PROTO_PATH := ${PROTO_PATH}:./vendor/github.com/gogo/protobuf/gogoproto

PACKAGES ?= $(shell go list ./...|grep -v vendor|grep -v example-)
BINARIES ?= $(shell go list -f "{{.Name}} {{.ImportPath}}" ./cmd/...|grep -v -e vendor|grep -e ^main|cut -f2 -d' ')
TEST_FLAGS ?= -v

.PHONY: all
all: test

.PHONY: install
install:
	go install $(BINARIES)

.PHONY: test
test:
	go $@ $(TEST_FLAGS) $(PACKAGES)

.PHONY: vet
vet:
	go $@ $(PACKAGES)

.PHONY: codecs
codecs: protobufs ffjson

.PHONY: protobufs
protobufs: clean-protobufs
	protoc --proto_path="${PROTO_PATH}" --gogo_out=. *.proto
	protoc --proto_path="${PROTO_PATH}" --gogo_out=. ./agent/*.proto
	protoc --proto_path="${PROTO_PATH}" --gogo_out=. ./allocator/*.proto
	protoc --proto_path="${PROTO_PATH}" --gogo_out=. ./executor/*.proto
	protoc --proto_path="${PROTO_PATH}" --gogo_out=. ./maintenance/*.proto
	protoc --proto_path="${PROTO_PATH}" --gogo_out=. ./master/*.proto
	protoc --proto_path="${PROTO_PATH}" --gogo_out=. ./scheduler/*.proto
	protoc --proto_path="${PROTO_PATH}" --gogo_out=. ./quota/*.proto

.PHONY: clean-protobufs
clean-protobufs:
	-rm *.pb.go **/*.pb.go

.PHONY: ffjson
ffjson: clean-ffjson
	ffjson *.pb.go
	ffjson agent/*.pb.go
	ffjson allocator/*.pb.go
	ffjson executor/*.pb.go
	ffjson maintenance/*.pb.go
	ffjson master/*.pb.go
	ffjson scheduler/*.pb.go
	ffjson quota/*.pb.go

.PHONY: clean-ffjson
clean-ffjson:
	rm -f *ffjson.go **/*ffjson.go

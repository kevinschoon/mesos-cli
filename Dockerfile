FROM quay.io/vektorcloud/go:dep AS build

COPY . /go/src/github.com/vektorlab/mesos-cli

RUN cd /go/src/github.com/vektorlab/mesos-cli \
  && go build \
  && chmod +x mesos-cli

FROM quay.io/vektorcloud/base:3.6

COPY --from=build /go/src/github.com/vektorlab/mesos-cli/mesos-cli /bin/mesos

ENTRYPOINT ["/bin/mesos"]

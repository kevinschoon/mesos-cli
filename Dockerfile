FROM quay.io/vektorcloud/go:1.9 AS build

COPY . /go/src/github.com/mesanine/mesos-cli

RUN cd /go/src/github.com/mesanine/mesos-cli \
  && go build \
  && chmod +x mesos-cli

FROM quay.io/vektorcloud/base:3.6

COPY --from=build /go/src/github.com/mesanine/mesos-cli/mesos-cli /bin/mesos

ENTRYPOINT ["/bin/mesos"]

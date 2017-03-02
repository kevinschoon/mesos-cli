FROM quay.io/vektorcloud/base:3.4

COPY release/mesos-cli-alpine /bin/mesos

CMD mesos

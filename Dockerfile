FROM quay.io/vektorcloud/base:3.4

COPY release/mesos-cli-linux-amd64 /bin/mesos

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

CMD mesos

# mesos-exec

Execute arbitrary commands against an [Apache Mesos]("http://mesos.apache.com") cluster.


**mesos-exec** implements the functionality of the [mesos-execute](https://github.com/apache/mesos/blob/master/src/cli/execute.cpp)
with some additional features. It communicates with Mesos via HTTP wire protocol so it does not require the cumbersome libmesos packages, 
nor language specific bindings to the Mesos C library.

### Usage

    $ mesos-exec --help

    Usage: mesos-exec [OPTIONS] [ARG...]

    Execute Commands on Apache Mesos

    Arguments:
      ARG=[]       Command Arguments

    Options:
      -v, --version               Show the version and exit
      -i, --image=""              Docker image to run
      --master="127.0.0.1:5050"   Master address <host:port>
      -c, --cpus="0.1"            CPU Resources to allocate
      -m, --mem="128.0"           Memory resources (mb) to allocate
      -d, --disk="64.0"           Memory resources (mb) to allocate
      --level="0"                 Logging level
      -n, --name="mesos-exec"     Task Name
      -s, --shell=false           Execute as shell command
      -u, --user="root"           User to run as

    # In native mesos executor
    $ mesos-exec --shell 'for i in $(seq 1 5); do echo $(date); sleep 1; done'
    ....
    Wed Dec 14 23:16:49 UTC 2016
    Wed Dec 14 23:16:50 UTC 2016
    Wed Dec 14 23:16:51 UTC 2016
    Wed Dec 14 23:16:52 UTC 2016
    Wed Dec 14 23:16:53 UTC 2016
    ....
    # Or Docker
    $ mesos-exec --image ubuntu:latest --shell 'for i in $(seq 1 5); do echo $(date); sleep 1; done'



### Installation

#### From Source

    go get -u github.com/vektorlab/mesos-exec


#### TODO

  * Support full TaskInfo object
  * Support multiple TaskInfos array
  * Support JSON config file with "profiles"
  * Improve logging output

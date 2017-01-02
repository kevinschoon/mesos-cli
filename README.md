# mesos-exec

Execute arbitrary commands against an [Apache Mesos]("http://mesos.apache.com") cluster.


**mesos-exec** implements the functionality of the [mesos-execute](https://github.com/apache/mesos/blob/master/src/cli/execute.cpp)
with some additional features. It communicates with Mesos via HTTP wire protocol so it does not require the cumbersome libmesos packages, 
nor language specific bindings to the Mesos C library.

### Usage

    Usage: mesos-exec [OPTIONS] [ARG...]

    Execute Commands on Apache Mesos

    Arguments:
      ARG=[]       Command Arguments

    Options:
      --profile="default"          Profile to load from ~/.mesos-exec.json
      --master="127.0.0.1:5050"    Master address <host:port>
      --task=""                    Path to a Mesos TaskInfo JSON file
      --param=[]                   Docker parameters
      -i, --image=""               Docker image to run
      -l, --level=0                Level of verbosity
      -v, --volume=[]              Volume mappings
      -p, --ports=[]               Port mappings
      -e, --env=[]                 Environment Variables
      -s, --shell=""               Shell command to execute
      -n, --name=mesos-exec        Task Name
      -u, --user=root              User to run as
      -c, --cpus=0.1               CPU Resources to allocate
      -m, --mem=128.0              Memory Resources (mb) to allocate
      -d, --disk=32.0              Disk Resources (mb) to allocate
      --privileged=false           Give extended privileges to this container
      -f, --forcePullImage=false   Always pull the container image


    # With Docker containerizer
    $ mesos-exec --image alpine:latest --shell 'for i in $(seq 1 5); do echo $(date); sleep 1; done'
    ....
    Wed Dec 14 23:16:49 UTC 2016
    Wed Dec 14 23:16:50 UTC 2016
    Wed Dec 14 23:16:51 UTC 2016
    Wed Dec 14 23:16:52 UTC 2016
    Wed Dec 14 23:16:53 UTC 2016
    ....
    # Or with native Mesos containerizer
    # Since native mesos containerizer doesn't redirect stdout/stderr by default you 
    # need to literally write to a file called `stdout`/`stderr` in the sandbox directory.
    $ mesos-exec --shell 'for i in $(seq 1 5); do echo $(date) >> stdout; sleep 1; done'
    
### Profiles
You can configure "profiles" by creating a JSON file at `~/.mesos-exec.json`.

    $ cat ~/.mesos-exec.json
    {
      "profiles": {
        "default": {
          "master": "localhost:5050"
        },
        "production": {
          "master": "production-host:5050"
        },
        "development": {
          "master": "development-host:5050"
        }
      }
    }





### Installation

#### From Source

    go get -u github.com/vektorlab/mesos-exec


#### TODO

  * Support full TaskInfo object
  * Support multiple TaskInfos array
  * Improve logging output

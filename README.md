# mesos-cli

Alternative command line interface for using an [Apache Mesos]("http://mesos.apache.com") cluster.

The existing CLI tools [bundled]("https://github.com/apache/mesos/tree/master/src/cli") with Mesos are a combintation
of C++ and Python scripts. The previously developed [mesos-cli]("https://github.com/mesosphere/mesos-cli") by Mesosphere has been integrated into their DC/OS platform. Both of those tools depend on the cubmersome libmesos packages which are difficult to install. `mesos-cli` is written in pure Go and communicates with Mesos via HTTP with the excellent [mesos-go]("https://github.com/mesos/mesos-go") library. `mesos-cli` additionally aims to add more convenient features than the original toolset.

`mesos-cli` is under development and not ready for use in a production environment.

### Installation 

Binary packages are not yet available so you need to install from source.

    go get -u github.com/vektorlab/mesos-cli
    
 If you don't mind potentially overriding the default `mesos` command you may add an alias:
 
     echo "alias mesos=mesos-cli" >> $HOME/.bashrc
    
### Usage

    mesos --help
    Usage: mesos-cli [OPTIONS] COMMAND [arg...]

    Alternative Apache Mesos CLI

    Options:
      --master="127.0.0.1:5050"   Master address <host:port>
      --profile="default"         Profile to load from ~/.mesos-cli.json
      --level=0                   Level of verbosity

    Commands:
      exec         Execute Arbitrary Commands Against a Cluster

    Run 'mesos-cli COMMAND --help' for more information on a command.
    
#### exec
`mesos exec` implements the functionality of the [mesos-execute](https://github.com/apache/mesos/blob/master/src/cli/execute.cpp)
with some additional features.

    Usage: mesos-cli exec [OPTIONS] [ARG...]

    Execute Arbitrary Commands Against a Cluster

    Arguments:
      ARG=[]       Command Arguments

    Options:
      --task=""                    Path to a Mesos TaskInfo JSON file
      --param=[]                   Docker parameters
      -i, --image=""               Docker image to run
      -v, --volume=[]              Volume mappings
      -p, --ports=[]               Port mappings
      -e, --env=[]                 Environment Variables
      -s, --shell=""               Shell command to execute
      -n, --name=mesos-cli        Task Name
      -u, --user=root              User to run as
      -c, --cpus=0.1               CPU Resources to allocate
      -m, --mem=128.0              Memory Resources (mb) to allocate
      -d, --disk=32.0              Disk Resources (mb) to allocate
      --privileged=false           Give extended privileges to this container
      -f, --forcePullImage=false   Always pull the container image

    # With Docker containerizer
    $ mesos exec --image alpine:latest --shell 'for i in $(seq 1 5); do echo $(date); sleep 1; done'
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
    $ mesos exec --shell 'for i in $(seq 1 5); do echo $(date) >> stdout; sleep 1; done'
    
### Profiles
You can configure "profiles" by creating a JSON file at `~/.mesos-cli.json`.

    $ cat ~/.mesos-cli.json
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

#### TODO

  * Support full TaskInfo object
  * Support multiple TaskInfos array
  * Improve logging output
  * mesos ps
  * mesos cat
  * mesos top
  * mesos tail
  * mesos head

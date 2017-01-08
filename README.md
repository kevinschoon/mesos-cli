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
    
### Usage

    Usage: mesos-cli [OPTIONS] COMMAND [arg...]

    Alternative Apache Mesos CLI

    Options:
      --profile="default"                      Profile to load
      --config="/home/kevin/.mesos-cli.json"   Path to load config from
      --level=0                                Level of verbosity

    Commands:
      cat          Output the contents of a file
      local        Launch a local Mesos cluster (requires Docker)
      ls           List the sandbox directory of a task
      ps           List currently running tasks on a cluster
      run          Run arbitrary commands against a cluster

Run 'mesos-cli COMMAND --help' for more information on a command.
    
#### run
`mesos run` implements the functionality of the [mesos-execute](https://github.com/apache/mesos/blob/master/src/cli/execute.cpp)
with some additional features.

    Usage: mesos-cli run [OPTIONS] [ARG...]

    Run arbitrary commands against a cluster

    Arguments:
      ARG=[]       Command Arguments

    Options:
      --master="127.0.0.1:5050"    Mesos Master
      --task=""                    Path to a Mesos TaskInfo JSON file
      --param=[]                   Docker parameters
      -i, --image=""               Docker image to run
      -v, --volume=[]              Volume mappings
      -p, --ports=[]               Port mappings
      -e, --env=[]                 Environment Variables
      -s, --shell=""               Shell command to execute
      -t, --tail=false             Tail command output
      -n, --name=mesos-cli         Task Name
      -u, --user=root              User to run as
      -c, --cpus=0.1               CPU Resources to allocate
      -m, --mem=128.0              Memory Resources (mb) to allocate
      -d, --disk=32.0              Disk Resources (mb) to allocate
      --privileged=false           Give extended privileges to this container
      -f, --forcePullImage=false   Always pull the container image


    # With Docker containerizer
    $ mesos run --tail --image alpine:latest --shell 'for i in $(seq 1 5); do echo $(date); sleep 1; done'
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
    
#### ps

    Usage: mesos-cli ps [OPTIONS]

    List currently running tasks on a cluster

    Arguments:
      order="desc"   accending or decending sort order [asc|desc]

    Options:
      --master="127.0.0.1:5050"   Mesos Master
      --limit=100                 maximum number of tasks to return per request
      --max=250                   maximum number of tasks to list
      --name=""                   regular expression to match the TaskId
      -a, --all=false             show all tasks
      -r, --running=true          show running tasks
      --fa, --failed=false        show failed tasks
      -k, --killed=false          show killed tasks
      -f, --finished=false        show finished tasks
    
    mesos ps
    
    ID                                        	FRAMEWORK	STATE       	CPUS	MEM	GPUS	DISK
    nginx.d6592dd7-d52a-11e6-bb61-6e9c129136b0	c654f3d1 	TASK_RUNNING	0.1 	64 	0   	0  
      
#### ls

    Usage: mesos-cli ls [OPTIONS] ID

    List the sandbox directory of a task

    Arguments:
      ID=""        Task to list

    Options:
      --master="127.0.0.1:5050"   Mesos Master
      -a, --absolute=false        Show absolute file paths
      
    mesos ls nginx.d6592dd7-d52a-11e6-bb61-6e9c129136b0
    
    UID 	GID 	MODE      	MODIFIED                     	SIZE  	PATH
    root	root	-rw-r--r--	2017-01-07 22:35:46 -0500 EST	1527  	stderr
    root	root	-rw-r--r--	2017-01-08 17:39:03 -0500 EST	642717	stdout

      
#### cat

    Usage: mesos-cli cat [OPTIONS] ID FILE

    Output the contents of a file

    Arguments:
      ID=""        TaskID
      FILE=""      Filename to retrieve

    Options:
      --master="127.0.0.1:5050"   Mesos Master
      -n, --lines=0               Output the last N lines
      -t, --tail=false            Tail output
      
      mesos cat nginx.d6592dd7-d52a-11e6-bb61-6e9c129136b0 stdout
      172.17.0.1 - - [08/Jan/2017:03:35:46 +0000] "GET / HTTP/1.1" 200 612 "http://localhost:8080/ui/" "Mozilla/5.0 (X11;     
      ...

#### local

`mesos local` provides a wrapper for launching a local Mesos cluster for development and testing purposes. It requires that you have Docker installed locally uses the [vektorcloud/mesos]("https://github.com/vektorcloud/mesos") image.

    Usage: mesos-cli local [OPTIONS] COMMAND [arg...]

    Launch a local Mesos cluster (requires Docker)

    Commands:
      up           Start the local cluster
      down         Stop the local cluster
      status       Display the status of the local cluster
      rm           Remove the local cluster

    Run 'mesos-cli local COMMAND --help' for more information on a command.


#### TODO

  * Support multiple TaskInfos array
  * Improve logging output
  * mesos top

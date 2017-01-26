# mesos-cli

Standalone commandline tool for interacting with an [Apache Mesos]("http://mesos.apache.com") cluster.

# Why?

Existing CLI tools for Mesos are tightly integrated into their parent projects(e.g. [[0]](https://github.com/apache/mesos/tree/master/src/cli), [[1]](https://github.com/mesosphere/mesos-cli)) and dependent on cumbersome libmesos packages.

`mesos-cli` is a lightweight alternative to these tools, leveraging the excellent [mesos-go]("https://github.com/mesos/mesos-go") library to communicate with Mesos via HTTP. `mesos-cli` additionally aims to add more convenient features than the original toolset.

`mesos-cli` is under development and not ready for use in a production environment.

# Installation

You can download binary packages for your platform (linux/darwin) from the [releases](https://github.com/vektorlab/mesos-cli/releases) section. 

    PLATFORM=linux
    wget https://github.com/vektorlab/mesos-cli/releases/download/v0.0.3/mesos-cli-v0.0.3-$PLATFORM-amd64 -o /usr/local/bin/mesos-cli

 If you don't mind potentially overriding the default `mesos` command you may add an alias:

     echo "alias mesos=mesos-cli" >> $HOME/.bashrc

# Profiles
You can configure "profiles" by creating a JSON file at `~/.mesos-cli.json`:

```json
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
```

# Usage

`mesos-cli` currently supports the following subcommands:

## run
`mesos run` implements the functionality of the [mesos-execute](https://github.com/apache/mesos/blob/master/src/cli/execute.cpp)
with some additional features. You can also specify a file containing a JSON encoded Mesos 
[TaskInfo](https://github.com/mesos/mesos-go/blob/master/mesosproto/mesos.proto#L1038-L1072) object with the `--task` flag.

```
mesos-cli run [OPTIONS] [ARG...]
```

### Example
With Docker containerizer:

```bash
mesos run --tail --image alpine:latest --shell 'date'
....
Wed Dec 14 23:16:50 UTC 2016
....
```

Or with native Mesos containerizer:
```bash
mesos run --shell 'echo $(date) >> stdout'
```
*note:* Since native mesos containerizer doesn't redirect stdout/stderr by default you need to literally write to a file called `stdout`/`stderr` in the sandbox directory.

### Options

Option | Description
--- | ---
--master="127.0.0.1:5050" | Mesos Master
--task="" | Path to a Mesos TaskInfo JSON file
--param=[] | Docker parameters
-i, --image="" | Docker image to run
-v, --volume=[] | Volume mappings
-p, --ports=[] | Port mappings
-e, --env=[] | Environment Variables
-s, --shell="" | Shell command to execute
-t, --tail=false | Tail command output
-n, --name=mesos-cli | Task Name
-u, --user=root | User to run as
-c, --cpus=0.1 | CPU Resources to allocate
-m, --mem=128.0 | Memory Resources (mb) to allocate
-d, --disk=32.0 | Disk Resources (mb) to allocate
--privileged=false | Give extended privileges to this container
-f, --forcePullImage=false | Always pull the container image

## ps

List currently running tasks on a cluster

```
mesos-cli ps [OPTIONS]
```

### Options

Option | Description
--- | ---
--master="127.0.0.1:5050" |  Mesos Master
--limit=2000              |  maximum number of tasks to return per request
--max=250                 |  maximum number of tasks to list
--truncate=true           |  truncate some values
--all=false               |  Show all tasks
--framework=""            |  Filter FrameworkID
--fuzzy=true              |  Fuzzy match Task name or Task ID prefix
--name=""                 |  Filter Task name
--id=""                   |  Filter Task ID
--state=["TASK_RUNNING"]  |  Filter based on Task state


### Example

```bash
mesos ps --state "TASK_FINISHED" --max=2
```

```
    ID        FRAMEWORK STATE         CPU MEM GPU DISK
    mesos-cli b620d6e2  TASK_FINISHED 0.1 128 0   32  
    mesos-cli b620d6e2  TASK_FINISHED 0.1 128 0   32  
```

## ls

List the sandbox directory of a task

```
Usage: mesos-cli ls [OPTIONS] TASKID
```

### Options

Option | Description
--- | ---
--master="127.0.0.1:5050"   | Mesos Master
-a, --absolute=false        | Show absolute file paths
--all=false                 | Show all tasks
--framework=""              | Filter FrameworkID
--fuzzy=true                | Fuzzy match Task name or Task ID prefix
--name=""                   | Filter Task name
--id=""                     | Filter Task ID
--state=[]                  | Filter based on Task state


### Example
```bash
mesos ls --id nginx.d6592dd7-d52a-11e6-bb61-6e9c129136b0
```

```
UID 	GID 	MODE      	MODIFIED                     	SIZE  	PATH
root	root	-rw-r--r--	2017-01-07 22:35:46 -0500 EST	1527  	stderr
root	root	-rw-r--r--	2017-01-08 17:39:03 -0500 EST	642717	stdout
```

## cat

Output the contents of a file

```
mesos-cli cat [OPTIONS] TASKID FILE
```

### Options

Option | Description
--- | ---
--master="127.0.0.1:5050"  | Mesos Master
-n, --lines=0              | Output the last N lines
-t, --tail=false           | Tail output
--all=false                | Show all tasks
--framework=""             | Filter FrameworkID
--fuzzy=true               | Fuzzy match Task name or Task ID prefix
--name=""                  | Filter Task name
--id=""                    | Filter Task ID
--state=[]                 | Filter based on Task state


### Example
```bash
mesos cat --id=nginx.d6592dd7-d52a-11e6-bb61-6e9c129136b0 stdout
```

```
172.17.0.1 - - [08/Jan/2017:03:35:46 +0000] "GET / HTTP/1.1" 200 612 "http://localhost:8080/ui/" "Mozilla/5.0 (X11;...
...
```

## agents
`mesos agents` lists all the agents running in the cluster

```bash
mesos agents
```

```
ID                                     	HOSTNAME             	VERSION	UPTIME                 	CPUS     	MEM         	GPUS     	DISK        
23d60c9d-dab0-4af9-8336-a7cb501ea2c1-S0	localhost.localdomain	1.1.0  	412626h27m31.603293767s	0.00/4.00	0.00/6867.00	0.00/0.00	0.00/5114.00
```

## top
`mesos top` provides a top-like overview of tasks, agent, and cluster status (work in progress)

```
mesos-cli top [OPTIONS] COMMAND [arg...]
```

### Example
```bash
mesos top
```

## local

`mesos local` provides a wrapper for launching a local Mesos cluster for development and testing purposes.
It requires that you have Docker installed locally, and uses the [vektorcloud/mesos]("https://github.com/vektorcloud/mesos") image.

```
mesos-cli local [OPTIONS] COMMAND [arg...]
```

### Commands
Command | Description
--- | ---
up | Start the local cluster
down | Stop the local cluster
status | Display the status of the local cluster
rm | Remove the local cluster

## Global Options

Option | Description
--- | ---
--profile | Profile to load
--config | Path to load config from
--level | Level of verbosity

## TODO

  * Support multiple TaskInfos array
  * Improve logging output
  * mesos top

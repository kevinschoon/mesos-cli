# mesos-cli

Standalone commandline tool for running containers on and interacting with [Apache Mesos]("http://mesos.apache.com").

`mesos-cli` interacts entirely with the new Mesos HTTP [scheduler](http://mesos.apache.org/documentation/latest/scheduler-http-api/) and [operator](http://mesos.apache.org/documentation/latest/operator-http-api/) APIs.

`mesos-cli` is under active development and not yet stable.

# Distinctive Features & Roadmap

| Feature                                                                             |ready|roadmap|
|-------------------------------------------------------------------------------------|-----|-------|
| Built ontop of the new Mesos HTTP V1 API                                            |✓    |       |
| Simple installation without platform specific libmesos drivers                      |✓    |       |
| Full support for latest Mesos features, e.g health checks, etc                      |     |✓      |
| Streamimg sandbox file content content to console (including task stdout and stderr)|✓    |       |
| Support for running Mesos TaskGroups (pods)                                         |     |✓      |
| Ability to search and filter across most Mesos types                                |     |✓      |
| Pure integration with Apache Mesos outside of the DC/OS ecosystem                   |✓    |       |
| Top-like interface for monitoring a cluster                                         |     |✓      |
| Subscribe to and monitor master event stream                                        |     |✓      |
| Lightweight scheduler for running arbitrary containers                              |✓    |       |
| Support for running multiple containers for batch-style jobs as a DAG               |     |✓      |
| Export to task to Kubernetes or Marathon                                            |     |✓      |
| Simple interface for lauguage agnostic executors                                    |     |✓      |
| Effortlessly run a local Mesos cluster                                              |✓    |       |

# Installation

## Binaries

You can download binary packages for your platform (linux/darwin) from the [releases](https://github.com/vektorlab/mesos-cli/releases) section. 

    PLATFORM=linux
    wget https://github.com/vektorlab/mesos-cli/releases/download/v0.0.5/mesos-cli-v0.0.5-$PLATFORM-amd64 -o /usr/local/bin/mesos-cli

 If you don't mind potentially overriding the default `mesos` command you may add an alias:

     echo "alias mesos=mesos-cli" >> $HOME/.bashrc
     
 ## Docker
 
    docker pull quay.io/vektorcloud/mesos-cli
    echo "alias mesos=docker run --rm -ti -v $HOME/.meoss-cli.json:/root/.mesos-cli.json --net host quay.io/vektorcloud/mesos-cli >> $HOME/.bashrc"

# Profiles
You can configure "profiles" by creating a JSON file at `~/.mesos-cli.json`. This will be created the first time you invoke command.

Example:

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
`mesos run` launches a lightweight scheduler for running a container on the target cluster.

```
mesos-cli run [OPTIONS] [ARG...]
Usage: mesos-cli run [OPTIONS] [CMD]
```
Option| Description
------------------- |--------------------
  --user="root"     |       User to run as
  --shell=false     |       Run as a shell command
  --master=""       |       Mesos master
  --path=""         |       Path to a JSON file containing a Mesos TaskInfo
  --json=false      |       Dump the task to JSON instead of running it
  --docker=false    |      Run as a Docker container
  --image=""        |      Image to run
  --restart=false   |     Restart container on failure
  --privileged=false|   Run in privileged mode [docker only]
  -e, --env=        |  Environment variables
  -v, --volume=     |   Container volume mappings
  --net=BRIDGE      |   Network Mode [Docker only]
  --param=          |   Freeform Docker parameters [Docker only]
  -p, --port=       |   Port mappings [Docker only]


### Examples

Launch a new task with the Mesos containerizer and restart it on exit

```bash
    mesos run --restart --shell 'echo $(date); sleep 2'
 ```
 
Run an app with Docker and keep it online

```bash
    mesos run --restart --docker -p 31000:80 --image nginx:latest 
    mesos tasks --state TASK_RUNNING # Check it's state
    curl localhost:31000
```


## Tasks

List currently running tasks on a cluster

```
mesos-cli tasks [OPTIONS]
```

### Options

Option | Description
--- | ---
--master=""    |    Mesos master
--truncate=true|   Truncate long values
--task=""      |  Filter by task id
--fuzzy=true   | Fuzzy matching on string values
--state=[]     |filter by task state



### Example

```bash
mesos ps --state "TASK_FINISHED"
```

```
    ID        FRAMEWORK STATE         CPU MEM GPU DISK
    mesos-cli b620d6e2  TASK_FINISHED 0.1 128 0   32  
    mesos-cli b620d6e2  TASK_FINISHED 0.1 128 0   32  
```

## List

List the sandbox directory of a task

```
Usage: mesos-cli list [OPTIONS] ID PATH
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
mesos list --master http://localhost:5050 5b7a7d41-1f2a-4b6d-93fd-48354d7fa785-S0  /opt/mesos/0/slaves/5b7a7d41-1f2a-4b6d-93fd-48354d7fa785-S0/frameworks/5b7a7d41-1f2a-4b6d-93fd-48354d7fa785-0115/executors/807fe9fa-55df-40fa-ab4f-20e359d51d43/runs/4295fd85-5023-4b0d-a9c8-1d9bdce309ba

UID 	GID 	MODE      	MODIFIED                     	SIZE  	PATH
root	root	-rw-r--r--	2017-01-07 22:35:46 -0500 EST	1527  	stderr
root	root	-rw-r--r--	2017-01-08 17:39:03 -0500 EST	642717	stdout
```

## Read

Output the contents of a file

```
mesos-cli read [OPTIONS] TASKID FILE
```

### Options

Option | Description
--- | ---
-f, --follow=false  | follow the content
-n, --nlines=0      | number of lines to read
-m, --master=""     | mesos master



### Example
```bash
mesos read 5b7a7d41-1f2a-4b6d-93fd-48354d7fa785-S0  /opt/mesos/0/slaves/5b7a7d41-1f2a-4b6d-93fd-48354d7fa785-S0/frameworks/5b7a7d41-1f2a-4b6d-93fd-48354d7fa785-0115/executors/807fe9fa-55df-40fa-ab4f-20e359d51d43/runs/4295fd85-5023-4b0d-a9c8-1d9bdce309ba/stdout

172.17.0.1 - - [08/Jan/2017:03:35:46 +0000] "GET / HTTP/1.1" 200 612 "http://localhost:8080/ui/" "Mozilla/5.0 (X11;...
...
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

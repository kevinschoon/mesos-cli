+++
title = "run"
next = "/usage/tasks"
prev = "/usage/read"
+++

Launch a lightweight scheduler for running containers on the target cluster.

```
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



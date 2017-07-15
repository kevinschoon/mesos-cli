+++
title = "task"
next = "/usage/tasks"
prev = "/usage/run"
+++

Task generates a valid mesos TaskInfo [protobuf object](https://github.com/apache/mesos/blob/master/include/mesos/v1/mesos.proto#L1551) and encodes it in JSON (or YAML). It (soon!) will support flags for the entire gamut of options available to a Mesos scheduler. It is intended that this command be convenient for use with the [run](/mesos-cli/usage/run) command.

```
Usage: mesos-cli task [OPTIONS] [CMD]
```

### Arguments
Argument        |   Description
----------------|------------------------------------------------
CMD=""          |   Command to run

### Options

Option              | Description
--------------------|---------------------------------------
  --name=mesos-cli  |   Friendly task name
  --user="root"     |   User to run as
  --shell=false     |   Run as a shell command
  --uri=            |   URIs to fetch
  -e, --env=        |   Environment variables
  -v, --volume=     |   Container volumes
  -i, --image=      |   Image to run
  --cpu=0.000000    |   CPU resources for this task
  --gpu=0.000000    |   GPU resources for this task
  --memory=0.000000 |   Memory resources for this task
  --disk=0.000000   |   Disk resources for this task
  --privileged=false|   Run Docker in privileged mode
  -p, --port=       |   Port mappings [Docker only]
  --param=          |   Docker parameters [Docker only]
  --net=BRIDGE      |   Network Mode [Docker only]
  --encoding="json" |   Output encoding [json/yaml]
  --docker=false    |   Run as a Docker container
  --role=""         |   Mesos role



### Examples

Launch a new task with the Mesos containerizer.

```bash
mesos task --shell 'echo $(date); sleep 2' | mesos run -
```

Generate a new task and save it to a file.

```bash
mesos task --encoding=yaml --shell 'echo $(date); sleep 2' > my-task.yaml
```

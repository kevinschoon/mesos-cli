+++
title = "run"
next = "/usage/task"
prev = "/usage/read"
+++

Launch a lightweight scheduler for running containers on the target cluster. Run accepts a JSON or YAML encoded document. 
By default run will attempt to load `Mesosfile` from the current directory.

```
Usage: mesos-cli run [OPTIONS] FILE
```
### Arguments
Argument        |   Description
----------------|------------------------------------------------
FILE="Mesosfile"|   File containing Mesos TaskInfos, - for stdin

### Options

Option | Description
-------|---------------------------------------
-m, --master=""  | Mesos Master
-s, --sync=false | Run containers synchronously
--restart=false  | Restart containers on failure


### Examples

Launch a new task with the Mesos containerizer and restart it on exit

```bash
mesos task --shell 'echo $(date); sleep 2' | mesos run --restart -
```
 
Run an app with Docker and keep it online

```bash
    mesos task --docker -p 31000:80 --image nginx:latest  | mesos run -
    curl localhost:31000
```



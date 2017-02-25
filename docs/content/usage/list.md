+++
draft = false
date = "2017-02-25T15:35:21+08:00"
title = "list"

+++

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



+++
date = "2017-02-25T15:35:05+08:00"
title = "local"
draft = false
+++

Launch a local Mesos cluster for development and testing purposes. This command requires that you have Docker installed locally, and uses the [vektorcloud/mesos]("https://github.com/vektorcloud/mesos") image.

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

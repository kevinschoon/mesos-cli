+++
title = "local"
next = "/usage/read"
prev = "/usage/list"
+++

Launch a local Mesos cluster for development and testing purposes. This command requires that you have Docker installed locally, and uses the [vektorcloud/mesos](https://github.com/vektorcloud/mesos) image.

```
mesos-cli local [OPTIONS] COMMAND [arg...]
```
## Options

Option    | Description
----------|------------------------
--profile | Profile to load
--config  | Path to load config from
--level   | Level of verbosity

### Subcommands
Command | Description
--------|---------------------------------------
up      | Start the local cluster
down    | Stop the local cluster
status  | Display the status of the local cluster
rm      | Remove the local cluster



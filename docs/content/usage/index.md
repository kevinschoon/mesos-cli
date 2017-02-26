+++
title = "Usage"
chapter = true
next = "/usage/agents"
prev = "/getting-started/configuration"
+++

**mesos-cli** is broken into several subcommands, you can find detailed options for each command by passing the `--help` flag.


    Usage: mesos-cli [OPTIONS] COMMAND [arg...]

### Options

Option                                | Description
--------------------------------------| --------------------------
--profile="default"                   |  Profile to load
--config="/home/kevin/.mesos-cli.json"|  Path to load config from
--debug=false                         |  Enable debugging
--version                             |  Show the version and exit

### Subcommands
Command| Description
------| --------------
agents| List Mesos Agents
list  | List files in a Mesos sandbox
local | Run a local Mesos cluster
read  | Read the contents of a file
run   | Run tasks on Mesos
tasks | List Mesos tasks
top   | Display a Mesos top interface


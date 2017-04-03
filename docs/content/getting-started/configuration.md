+++
title = "configuration"
weight = 10
next = "/usage"
prev = "/getting-started/installation"
+++

It is recommended that you create a bash alias for mesos-cli, note that if the default `mesos` executable is installed on your system this will override it.

     # With a binary installation
     echo "alias mesos=mesos-cli" >> $HOME/.bashrc
     # With a Docker installation
     echo "alias mesos=docker run --rm -ti -v $HOME/.mesos-cli.json:/root/.mesos-cli.json --net host quay.io/vektorlab/mesos-cli >> $HOME/.bashrc"


### Test Your Installation

**mesos-cli** has built-in support for launching a local cluster for testing and development.

    # Launch a local Mesos cluster
    mesos local up
    # List agents running on your cluster
    mesos agents

    ID                                     	HOSTNAME             	CPUS	MEM    	GPUS	DISK   
    48bf0171-3c61-49ed-9e05-3ac1a9274478-S0	localhost.localdomain	4.00	6867.00	0.00	5114.0


### Profiles
You can configure "profiles" by creating a JSON file at `~/.mesos-cli.json`. This file is automatically created for you the first time you invoke **mesos-cli**. You can choose an alternative profile for use with any command. All options specified in a profile may be overriden by specifying the same option on the command-line.

Example:

```json
{
  "profiles": {
    "default": {
      "master": "http://localhost:5050",
      "debug": false,
      "restart": false
    },
    "production": {
      "master": "http://localhost:5050",
      "debug": false,
      "restart": true
    }
  }
}
```


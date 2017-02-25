+++
date = "2017-02-25T15:35:12+08:00"
title = "Getting Started"
draft = false

+++

**mesos-cli** is a command-line tool for running containers on and interacting with [Apache Mesos](http://mesos.apache.com). It is designed to be lightweight and more featureful than the [native](https://github.com/apache/mesos/tree/master/src/cli) CLI tool provided with Apache Mesos. Additionally, it seeks to act as a stand-alone and extensible task scheduler on Mesos.

By interacting entirely with the new Mesos HTTP [scheduler](http://mesos.apache.org/documentation/latest/scheduler-http-api/) API as a general purpose Mesos framework, **mesos-cli** does not require the Mesos Master server to establish a direct network connection to it, which makes it more flexible than other frameworks.

**mesos-cli is under active development and not yet considered stable!**


# Distinctive Features & Roadmap

| Feature                                                                             |ready|roadmap|
|-------------------------------------------------------------------------------------|-----|-------|
| Built ontop of the new Mesos HTTP V1 API                                            |✓    |       |
| Simple installation without platform specific libmesos drivers                      |✓    |       |
| Streamimg sandbox file content content to console (including task stdout and stderr)|✓    |       |
| Effortlessly run a local Mesos cluster                                              |✓    |       |
| Pure integration with Apache Mesos outside of the DC/OS ecosystem                   |✓    |       |
| Lightweight scheduler for running arbitrary containers                              |✓    |       |
| Ability to search and filter across most Mesos types                                |     |✓      |
| Support for running Mesos TaskGroups (pods)                                         |     |✓      |
| Ability to search and filter across most Mesos types                                |     |✓      |
| Top-like interface for monitoring a cluster                                         |     |✓      |
| Subscribe to and monitor master event stream                                        |     |✓      |
| Export to task to Kubernetes or Marathon                                            |     |✓      |
| Run support running docker-compose files directly against Mesos                     |     |✓      |
| Full support for latest Mesos features, e.g health checks, etc                      |     |✓      |


# Installation

**mesos-cli** can be installed by downloading the latest release for your architecture or it can be run from a Docker container.

## Binaries

You can download binary packages for your platform (linux/darwin) from the [releases](https://github.com/vektorlab/mesos-cli/releases) section on Github or below: 

  - [Linux](https://github.com/vektorlab/mesos-cli/releases/download/v0.0.5/mesos-cli-linux-amd64-v0.0.5)
  - [OSX/Darwin](https://github.com/vektorlab/mesos-cli/releases/download/v0.0.5/mesos-cli-darwin-amd64-v0.0.5)


     
## Docker

A Docker container is also available for download:
 
    docker pull quay.io/vektorcloud/mesos-cli
    docker run --rm -ti quay.io/vektorcloud/mesos-cli tasks --master http://your-mesos-server:5050

## Configuration

 If you don't mind potentially overriding the default `mesos` command you may add an alias:

     # With a binary installation
     echo "alias mesos=mesos-cli" >> $HOME/.bashrc
     # With a Docker installation
     echo "alias mesos=docker run --rm -ti -v $HOME/.meoss-cli.json:/root/.mesos-cli.json --net host quay.io/vektorcloud/mesos-cli >> $HOME/.bashrc"


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


+++
title = "installation"
weight = 0
next = "/getting-started/configuration"
prev = "/getting-started"
+++

**mesos-cli** can be installed by downloading the latest release for your architecture or it can be run from a Docker container.

## Binaries

You can download binary packages for your platform (linux/darwin) from the [releases](https://github.com/vektorlab/mesos-cli/releases) section on Github or below: 

  - [Linux](https://github.com/vektorlab/mesos-cli/releases/download/v0.0.7/mesos-cli-linux-amd64-v0.0.7)
  - [OSX/Darwin](https://github.com/vektorlab/mesos-cli/releases/download/v0.0.7/mesos-cli-darwin-amd64-v0.0.7)

     
## Docker

A Docker container is also available for download:
 
    docker pull quay.io/vektorlab/mesos-cli
    docker run --rm -ti quay.io/vektorlab/mesos-cli tasks --master http://your-mesos-server:5050

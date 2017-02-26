+++
title = "Getting Started"
weight = -10
chapter = true
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



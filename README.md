# mesos-cli
[![Build](https://img.shields.io/circleci/project/github/mesanine/mesos-cli.svg)](https://circleci.com/gh/mesanine/mesos-cli)

Featureful commandline interface for [Apache Mesos](http://mesos.apache.com).

**NOTE**: *mesos-cli is under active development and not yet considered stable!*

`mesos-cli` is designed to be a lightweight alternative to the [native tool](https://github.com/apache/mesos/tree/master/src/cli) provided with Mesos, with extended features for orchestration, management, and task scheduling.

By interacting entirely with the new Mesos [HTTP scheduler API](http://mesos.apache.org/documentation/latest/scheduler-http-api/), `mesos-cli` does not require a direct network connection to the Mesos Master server, which makes it more flexible than other frameworks.

Check out the documentation for mesos-cli [here](https://mesanine.co/mesos-cli).

## Quickstart
Simply run `mesos-cli` via the official Docker image to get started:
```bash
docker run --rm -ti quay.io/mesanine/mesos-cli:latest tasks --master http://your-mesos-server:5050
```
Full install and configuration documentation is available [here](https://mesanine.co/mesos-cli/getting-started/)

## Commands
cmd | description
--- | ---
[agents][usage_agents] | List Mesos Agents
[list][usage_list] | List files in a Mesos sandbox
[read][usage_read] | Read the contents of a file
[run][usage_run] | Run tasks on Mesos
[tasks][usage_tasks] | List Mesos tasks
[top][usage_top] | Display a Mesos top interface

## Distinctive Features & Roadmap

| Feature                                                                             |ready|roadmap|
|-------------------------------------------------------------------------------------|-----|-------|
| Built ontop of the new Mesos HTTP V1 API                                            |✓    |       |
| Simple installation without platform specific libmesos drivers                      |✓    |       |
| Streamimg sandbox file content content to console (including task stdout and stderr)|✓    |       |
| Pure integration with Apache Mesos outside of the DC/OS ecosystem                   |✓    |       |
| Lightweight scheduler for running arbitrary containers                              |✓    |       |
| Top-like interface for monitoring a cluster                                         |✓    |       |
| Full support for latest Mesos features, e.g health checks, etc                      |     |✓      |
| Ability to search and filter across most Mesos types                                |     |✓      |
| Subscribe to and monitor master event stream                                        |     |✓      |
| Run docker-compose files directly against Mesos                                     |     |✓      |

[usage_agents]: https://mesanine.co/mesos-cli/usage#agents
[usage_list]: https://mesanine.co/mesos-cli/usage#list
[usage_read]: https://mesanine.co/mesos-cli/usage#read
[usage_run]: https://mesanine.co/mesos-cli/usage#run
[usage_tasks]: https://mesanine.co/mesos-cli/usage#tasks
[usage_top]: https://mesanine.co/mesos-cli/usage#top

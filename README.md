# mesos-cli

Featureful commandline interface for [Apache Mesos](http://mesos.apache.com).

**NOTE**: *mesos-cli is under active development and not yet considered stable!*

`mesos-cli` is designed to be a lightweight alternative to the [native tool](https://github.com/apache/mesos/tree/master/src/cli) provided with Mesos, with extended features for orchestration, management, and task scheduling.

By interacting entirely with the new Mesos [HTTP scheduler API](http://mesos.apache.org/documentation/latest/scheduler-http-api/), `mesos-cli` does not require a direct network connection to the Mesos Master server, which makes it more flexible than other frameworks.

Check out the documentation for mesos-cli [here](https://vektorlab.github.io/mesos-cli).

## Quickstart
Simply run `mesos-cli` via the official Docker image to get started:
```bash
docker run --rm -ti quay.io/vektorcloud/mesos-cli:latest tasks --master http://your-mesos-server:5050
```
Full install and configuration documentation is available [here](https://vektorlab.github.io/mesos-cli/getting-started/)

## Commands
--- | ---
[agents][usage_agents] | List Mesos Agents
[list][usage_list] | List files in a Mesos sandbox
[local][usage_local] | Run a local Mesos cluster
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
| Effortlessly run a local Mesos cluster                                              |✓    |       |
| Pure integration with Apache Mesos outside of the DC/OS ecosystem                   |✓    |       |
| Lightweight scheduler for running arbitrary containers                              |✓    |       |
| Top-like interface for monitoring a cluster                                         |✓    |       |
| Full support for latest Mesos features, e.g health checks, etc                      |     |✓      |
| Ability to search and filter across most Mesos types                                |     |✓      |
| Subscribe to and monitor master event stream                                        |     |✓      |
| Run docker-compose files directly against Mesos                                     |     |✓      |

[usage_agents]: https://vektorlab.github.io/mesos-cli/usage/agents
[usage_list]: https://vektorlab.github.io/mesos-cli/usage/list
[usage_local]: https://vektorlab.github.io/mesos-cli/usage/local
[usage_read]: https://vektorlab.github.io/mesos-cli/usage/read
[usage_run]: https://vektorlab.github.io/mesos-cli/usage/run
[usage_tasks]: https://vektorlab.github.io/mesos-cli/usage/tasks
[usage_top]: https://vektorlab.github.io/mesos-cli/usage/top

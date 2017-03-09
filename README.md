# mesos-cli

Featureful commandline interface for [Apache Mesos](http://mesos.apache.com).

`mesos-cli` is designed to be a lightweight alternative to the [native tool](https://github.com/apache/mesos/tree/master/src/cli) provided with Mesos, with extended features for orchestration, management, and task scheduling.

By interacting entirely with the new Mesos [HTTP scheduler API](http://mesos.apache.org/documentation/latest/scheduler-http-api/), `mesos-cli` does not require a direct network connection to the Mesos Master server, which makes it more flexible than other frameworks.

**mesos-cli is under active development and not yet considered stable!**

Check out the documentation for mesos-cli [here](https://vektorlab.github.io/mesos-cli).

## Quickstart
Simply run `mesos-cli` via the official Docker image to get started:
```bash
docker run --rm -ti quay.io/vektorcloud/mesos-cli:latest tasks --master http://your-mesos-server:5050
```
Full install and configuration documentation is available [here](https://vektorlab.github.io/mesos-cli/getting-started/)


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


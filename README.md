# mesos-cli

**mesos-cli** is a command-line tool for running containers on and interacting with [Apache Mesos](http://mesos.apache.com). It is designed to be lightweight and more featureful than the [native](https://github.com/apache/mesos/tree/master/src/cli) CLI tool provided with Mesos. Additionally, it seeks to act as a stand-alone and extensible task scheduler on Mesos.

By interacting entirely with the new Mesos HTTP [scheduler](http://mesos.apache.org/documentation/latest/scheduler-http-api/) API as a general purpose Mesos framework, **mesos-cli** does not require the Mesos Master server to establish a direct network connection to it, which makes it more flexible than other frameworks.

**mesos-cli is under active development and not yet considered stable!**

Check out the documentation for mesos-cli [here](https://vektorlab.github.io/mesos-cli).


# Distinctive Features & Roadmap

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


# mesos-cli
[![Build](https://img.shields.io/circleci/project/github/mesanine/mesos-cli.svg)](https://circleci.com/gh/mesanine/mesos-cli)

Featureful commandline interface for [Apache Mesos](http://mesos.apache.com).

`mesos-cli` is designed to be a lightweight alternative to the [native tool](https://github.com/apache/mesos/tree/master/src/cli) provided with Mesos, with extended features for orchestration, management, and task scheduling.

By interacting entirely with the new Mesos [HTTP scheduler API](http://mesos.apache.org/documentation/latest/scheduler-http-api/), `mesos-cli` does not require a direct network connection to the Mesos Master server, which makes it more flexible than other frameworks.

**NOTE**: *mesos-cli is under active development and not yet considered stable!*

# Distinctive Features & Roadmap

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

# Configuration

It is recommended that you create a bash alias for mesos-cli, note that if the default `mesos` executable is installed on your system this will override it.

     # With a binary installation
     echo "alias mesos=mesos-cli" >> $HOME/.bashrc
     # With a Docker installation
     echo "alias mesos=docker run --rm -ti -v $HOME/.mesos-cli.json:/root/.mesos-cli.json --net host mesanine/mesos-cli >> $HOME/.bashrc"

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

# Usage

```

Usage: mesos-cli [OPTIONS] COMMAND [arg...]

Alternative Apache Mesos CLI

Options:
  --profile="default"                      Profile to load
  --config="~/.mesos-cli.json"             Path to load config from
  --debug=false                            Enable debugging
  --version                                Show the version and exit

Commands:
  agents       List Mesos Agents
  list         List files in a Mesos sandbox
  read         Read the contents of a file
  run          Run tasks on Mesos
  task         Generate a Mesos Task
  tasks        List Mesos tasks
  top          Display a Mesos top interface

Run 'mesos-cli COMMAND --help' for more information on a command.
```

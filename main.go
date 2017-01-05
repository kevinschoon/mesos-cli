package main

import (
	"flag"
	"fmt"
	"github.com/jawher/mow.cli"
	"os"
	"os/user"
)

const Version = "0.0.1"

// Singleton config object
var config *Config

func main() {
	app := cli.App("mesos-cli", "Alternative Apache Mesos CLI")
	app.Spec = "[OPTIONS]"
	app.Command("ps", "List currently running tasks on a cluster", ps)
	app.Command("agents", "List registered agents on a cluster", agents)
	app.Command("exec", "Execute arbitrary commands against a cluster", exec)
	app.Command("local", "Launch a local Mesos cluster (requires Docker)", local)
	var (
		profile    = app.StringOpt("profile", "default", "Profile to load")
		configPath = app.StringOpt("config", fmt.Sprintf("%s/.mesos-cli.json", homeDir()),
			"Path to load config from")
		level = app.IntOpt("level", 0, "Level of verbosity")
		err   error
	)

	// This is done to satisfy the presumptuous golang/glog package
	// which assumes I am using flag and insists it be configured
	// with such. Since glog is used in go-mesos it is easiest to use
	// the same library for the moment.
	flag.CommandLine.Set("v", string(*level))
	flag.CommandLine.Set("logtostderr", "1")
	flag.CommandLine.Parse([]string{})

	app.Before = func() {
		config, err = LoadConfig(*configPath, *profile)
		failOnErr(err)
	}
	app.Run(os.Args)
}

func homeDir() string {
	u, err := user.Current()
	if err != nil {
		cli.Exit(1)
	}
	return u.HomeDir
}

func failOnErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		cli.Exit(1)
	}
}

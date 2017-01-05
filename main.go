package main

import (
	"flag"
	"fmt"
	"github.com/jawher/mow.cli"
	"os"
)

const Version = "0.0.1"

// Singleton config object
var config *Config

func main() {
	app := cli.App("mesos-cli", "Alternative Apache Mesos CLI")
	app.Spec = "[OPTIONS]"
	app.Command("ps", "List currently running tasks on a cluster", ps)
	app.Command("exec", "Execute Arbitrary Commands Against a Cluster", exec)
	app.Command("local", "Launch a local Mesos cluster (requires Docker)", local)
	var (
		master  = app.StringOpt("master", "127.0.0.1:5050", "Master address <host:port>")
		profile = app.StringOpt("profile", "default", "Profile to load from ~/.mesos-cli.json")
		level   = app.IntOpt("level", 0, "Level of verbosity")
		err     error
	)

	// This is done to satisfy the presumptuous golang/glog package
	// which assumes I am using flag and insists it be configured
	// with such. Since glog is used in go-mesos it is easiest to use
	// the same library for the moment.
	flag.CommandLine.Set("v", string(*level))
	flag.CommandLine.Set("logtostderr", "1")
	flag.CommandLine.Parse([]string{})

	config, err = LoadConfig(*profile, *master)
	failOnErr(err)
	app.Run(os.Args)
}

func failOnErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		cli.Exit(1)
	}
}

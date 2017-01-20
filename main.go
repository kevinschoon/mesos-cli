package main

import (
	"flag"
	"fmt"
	"github.com/jawher/mow.cli"
	"os"
)

var (
	Version = "undefined"
	GitSHA  = "undefied"
)

// Singleton config object
var config *Config

func main() {
	app := cli.App("mesos-cli", "Alternative Apache Mesos CLI")
	app.Spec = "[OPTIONS]"
	app.Command("cat", "Output the contents of a file", cat)
	app.Command("local", "Launch a local Mesos cluster (requires Docker)", local)
	app.Command("ls", "List the sandbox directory of a task", ls)
	app.Command("ps", "List currently running tasks on a cluster", ps)
	app.Command("run", "Run arbitrary commands against a cluster", run)
	app.Command("top", "Display a Mesos top dialog", topCmd)

	var (
		profile    = app.StringOpt("profile", "default", "Profile to load")
		configPath = app.StringOpt("config", fmt.Sprintf("%s/.mesos-cli.json", homeDir()),
			"Path to load config from")
		level = app.IntOpt("level", 0, "Logging level (higher is more verbose)")
		err   error
	)

	app.Version("version", fmt.Sprintf("%s (%s)", Version, GitSHA))

	app.Before = func() {
		// This is done to satisfy the presumptuous golang/glog package
		// which assumes I am using flag and insists it be configured
		// with such. Since glog is used in go-mesos it is easiest to use
		// the same library for the moment.
		flag.CommandLine.Parse([]string{fmt.Sprintf("-v=%d", *level), "-logtostderr=true"})

		config, err = LoadConfig(*configPath, *profile)
		failOnErr(err)
	}
	app.Run(os.Args)
}

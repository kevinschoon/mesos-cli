package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/commands"
	"github.com/vektorlab/mesos-cli/config"
	"os"
)

var (
	Version = "undefined"
	GitSHA  = "undefied"
)

func main() {
	app := cli.App("mesos-cli", "Alternative Apache Mesos CLI")
	app.Spec = "[OPTIONS]"
	/*
		app.Command("agents", "List Mesos agents", agents)
		app.Command("cat", "Output the contents of a file", cat)
		app.Command("local", "Launch a local Mesos cluster (requires Docker)", local)
		app.Command("ls", "List the sandbox directory of a task", ls)
		app.Command("tasks", "List currently running tasks on a cluster", tasks)
		app.Command("run", "Run arbitrary commands against a cluster", run)
		app.Command("top", "Display a Mesos top dialog", topCmd)
	*/
	var (
		cfg         *config.Config
		profileName = app.StringOpt("profile", "default", "Profile to load")
		configPath  = app.StringOpt("config", fmt.Sprintf("%s/.mesos-cli.json", config.HomeDir()),
			"Path to load config from")
		err error
	)

	app.Version("version", fmt.Sprintf("%s (%s)", Version, GitSHA))

	app.Before = func() {
		cfg, err = config.LoadConfig(*configPath, *profileName)
		failOnErr(err)
		fmt.Println("Profile: ", cfg.Profile())
	}
	cfgFn := func() *config.Config { return cfg }
	for _, cmd := range commands.Commands {
		app.Command(cmd.Name(), cmd.Desc(), cmd.Init(cfgFn))
	}

	app.Run(os.Args)
}

func failOnErr(err error) {
	if err != nil {
		fmt.Printf("Encountered Error: %v", err)
		os.Exit(2)
	}
}

package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/mesanine/mesos-cli/commands"
	"github.com/mesanine/mesos-cli/config"

	"os"
)

var (
	Version = "undefined"
	GitSHA  = "undefied"
)

func main() {
	app := cli.App("mesos-cli", "Alternative Apache Mesos CLI")
	app.Spec = "[OPTIONS]"
	var (
		name = app.StringOpt("profile", "default", "Profile to load")
		path = app.StringOpt("config", fmt.Sprintf("%s/.mesos-cli.json", config.HomeDir()),
			"Path to load config from")
		debug   = app.BoolOpt("debug", false, "Enable debugging")
		profile *config.Profile
	)
	app.Version("version", fmt.Sprintf("%s (%s)", Version, GitSHA))
	app.Before = func() {
		p, err := config.LoadProfile(*path, *name)
		if err != nil {
			fmt.Printf("Could not load configuration profile %s: %s\n", *name, err)
			os.Exit(2)
		}
		profile = p.With(config.Debug(*debug))
	}
	for _, cmd := range commands.Commands {
		app.Command(cmd.Name(), cmd.Desc(), cmd.Init(func() *config.Profile { return profile }))
	}
	app.Run(os.Args)
}

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
		if err != nil {
			fmt.Println("Unable to load configuration: %v", err)
			os.Exit(2)
		}
	}
	for _, cmd := range commands.Commands {
		cmd.SetConfig(func() *config.Config { return cfg })
		app.Command(cmd.Name(), cmd.Desc(), cmd.Init())
	}

	app.Run(os.Args)
}

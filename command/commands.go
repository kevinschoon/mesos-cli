package command

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/mesanine/mesos-cli/config"

	"os"
)

var (
	Version = "undefined"
	GitSHA  = "undefied"
	profile *config.Profile
)

func MesosCLI() {
	app := cli.App("mesos-cli", "Alternative Apache Mesos CLI")
	app.Spec = "[OPTIONS]"
	var (
		name  = app.StringOpt("profile", "default", "profile to load")
		path  = app.StringOpt("config", fmt.Sprintf("%s/.mesos-cli.json", config.HomeDir()), "config path")
		debug = app.BoolOpt("debug", false, "enable debugging")
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
	app.Command("generate", "generate Mesos API calls", func(cmd *cli.Cmd) {
		cmd.Command("master", "generate Mesos master API calls", MasterCMDs)
	})
	app.Run(os.Args)
}

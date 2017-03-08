package commands

import (
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/top"
)

type Top struct{}

func (_ Top) Name() string { return "top" }
func (_ Top) Desc() string { return "Display a Mesos top interface" }

func (_ Top) Init(profile config.ProfileFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		hostname := cmd.StringOpt("master", "", "Mesos Master")
		cmd.Action = func() {
			failOnErr(top.Run(profile().With(config.Master(*hostname))))
		}
	}
}

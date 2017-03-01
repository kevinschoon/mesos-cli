package commands

import (
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/config"
)

type Top struct{}

func (_ Top) Name() string { return "top" }
func (_ Top) Desc() string { return "Display a Mesos top interface" }

func (top Top) Init(_ config.ProfileFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {}
}

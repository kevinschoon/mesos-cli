package commands

import (
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/config"
)

type Top struct {
	*command
}

func NewTop() Command {
	return Top{
		&command{"top", "Display a Mesos top interface"},
	}
}

func (top Top) Init(cfg config.CfgFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {}
}

package commands

import (
	"github.com/jawher/mow.cli"
)

type Top struct {
	*command
}

func NewTop() Command {
	return Top{
		command: &command{
			name: "top",
			desc: "Display a Mesos top interface",
		},
	}
}

func (top Top) Init() func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {}
}

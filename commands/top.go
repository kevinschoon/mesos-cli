package commands

import (
	"github.com/jawher/mow.cli"
)

type Top struct{}

func (_ Top) Name() string { return "top" }
func (_ Top) Desc() string { return "Display a Mesos top interface" }

func (top Top) Init(_ Profile) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {}
}

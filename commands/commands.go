package commands

import (
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/config"
)

var Commands = []Command{
	&Agents{},
	&List{},
	&Local{},
	&Read{},
	&Run{},
	&Tasks{},
	&Top{},
}

type Profile func() *config.Profile

type Command interface {
	Name() string
	Desc() string
	SetProfile(Profile)
	Init() func(*cli.Cmd)
}

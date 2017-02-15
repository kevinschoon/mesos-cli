package commands

import (
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/config"
)

var Commands = []Command{
	NewAgents(),
	NewList(),
	NewLocal(),
	NewRead(),
	NewRun(),
	NewTasks(),
	NewTop(),
}

type Command interface {
	Name() string
	Desc() string
	Init(config.CfgFn) func(*cli.Cmd)
}

type command struct {
	name string
	desc string
}

func (c command) Name() string { return c.name }
func (c command) Desc() string { return c.desc }

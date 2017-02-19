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
	SetConfig(config.ConfigFn)
	Init() func(*cli.Cmd)
}

type command struct {
	name   string
	desc   string
	config config.ConfigFn
}

func (c command) Name() string                  { return c.name }
func (c command) Desc() string                  { return c.desc }
func (c *command) SetConfig(fn config.ConfigFn) { c.config = fn }

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

package commands

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	master "github.com/mesos/mesos-go/master/calls"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/filter"
)

type Agents struct {
	*command
	Hostname *string
}

func NewAgents() Command {
	return &Agents{
		command: &command{
			name: "agents",
			desc: "List Mesos Agents",
		},
	}
}

func (a Agents) Action() {
	resp, err := NewCaller(a.config().Profile()).CallMaster(master.GetAgents())
	failOnErr(err)
	table := uitable.New()
	table.AddRow("ID", "HOSTNAME", "CPUS", "MEM", "GPUS", "DISK")

	for _, agent := range filter.AsAgents(filter.FromMaster(resp).FindMany()) {
		table.AddRow(
			agent.GetID().GetValue(),
			agent.GetHostname(),
			fmt.Sprintf("%.2f", Scalar("cpus", agent.Resources)),
			fmt.Sprintf("%.2f", Scalar("mem", agent.Resources)),
			fmt.Sprintf("%.2f", Scalar("gpus", agent.Resources)),
			fmt.Sprintf("%.2f", Scalar("disk", agent.Resources)),
		)
	}
	fmt.Println(table)
}

func (a *Agents) Init() func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		a.Hostname = cmd.StringOpt("hostname", config.DefaultProfile().Master, "Mesos Master")
		cmd.Action = a.Action
	}
}

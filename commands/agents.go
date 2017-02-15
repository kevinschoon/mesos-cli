package commands

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/client"
	"github.com/vektorlab/mesos-cli/config"
)

type Agents struct {
	*command
}

func NewAgents() Command {
	return Agents{
		&command{"agents", "List Mesos Agents"},
	}
}

func (_ Agents) Init(cfg config.CfgFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			defaults = config.DefaultProfile()
			hostname = cmd.StringOpt("hostname", defaults.Master, "Mesos Master")
		)
		cmd.Action = func() {
			client := client.New(
				cfg().Profile(
					config.WithMaster(*hostname),
				),
			)
			agents, err := client.Agents()
			failOnErr(err)
			table := uitable.New()
			table.AddRow("ID", "HOSTNAME", "CPUS", "MEM", "GPUS", "DISK")
			for _, agent := range agents {
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
	}
}

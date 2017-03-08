package commands

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	master "github.com/mesos/mesos-go/master/calls"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/filter"
	"github.com/vektorlab/mesos-cli/helper"
)

type Agents struct{}

func (_ Agents) Name() string { return "agents" }
func (_ Agents) Desc() string { return "List Mesos Agents" }

func (_ Agents) Init(profile config.ProfileFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		hostname := cmd.StringOpt("master", "", "Mesos Master")
		cmd.Action = func() {
			resp, err := helper.NewCaller(profile().With(
				config.Master(*hostname),
			)).CallMaster(master.GetAgents())
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
	}
}

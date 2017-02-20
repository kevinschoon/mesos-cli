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
	Hostname *string
	profile  Profile
}

func (_ Agents) Name() string { return "agents" }
func (_ Agents) Desc() string { return "List Mesos Agents" }

func (a *Agents) SetProfile(p Profile) {
	a.profile = func() *config.Profile {
		profile := p()
		if *a.Hostname != "" {
			profile = profile.With(
				config.Master(*a.Hostname),
			)
		}
		return profile
	}
}

func (a *Agents) Init() func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		a.Hostname = cmd.StringOpt("hostname", "", "Mesos Master")
		cmd.Action = func() {
			resp, err := NewCaller(a.profile()).CallMaster(master.GetAgents())
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

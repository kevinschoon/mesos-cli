package commands

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/client"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/filters"
	"github.com/vektorlab/mesos-cli/options"
)

type Tasks struct {
	*command
}

func NewTasks() Command {
	return Tasks{
		&command{"tasks", "List Mesos Tasks"},
	}
}

func (t Tasks) Init(cfg config.CfgFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			defaults = config.DefaultProfile()
			hostname = cmd.StringOpt("master", defaults.Master, "Mesos Master")
			truncate = cmd.BoolOpt("truncate", true, "truncate some values")
			// TODO: Reduce duplicated code
			fuzzy  = cmd.BoolOpt("fuzzy", true, "fuzzy match")
			name   = cmd.StringOpt("name", "", "filter by task name")
			states = options.NewStates()
		)
		cmd.VarOpt("state", &states, "Filter by task state")
		cmd.Action = func() {
			client := client.New(
				cfg().Profile(
					config.WithMaster(*hostname),
				),
			)
			filter := filters.TaskFilter{
				Name:   *name,
				Fuzzy:  *fuzzy,
				States: states,
			}
			tasks, err := client.Tasks()
			failOnErr(err)
			table := uitable.New()
			table.AddRow("ID", "FRAMEWORK", "STATE", "CPU", "MEM", "GPU", "DISK")
			for _, task := range tasks {
				if filter.Match(task) {
					frameworkID := task.FrameworkID.Value
					if *truncate {
						frameworkID = truncStr(frameworkID, 8)
					}
					table.AddRow(
						task.TaskID.Value,
						frameworkID,
						task.GetState().String(),
						Scalar("cpus", task.Resources),
						Scalar("mem", task.Resources),
						Scalar("gpus", task.Resources),
						Scalar("disk", task.Resources),
					)
				}
			}
			fmt.Println(table)
		}

	}
}

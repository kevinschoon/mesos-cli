package commands

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	master "github.com/mesos/mesos-go/master/calls"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/options"
	"github.com/vektorlab/mesos-cli/state"
)

type Tasks struct {
	command
	Hostname *string
	TaskID   *string
	Truncate *bool
	Fuzzy    *bool
	States   options.States
}

func NewTasks() Command {
	return &Tasks{
		command: command{
			name: "tasks",
			desc: "List Mesos Tasks",
		},
		States: options.NewStates(),
	}
}

func (t *Tasks) filters() []state.Filter {
	filters := []state.Filter{
		state.TaskStateFilter(t.States),
	}
	if *t.TaskID != "" {
		filters = append(filters, state.TaskIDFilter(*t.TaskID, *t.Fuzzy))
	}
	return filters
}

func (t *Tasks) Action() {
	resp, err := NewCaller(t.configFn().Profile()).CallMaster(master.GetState())
	failOnErr(err)

	table := uitable.New()
	table.AddRow("ID", "FRAMEWORK", "STATE", "CPU", "MEM", "GPU", "DISK")

	for _, task := range state.AsTasks(state.StateFromMaster(resp.GetState).FindMany(t.filters()...)) {
		frameworkID := task.FrameworkID.Value
		if *t.Truncate {
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

	fmt.Println(table)

}

func (t *Tasks) Init() func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		t.Hostname = cmd.StringOpt("master", config.DefaultProfile().Master, "Mesos Master")
		t.Truncate = cmd.BoolOpt("truncate", true, "truncate long values")
		t.TaskID = cmd.StringOpt("task", "", "filter by task id")
		t.Fuzzy = cmd.BoolOpt("fuzzy", true, "fuzzy matching on string values")
		cmd.VarOpt("state", &t.States, "filter by task state")
		cmd.Action = t.Action
	}
}

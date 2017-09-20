package commands

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	"github.com/mesanine/mesos-cli/config"
	"github.com/mesanine/mesos-cli/filter"
	"github.com/mesos/mesos-go"
)

type Tasks struct{}

func (_ Tasks) Name() string { return "tasks" }
func (_ Tasks) Desc() string { return "List Mesos tasks" }

func (_ Tasks) Init(profile config.ProfileFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			hostname = cmd.StringOpt("master", "", "Mesos master")
			truncate = cmd.BoolOpt("truncate", true, "Truncate long values")
			taskID   = cmd.StringOpt("task", "", "Filter by task id")
			name     = cmd.StringOpt("name", "", "Filter by Task name")
			fuzzy    = cmd.BoolOpt("fuzzy", true, "Fuzzy matching on string values")
			all      = cmd.BoolOpt("a all", false, "Show all tasks")
			states   = states([]*mesos.TaskState{})
		)
		cmd.VarOpt("state", &states, "filter by task state")
		cmd.Action = func() {

			criteria := filter.Criteria{
				Target:     filter.TASKS,
				TaskID:     *taskID,
				TaskName:   *name,
				Fuzzy:      *fuzzy,
				TaskStates: []*mesos.TaskState{},
			}

			if len(criteria.TaskStates) == 0 && !*all {
				criteria.TaskStates = append(criteria.TaskStates, mesos.TASK_RUNNING.Enum())
			}

			table := uitable.New()
			table.AddRow("ID", "NAME", "FRAMEWORK", "STATE", "CPU", "MEM", "GPU", "DISK")

			results, err := filter.Find(profile().With(config.Master(*hostname)), criteria)
			failOnErr(err)

			for _, task := range filter.AsTasks(results) {
				frameworkID := task.FrameworkID.Value
				if *truncate {
					frameworkID = truncStr(frameworkID, 8)
				}
				table.AddRow(
					task.TaskID.Value,
					task.Name,
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
	}
}

type states []*mesos.TaskState

func (o *states) String() string {
	return fmt.Sprintf("%v", *o)
}

func (o *states) Set(name string) error {
	v, ok := mesos.TaskState_value[name]
	if !ok {
		return fmt.Errorf("Invalid state %s", name)
	}
	*o = append(*o, mesos.TaskState(v).Enum())
	return nil
}

func (o *states) Clear() {
	*o = []*mesos.TaskState{}
}

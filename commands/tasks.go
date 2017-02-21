package commands

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	"github.com/mesos/mesos-go"
	master "github.com/mesos/mesos-go/master/calls"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/filter"
)

type Tasks struct{}

func (_ Tasks) Name() string { return "tasks" }
func (_ Tasks) Desc() string { return "List Mesos tasks" }

func (t *Tasks) Init(profile Profile) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS]"
		var (
			hostname = cmd.StringOpt("master", "", "Mesos master")
			truncate = cmd.BoolOpt("truncate", true, "Truncate long values")
			taskID   = cmd.StringOpt("task", "", "Filter by task id")
			fuzzy    = cmd.BoolOpt("fuzzy", true, "Fuzzy matching on string values")
			states   = states([]*mesos.TaskState{})
		)
		cmd.VarOpt("state", &states, "filter by task state")
		cmd.Action = func() {

			resp, err := NewCaller(
				profile().With(config.Master(*hostname)),
			).CallMaster(master.GetTasks())
			failOnErr(err)

			filters := []filter.Filter{
				filter.TaskStateFilter(states),
			}

			if *taskID != "" {
				filters = append(filters, filter.TaskIDFilter(*taskID, *fuzzy))
			}

			table := uitable.New()
			table.AddRow("ID", "FRAMEWORK", "STATE", "CPU", "MEM", "GPU", "DISK")

			for _, task := range filter.AsTasks(filter.FromMaster(resp).FindMany(filters...)) {
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

			fmt.Println(table)
		}
	}
}

type states []*mesos.TaskState

//func NewStates() States {
//	return states{mesos.TASK_RUNNING.Enum()}
//}

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

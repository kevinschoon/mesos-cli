package commands

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/jawher/mow.cli"
	"github.com/mesos/mesos-go"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/options"
	"github.com/vektorlab/mesos-cli/runner"
	"os"
)

type Run struct {
	*command
}

func NewRun() Command {
	return Run{
		&command{"run", "Run Tasks on Mesos"},
	}
}

func (r Run) Init(cfg config.CfgFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [CMD]"
		var (
			resources = mesos.Resources{
				mesos.Resource{
					Name:   "cpus",
					Type:   mesos.SCALAR.Enum(),
					Role:   proto.String("*"),
					Scalar: &mesos.Value_Scalar{Value: 0.1},
				},
			}
			defaults = config.DefaultProfile()
			master   = cmd.StringOpt("master", defaults.Master, "Mesos Master")
			command  = cmd.StringArg("CMD", "", "Command to run")
			name     = cmd.StringOpt("name", "mesos-cli", "Task Name")
			user     = cmd.StringOpt("user", "root", "User to run as")
			shell    = cmd.BoolOpt("shell", true, "Run as a shell command")
			toJson   = cmd.BoolOpt("json", false, "Write task to JSON instead of running")
			taskID   = options.NewTaskID()
		)

		cmd.VarOpt("TaskID", taskID, "Mesos TaskID")
		cmd.VarOpt("cpus", options.NewScalarResources("cpus", resources), "CPU Resources")
		cmd.VarOpt("memory", options.NewScalarResources("memory", resources), "Memory Resources")
		cmd.VarOpt("disk", options.NewScalarResources("disk", resources), "Disk Resources")

		cmd.Before = func() {}

		cmd.Action = func() {
			info := &mesos.TaskInfo{
				TaskID: *taskID.ID,
				Name:   *name,
				Command: &mesos.CommandInfo{
					Shell: shell,
					User:  user,
					Value: command,
				},
				Resources: resources,
			}
			if *toJson {
				raw, err := json.Marshal(info)
				failOnErr(err)
				fmt.Println(string(raw))
				os.Exit(0)
			}
			failOnErr(
				runner.New(
					cfg().Profile(
						config.WithMaster(*master),
					),
				).Run(info))
		}
	}
}

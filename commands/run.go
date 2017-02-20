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
	Hostname *string
	profile  Profile
}

func (_ Run) Name() string { return "run" }
func (_ Run) Desc() string { return "Run tasks on Mesos" }
func (r *Run) SetProfile(p Profile) {
	r.profile = func() *config.Profile {
		profile := p()
		if *r.Hostname != "" {
			profile = profile.With(
				config.Master(*r.Hostname),
			)
		}
		return profile
	}
}

func (r *Run) Init() func(*cli.Cmd) {
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
			command = cmd.StringArg("CMD", "", "Command to run")
			name    = cmd.StringOpt("name", "mesos-cli", "Task Name")
			user    = cmd.StringOpt("user", "root", "User to run as")
			shell   = cmd.BoolOpt("shell", true, "Run as a shell command")
			toJson  = cmd.BoolOpt("json", false, "Write task to JSON instead of running")
			taskID  = options.NewTaskID()
		)

		r.Hostname = cmd.StringOpt("master", "", "Mesos Master")
		cmd.VarOpt("TaskID", taskID, "Mesos TaskID")
		cmd.VarOpt("cpus", options.NewScalarResources("cpus", resources), "CPU Resources")
		cmd.VarOpt("memory", options.NewScalarResources("memory", resources), "Memory Resources")
		cmd.VarOpt("disk", options.NewScalarResources("disk", resources), "Disk Resources")

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
			failOnErr(runner.New(r.profile()).Run(info))
		}
	}
}

package commands

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/jawher/mow.cli"
	"github.com/mesos/mesos-go"
	"github.com/vektorlab/mesos-cli/commands/flags"
	"github.com/vektorlab/mesos-cli/commands/options"
	"github.com/vektorlab/mesos-cli/config"
)

type Task struct{}

func (_ Task) Name() string { return "task" }
func (_ Task) Desc() string { return "Generate a Mesos Task" }

func (_ Task) Init(_ config.ProfileFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {

		cmd.Spec = "[OPTIONS] [CMD]"

		task := &mesos.TaskInfo{
			Name: "mesos-cli",
			Container: &mesos.ContainerInfo{
				Mesos: &mesos.ContainerInfo_MesosInfo{},
				Docker: &mesos.ContainerInfo_DockerInfo{
					Network:      mesos.ContainerInfo_DockerInfo_BRIDGE.Enum(),
					PortMappings: []mesos.ContainerInfo_DockerInfo_PortMapping{},
					Parameters:   []mesos.Parameter{},
				},
			},
			Resources: []mesos.Resource{},
		}

		for _, flag := range flags.Flags {
			flag(task, cmd)
		}

		var (
			encoding = cmd.StringOpt("encoding", "json", "Output encoding [json/yaml]")
			asDocker = cmd.BoolOpt("docker", false, "Run as a Docker container")
			role     = cmd.StringOpt("role", "*", "Mesos role")
		)

		cmd.Action = func() {

			options.Apply(
				task,
				options.WithContainerizer(*asDocker),
				options.WithPorts(),
				options.WithDefaultResources(),
				options.WithRole(*role),
			)

			out := []*mesos.TaskInfo{task}
			switch *encoding {
			case "json":
				raw, err := json.Marshal(&out)
				failOnErr(err)
				fmt.Println(string(raw))
			case "yaml":
				raw, err := yaml.Marshal(&out)
				failOnErr(err)
				fmt.Println(string(raw))
			default:
				failOnErr(fmt.Errorf("Bad encoding %s", *encoding))
			}
		}
	}
}

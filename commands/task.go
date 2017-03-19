package commands

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/jawher/mow.cli"
	"github.com/mesos/mesos-go"
	"github.com/vektorlab/mesos-cli/commands/flags"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/mesosfile"
)

type Task struct{}

func (_ Task) Name() string { return "task" }
func (_ Task) Desc() string { return "Generate a Mesos Task" }

func (_ Task) Init(_ config.ProfileFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {

		cmd.Spec = "[OPTIONS] [CMD]"

		task := mesosfile.NewTask()

		for _, flag := range flags.Flags {
			flag(task, cmd)
		}

		var (
			encoding = cmd.StringOpt("encoding", "json", "Output encoding [json/yaml]")
			docker   = cmd.BoolOpt("docker", false, "Run as a Docker container")
			role     = cmd.StringOpt("role", "*", "Mesos role")
		)

		cmd.Action = func() {
			group := &mesosfile.Group{Tasks: []*mesos.TaskInfo{task}}
			group = group.With(
				mesosfile.Init(),
				mesosfile.Role(*role),
				mesosfile.Docker(*docker),
			)

			out := mesosfile.Mesosfile{group}

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

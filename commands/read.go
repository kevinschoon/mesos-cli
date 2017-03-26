package commands

import (
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/filter"
	"github.com/vektorlab/mesos-cli/helper"
	"github.com/vektorlab/mesos-cli/pailer"
	"os"
)

type Read struct{}

func (_ Read) Name() string { return "read" }
func (_ Read) Desc() string { return "Read the contents of a file" }

func (_ Read) Init(profile config.ProfileFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID PATH"
		var (
			agentID  = cmd.StringArg("ID", "", "AgentID")
			path     = cmd.StringArg("PATH", "", "path to read")
			follow   = cmd.BoolOpt("f follow", false, "follow the content")
			lines    = cmd.IntOpt("n nlines", 0, "number of lines to read")
			hostname = cmd.StringOpt("m master", "", "mesos master")
		)
		cmd.Action = func() {
			msgs, err := filter.Find(
				profile().With(config.Master(*hostname)),
				filter.Criteria{Target: filter.AGENTS, AgentID: *agentID},
			)
			failOnErr(err)
			agent, err := filter.AsAgent(msgs.FindOne())
			failOnErr(err)
			pg := &pailer.FilePaginator{
				Path:   *path,
				Follow: *follow,
				Max:    uint64(*lines),
			}
			failOnErr(pailer.Monitor(helper.NewAgentCaller(profile().With(config.Master(*hostname)), agent), os.Stdout, -1, pg))
		}
	}
}

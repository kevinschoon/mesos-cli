package commands

import (
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/pailer"
	"os"
)

type Read struct{}

func (_ Read) Name() string { return "read" }
func (_ Read) Desc() string { return "Read the contents of a file" }

func (r *Read) Init(profile Profile) func(*cli.Cmd) {
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
			caller, err := NewAgentCaller(profile().With(config.Master(hostname)), *agentID)
			failOnErr(err)
			pg := &pailer.FilePaginator{
				Path:   *path,
				Follow: *follow,
				Max:    uint64(*lines),
			}
			failOnErr(pailer.Monitor(caller, os.Stdout, -1, pg))
		}
	}
}

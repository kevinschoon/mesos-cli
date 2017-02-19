package commands

import (
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/pailer"
	"os"
)

type Read struct {
	*command
	AgentID *string
	Path    *string
	Follow  *bool
	Lines   *int
}

func NewRead() Command {
	return &Read{
		command: &command{
			name: "read",
			desc: "Read the contents of a file",
		},
	}
}

func (r *Read) Action() {
	caller, err := NewAgentCaller(r.config().Profile(), *r.AgentID)
	failOnErr(err)
	pg := &pailer.FilePaginator{
		Path:   *r.Path,
		Follow: *r.Follow,
		Max:    uint64(*r.Lines),
	}
	failOnErr(pailer.Monitor(caller, os.Stdout, -1, pg))
}

func (r *Read) Init() func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID PATH"
		cmd.Action = r.Action
		r.AgentID = cmd.StringArg("ID", "", "AgentID")
		r.Path = cmd.StringArg("PATH", "", "path to read")
		r.Follow = cmd.BoolOpt("f follow", false, "follow the content")
		r.Lines = cmd.IntOpt("n nlines", 0, "number of lines to read")
	}
}

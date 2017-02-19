package commands

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	agent "github.com/mesos/mesos-go/agent/calls"
	"strings"
)

type List struct {
	*command
	Path     *string
	AgentID  *string
	Relative *bool
}

func NewList() Command {
	return &List{
		command: &command{
			name: "list",
			desc: "List files in a Mesos sandbox",
		},
	}
}

// TODO: The HTTP operator API does not provide a way to pull down the sandbox
// paths of tasks that are not currently running. Once I work around this I will implement
// a way to search across all agents from a root path like /<agentid>/<framework>/<executor>/<containerid>/...
func (l *List) Action() {
	caller, err := NewAgentCaller(l.config().Profile(), *l.AgentID)
	failOnErr(err)
	resp, err := caller.CallAgent(agent.ListFiles(*l.Path))
	failOnErr(err)
	table := uitable.New()
	table.AddRow("UID", "GID", "MODE", "MODIFIED", "SIZE", "PATH")
	for _, info := range resp.ListFiles.FileInfos {
		path := info.Path
		if *l.Relative {
			split := strings.Split(path, "/")
			if len(split) > 0 {
				path = split[len(split)-1]
			}
		}
		table.AddRow(*info.UID, *info.GID, *info.Mode, "TODO", fmt.Sprintf("%d", info.Size()), path)
	}
	fmt.Println(table)
}

func (l *List) Init() func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID PATH"
		cmd.Action = l.Action
		l.AgentID = cmd.StringArg("ID", "", "AgentID")
		l.Path = cmd.StringArg("PATH", "", "path to list")
		l.Relative = cmd.BoolOpt("relative", true, "Display the relative sandbox path")
	}
}

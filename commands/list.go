package commands

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	agent "github.com/mesos/mesos-go/agent/calls"
	"github.com/vektorlab/mesos-cli/config"
	"strings"
)

type List struct {
	Path     *string
	AgentID  *string
	Relative *bool
	Hostname *string
	profile  Profile
}

func (_ List) Name() string { return "list" }
func (_ List) Desc() string { return "List files in a Mesos sandbox" }
func (l *List) SetProfile(p Profile) {
	l.profile = func() *config.Profile {
		profile := p()
		if *l.Hostname != "" {
			profile = profile.With(
				config.Master(*l.Hostname),
			)
		}
		return profile
	}
}

// TODO: The HTTP operator API does not provide a way to pull down the sandbox
// paths of tasks that are not currently running. Once I work around this I will implement
// a way to search across all agents from a root path like /<agentid>/<framework>/<executor>/<containerid>/...
func (l *List) Init() func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID PATH"
		l.AgentID = cmd.StringArg("ID", "", "AgentID")
		l.Path = cmd.StringArg("PATH", "", "path to list")
		l.Relative = cmd.BoolOpt("relative", true, "Display the relative sandbox path")
		l.Hostname = cmd.StringOpt("master", "", "Mesos master hostname")
		cmd.Action = func() {

			caller, err := NewAgentCaller(l.profile(), *l.AgentID)
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
	}
}

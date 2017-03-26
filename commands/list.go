package commands

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/filter"
	"strings"
)

type List struct{}

func (_ List) Name() string { return "list" }
func (_ List) Desc() string { return "List files in a Mesos sandbox" }

// TODO: The HTTP operator API does not provide a way to pull down the sandbox
// paths of tasks that are not currently running. Once I work around this I will implement
// a way to search across all agents from a root path like /<agentid>/<framework>/<executor>/<containerid>/...
func (_ List) Init(profile config.ProfileFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] ID PATH"
		var (
			agentID  = cmd.StringArg("ID", "", "AgentID")
			filePath = cmd.StringArg("PATH", "", "path to list")
			relative = cmd.BoolOpt("relative", true, "Display the relative sandbox path")
			hostname = cmd.StringOpt("master", "", "Mesos master")
		)
		cmd.Action = func() {
			msgs, err := filter.Find(profile().With(config.Master(*hostname)), filter.Criteria{Target: filter.FILES, AgentID: *agentID, FilePath: *filePath})
			failOnErr(err)
			table := uitable.New()
			table.AddRow("UID", "GID", "MODE", "MODIFIED", "SIZE", "PATH")
			for _, info := range filter.AsFileInfos(msgs) {
				path := info.Path
				if *relative {
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

package commands

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/client"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/filters"
	"github.com/vektorlab/mesos-cli/options"
)

type List struct {
	*command
}

func NewList() Command {
	return List{
		&command{"list", "List files in a Mesos sandbox"},
	}
}

func (_ List) Init(cfg config.CfgFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [ID] [PATH]"
		var (
			defaults = config.DefaultProfile()
			hostname = cmd.StringOpt("hostname", defaults.Master, "Mesos Master")
			agentID  = cmd.StringArg("ID", "", "AgentID")
			path     = cmd.StringArg("PATH", "", "Directory Path")
			// TODO: Reduce duplicated code
			fuzzy  = cmd.BoolOpt("fuzzy", true, "fuzzy match")
			name   = cmd.StringOpt("name", "", "filter by task name")
			states = options.NewStates()
		)
		cmd.VarOpt("state", &states, "Filter by task state")
		cmd.Action = func() {
			client := client.New(
				cfg().Profile(
					config.WithMaster(*hostname),
				),
			)
			if *agentID == "" {
				result, err := client.SearchByTask(
					&filters.TaskFilter{
						Name:   *name,
						Fuzzy:  *fuzzy,
						States: states,
					},
				)
				failOnErr(err)
				// TODO
				fmt.Println(result)
				fmt.Println(path)
			}
			/*
				filters, err := NewTaskFilters(&TaskFilterOptions{
					All:         *all,
					FrameworkID: *frameworkID,
					Fuzzy:       *fuzzy,
					ID:          *id,
					Name:        *name,
					States:      *state,
				})
				failOnErr(err)
				// First attempt to resolve the task by ID
				task, err := client.Task(filters...)
				failOnErr(err)
				// Attempt to get the full agent state
				agent, err := client.Agent(AgentFilterId(task.GetAgentId().GetValue()))
				failOnErr(err)
				// Lookup executor information in agent state
				client = &Client{Hostname: FQDN(agent)}
				executor, err := client.Executor(ExecutorFilterId(task.GetExecutorId().GetValue()))
				failOnErr(err)
				fmt.Println(executor)
				files, err := client.Files()
				failOnErr(err)
				table := uitable.New()
				table.AddRow("UID", "GID", "MODE", "MODIFIED", "SIZE", "PATH")
				for _, file := range files {
					path := file.Relative()
					if *absolute {
						path = file.Path
					}
					table.AddRow(file.UID, file.GID, file.Mode, file.Modified().String(), fmt.Sprintf("%d", file.Size), path)
				}
				fmt.Println(table)
			*/
		}
	}
}

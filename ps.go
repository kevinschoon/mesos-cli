package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	mesos "github.com/mesos/mesos-go/mesosproto"
)

func ps(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS]"
	defaults := DefaultProfile()
	var (
		master = cmd.StringOpt("master", defaults.Master, "Mesos Master")
		limit  = cmd.IntArg("limit", 100, "maximum number of tasks to return per request")
		max    = cmd.IntArg("max", 250, "maximum number of tasks to list")
		order  = cmd.StringArg("order", "desc", "accending or decending sort order [asc|desc]")
	)
	cmd.Action = func() {
		profile := config.Profile(WithMaster(*master))
		tasks := make(chan *mesos.Task)
		client := &Client{handler: DefaultHandler{hostname: profile.Master}}
		paginator := &TaskPaginator{
			limit: *limit,
			max:   *max,
			order: *order,
			tasks: tasks,
		}
		go func() {
			failOnErr(client.Paginate(paginator))
		}()
		// TODO: make pretty
		for task := range tasks {
			fmt.Printf("%s   %s   %s\n", *task.TaskId.Value, *task.FrameworkId.Value, task.State.String())
		}
	}
}

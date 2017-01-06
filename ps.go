package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"net/url"
	"regexp"
)

// TaskPaginator paginates requests from /master/tasks
// TODO: It appears that we should be able to
// unmarshal responses from the non-scheduler API
// with protobuf code generated from /include/mesos/master/master.proto
// however I was unsuccesful after several attempts. Additionally
// we would ideally want to match the vendored mesos-go protobufs.
type TaskPaginator struct {
	tasks     chan *mesos.Task
	processed int    // Total number of tasks proessed
	count     int    // Total number of matching tasks
	limit     int    // Limit of tasks per request
	max       int    // Maximum amount of matching tasks
	order     string // Order of tasks
}

func (t *TaskPaginator) Close() { close(t.tasks) }

func (t *TaskPaginator) Next(c *Client, f ...Filter) error {
	u := &url.URL{
		Path: "/master/tasks",
		RawQuery: url.Values{
			"offset": []string{fmt.Sprintf("%d", t.processed)},
			"limit":  []string{fmt.Sprintf("%d", t.limit)},
		}.Encode(),
	}
	tasks := struct {
		Tasks []struct {
			ID          *string          `json:"id"`
			Name        *string          `json:"name"`
			FrameworkID *string          `json:"framework_id"`
			AgentID     *string          `json:"slave_id"`
			State       *mesos.TaskState `json:"state"`
		} `json:"tasks"`
	}{}
	if err := c.Get(u, &tasks); err != nil {
		return err
	}
loop:
	for _, other := range tasks.Tasks {
		t.processed++
		task := &mesos.Task{
			TaskId:      &mesos.TaskID{Value: other.ID},
			Name:        other.Name,
			FrameworkId: &mesos.FrameworkID{Value: other.FrameworkID},
			SlaveId:     &mesos.SlaveID{Value: other.AgentID},
			State:       other.State,
		}
		for _, filter := range f {
			// If any filter does not match discard the task
			if !filter(task) {
				continue loop
			}
		}
		t.count++
		// Check if we've exceeded the maximum tasks
		if t.count >= t.max {
			return ErrMaxExceeded
		}
		t.tasks <- task
	}
	// If the response is smaller than the limit
	// we have finished this request
	if len(tasks.Tasks) < t.limit {
		return ErrEndPagination
	}
	return nil
}

func ps(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS]"
	defaults := DefaultProfile()
	var (
		master = cmd.StringOpt("master", defaults.Master, "Mesos Master")
		limit  = cmd.IntOpt("limit", 100, "maximum number of tasks to return per request")
		max    = cmd.IntOpt("max", 250, "maximum number of tasks to list")
		order  = cmd.StringArg("order", "desc", "accending or decending sort order [asc|desc]")
		name   = cmd.StringOpt("name", "", "regular expression to match the TaskId")
		state  = cmd.StringOpt("state", "", "regular expression to match the State")
	)
	cmd.Action = func() {
		filters := []Filter{}
		if *name != "" {
			exp, err := regexp.Compile(*name)
			failOnErr(err)
			filters = append(filters, func(t interface{}) bool {
				task := t.(*mesos.Task)
				return exp.MatchString(*task.Name)
			})
		}
		if *state != "" {
			exp, err := regexp.Compile(*state)
			failOnErr(err)
			filters = append(filters, func(t interface{}) bool {
				task := t.(*mesos.Task)
				return exp.MatchString(task.State.String())
			})
		}
		tasks := make(chan *mesos.Task)
		client := &Client{
			Hostname: config.Profile(WithMaster(*master)).Master,
		}
		paginator := &TaskPaginator{
			limit: *limit,
			max:   *max,
			order: *order,
			tasks: tasks,
		}
		go func() {
			failOnErr(Paginate(client, paginator, filters...))
		}()
		// TODO: make pretty
		for task := range tasks {
			fmt.Printf("%s   %s   %s\n", *task.TaskId.Value, *task.FrameworkId.Value, task.State.String())
		}
	}
}

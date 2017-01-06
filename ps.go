package main

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"net/url"
	"regexp"
)

// TODO: It appears that we should be able to
// unmarshal responses from the non-scheduler API
// with protobuf code generated from /include/mesos/master/master.proto
// however I was unsuccesful after several attempts.
type taskInfo struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	FrameworkID string           `json:"framework_id"`
	AgentID     string           `json:"slave_id"`
	State       *mesos.TaskState `json:"state"`
	Resources   struct {
		CPU  float64 `json:"cpus"`
		Mem  float64 `json:"mem"`
		Disk float64 `json:"disk"`
		GPUs float64 `json:"gpus"`
	}
	//Resources map[string]float64 `json:"resources"`
}

// TaskPaginator paginates requests from /master/tasks
type TaskPaginator struct {
	tasks     chan *taskInfo
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
		Tasks []*taskInfo `json:"tasks"`
	}{}
	if err := c.Get(u, &tasks); err != nil {
		return err
	}
loop:
	for _, task := range tasks.Tasks {
		t.processed++
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

func stateFilter(name string) Filter {
	return func(t interface{}) bool {
		task := t.(*taskInfo)
		return task.State.String() == name
	}
}

func ps(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS]"
	defaults := DefaultProfile()
	var (
		master   = cmd.StringOpt("master", defaults.Master, "Mesos Master")
		limit    = cmd.IntOpt("limit", 100, "maximum number of tasks to return per request")
		max      = cmd.IntOpt("max", 250, "maximum number of tasks to list")
		order    = cmd.StringArg("order", "desc", "accending or decending sort order [asc|desc]")
		name     = cmd.StringOpt("name", "", "regular expression to match the TaskId")
		all      = cmd.BoolOpt("a all", false, "show all tasks")
		running  = cmd.BoolOpt("r running", true, "show running tasks")
		failed   = cmd.BoolOpt("fa failed", false, "show failed tasks")
		killed   = cmd.BoolOpt("k killed", false, "show killed tasks")
		finished = cmd.BoolOpt("f finished", false, "show finished tasks")
	)
	Filters := func() []Filter {
		filters := []Filter{}
		if *name != "" {
			exp, err := regexp.Compile(*name)
			failOnErr(err)
			filters = append(filters, func(t interface{}) bool {
				task := t.(*taskInfo)
				return exp.MatchString(task.Name)
			})
		}
		if *all {
			filters = append(filters, func(t interface{}) bool { return true })
			return filters
		}
		if *running {
			filters = append(filters, stateFilter("TASK_RUNNING"))
		}
		if *failed {
			filters = append(filters, stateFilter("TASK_FAILED"))
		}
		if *killed {
			filters = append(filters, stateFilter("TASK_KILLED"))
		}
		if *finished {
			filters = append(filters, stateFilter("TASK_FINISHED"))
		}
		return filters
	}
	cmd.Action = func() {
		tasks := make(chan *taskInfo)
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
			failOnErr(Paginate(client, paginator, Filters()...))
		}()
		table := uitable.New()
		table.AddRow("ID", "FRAMEWORK", "STATE", "CPUS", "MEM", "GPUS", "DISK")
		for task := range tasks {
			table.AddRow(task.ID, truncStr(task.FrameworkID, 8), task.State.String(), task.Resources.CPU, task.Resources.Mem, task.Resources.GPUs, task.Resources.Disk)
		}
		fmt.Println(table)
	}
}

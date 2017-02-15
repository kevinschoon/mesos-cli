package client

import (
	"bytes"
	"fmt"
	"github.com/mesos/mesos-go"
	"github.com/mesos/mesos-go/httpcli"
	"github.com/mesos/mesos-go/httpcli/operator"
	master "github.com/mesos/mesos-go/master/calls"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/filters"
	"io/ioutil"
	"net/http"
	"net/url"
)

func RequestLogger(req *http.Request) {
	buf, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	fmt.Println(string(buf))
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
}

// Result is returned from a Resolve* method
type Result struct {
	Agent      *mesos.AgentInfo
	ExecutorID *mesos.ExecutorID
}

type Client struct {
	caller operator.Caller
}

func New(profile *config.Profile) *Client {
	endpoint := url.URL{
		Scheme: profile.Scheme,
		Host:   profile.Master,
		Path:   config.OperatorAPIPath,
	}
	return &Client{
		caller: operator.NewCaller(
			httpcli.New(
				httpcli.Endpoint(endpoint.String()),
				httpcli.RequestOptions(RequestLogger),
			),
		)}
}

func (c *Client) Agents() ([]*mesos.AgentInfo, error) {
	resp, err := c.caller.CallMaster(master.GetAgents())
	if err != nil {
		return nil, err
	}
	agents := []*mesos.AgentInfo{}
	for _, agnt := range resp.GetGetAgents().GetAgents() {
		agents = append(agents, agnt.GetAgentInfo())
	}
	return agents, nil
}

// Tasks returns ALL tasks in a Mesos cluster
// NOTE: The old /master/tasks endpoint allowed for pagination
// but the new operator API does not. This call will
// be slower the larger the cluster. Only have tested this
// against a cluster with ~20,000 tasks.
func (c *Client) Tasks() ([]*mesos.Task, error) {
	resp, err := c.caller.CallMaster(master.GetTasks())
	if err != nil {
		return nil, err
	}
	tr := resp.GetGetTasks()
	tasks := []*mesos.Task{}
	for _, task := range tr.GetTasks() {
		tasks = append(tasks, task)
	}
	for _, task := range tr.GetOrphanTasks() {
		tasks = append(tasks, task)
	}
	for _, task := range tr.GetPendingTasks() {
		tasks = append(tasks, task)
	}
	for _, task := range tr.GetCompletedTasks() {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// ResolveTask attempts to return a Result based on
// all of the tasks in a Mesos cluster. If more than one task is
// is found we will return ErrTooManyResults which indicates the caller
// needs to be more specific with their search criteria.
func (c *Client) SearchByTask(filter *filters.TaskFilter) (*Result, error) {
	tasks, err := c.Tasks()
	if err != nil {
		return nil, err
	}
	task, err := filter.Find(tasks)
	if err != nil {
		return nil, err
	}
	agents, err := c.Agents()
	if err != nil {
		return nil, err
	}
	info, err := filters.AgentFilter{ID: &task.AgentID}.Find(agents)
	if err != nil {
		return nil, err
	}
	return &Result{
		Agent:      info,
		ExecutorID: task.ExecutorID,
	}, nil
}

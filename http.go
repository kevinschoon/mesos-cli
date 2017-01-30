package main

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	mesos "github.com/vektorlab/mesos/v1"
	agent "github.com/vektorlab/mesos/v1/agent"
	master "github.com/vektorlab/mesos/v1/master"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	ErrMaxExceeded      = errors.New("max exceeded")
	ErrEndPagination    = errors.New("no more items to paginate")
	ErrTooManyTasks     = errors.New("too many tasks")
	ErrTaskNotFound     = errors.New("task not found")
	ErrTooManyAgents    = errors.New("too many agents")
	ErrAgentNotFound    = errors.New("agent not found")
	ErrTooManyExecutors = errors.New("too many executors")
	ErrExecutorNotFound = errors.New("executor not found")
	ErrTooManyFiles     = errors.New("too many files")
	ErrFileNotFound     = errors.New("file not foun")
)

// Client implements a simple HTTP client for
// interacting with Mesos API endpoints.
type Client struct {
	Hostname string
}

func (c Client) handle(u *url.URL, method string, body io.ReadCloser) (*http.Response, error) {
	u.Scheme = "http"
	u.Host = c.Hostname
	resp, err := http.DefaultClient.Do(&http.Request{
		Method:     method,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: body,
		Host: c.Hostname,
		URL:  u,
	})
	if err != nil {
		return nil, err
	}
	log.V(1).Infof("%s[%d] %s", method, resp.StatusCode, u.String())
	return resp, nil
}

func (c Client) GetBytes(u *url.URL) ([]byte, error) {
	resp, err := c.handle(u, "GET", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c Client) Get(u *url.URL, o interface{}) error {
	resp, err := c.handle(u, "GET", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, o)
}

func (c Client) call(call, o proto.Message, u *url.URL) error {
	buf := bytes.NewBuffer(nil)
	m := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
	}
	err := m.Marshal(buf, call)
	if err != nil {
		return err
	}
	//log.V(1).Infof("REQUEST: %s", buf.String())
	resp, err := c.handle(u, "POST", ioutil.NopCloser(buf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return jsonpb.Unmarshal(resp.Body, o)

}
func (c Client) callMaster(call *master.Call) (*master.Response, error) {
	msg := &master.Response{}
	return msg, c.call(call, msg, &url.URL{Path: "/master/api/v1"})

}
func (c Client) callAgent(call *agent.Call) (*agent.Response, error) {
	msg := &agent.Response{}
	return msg, c.call(call, msg, &url.URL{Path: "/slave(1)/api/v1"})
}

func (c *Client) PaginateFile(pag *FilePaginator) (err error) {
	defer pag.Close()
	for err == nil {
		err = pag.Next(c)
	}
	switch err {
	case ErrMaxExceeded:
		return nil
	case ErrEndPagination:
		return nil
	}
	return err
}

// TODO: Test against huge Task history.
// New Mesos v1 API does not seem to support
// pagination for tasks. Have only tested with
// against a cluster with Task history of around
// ~25k active/completed.
func (c Client) Tasks(f ...TaskFilter) ([]*mesos.Task, error) {
	resp, err := c.callMaster(NewMasterCall("GET_TASKS"))
	if err != nil {
		return nil, err
	}
	tasks := []*mesos.Task{}
	for _, task := range resp.GetGetTasks().GetTasks() {
		if FilterTask(task, f, false) {
			tasks = append(tasks, task)
		}
	}
	for _, task := range resp.GetGetTasks().GetCompletedTasks() {
		if FilterTask(task, f, false) {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (c Client) Task(f ...TaskFilter) (*mesos.Task, error) {
	tasks, err := c.Tasks(f...)
	if err != nil {
		return nil, err
	}
	if len(tasks) > 1 {
		return nil, ErrTooManyTasks
	}
	if len(tasks) == 0 {
		return nil, ErrTaskNotFound
	}
	return tasks[0], nil
}

func (c Client) Agents(f ...AgentFilter) ([]*mesos.AgentInfo, error) {
	resp, err := c.callMaster(NewMasterCall("GET_AGENTS"))
	if err != nil {
		return nil, err
	}
	infos := []*mesos.AgentInfo{}
	for _, agent := range resp.GetGetAgents().GetAgents() {
		info := agent.GetAgentInfo()
		if FilterAgent(info, f, false) {
			infos = append(infos, info)
		}
	}
	return infos, nil
}

func (c Client) Agent(f ...AgentFilter) (*mesos.AgentInfo, error) {
	agents, err := c.Agents()
	if err != nil {
		return nil, err
	}
	if len(agents) > 1 {
		return nil, ErrTooManyAgents
	}
	if len(agents) == 0 {
		return nil, ErrAgentNotFound
	}
	return agents[0], nil
}

// Calls must be made against agent API
func (c Client) Executors(f ...ExecutorFilter) ([]*mesos.ExecutorInfo, error) {
	resp, err := c.callAgent(NewAgentCall("GET_EXECUTORS"))
	if err != nil {
		return nil, err
	}
	infos := []*mesos.ExecutorInfo{}
	for _, executor := range resp.GetGetExecutors().GetExecutors() {
		info := executor.GetExecutorInfo()
		if FilterExecutor(info, f, false) {
			infos = append(infos, info)
		}
	}
	return infos, nil
}

func (c Client) Executor(f ...ExecutorFilter) (*mesos.ExecutorInfo, error) {
	executors, err := c.Executors(f...)
	if err != nil {
		return nil, err
	}
	if len(executors) > 1 {
		return nil, ErrTooManyExecutors
	}
	if len(executors) == 0 {
		return nil, ErrExecutorNotFound
	}
	return executors[0], nil
}

func (c Client) Files(path string, f ...FileFilter) ([]*mesos.FileInfo, error) {
	call := NewAgentCall("LIST_FILES")
	call.ListFiles.Path = proto.String(path)
	resp, err := c.callAgent(call)
	if err != nil {
		return nil, err
	}
	infos := []*mesos.FileInfo{}
	for _, info := range resp.GetListFiles().GetFileInfos() {
		if FilterFile(info, f, false) {
			infos = append(infos, info)
		}
	}
	return infos, nil
}

func (c Client) File(path string, f ...FileFilter) (*mesos.FileInfo, error) {
	files, err := c.Files(path, f...)
	if err != nil {
		return nil, err
	}
	if len(files) > 1 {
		return nil, ErrTooManyFiles
	}
	if len(files) == 0 {
		return nil, ErrFileNotFound
	}
	return files[0], nil
}

func (c Client) Read(path string, offset uint64) (*ReadData, error) {
	call := NewAgentCall("READ_FILE")
	call.ReadFile.Path = proto.String(path)
	call.ReadFile.Offset = proto.Uint64(offset)
	resp, err := c.callAgent(call)
	if err != nil {
		return nil, err
	}
	return &ReadData{
		Data: resp.GetReadFile().GetData(),
		Size: resp.GetReadFile().GetSize(),
		//Offset: int64(resp.GetReadFile().GetSize()),
	}, nil
}

func NewAgentCall(name string) *agent.Call {
	ct := agent.Call_Type(agent.Call_Type_value[name])
	return &agent.Call{Type: &ct}
}

func NewMasterCall(name string) *master.Call {
	ct := master.Call_Type(master.Call_Type_value[name])
	return &master.Call{Type: &ct}
}

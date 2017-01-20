package main

import (
	"encoding/json"
	"fmt"
	log "github.com/golang/glog"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

// TaskPaginator paginates requests from /master/tasks
type TaskPaginator struct {
	tasks     chan *taskInfo
	processed int    // Total number of tasks proessed
	count     int    // Total number of matching tasks
	limit     int    // Limit of tasks per request
	max       int    // Maximum amount of matching tasks
	order     string // Order of tasks
	Any       bool   // Any match
}

func (t *TaskPaginator) Close() { close(t.tasks) }

func (t *TaskPaginator) Next(c *Client, f ...TaskFilter) error {
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
		if !FilterTask(task, f, t.Any) {
			continue loop
		}
		t.count++
		// Check if we've exceeded the maximum tasks
		// If the maximum tasks is less than zero
		// continue forever.
		if t.count >= t.max && t.max > 0 {
			t.tasks <- task // Send the last task
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

// Client implements a simple HTTP client for
// interacting with Mesos API endpoints.
type Client struct {
	Hostname string
}

func (c Client) handle(u *url.URL, method string) ([]byte, error) {
	u.Scheme = "http"
	u.Host = c.Hostname
	resp, err := http.DefaultClient.Do(&http.Request{
		Method:     method,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
		Body:       nil,
		Host:       c.Hostname,
		URL:        u,
	})
	log.V(1).Infof("%s[%d] %s", method, resp.StatusCode, u.String())
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = resp.Body.Close(); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c Client) GetBytes(u *url.URL) ([]byte, error) {
	return c.handle(u, "GET")
}

func (c Client) Get(u *url.URL, o interface{}) error {
	raw, err := c.handle(u, "GET")
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, o)
}

// Browse returns all of the files on an agent at the given path
func (c Client) Browse(path string) ([]*fileInfo, error) {
	//client := &Client{Hostname: agent.FQDN()}
	files := []*fileInfo{}
	err := c.Get(&url.URL{
		Path: "/files/browse",
		RawQuery: url.Values{
			"path": []string{path},
		}.Encode(),
	}, &files)
	if err != nil {
		return nil, err
	}
	return files, nil
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

// Attempt to monitor one or more files
func (c *Client) ReadFiles(w io.Writer, targets ...*fileInfo) error {
	var (
		wg  sync.WaitGroup
		err error
	)
	for _, target := range targets {
		wg.Add(2)
		fp := &FilePaginator{
			data:   make(chan *fileData),
			cancel: make(chan bool),
			path:   target.Path,
			tail:   true,
		}
		err := fp.init(c)
		if err != nil {
			return err
		}
		// TODO: Need to bubble these errors back properly
		go func() {
			defer wg.Done()
			err = c.PaginateFile(fp)
		}()
		go func() {
			defer wg.Done()
			err = Pailer(fp.data, fp.cancel, 0, w)
		}()
	}
	wg.Wait()
	return err
}

func (c *Client) PaginateTasks(pag *TaskPaginator, f ...TaskFilter) (err error) {
	defer pag.Close()
	for err == nil {
		err = pag.Next(c, f...)
	}
	switch err {
	case ErrMaxExceeded:
		return nil
	case ErrEndPagination:
		return nil
	}
	return err
}

// FindTask attempts to find a single task
func (c *Client) FindTask(f ...TaskFilter) (*taskInfo, error) {
	var err error
	results := []*taskInfo{}
	tasks := make(chan *taskInfo)
	paginator := &TaskPaginator{
		limit: 2000,
		max:   -1,
		order: "asc",
		tasks: tasks,
	}
	go func() {
		err = c.PaginateTasks(paginator, f...)
	}()
	for task := range tasks {
		results = append(results, task)
	}
	if err != nil {
		return nil, err
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("too many results")
	}
	if len(results) != 1 {
		return nil, fmt.Errorf("task not found")
	}
	return results[0], nil
}

// Agents will return an array of agents as reported by the master
func (c *Client) Agents() ([]*agentInfo, error) {
	agents := struct {
		Agents []*agentInfo `json:"slaves"`
	}{}
	err := c.Get(&url.URL{Path: "/master/slaves"}, &agents)
	if err != nil {
		return nil, err
	}
	return agents.Agents, nil
}

// Agent returns an agent with it's full state
func (c *Client) Agent(agentID string) (*agentState, error) {
	agents, err := c.Agents()
	if err != nil {
		return nil, err
	}
	var info *agentInfo
	for _, a := range agents {
		if a.ID == agentID {
			info = a
			break
		}
	}
	if info == nil {
		return nil, fmt.Errorf("agent not found")
	}
	agent := &agentState{}
	err = Client{Hostname: info.FQDN()}.Get(&url.URL{Path: "/slave(1)/state"}, agent)
	return agent, err
}

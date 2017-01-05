package main

import (
	"encoding/json"
	"errors"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	ErrMaxExceeded   = errors.New("max exceeded")
	ErrEndPagination = errors.New("no more items to paginate")
)

// Handler resolves some endpoint into interface
type Handler interface {
	Handle(*url.URL, interface{}) error
}

// DefaultHandler unmarshals the JSON response at the given endpoint
type DefaultHandler struct {
	hostname string
}

func (h DefaultHandler) Handle(u *url.URL, o interface{}) error {
	u.Scheme = "http"
	u.Host = h.hostname
	resp, err := http.DefaultClient.Do(&http.Request{
		Method:     "GET",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
		Body:       nil,
		Host:       h.hostname,
		URL:        u,
	})
	if err != nil {
		return err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = resp.Body.Close(); err != nil {
		return err
	}
	return json.Unmarshal(raw, o)
}

// Paginator is handles some stateful request
type Paginator interface {
	Next(Handler) error // Make the next HTTP request
	Close()             // Close any open channels
}

// TaskPaginator paginates requests from /master/tasks
// TODO: It appears that we should be able to
// unmarshal responses from the non-scheduler API
// with protobuf code generated from /include/mesos/master/master.proto
// however I was unsuccesful after several attempts. Additionally
// we would ideally want to match the vendored mesos-go protobufs.
type TaskPaginator struct {
	tasks chan *mesos.Task
	count int
	limit int
	max   int
	order string
}

func (t *TaskPaginator) Close() { close(t.tasks) }

func (t *TaskPaginator) Next(h Handler) error {
	u := &url.URL{
		Path: "/master/tasks",
		RawQuery: url.Values{
			"offset": []string{string(t.count)},
			"limit":  []string{string(t.limit)},
		}.Encode(),
	}
	tasks := struct {
		Tasks []struct {
			ID          *string          `json:"id"`
			FrameworkID *string          `json:"framework_id"`
			AgentID     *string          `json:"slave_id"`
			State       *mesos.TaskState `json:"state"`
		} `json:"tasks"`
	}{}
	if err := h.Handle(u, &tasks); err != nil {
		return err
	}
	for _, task := range tasks.Tasks {
		t.count++
		// Check if we've exceeded the maximum tasks
		if t.count >= t.max {
			return ErrMaxExceeded
		}
		t.tasks <- &mesos.Task{
			TaskId:      &mesos.TaskID{Value: task.ID},
			FrameworkId: &mesos.FrameworkID{Value: task.FrameworkID},
			SlaveId:     &mesos.SlaveID{Value: task.AgentID},
			State:       task.State,
		}
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
	handler Handler
}

func (c *Client) Paginate(pag Paginator) (err error) {
	defer pag.Close()
	for err == nil {
		err = pag.Next(c.handler)
	}
	switch err {
	case ErrMaxExceeded:
		return nil
	case ErrEndPagination:
		return nil
	}
	return err
}

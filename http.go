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

// Filter filters the results of a
// paginator based on some criteria.
type Filter func(interface{}) bool

// Paginator is handles some stateful request
type Paginator interface {
	Next(Handler, ...Filter) error // Make the next HTTP request
	Close()                        // Close any open channels
}

// Client implements a simple HTTP client for
// interacting with Mesos API endpoints.
type Client struct {
	handler Handler
}

func (c *Client) Agents() ([]*mesos.SlaveInfo, error) {
	agents := []*mesos.SlaveInfo{}
	agnts := struct {
		Agents []struct {
			ID       *string `json:"id"`
			Hostname *string `json:"hostname"`
		} `json:"slaves"`
	}{}
	if err := c.handler.Handle(&url.URL{Path: "/master/slaves"}, &agnts); err != nil {
		return nil, err
	}
	for _, agent := range agnts.Agents {
		agents = append(agents, &mesos.SlaveInfo{
			Id:       &mesos.SlaveID{Value: agent.ID},
			Hostname: agent.Hostname,
		})
	}
	return agents, nil
}

func (c *Client) Paginate(pag Paginator, f ...Filter) (err error) {
	defer pag.Close()
	for err == nil {
		err = pag.Next(c.handler, f...)
	}
	switch err {
	case ErrMaxExceeded:
		return nil
	case ErrEndPagination:
		return nil
	}
	return err
}

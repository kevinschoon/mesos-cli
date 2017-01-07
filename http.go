package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	ErrMaxExceeded   = errors.New("max exceeded")
	ErrEndPagination = errors.New("no more items to paginate")
)

// Filter filters the results of a
// paginator based on some criteria.
type Filter func(interface{}) bool

// Paginator handles some stateful request
type Paginator interface {
	Next(*Client, ...Filter) error // Make the next HTTP request
	Close()                        // Close any open channels
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

func Paginate(client *Client, pag Paginator, f ...Filter) (err error) {
	defer pag.Close()
	for err == nil {
		err = pag.Next(client, f...)
	}
	switch err {
	case ErrMaxExceeded:
		return nil
	case ErrEndPagination:
		return nil
	}
	return err
}

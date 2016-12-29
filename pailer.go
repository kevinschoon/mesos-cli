package main

import (
	"bufio"
	"fmt"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const PollInterval = 100 * time.Millisecond

type Pailer struct {
	Hostname string
	Path     string
	offset   int
}

func (p *Pailer) url() *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   p.Hostname,
		Path:   "/files/read",
		RawQuery: url.Values{
			"path":   []string{p.Path},
			"length": []string{"50000"},
			"offset": []string{fmt.Sprintf("%d", p.offset)},
		}.Encode(),
	}
}

func (p *Pailer) Monitor(target io.Writer) (err error) {
	writer := bufio.NewWriter(target)
	for {
		resp, err := http.Get(p.url().String())
		if err != nil {
			break
		}
		raw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			break
		}
		resp.Body.Close()
		data := gjson.GetBytes(raw, "data").Str
		_, err = writer.WriteString(data)
		if err != nil {
			break
		}
		err = writer.Flush()
		if err != nil {
			break
		}
		p.offset += len(data)
		time.Sleep(PollInterval)
	}
	return err
}

// LogTask monitors a task redirecting it's
// stdout and stderr log files to the operator
func LogTask(master string, status *mesos.TaskStatus) error {
	agents, err := Agents(master)
	if err != nil {
		return err
	}
	hostname, ok := agents[*status.SlaveId.Value]
	if !ok {
		return fmt.Errorf("Cannot find agent host")
	}
	logDir, err := LogDir(hostname, *status.ExecutorId.Value)
	if err != nil {
		return err
	}
	stdout := &Pailer{
		Hostname: hostname,
		Path:     fmt.Sprintf("%s/stdout", logDir),
	}
	stderr := &Pailer{
		Hostname: hostname,
		Path:     fmt.Sprintf("%s/stderr", logDir),
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() { defer wg.Done(); err = stdout.Monitor(os.Stdout) }()
	go func() { defer wg.Done(); err = stderr.Monitor(os.Stderr) }()
	wg.Wait()
	return err
}

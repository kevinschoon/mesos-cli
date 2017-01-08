package main

import (
	//"bytes"
	"fmt"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"strconv"
	"strings"
	"time"
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
}

type agentInfo struct {
	ID               string            `json:"id"`
	Pid              string            `json:"pid"`
	Hostname         string            `json:"hostname"`
	RegisteredTime   float64           `json:"registered_time"`
	ReRegisteredTime float64           `json:"reregistered_time"`
	Version          string            `json:"version"`
	Flags            map[string]string `json:"flags"`
	Frameworks       []*frameworkInfo  `json:"frameworks"`
	Resources        struct {
		CPU  float64 `json:"cpus"`
		Mem  float64 `json:"mem"`
		Disk float64 `json:"disk"`
		GPUs float64 `json:"gpus"`
	} `json:"resources"`
	UsedResources struct {
		CPU  float64 `json:"cpus"`
		Mem  float64 `json:"mem"`
		Disk float64 `json:"disk"`
		GPUs float64 `json:"gpus"`
	} `json:"used_resources"`
}

func (a *agentInfo) Registered() time.Time {
	return time.Unix(int64(a.RegisteredTime), 0)
}

func (a *agentInfo) ReRegistered() time.Time {
	return time.Unix(int64(a.ReRegisteredTime), 0)
}

func (a *agentInfo) Uptime() time.Duration {
	return time.Since(a.ReRegistered())
}

// Detect port agent is listening on
func (a *agentInfo) Port() int64 {
	split := strings.Split(a.Pid, ":")
	if len(split) != 2 {
		panic(fmt.Errorf("cannot detect port"))
	}
	port, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		panic(err)
	}
	return port
}

func (a *agentInfo) FQDN() string {
	return fmt.Sprintf("%s:%d", a.Hostname, a.Port())
}

type frameworkInfo struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Executors []*executorInfo `json:"executors"`
}

type executorInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Directory string `json:"directory"`
}

type fileInfo struct {
	GID   string  `json:"gid"`
	UID   string  `json:"uid"`
	Path  string  `json:"path"`
	Mode  string  `json:"mode"`
	MTime float64 `json:"mtime"`
	Nlink int64   `json:"nlink"`
	Size  int64   `json:"size"`
}

func (f fileInfo) Modified() time.Time {
	return time.Unix(int64(f.MTime), 0)
}

// Return relative file path
func (f fileInfo) Relative() string {
	path := f.Path
	split := strings.Split(f.Path, "/")
	if len(split) > 0 {
		path = split[len(split)-1]
	}
	return path
}

type fileData struct {
	Data   string `json:"data"`
	Offset int    `json:"offset"`
}

func (f fileData) Length() int {
	return len([]byte(f.Data))
}

package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Agents returns a map of IDs to hostnames
func Agents(master string) (map[string]string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/master/slaves", master))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	agents := map[string]string{}
	for _, agent := range gjson.GetBytes(raw, "slaves").Array() {
		agents[agent.Get("id").Str] = agent.Get("hostname").Str
	}
	return agents, nil
}

// LogDir returns the directory path for following task output
func LogDir(hostname, executorId string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:5051/slave(1)/state", hostname))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	for _, framework := range gjson.GetBytes(raw, "frameworks").Array() {
		for _, executor := range framework.Get("executors").Array() {
			if executor.Get("id").Str == executorId {
				return executor.Get("directory").Str, nil
			}
		}
	}
	return "", fmt.Errorf("Unable to find log directory")
}

// Resource returns the value of a resource
func Resource(name string, resources []*mesos.Resource) float64 {
	var value float64
	for _, resource := range resources {
		if resource.GetName() == name {
			value = resource.GetScalar().GetValue()
		}
	}
	return value
}

// Check if a Mesos resource offer can satisfy the Task
func Sufficent(task *mesos.TaskInfo, offer *mesos.Offer) bool {
	for _, resource := range offer.Resources {
		value := resource.GetScalar().GetValue()
		switch resource.GetName() {
		case "cpus":
			if value < Resource("cpus", task.Resources) {
				return false
			}
		case "mem":
			if value < Resource("mem", task.Resources) {
				return false
			}
		case "disk":
			if value < Resource("disk", task.Resources) {
				return false
			}
		}
	}
	return true
}

// NewTask returns a mesos.TaskInfo with sensibly
// populated default values.
func NewTask() *mesos.TaskInfo {
	task := &mesos.TaskInfo{
		// TODO: Generate unique taskid
		TaskId: &mesos.TaskID{Value: proto.String("mesos-exec")},
		Name:   proto.String("mesos-exec"),
		Command: &mesos.CommandInfo{
			Shell: proto.Bool(false),
			User:  proto.String("root"),
		},
		Container: &mesos.ContainerInfo{
			// Default to Mesos Containerizer
			Type: mesos.ContainerInfo_MESOS.Enum(),
			// Docker specific settings
			//Docker: &mesos.ContainerInfo_DockerInfo{
			//	Network: mesos.ContainerInfo_DockerInfo_BRIDGE.Enum(),
			//},
			// Mesos settings
			Mesos: &mesos.ContainerInfo_MesosInfo{},
		},
		Resources: []*mesos.Resource{
			&mesos.Resource{
				Name: proto.String("cpus"),
				Type: mesos.Value_SCALAR.Enum(),
				Scalar: &mesos.Value_Scalar{
					Value: proto.Float64(0.1),
				},
			},
			&mesos.Resource{
				Name: proto.String("mem"),
				Type: mesos.Value_SCALAR.Enum(),
				Scalar: &mesos.Value_Scalar{
					Value: proto.Float64(128.0),
				},
			},
			&mesos.Resource{
				Name: proto.String("disk"),
				Type: mesos.Value_SCALAR.Enum(),
				Scalar: &mesos.Value_Scalar{
					Value: proto.Float64(32.0),
				},
			},
		},
	}
	return task
}

// Convenience types for cli so we may
// specify default values in one place
// as pass them to the cli parser.
type str struct {
	pt *string
}

func (s str) String() string {
	return *s.pt
}

func (s str) Set(other string) error {
	*s.pt = other
	return nil
}

type bl struct {
	pt *bool
}

func (b bl) String() string {
	if *b.pt {
		return "true"
	}
	return "false"
}

func (b bl) Set(other string) error {
	if other == "true" {
		*b.pt = true
	}
	return nil
}

type flt struct {
	pt *float64
}

func (f flt) String() string {
	return fmt.Sprintf("%.1f", *f.pt)
}

func (f flt) Set(other string) error {
	v, err := strconv.ParseFloat(other, 64)
	if err != nil {
		return err
	}
	*f.pt = v
	return nil
}

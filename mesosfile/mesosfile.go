package mesosfile

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	mesos "github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/scheduler/calls"
	"io/ioutil"
	"os"
)

// Mesosfile is a structure for describing tasks to run on Mesos
// In the future it might hold other other things like variables
// for templating, etc.
type Mesosfile []*Group

// Group is a pair of mesos.ExecutorInfo and []mesos.TaskInfo
// This exists because Mesos does not have an explicit LaunchGroup
// type. See https://github.com/apache/mesos/blob/master/include/mesos/v1/mesos.proto#L1459-L1462
type Group struct {
	Networks    []mesos.NetworkInfo `json:"networks"`
	Tasks       []*mesos.TaskInfo   `json:"tasks"`
	executor    *mesos.ExecutorInfo
	frameworkID string
}

func (g *Group) UnmarshalJSON(data []byte) error {
	defer func() { initTasks(g.Tasks) }()
	g.executor = &mesos.ExecutorInfo{
		FrameworkID: &mesos.FrameworkID{},
		ExecutorID:  mesos.ExecutorID{},
		Type:        mesos.ExecutorInfo_DEFAULT.Enum(),
		Container: &mesos.ContainerInfo{
			Type:         mesos.ContainerInfo_MESOS.Enum(),
			NetworkInfos: g.Networks,
		},
		Resources: []mesos.Resource{
			scalar("cpus", CPU),
			scalar("mem", MEMORY),
			scalar("disk", DISK),
		},
	}
	temp := struct {
		Networks []mesos.NetworkInfo `json:"networks"`
		Tasks    []json.RawMessage   `json:"tasks"`
	}{}
	if json.Unmarshal(data, &temp) == nil {
		if len(temp.Tasks) > 0 {
			for _, data := range temp.Tasks {
				task := NewTask()
				err := json.Unmarshal(data, task)
				if err != nil {
					return err
				}
				g.Tasks = append(g.Tasks, task)
			}
			g.Networks = temp.Networks
			g.executor.Container.NetworkInfos = temp.Networks
			return nil
		}
	}
	g.Tasks = []*mesos.TaskInfo{NewTask()}
	err := json.Unmarshal(data, g.Tasks[0])
	if err != nil {
		return err
	}
	return nil
}

func (g *Group) Resources() mesos.Resources {
	resources := mesos.Resources{}
	for _, task := range g.Tasks {
		for _, resource := range task.Resources {
			resources.Add(resource)
		}
	}
	// Executor resources
	if g.executor != nil {
		for _, resource := range g.executor.Resources {
			resources.Add(resource)
		}
	}
	return resources
}

// Find returns a matching TaskInfo
func (g *Group) Find(id string) *mesos.TaskInfo {
	for _, task := range g.Tasks {
		if task.TaskID.Value == id {
			return task
		}
	}
	return nil
}

func (g *Group) Reset() {
	g = g.With(Reset())
}

func (g *Group) With(opts ...Option) *Group {
	for _, opt := range opts {
		g = opt(g)
	}
	return g
}

func (g *Group) LaunchOp() mesos.Offer_Operation {
	if len(g.Tasks) > 1 {
		return calls.OpLaunchGroup(*g.executor, mesos.TaskGroupInfo{Tasks: g.Tasks})
	}
	if len(g.Tasks) == 1 {
		return calls.OpLaunch(*g.Tasks[0])
	}
	return mesos.Offer_Operation{}
}

// Load reads a Mesosfile from the given path
// if path == "-" it reads from stdin instead
func Load(path string) (Mesosfile, error) {
	var (
		raw []byte
		err error
	)
	if path == "-" {
		raw, err = ioutil.ReadAll(os.Stdin)
	} else {
		raw, err = ioutil.ReadFile(path)
	}
	if err != nil {
		return nil, err
	}
	mf := Mesosfile{}
	return mf, yaml.Unmarshal(raw, &mf)
}

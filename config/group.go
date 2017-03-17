package config

import (
	"github.com/mesos/mesos-go"
)

// Mesosfile is a container for holding pairs of Groups (pods).
// In the future we might consider allowing key/value pairs here to use
// to parameterize environment variables and labels globally.
type Mesosfile struct {
	Groups []*Group `json:"groups"`
}

// Tasks returns a flattened array of all tasks in each Group.
func (m Mesosfile) Tasks() []*mesos.TaskInfo {
	tasks := []*mesos.TaskInfo{}
	for _, group := range m.Groups {
		for _, task := range group.Tasks {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// Executors returns a flattened array of all executors in each Group.
func (m Mesosfile) Executors() []*mesos.ExecutorInfo {
	executors := []*mesos.ExecutorInfo{}
	for _, group := range m.Groups {
		executors = append(executors, group.Executor)
	}
	return executors
}

// Group is a pair of mesos.ExecutorInfo and []mesos.TaskInfo
// This exists because Mesos does not have an explicit LaunchGroup
// type. See https://github.com/apache/mesos/blob/master/include/mesos/v1/mesos.proto#L1459-L1462
type Group struct {
	Executor *mesos.ExecutorInfo `json:"executor"`
	Tasks    []*mesos.TaskInfo   `json:"tasks"`
}

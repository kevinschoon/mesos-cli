package main

import (
	"errors"
	log "github.com/golang/glog"
	mesos "github.com/vektorlab/mesos/v1"
	"strings"
)

var ErrInvalidTaskState = errors.New("Invalid Task State")

type AgentFilterOptions struct {
	All   bool
	Id    string
	Fuzzy bool
}

type AgentFilter func(*mesos.AgentInfo) bool

func NewAgentFilters(opts *AgentFilterOptions) ([]AgentFilter, error) {
	filters := []AgentFilter{}
	if opts.All {
		filters = append(filters, AgentFilterAll)
		return filters, nil
	}
	if opts.Id != "" {
		if opts.Fuzzy {
			filters = append(filters, AgentFilterIdFuzzy(opts.Id))
		} else {
			filters = append(filters, AgentFilterId(opts.Id))
		}
	}
	return filters, nil
}

func FilterAgent(agent *mesos.AgentInfo, filters []AgentFilter, any bool) bool {
	var matches int
	for _, filter := range filters {
		if filter(agent) {
			if any {
				log.V(1).Infof("Agent %s matched: true", agent.GetId().GetValue())
				return true
			}
			matches++
		}
	}
	log.V(1).Infof("Agent %s matched: %t", agent.GetId().GetValue(), matches == len(filters))
	return matches == len(filters)
}

func AgentFilterAll(*mesos.AgentInfo) bool { return true }

func AgentFilterId(id string) AgentFilter {
	return func(a *mesos.AgentInfo) bool {
		return a.GetId().GetValue() == id
	}
}

func AgentFilterIdFuzzy(id string) AgentFilter {
	return func(a *mesos.AgentInfo) bool {
		return strings.HasPrefix(a.GetId().GetValue(), id)
	}
}

type ExecutorFilterOptions struct {
	Id string
}

type ExecutorFilter func(*mesos.ExecutorInfo) bool

func NewExecutorFilters(opts *ExecutorFilterOptions) ([]ExecutorFilter, error) {
	filters := []ExecutorFilter{}
	if opts.Id != "" {
		filters = append(filters, ExecutorFilterId(opts.Id))
	}
	return filters, nil
}

func ExecutorFilterId(id string) ExecutorFilter {
	return func(e *mesos.ExecutorInfo) bool {
		return e.GetExecutorId().GetValue() == id
	}
}

func FilterExecutor(executor *mesos.ExecutorInfo, filters []ExecutorFilter, any bool) bool {
	var matches int
	for _, filter := range filters {
		if filter(executor) {
			if any {
				log.V(1).Infof("Executor %s matched: true", executor.GetExecutorId().GetValue())
				return true
			}
			matches++
		}
	}
	log.V(1).Infof("Executor %s matched: %t", executor.GetExecutorId().GetValue(), matches == len(filters))
	return matches == len(filters)
}

type FileFilter func(*mesos.FileInfo) bool

func FilterFile(file *mesos.FileInfo, filters []FileFilter, any bool) bool {
	var matches int
	for _, filter := range filters {
		if filter(file) {
			if any {
				log.V(1).Infof("File %s matched: true", file.GetPath())
				return true
			}
			matches++
		}
	}
	log.V(1).Infof("File %s matched: %t", file.GetPath(), matches == len(filters))
	return matches == len(filters)
}

func FileFilterPathRelative(path string) FileFilter {
	return func(f *mesos.FileInfo) bool {
		return Relative(f) == path
	}
}

func FileFilterPath(path string) FileFilter {
	return func(f *mesos.FileInfo) bool {
		return f.GetPath() == path
	}
}

func FileFilterPathFuzzy(path string) FileFilter {
	return func(f *mesos.FileInfo) bool {
		return strings.HasPrefix(f.GetPath(), path)
	}
}

type TaskFilterOptions struct {
	All         bool
	FrameworkID string
	Fuzzy       bool
	Name        string
	ID          string
	States      []string
}

// TaskFilter filters tasks based on some condition
type TaskFilter func(*mesos.Task) bool

func NewTaskFilters(opts *TaskFilterOptions) ([]TaskFilter, error) {
	filters := []TaskFilter{}
	if opts.All {
		filters = append(filters, TaskFilterAll)
		return filters, nil
	}
	if opts.ID != "" {
		if opts.Fuzzy {
			filters = append(filters, TaskFilterIDFuzzy(opts.ID))
		} else {
			filters = append(filters, TaskFilterID(opts.ID))
		}
	}

	if opts.Name != "" {
		if opts.Fuzzy {
			filters = append(filters, TaskFilterNameFuzzy(opts.Name))
		} else {
			filters = append(filters, TaskFilterName(opts.Name))
		}
	}

	if opts.FrameworkID != "" {
		filters = append(filters, TaskFilterFrameworkID(opts.FrameworkID))
	}

	for _, name := range opts.States {
		filter, err := TaskFilterState(name)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}

	return filters, nil
}

func FilterTask(task *mesos.Task, filters []TaskFilter, any bool) bool {
	var matches int
	for _, filter := range filters {
		if filter(task) {
			if any {
				log.V(1).Infof("Task %s matched: true", task.GetTaskId().GetValue())
				return true
			}
			matches++
		}
	}
	log.V(1).Infof("Task %s matched: %t", task.GetTaskId().GetValue(), matches == len(filters))
	return matches == len(filters)
}

func TaskFilterAll(*mesos.Task) bool { return true }

func TaskFilterID(id string) TaskFilter {
	return func(task *mesos.Task) bool {
		return task.GetTaskId().GetValue() == id
	}
}

func TaskFilterIDFuzzy(id string) TaskFilter {
	return func(task *mesos.Task) bool {
		return strings.HasPrefix(task.GetTaskId().GetValue(), id)
	}
}

func TaskFilterName(name string) TaskFilter {
	return func(task *mesos.Task) bool {
		return task.GetName() == name
	}
}

func TaskFilterNameFuzzy(name string) TaskFilter {
	return func(task *mesos.Task) bool {
		return strings.HasPrefix(task.GetName(), name)
	}
}

func TaskFilterFrameworkID(id string) TaskFilter {
	return func(task *mesos.Task) bool {
		return task.GetFrameworkId().String() == id
	}
}

func TaskFilterState(name string) (TaskFilter, error) {
	s, ok := mesos.TaskState_value[name]
	if !ok {
		return nil, ErrInvalidTaskState
	}
	state := mesos.TaskState(s)
	return func(task *mesos.Task) bool {
		return *task.GetState().Enum() == state
	}, nil
}

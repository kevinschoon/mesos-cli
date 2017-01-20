package main

import (
	"errors"
	log "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"strings"
)

var ErrInvalidTaskState = errors.New("Invalid Task State")

type TaskFilterOptions struct {
	All         bool
	FrameworkID string
	Fuzzy       bool
	Name        string
	ID          string
	States      []string
}

// TaskFilter filters tasks based on some condition
type TaskFilter func(*taskInfo) bool

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

func FilterTask(task *taskInfo, filters []TaskFilter, any bool) bool {
	var matches int
	for _, filter := range filters {
		if filter(task) {
			if any {
				log.V(1).Infof("Task %s matched: true", task.ID)
				return true
			}
			matches++
		}
	}
	log.V(1).Infof("Task %s matched: %t", task.ID, matches == len(filters))
	return matches == len(filters)
}

func TaskFilterAll(*taskInfo) bool { return true }

func TaskFilterID(id string) TaskFilter {
	return func(task *taskInfo) bool {
		return task.ID == id
	}
}

func TaskFilterIDFuzzy(id string) TaskFilter {
	return func(task *taskInfo) bool {
		return strings.HasPrefix(task.ID, id)
	}
}

func TaskFilterName(name string) TaskFilter {
	return func(task *taskInfo) bool {
		return task.Name == name
	}
}

func TaskFilterNameFuzzy(name string) TaskFilter {
	return func(task *taskInfo) bool {
		return strings.HasPrefix(task.Name, name)
	}
}

func TaskFilterFrameworkID(id string) TaskFilter {
	return func(task *taskInfo) bool {
		return task.FrameworkID == id
	}
}

func TaskFilterState(name string) (TaskFilter, error) {
	s, ok := mesos.TaskState_value[name]
	if !ok {
		return nil, ErrInvalidTaskState
	}
	state := mesos.TaskState(s)
	return func(task *taskInfo) bool {
		return task.State == state
	}, nil
}

package main

import (
	mesos "github.com/mesos/mesos-go/mesosproto"
	"regexp"
)

type TaskFilter func(*taskInfo) bool

func MatchTaskAll(*taskInfo) bool { return true }

func MatchTaskID(name string) TaskFilter {
	return func(task *taskInfo) bool {
		return task.Name == name
	}
}

func MatchTaskIDFuzzy(name string) (TaskFilter, error) {
	expr, err := regexp.Compile(name)
	if err != nil {
		return nil, err
	}
	return func(task *taskInfo) bool {
		return expr.MatchString(task.ID)
	}, nil
}

func MatchTaskState(state mesos.TaskState) TaskFilter {
	return func(task *taskInfo) bool {
		return *task.State == state
	}
}

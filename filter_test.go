package main

import (
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaskFilter(t *testing.T) {
	tasks := []*taskInfo{
		&taskInfo{
			ID:    "T1",
			Name:  "T1 Task",
			State: mesos.TaskState_TASK_FAILED,
		},
		&taskInfo{
			ID:    "T2",
			Name:  "T2 Task",
			State: mesos.TaskState_TASK_RUNNING,
		},
		&taskInfo{
			ID:    "T3",
			Name:  "T3 Task",
			State: mesos.TaskState_TASK_RUNNING,
		},
		&taskInfo{
			ID:    "T3-1",
			Name:  "T3 Task (1)",
			State: mesos.TaskState_TASK_RUNNING,
		},
	}
	filters, err := NewTaskFilters(&TaskFilterOptions{
		States: []string{"nope"},
	})
	assert.Error(t, err)
	filters, _ = NewTaskFilters(&TaskFilterOptions{
		All:    true,
		ID:     "T1",
		States: []string{"TASK_RUNNING"},
	})
	assert.Len(t, filters, 1)
	assert.True(t, FilterTask(tasks[0], filters, false))
	filters, _ = NewTaskFilters(&TaskFilterOptions{
		ID:     "T1",
		States: []string{"TASK_RUNNING"},
	})
	assert.Len(t, filters, 2)
	assert.False(t, FilterTask(tasks[0], filters, false))
	assert.True(t, FilterTask(tasks[0], filters, true))
	filters, _ = NewTaskFilters(&TaskFilterOptions{
		ID:     "T2",
		Name:   "T2 Task",
		States: []string{"TASK_RUNNING"},
	})
	assert.True(t, FilterTask(tasks[1], filters, false))
	filters, _ = NewTaskFilters(&TaskFilterOptions{
		ID:    "T3",
		Name:  "T3",
		Fuzzy: true,
	})
	assert.False(t, FilterTask(tasks[0], filters, true))
	assert.False(t, FilterTask(tasks[1], filters, true))
	assert.True(t, FilterTask(tasks[2], filters, true))
	assert.True(t, FilterTask(tasks[3], filters, true))
}

package state

/*
import (
	mesos "github.com/mesos/mesos-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAgent(t *testing.T) {
	assert.True(t, AgentFilter{
		ID: &mesos.AgentID{
			Value: "agent-1234",
		},
	}.Match(&mesos.AgentInfo{
		ID: &mesos.AgentID{
			Value: "agent-1234",
		},
	},
	))
}

func TestExecutor(t *testing.T) {
	assert.True(t, ExecutorFilter{
		ExecutorID: &mesos.ExecutorID{
			Value: "exec-1234",
		},
	}.Match(&mesos.ExecutorInfo{
		ExecutorID: mesos.ExecutorID{
			Value: "exec-1234",
		},
	},
	))
}

func TestFile(t *testing.T) {
	assert.True(t, FileFilter{
		Path: "/fuu",
	}.Match(&mesos.FileInfo{
		Path: "/fuu",
	}))
	assert.True(t, FileFilter{
		RelativePath: "bar",
	}.Match(&mesos.FileInfo{
		Path: "/fuu/bar",
	}))
}

func TestTask(t *testing.T) {
	assert.True(t, TaskFilter{
		FrameworkID: &mesos.FrameworkID{Value: "framework-1234"},
		TaskID:      &mesos.TaskID{Value: "task-1234"},
	}.Match(
		&mesos.Task{
			FrameworkID: mesos.FrameworkID{Value: "framework-1234"},
			TaskID:      mesos.TaskID{Value: "task-1234"},
		},
	))
	assert.False(t, TaskFilter{
		FrameworkID: &mesos.FrameworkID{Value: "framework-1234"},
		TaskID:      &mesos.TaskID{Value: "task-1234"},
	}.Match(
		&mesos.Task{
			FrameworkID: mesos.FrameworkID{Value: "framework-1234-1"},
			TaskID:      mesos.TaskID{Value: "task-1234"},
		},
	))
	assert.True(t, TaskFilter{
		FrameworkID: &mesos.FrameworkID{Value: "framework-1234"},
		TaskID:      &mesos.TaskID{Value: "task-1234"},
		States:      []*mesos.TaskState{mesos.TASK_RUNNING.Enum()},
	}.Match(
		&mesos.Task{
			FrameworkID: mesos.FrameworkID{Value: "framework-1234"},
			TaskID:      mesos.TaskID{Value: "task-1234"},
			State:       mesos.TASK_RUNNING.Enum(),
		},
	))
	assert.False(t, TaskFilter{
		FrameworkID: &mesos.FrameworkID{Value: "framework-1234"},
		TaskID:      &mesos.TaskID{Value: "task-1234"},
		States:      []*mesos.TaskState{mesos.TASK_FAILED.Enum()},
	}.Match(
		&mesos.Task{
			FrameworkID: mesos.FrameworkID{Value: "framework-1234"},
			TaskID:      mesos.TaskID{Value: "task-1234"},
			State:       mesos.TASK_RUNNING.Enum(),
		},
	))
}
*/

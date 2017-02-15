package filters

import (
	"errors"
	"github.com/mesos/mesos-go"
	"strings"
)

var (
	ErrTooManyResults = errors.New("Too many results")
)

type AgentFilter struct {
	ID *mesos.AgentID
}

func (f AgentFilter) Find(infos []*mesos.AgentInfo) (*mesos.AgentInfo, error) {
	var match *mesos.AgentInfo
	for _, info := range infos {
		if f.Match(info) {
			if match != nil {
				return nil, ErrTooManyResults
			}
			match = info
		}
	}
	return match, nil
}

func (f AgentFilter) Match(info *mesos.AgentInfo) bool {
	return f.ID.Equal(info.ID)
}

type ExecutorFilter struct {
	ExecutorID *mesos.ExecutorID
}

func (f ExecutorFilter) Match(info *mesos.ExecutorInfo) bool {
	return f.ExecutorID.Equal(&info.ExecutorID)
}

type FileFilter struct {
	Path         string
	RelativePath string
}

func (f FileFilter) Match(info *mesos.FileInfo) bool {
	if f.Path != "" {
		if info.Path != f.Path {
			return false
		}
	}
	if f.RelativePath != "" {
		path := info.Path
		split := strings.Split(info.Path, "/")
		if len(split) > 0 {
			path = split[len(split)-1]
		}
		if path != f.RelativePath {
			return false
		}
	}
	return true
}

type TaskFilter struct {
	FrameworkID *mesos.FrameworkID
	Fuzzy       bool
	Name        string
	TaskID      *mesos.TaskID
	States      []*mesos.TaskState
}

// Find attempts to resolve a single mesos.Task. If more than one
// result is found we return ErrTooManyResults.
func (f TaskFilter) Find(tasks []*mesos.Task) (*mesos.Task, error) {
	var match *mesos.Task
	for _, task := range tasks {
		if f.Match(task) {
			if match != nil {
				return nil, ErrTooManyResults
			}
			match = task
		}
	}
	return match, nil
}

func (f TaskFilter) Match(task *mesos.Task) bool {
	if f.Name != "" {
		if f.Fuzzy {
			if !strings.HasPrefix(task.Name, f.Name) {
				return false
			} else {
				if task.Name != f.Name {
					return false
				}
			}
		}
	}
	if f.TaskID != nil {
		if f.Fuzzy {
			if !strings.HasPrefix(task.TaskID.Value, f.TaskID.Value) {
				return false
			}
		} else if !f.TaskID.Equal(&task.TaskID) {
			return false
		}
	}
	if f.FrameworkID != nil {
		if f.Fuzzy {
			return strings.HasPrefix(task.FrameworkID.Value, f.FrameworkID.Value)
		} else if !f.FrameworkID.Equal(&task.FrameworkID) {
			return false
		}
	}
	if len(f.States) > 0 {
		var matched bool
		for _, state := range f.States {
			if task.State.String() == state.String() {
				matched = true
			}
		}
		return matched
	}
	return true
}

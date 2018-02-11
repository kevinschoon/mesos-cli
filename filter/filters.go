package filter

import (
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/api/v1/lib"
	"strings"
)

// Filter is used to query a State and match
// a single protobuf message.
type Filter func(proto.Message) bool

func TaskNameFilter(name string, fuzzy bool) Filter {
	return func(msg proto.Message) bool {
		task, _ := AsTask(msg, nil)
		if task == nil {
			return false
		}
		if fuzzy {
			return strings.HasPrefix(task.GetName(), name)
		}
		return task.GetName() == name
	}
}

func TaskIDFilter(taskID string, fuzzy bool) Filter {
	return func(msg proto.Message) bool {
		task, _ := AsTask(msg, nil)
		if task == nil {
			return false
		}
		if fuzzy {
			return strings.HasPrefix(task.GetTaskID().Value, taskID)
		}
		return task.GetTaskID().Value == taskID
	}
}

func TaskStateFilter(states []*mesos.TaskState) Filter {
	return func(msg proto.Message) bool {
		task, _ := AsTask(msg, nil)
		if task == nil {
			return false
		}
		for _, state := range states {
			if *task.State == *state {
				return true
			}
		}
		return false
	}
}

func FrameworkIDFilter(frameworkID string, fuzzy bool) Filter {
	return func(msg proto.Message) bool {
		framework, _ := AsFramework(msg, nil)
		if framework == nil {
			return false
		}
		if fuzzy {
			return strings.HasPrefix(framework.GetID().Value, frameworkID)
		}
		return framework.GetID().Value == frameworkID
	}
}

func ExecutorIDFilter(executorID string, fuzzy bool) Filter {
	return func(msg proto.Message) bool {
		executor, _ := AsExecutor(msg, nil)
		if executor == nil {
			return false
		}
		if fuzzy {
			return strings.HasPrefix(executor.GetExecutorID().Value, executorID)
		}
		return executor.GetExecutorID().Value == executorID
	}
}

func AgentIDFilter(agentID string, fuzzy bool) Filter {
	return func(msg proto.Message) bool {
		agent, _ := AsAgent(msg, nil)
		if agent == nil {
			return false
		}
		if fuzzy {
			return strings.HasPrefix(agent.GetID().Value, agentID)
		}
		return agent.GetID().Value == agentID
	}
}

func AgentAnyFilter() Filter {
	return func(msg proto.Message) bool {
		agent, _ := AsAgent(msg, nil)
		if agent == nil {
			return false
		}
		return true
	}
}

package state

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go"
	"strings"
)

// ErrInvalidMessage is returned when we expect
// one type of proto.Message but get another.
type ErrInvalidMessage struct {
	msg proto.Message
}

func (e ErrInvalidMessage) Error() string {
	return fmt.Sprintf("Unexpected message: %s", e.msg.String())
}

func AsTask(msg proto.Message, err error) (*mesos.Task, error) {
	if err != nil {
		return nil, err
	}
	task, ok := msg.(*mesos.Task)
	if !ok {
		return nil, ErrInvalidMessage{msg}
	}
	return task, nil
}

func AsTasks(msgs []proto.Message) []*mesos.Task {
	tasks := []*mesos.Task{}
	for _, msg := range msgs {
		if task, ok := msg.(*mesos.Task); ok {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func AsFramework(msg proto.Message, err error) (*mesos.FrameworkInfo, error) {
	framework, ok := msg.(*mesos.FrameworkInfo)
	if !ok {
		return nil, ErrInvalidMessage{msg}
	}
	return framework, nil
}

func AsFrameworks(msgs []proto.Message) []*mesos.FrameworkInfo {
	frameworks := []*mesos.FrameworkInfo{}
	for _, msg := range msgs {
		if framework, ok := msg.(*mesos.FrameworkInfo); ok {
			frameworks = append(frameworks, framework)
		}
	}
	return frameworks
}

func AsExecutor(msg proto.Message, err error) (*mesos.ExecutorInfo, error) {
	executor, ok := msg.(*mesos.ExecutorInfo)
	if !ok {
		return nil, ErrInvalidMessage{msg}
	}
	return executor, nil
}

func AsExecutors(msgs []proto.Message) []*mesos.ExecutorInfo {
	executors := []*mesos.ExecutorInfo{}
	for _, msg := range msgs {
		if executor, ok := msg.(*mesos.ExecutorInfo); ok {
			executors = append(executors, executor)
		}
	}
	return executors
}

func AsAgent(msg proto.Message, err error) (*mesos.AgentInfo, error) {
	if err != nil {
		return nil, err
	}
	agent, ok := msg.(*mesos.AgentInfo)
	if !ok {
		return nil, ErrInvalidMessage{msg}
	}
	return agent, nil
}

func AsAgents(msgs []proto.Message) []*mesos.AgentInfo {
	agents := []*mesos.AgentInfo{}
	for _, msg := range msgs {
		if agent, ok := msg.(*mesos.AgentInfo); ok {
			agents = append(agents, agent)
		}
	}
	return agents
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

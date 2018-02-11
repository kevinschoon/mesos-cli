package filter

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/agent"
	"github.com/mesos/mesos-go/api/v1/lib/master"
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

func AsFileInfo(msg proto.Message, err error) (*mesos.FileInfo, error) {
	if err != nil {
		return nil, err
	}
	info, ok := msg.(*mesos.FileInfo)
	if !ok {
		return nil, ErrInvalidMessage{msg}
	}
	return info, nil
}

func AsFileInfos(msgs []proto.Message) []*mesos.FileInfo {
	infos := []*mesos.FileInfo{}
	for _, msg := range msgs {
		if info, ok := msg.(*mesos.FileInfo); ok {
			infos = append(infos, info)
		}
	}
	return infos
}

// TODO: Consider using reflect to reduce typing.
func FromAgent(resp *agent.Response, err error) (Messages, error) {
	if err != nil {
		return nil, err
	}
	messages := Messages{}
	if tasks := resp.GetTasks; tasks != nil {
		for _, task := range tasks.QueuedTasks {
			messages = append(messages, task)
		}
		for _, task := range tasks.PendingTasks {
			messages = append(messages, task)
		}
		for _, task := range tasks.LaunchedTasks {
			messages = append(messages, task)
		}
		for _, task := range tasks.CompletedTasks {
			messages = append(messages, task)
		}
		for _, task := range tasks.TerminatedTasks {
			messages = append(messages, task)
		}
	}
	if frameworks := resp.GetFrameworks; frameworks != nil {
		for _, framework := range frameworks.Frameworks {
			messages = append(messages, framework.FrameworkInfo)
		}
		for _, framework := range frameworks.CompletedFrameworks {
			messages = append(messages, framework.FrameworkInfo)
		}
	}
	if executors := resp.GetExecutors; executors != nil {
		for _, executor := range executors.Executors {
			messages = append(messages, executor.ExecutorInfo)
		}
		for _, executor := range executors.CompletedExecutors {
			messages = append(messages, executor.ExecutorInfo)
		}
	}
	if files := resp.ListFiles; files != nil {
		for _, file := range files.FileInfos {
			messages = append(messages, file)
		}
	}
	return messages, nil
}

func FromMaster(resp *master.Response, err error) (Messages, error) {
	if err != nil {
		return nil, err
	}
	messages := Messages{}
	if tasks := resp.GetTasks; tasks != nil {
		for _, task := range tasks.Tasks {
			messages = append(messages, task)
		}
		for _, task := range tasks.OrphanTasks {
			messages = append(messages, task)
		}
		for _, task := range tasks.PendingTasks {
			messages = append(messages, task)
		}
		for _, task := range tasks.CompletedTasks {
			messages = append(messages, task)
		}
	}
	if frameworks := resp.GetFrameworks; frameworks != nil {
		for _, framework := range frameworks.Frameworks {
			messages = append(messages, framework.FrameworkInfo)
		}
		for _, framework := range frameworks.CompletedFrameworks {
			messages = append(messages, framework.FrameworkInfo)
		}
		for _, framework := range frameworks.RecoveredFrameworks {
			messages = append(messages, framework)
		}
	}
	if executors := resp.GetExecutors; executors != nil {
		for _, executor := range executors.Executors {
			messages = append(messages, executor.ExecutorInfo)
		}
	}
	if agents := resp.GetAgents; agents != nil {
		for _, agent := range agents.Agents {
			messages = append(messages, agent.AgentInfo)
		}
	}
	return messages, nil
}

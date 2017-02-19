package filter

import (
	"github.com/mesos/mesos-go/agent"
	"github.com/mesos/mesos-go/master"
)

// TODO: Consider using reflect to reduce typing.

func FromAgent(resp *agent.Response) Messages {
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

	return messages
}

func FromMaster(resp *master.Response) Messages {
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
			messages = append(messages, framework)
		}
		for _, framework := range frameworks.RecoveredFrameworks {
			messages = append(messages, framework)
		}
		for _, framework := range frameworks.CompletedFrameworks {
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
	return messages
}

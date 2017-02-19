package state

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/agent"
	"github.com/mesos/mesos-go/master"
)

var (
	ErrNotFound       = errors.New("Not found")
	ErrTooManyResults = errors.New("Too many results")
)

// Filter is used to query a State and match
// a single protobuf message.
type Filter func(proto.Message) bool

// State represents the state of a Mesos cluster.
// Mesos protobufs do not immplement an explicit
// message for this because it is a combination
// of different Mesos types.
// NOTE: This could be improved with fancier indexing
// but I am unsure if that is necessary.
type State struct {
	messages []proto.Message
}

// StateFromAgent loads a State with proto.messages
// from a GET_STATE call.
// TODO: Consider using reflect to reduce typing.
func StateFromAgent(resp *agent.Response_GetState) *State {
	state := &State{messages: []proto.Message{}}
	tasks := resp.GetTasks
	for _, task := range tasks.QueuedTasks {
		state.messages = append(state.messages, task)
	}
	for _, task := range tasks.PendingTasks {
		state.messages = append(state.messages, task)
	}
	for _, task := range tasks.LaunchedTasks {
		state.messages = append(state.messages, task)
	}
	for _, task := range tasks.CompletedTasks {
		state.messages = append(state.messages, task)
	}
	for _, task := range tasks.TerminatedTasks {
		state.messages = append(state.messages, task)
	}
	frameworks := resp.GetFrameworks
	for _, framework := range frameworks.Frameworks {
		state.messages = append(state.messages, framework.FrameworkInfo)
	}
	for _, framework := range frameworks.CompletedFrameworks {
		state.messages = append(state.messages, framework.FrameworkInfo)
	}
	executors := resp.GetExecutors
	for _, executor := range executors.Executors {
		state.messages = append(state.messages, executor.ExecutorInfo)
	}
	for _, executor := range executors.CompletedExecutors {
		state.messages = append(state.messages, executor.ExecutorInfo)
	}
	return state
}

// StateFromMaster loads a State with proto.messages
// from a GET_STATE call.
// TODO: Consider using reflect to reduce typing.
func StateFromMaster(resp *master.Response_GetState) *State {
	state := &State{messages: []proto.Message{}}
	tasks := resp.GetTasks
	for _, task := range tasks.Tasks {
		state.messages = append(state.messages, task)
	}
	for _, task := range tasks.OrphanTasks {
		state.messages = append(state.messages, task)
	}
	for _, task := range tasks.PendingTasks {
		state.messages = append(state.messages, task)
	}
	for _, task := range tasks.CompletedTasks {
		state.messages = append(state.messages, task)
	}
	frameworks := resp.GetFrameworks
	for _, framework := range frameworks.Frameworks {
		state.messages = append(state.messages, framework)
	}
	executors := resp.GetExecutors
	for _, executor := range executors.Executors {
		state.messages = append(state.messages, executor.ExecutorInfo)
	}
	agents := resp.GetAgents
	for _, agent := range agents.Agents {
		state.messages = append(state.messages, agent.AgentInfo)
	}
	return state
}

func (s *State) Add(msg proto.Message) {
	s.messages = append(s.messages, msg)
}

// FindAny will return the first message
// where all filters return true. If no
// messages match we will return ErrNotFound.
func (s State) FindAny(filters ...Filter) (proto.Message, error) {
	var match bool
loop:
	for _, message := range s.messages {
		for _, filter := range filters {
			match = filter(message)
			if !match {
				continue loop
			}
		}
		if match {
			return message, nil
		}
	}
	return nil, ErrNotFound
}

// FindOne will a single message where all
// filters return true. If more than one
// message matches return ErrTooManyResults.
// If no messages match all filters we will
// return ErrNotFound.
func (s State) FindOne(filters ...Filter) (proto.Message, error) {
	var (
		position int
		match    bool
	)
loop:
	for i, message := range s.messages {
		for _, filter := range filters {
			if !filter(message) {
				continue loop
			}
		}
		// Already matched a message
		if match {
			return nil, ErrTooManyResults
		}
		// Mark this message as matched
		match = true
		position = i
	}
	if !match {
		return nil, ErrNotFound
	}
	return s.messages[position], nil
}

// FindMany will return as many messages
// in which all of the filters match true.
// If no messages match all the filters
// we will return an empty []proto.Message.
func (s State) FindMany(filters ...Filter) []proto.Message {
	matches := []proto.Message{}
loop:
	for _, message := range s.messages {
		for _, filter := range filters {
			if !filter(message) {
				continue loop
			}
		}
		matches = append(matches, message)
	}
	return matches
}

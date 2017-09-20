package filter

import (
	"errors"
	"github.com/mesanine/mesos-cli/config"
	"github.com/mesanine/mesos-cli/helper"
	"github.com/mesos/mesos-go"
	agent "github.com/mesos/mesos-go/agent/calls"
	master "github.com/mesos/mesos-go/master/calls"
)

// Target represents the expected type of result, e.g. TaskInfo, AgentInfo, etc.
type Target int

const (
	AGENTS Target = iota
	TASKS
	FILES
)

// ErrInvalidCriteria indicates there is a problem with the search critiera
var ErrInvalidCriteria = errors.New("invalid search criteria")

// Criteria represent search criteria for finding Mesos objects
type Criteria struct {
	Target      Target
	Fuzzy       bool
	AgentID     string
	FrameworkID string
	TaskID      string
	TaskName    string
	FilePath    string
	TaskStates  []*mesos.TaskState
}

// FindAgents finds AgentInfos which match the given Criteria
func FindAgents(profile *config.Profile, criteria Criteria) (Messages, error) {
	caller := helper.NewCaller(profile)
	switch {
	// Caller is attempting to resolve the agent based on a TaskID
	case criteria.TaskID != "":
		// Get all tasks running on the master
		msgs, err := FromMaster(caller.CallMaster(master.GetTasks()))
		if err != nil {
			return nil, err
		}
		// Lookup the AgentID from the task if there is a match
		task, err := AsTask(msgs.FindOne(TaskIDFilter(criteria.TaskID, criteria.Fuzzy)))
		if err != nil {
			return nil, err
		}
		// Return the results of a Find based on the resolved AgentID
		return FindAgents(profile, Criteria{AgentID: task.AgentID.Value})
	// Caller specified an AgentID, match based on GET_AGENTS call
	case criteria.AgentID != "":
		msgs, err := FromMaster(caller.CallMaster(master.GetAgents()))
		if err != nil {
			return nil, err
		}
		return msgs.FindMany(AgentIDFilter(criteria.AgentID, criteria.Fuzzy)), nil
	default:
		// Return all agents
		msgs, err := FromMaster(caller.CallMaster(master.GetAgents()))
		if err != nil {
			return nil, err
		}
		return msgs.FindMany(), nil
	}
}

// FindTasks finds Tasks which match the given Criteria
func FindTasks(profile *config.Profile, criteria Criteria) (Messages, error) {
	caller := helper.NewCaller(profile)
	filters := []Filter{TaskIDFilter(criteria.TaskID, criteria.Fuzzy)}
	if len(criteria.TaskStates) > 0 {
		filters = append(filters, TaskStateFilter(criteria.TaskStates))
	}
	switch {
	// Caller is seaching by TaskName
	case criteria.TaskName != "":
		msgs, err := FromMaster(caller.CallMaster(master.GetTasks()))
		if err != nil {
			return nil, err
		}
		filters = append(filters, TaskNameFilter(criteria.TaskName, criteria.Fuzzy))
		return msgs.FindMany(filters...), nil
	// Caller is explicitly searching for a task or tasks that match a string
	case criteria.TaskID != "":
		msgs, err := FromMaster(caller.CallMaster(master.GetTasks()))
		if err != nil {
			return nil, err
		}
		return msgs.FindMany(filters...), nil
		// Caller wants to return all tasks on a particular agent
	case criteria.AgentID != "":
		// Search for agents with the matching ID
		msgs, err := FindAgents(profile, Criteria{AgentID: criteria.AgentID})
		if err != nil {
			return nil, err
		}
		// Only one agent may match
		agnt, err := AsAgent(msgs.FindOne())
		if err != nil {
			return nil, err
		}
		// Build a new Caller for the resolved agent
		msgs, err = FromAgent(helper.NewAgentCaller(profile, agnt).CallAgent(agent.GetTasks()))
		if err != nil {
			return nil, err
		}
		return msgs.FindMany(filters...), nil
	default:
		// Return all the tasks still matching any state filters
		msgs, err := FromMaster(caller.CallMaster(master.GetTasks()))
		if err != nil {
			return nil, err
		}
		return msgs.FindMany(filters...), nil
	}
}

// FindFiles files FileInfos on an Agent
func FindFiles(profile *config.Profile, criteria Criteria) (Messages, error) {
	// Resolve the AgentInfo
	msgs, err := FindAgents(profile, Criteria{AgentID: criteria.AgentID})
	if err != nil {
		return nil, err
	}
	agnt, err := AsAgent(msgs.FindOne())
	if err != nil {
		return nil, err
	}
	return FromAgent(helper.NewAgentCaller(profile, agnt).CallAgent(agent.ListFiles(criteria.FilePath)))
}

// Find provides basic searching functionality across a Mesos cluster
// enabling searching for supported Mesos objects (Targets). Find may perform
// multiple calls to the Mesos master and agents depending on the provided criteria.
func Find(profile *config.Profile, criteria Criteria) (msgs Messages, err error) {
	switch criteria.Target {
	case AGENTS:
		return FindAgents(profile, criteria)
	case TASKS:
		return FindTasks(profile, criteria)
	case FILES:
		switch {
		case criteria.AgentID == "":
			err = ErrInvalidCriteria
		case criteria.FilePath == "":
			err = ErrInvalidCriteria
		default:
			return FindFiles(profile, criteria)
		}
	default:
		err = ErrInvalidCriteria
	}
	return msgs, err
}

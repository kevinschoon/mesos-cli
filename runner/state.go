package runner

import (
	"fmt"
	"github.com/mesos/mesos-go"
	"github.com/vektorlab/mesos-cli/mesosfile"
)

type ErrTaskTerminal struct {
	task   mesos.TaskID
	status mesos.TaskStatus
}

func (e ErrTaskTerminal) Error() string {
	return fmt.Sprintf(
		"Task %s entered a terminal state %s %s",
		e.task.Value,
		e.status.State.String(),
		e.status.Message,
	)
}

// State is an in-memory structure for keeping
// track of tasks. The State is updated when
// provided a mesos.TaskStatus via the Update
// method. If a task is encounted that was not
// specified during the creation of State it will
// panic.
type State struct {
	groups  []*mesosfile.Group
	states  map[string]mesos.TaskState
	pending chan *mesosfile.Group
	updates chan mesos.TaskStatus
	last    mesos.TaskState
	restart bool
	sync    bool
	done    bool
}

func NewState(groups []*mesosfile.Group, restart, sync bool) *State {
	state := &State{
		groups:  groups,
		states:  map[string]mesos.TaskState{},
		pending: make(chan *mesosfile.Group, len(groups)),
		updates: make(chan mesos.TaskStatus),
		restart: restart,
		sync:    sync,
		last:    mesos.TASK_FINISHED,
	}
	for _, group := range groups {
		// Reset Task and executor UUIDs
		group.Reset()
		for _, task := range group.Tasks {
			// Set state to TASK_STARTING
			state.states[task.TaskID.Value] = mesos.TaskState(0)
		}
		// Push the group into pending chan
		state.pending <- group
	}
	return state
}

// Total returns the total number of task groups.
func (s *State) Total() int {
	return len(s.groups)
}

// Pending returns the next task group waiting
// to be scheduled. If a returned group is not
// scheduled the caller must return it via
// Append or it will be lost.
func (s *State) Pop() *mesosfile.Group {
	// If running task groups synchronously each subsequent group (and task) must reach
	// TASK_FINISHED before the next task is available for scheduling.
	if s.sync && s.last != mesos.TASK_FINISHED {
		return nil
	}
	select {
	case group := <-s.pending:
		if s.sync {
			// Reset last to TASK_STARTING to avoid a race condition
			s.last = mesos.TaskState(0)
		}
		return group
	default:
	}
	return nil
}

// Append pushes the group into the pending chan.
func (s *State) Append(group *mesosfile.Group) {
	s.pending <- group
}

// Update places the TaskStatus into the updates channel
func (s *State) Update(status mesos.TaskStatus) {
	s.updates <- status
}

func (s *State) Monitor() (err error) {
loop:
	for {
		select {
		case status := <-s.updates:
			group := s.find(status.TaskID.Value)
			// Discard updates from orphaned tasks
			if group == nil {
				continue loop
			}
			// Check if the state is "terminal" or if the task exited normally but should be restarted
			if terminal(*status.State) || s.restart && *status.State == mesos.TASK_FINISHED {
				// TODO: Need to "backoff"
				if !s.restart {
					// If we will not restart the task return with ErrTaskTerminal
					err = ErrTaskTerminal{status.TaskID, status}
					break loop
				}
				// Remove all recorded states
				s.purge(group)
				// Reset UUIDs
				group.Reset()
				// Insert new states
				s.insert(group)
				// Push the group back into the pending chan
				s.pending <- group
				continue loop
			}
			// Update the state of this task
			s.states[status.TaskID.Value] = *status.State
			// Store the state of latest update
			s.last = *status.State
		}
		if s.finished() { // All tasks have finished
			break
		}
		if s.done {
			break // Done was toggled
		}
	}
	return err
}

// Toggle shutdown
func (s *State) Done() {
	s.done = true
}

func (s *State) find(id string) (g *mesosfile.Group) {
	for _, group := range s.groups {
		if task := group.Find(id); task != nil {
			g = group
		}
	}
	return g
}

func (s *State) purge(g *mesosfile.Group) {
	for _, task := range g.Tasks {
		delete(s.states, task.TaskID.Value)
	}
}

func (s *State) insert(g *mesosfile.Group) {
	for _, task := range g.Tasks {
		s.states[task.TaskID.Value] = mesos.TaskState(0)
	}
}

// Finished checks if all tasks are in state TASK_FINISHED
func (s *State) finished() bool {
	for _, state := range s.states {
		if state != mesos.TASK_FINISHED {
			return false
		}
	}
	return true
}

package runner

import (
	"fmt"
	"github.com/mesos/mesos-go"
	"github.com/satori/go.uuid"
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
	tasks   map[string]*mesos.TaskInfo
	states  map[string]mesos.TaskState
	pending chan *mesos.TaskInfo
	updates chan mesos.TaskStatus
	last    mesos.TaskState
	restart bool
	sync    bool
	done    bool
}

func NewState(tasks []*mesos.TaskInfo, restart, sync bool) *State {
	state := &State{
		states:  map[string]mesos.TaskState{},
		tasks:   map[string]*mesos.TaskInfo{},
		pending: make(chan *mesos.TaskInfo, len(tasks)),
		updates: make(chan mesos.TaskStatus),
		restart: restart,
		sync:    sync,
		last:    mesos.TASK_FINISHED,
	}
	for _, task := range tasks {
		// Assign a random task id
		task.TaskID.Value = uuid.NewV4().String()
		// Push the task into pending chan
		state.pending <- task
		// Record the TaskID
		state.tasks[task.TaskID.Value] = task
		// Sets the task as state "TASK_STARTING"
		state.states[task.TaskID.Value] = mesos.TaskState(0)
	}
	return state
}

// Total returns the total number of tasks.
func (s *State) Total() int {
	return len(s.tasks)
}

// Pending returns the next task waiting to
// be scheduled. If a returned task is not
// scheduled the caller must return it via
// Append or the Task will be lost.
func (s *State) Pop() *mesos.TaskInfo {
	// If running tasks synchronously each subsequent task must reach
	// TASK_FINISHED before the next task is available for scheduling.
	if s.sync && s.last != mesos.TASK_FINISHED {
		return nil
	}
	select {
	case task := <-s.pending:
		if s.sync {
			// Reset last to TASK_STARTING to avoid a race condition
			s.last = mesos.TaskState(0)
		}
		return task
	default:
	}
	return nil
}

// Append pushes the task into the pending chan.
func (s *State) Append(task *mesos.TaskInfo) {
	s.pending <- task
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
			// Check if the state is "terminal" or if the task exited normally but should be restarted
			if terminal(*status.State) || s.restart && *status.State == mesos.TASK_FINISHED {
				// TODO: Need to "backoff"
				if !s.restart {
					// If we will not restart the task return with ErrTaskTerminal
					err = ErrTaskTerminal{status.TaskID, status}
					break loop
				}

				task := s.tasks[status.TaskID.Value]
				// Remove the old ID
				delete(s.tasks, task.TaskID.Value)
				delete(s.states, task.TaskID.Value)
				// Generate a new ID
				task.TaskID.Value = uuid.NewV4().String()
				// Reset the task state
				s.tasks[task.TaskID.Value] = task
				s.states[task.TaskID.Value] = mesos.TaskState(0)
				// Push the task back into the pending chan
				s.pending <- s.tasks[task.TaskID.Value]
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

// Finished checks if all tasks are in state TASK_FINISHED
func (s *State) finished() bool {
	for _, task := range s.tasks {
		if s.states[task.TaskID.Value] != mesos.TASK_FINISHED {
			return false
		}
	}
	return true
}

func terminal(state mesos.TaskState) bool {
	switch state {
	case mesos.TASK_FAILED:
		return true
	case mesos.TASK_KILLED:
		return true
	case mesos.TASK_ERROR:
		return true
	case mesos.TASK_LOST:
		return true
	case mesos.TASK_DROPPED:
		return true
	//case mesos.TASK_UNREACHABLE:
	//	return true
	case mesos.TASK_GONE:
		return true
	case mesos.TASK_GONE_BY_OPERATOR:
		return true
	case mesos.TASK_UNKNOWN:
		return true
	}
	return false
}
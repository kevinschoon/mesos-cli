package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
)

/*
Runner is a task runner which implements the mesos.Scheduler interface.
*/
type Runner struct {
	task     *mesos.TaskInfo
	status   chan *mesos.TaskStatus
	errors   chan error
	launched bool
}

func (r *Runner) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {
	for _, offer := range offers {
		if r.launched {
			if _, err := driver.DeclineOffer(offer.Id, &mesos.Filters{RefuseSeconds: proto.Float64(5)}); err != nil {
				r.errors <- err
			}
			continue
		}
		if Sufficent(r.task, offer) {
			r.task.SlaveId = offer.SlaveId
			if _, err := driver.LaunchTasks(
				[]*mesos.OfferID{offer.Id},
				[]*mesos.TaskInfo{r.task},
				&mesos.Filters{RefuseSeconds: proto.Float64(5)},
			); err != nil {
				r.errors <- err
			}
			r.launched = true
		}
	}
}

func (r *Runner) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	if status.GetTaskId().Equal(r.task.TaskId) {
		r.status <- status
	}
}

func (r *Runner) Registered(_ sched.SchedulerDriver, _ *mesos.FrameworkID, _ *mesos.MasterInfo) {}
func (r *Runner) Reregistered(_ sched.SchedulerDriver, _ *mesos.MasterInfo)                     {}
func (r *Runner) Disconnected(_ sched.SchedulerDriver)                                          {}
func (r *Runner) OfferRescinded(_ sched.SchedulerDriver, _ *mesos.OfferID)                      {}

func (r *Runner) FrameworkMessage(_ sched.SchedulerDriver, _ *mesos.ExecutorID, _ *mesos.SlaveID, _ string) {

}

func (r *Runner) SlaveLost(_ sched.SchedulerDriver, id *mesos.SlaveID) {
	r.errors <- fmt.Errorf("Slave lost: %s", *id.Value)
}

func (r *Runner) ExecutorLost(_ sched.SchedulerDriver, executor *mesos.ExecutorID, agent *mesos.SlaveID, _ int) {
	r.errors <- fmt.Errorf("Executor lost: %s, %s", *executor.Value, *agent.Value)
}

func (r *Runner) Error(_ sched.SchedulerDriver, message string) {
	r.errors <- fmt.Errorf(message)
}

func NewDriver(master string, runner *Runner) (*sched.MesosSchedulerDriver, error) {
	config := sched.DriverConfig{
		Scheduler: runner,
		Framework: &mesos.FrameworkInfo{
			User: proto.String(""),
			Name: proto.String("mesos-exec"),
		},
		Master: master,
	}
	return sched.NewMesosSchedulerDriver(config)
}

func RunTask(master string, task *mesos.TaskInfo) (err error) {
	status := make(chan *mesos.TaskStatus)
	errors := make(chan error)
	runner := &Runner{
		task:   task,
		status: status,
		errors: errors,
	}
	driver, err := NewDriver(master, runner)
	if err != nil {
		return err
	}
	go func() { _, err = driver.Run() }()
loop:
	for {
		select {
		case s := <-status:
			state := s.State
			switch *state {
			case mesos.TaskState_TASK_RUNNING:
				go func() {
					fmt.Println(LogTask(master, s))
				}()
			case mesos.TaskState_TASK_LOST:
				driver.Stop(false)
				err = fmt.Errorf(state.String())
				break loop
			case mesos.TaskState_TASK_ERROR:
				err = fmt.Errorf(state.String())
				driver.Stop(false)
				break loop
			case mesos.TaskState_TASK_KILLED:
				err = fmt.Errorf(state.String())
				driver.Stop(false)
				break loop
			case mesos.TaskState_TASK_FAILED:
				err = fmt.Errorf(state.String())
				driver.Stop(false)
				break loop
			case mesos.TaskState_TASK_FINISHED:
				driver.Stop(false)
				break loop
			}
		case err = <-errors:
			break loop
		}
	}
	return err
}

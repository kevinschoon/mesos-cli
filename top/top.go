package top

import (
	"github.com/mesos/mesos-go"
	"github.com/mesos/mesos-go/httpcli/operator"
	master "github.com/mesos/mesos-go/master/calls"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/filter"
	"github.com/vektorlab/mesos-cli/helper"
	"github.com/vektorlab/toplib"
	"github.com/vektorlab/toplib/sample"
	"github.com/vektorlab/toplib/section"
	"time"
)

func collect(caller operator.Caller) ([]*sample.Sample, error) {
	samples := []*sample.Sample{}
	// NOTE: On large clusters this may be quite slow because the entire
	// Task state must be downloaded.
	resp, err := caller.CallMaster(master.GetTasks())
	if err != nil {
		return nil, err
	}
	filters := []filter.Filter{filter.TaskStateFilter([]*mesos.TaskState{mesos.TASK_RUNNING.Enum()})}
	for _, task := range filter.AsTasks(filter.FromMaster(resp).FindMany(filters...)) {
		smpl := sample.NewSample(task.TaskID.Value)
		smpl.SetString("AGENT", task.AgentID.Value)
		resources := mesos.Resources(task.Resources)
		cpus, _ := resources.CPUs()
		mem, _ := resources.Memory()
		disk, _ := resources.Disk()
		smpl.SetFloat64("CPU", cpus)
		smpl.SetFloat64("MEM", float64(mem))
		smpl.SetFloat64("DISK", float64(disk))
		samples = append(samples, smpl)
	}
	return samples, nil
}

func Run(profile *config.Profile) error {
	caller := helper.NewCaller(profile)
	// TODO: Add more sections like Agents, framework, etc.
	sections := []toplib.Section{
		section.NewSamples("tasks", "ID", "AGENT", "CPU", "MEM", "DISK"),
	}
	top := toplib.NewTop(sections)
	tick := time.NewTicker(1500 * time.Millisecond)
	go func() {
	loop:
		for {
			select {
			case <-top.Exit:
				close(top.Samples)
				break loop
			case <-tick.C:
				samples, err := collect(caller)
				if err != nil {
					break loop
				}
				top.Samples <- samples
			}
		}
		tick.Stop()
	}()
	return toplib.Run(top)
}

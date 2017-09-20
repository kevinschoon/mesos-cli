package top

import (
	"github.com/mesanine/mesos-cli/config"
	"github.com/mesanine/mesos-cli/filter"
	"github.com/mesanine/mesos-cli/helper"
	"github.com/mesos/mesos-go"
	agent "github.com/mesos/mesos-go/agent/calls"
	"github.com/mesos/mesos-go/httpcli/operator"
	master "github.com/mesos/mesos-go/master/calls"
	"github.com/vektorlab/toplib"
	"github.com/vektorlab/toplib/sample"
	"github.com/vektorlab/toplib/section"
)

func getContainers(profile *config.Profile, caller operator.Caller, namespace sample.Namespace) ([]*sample.Sample, error) {
	msgs, err := filter.FromMaster(caller.CallMaster(master.GetAgents()))
	if err != nil {
		return nil, err
	}
	samples := []*sample.Sample{}
	// TODO: Reduce several redundant calls here
	for _, agnt := range filter.AsAgents(msgs.FindMany()) {
		ac := helper.NewAgentCaller(profile, agnt)
		// TODO: Handle concurrently, will choke with many agents
		resp, err := ac.CallAgent(agent.GetContainers())
		if err != nil {
			return nil, err
		}
		if resp.GetContainers != nil {
			for _, container := range resp.GetContainers.Containers {
				smpl := sample.NewSample(container.ContainerId.Value, namespace)
				if stats := container.GetResourceStatistics(); stats != nil {
					smpl.SetFloat64("CPU_USR", *stats.CpusUserTimeSecs)
					smpl.SetFloat64("CPU_SYS", *stats.CpusSystemTimeSecs)
					smpl.SetFloat64("MEM_RSS", float64(*stats.MemRssBytes))
					smpl.SetFloat64("MEM_AVL", float64(*stats.MemTotalBytes))
				}
				samples = append(samples, smpl)
			}
		}
	}
	return samples, nil
}

func getTasks(caller operator.Caller, namespace sample.Namespace) ([]*sample.Sample, error) {
	msgs, err := filter.FromMaster(caller.CallMaster(master.GetTasks()))
	if err != nil {
		return nil, err
	}
	samples := []*sample.Sample{}
	filters := []filter.Filter{filter.TaskStateFilter([]*mesos.TaskState{mesos.TASK_RUNNING.Enum()})}
	for _, task := range filter.AsTasks(msgs.FindMany(filters...)) {
		smpl := sample.NewSample(task.TaskID.Value, namespace)
		smpl.SetString("NAME", task.Name)
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

func getAgents(caller operator.Caller, namespace sample.Namespace) ([]*sample.Sample, error) {
	msgs, err := filter.FromMaster(caller.CallMaster(master.GetAgents()))
	if err != nil {
		return nil, err
	}
	samples := []*sample.Sample{}
	for _, agent := range filter.AsAgents(msgs.FindMany()) {
		smpl := sample.NewSample(agent.ID.Value, namespace)
		resources := mesos.Resources(agent.Resources)
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

func getFrameworks(caller operator.Caller, namespace sample.Namespace) ([]*sample.Sample, error) {
	msgs, err := filter.FromMaster(caller.CallMaster(master.GetFrameworks()))
	if err != nil {
		return nil, err
	}
	samples := []*sample.Sample{}
	for _, framework := range filter.AsFrameworks(msgs.FindMany()) {
		smpl := sample.NewSample(framework.ID.Value, namespace)
		smpl.SetString("NAME", framework.Name)
		smpl.SetString("ROLE", *framework.Role)
		smpl.SetString("HOSTNAME", *framework.Hostname)
		samples = append(samples, smpl)
	}
	return samples, nil
}

func Run(profile *config.Profile) error {
	caller := helper.NewCaller(profile)
	// TODO: Add more sections like Agents, framework, etc.
	namespaces := map[string]sample.Namespace{
		"containers": sample.Namespace("containers"),
		"tasks":      sample.Namespace("tasks"),
		"agents":     sample.Namespace("agents"),
		"frameworks": sample.Namespace("frameworks"),
	}
	sections := []toplib.Section{
		section.NewSamples(namespaces["containers"], "ID", "CPU_USR", "CPU_SYS", "MEM_RSS", "MEM_AVL"),
		section.NewSamples(namespaces["tasks"], "ID", "NAME", "AGENT", "CPU", "MEM", "DISK"),
		section.NewSamples(namespaces["agents"], "ID", "CPU", "MEM", "DISK"),
		section.NewSamples(namespaces["frameworks"], "ID", "NAME", "ROLE", "HOSTNAME"),
	}
	top := toplib.NewTop(sections)
	return toplib.Run(
		top,
		func() ([]*sample.Sample, error) { return getContainers(profile, caller, namespaces["containers"]) },
		func() ([]*sample.Sample, error) { return getTasks(caller, namespaces["tasks"]) },
		func() ([]*sample.Sample, error) { return getAgents(caller, namespaces["agents"]) },
		func() ([]*sample.Sample, error) { return getAgents(caller, namespaces["frameworks"]) },
	)
}

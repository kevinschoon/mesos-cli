package mesosfile

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go"
)

const (
	CPU    float64 = 0.1
	MEMORY float64 = 32.0
	DISK   float64 = 32.0
)

func NewTask() *mesos.TaskInfo {
	return &mesos.TaskInfo{
		Name: "mesos-cli",
		Container: &mesos.ContainerInfo{
			Type: mesos.ContainerInfo_MESOS.Enum(),
			Docker: &mesos.ContainerInfo_DockerInfo{
				Network: mesos.ContainerInfo_DockerInfo_HOST.Enum(),
			},
			Mesos: &mesos.ContainerInfo_MesosInfo{},
		},
		Resources: []mesos.Resource{
			scalar("cpus", CPU),
			scalar("mem", MEMORY),
			scalar("disk", DISK),
		},
	}
}

// populate defaults and handle incompatible options
func initTasks(tasks []*mesos.TaskInfo) {
	for _, task := range tasks {
		def := NewTask()
		if len(task.Resources) == 0 {
			task.Resources = def.Resources
		} else {
			resources := mesos.Resources(task.Resources)
			if cpus := resources.SumScalars(mesos.NamedResources("cpus")); cpus == nil {
				resources.Add(
					scalar(
						"cpus",
						mesos.Resources(def.Resources).
							SumScalars(mesos.NamedResources("cpus")).Value,
					),
				)
			}
			if mem := resources.SumScalars(mesos.NamedResources("mem")); mem == nil {
				resources.Add(
					scalar(
						"mem", mesos.Resources(def.Resources).
							SumScalars(mesos.NamedResources("mem")).Value,
					),
				)
			}
			if disk := resources.SumScalars(mesos.NamedResources("disk")); disk == nil {
				resources.Add(
					scalar(
						"disk",
						mesos.Resources(def.Resources).
							SumScalars(mesos.NamedResources("disk")).Value,
					),
				)
			}
		}
		if *task.Container.Type == mesos.ContainerInfo_DOCKER {
			task.Container.Mesos = nil
		}
		if *task.Container.Type == mesos.ContainerInfo_MESOS {
			task.Container.Docker = nil
		}
	}
}

func scalar(name string, value float64) mesos.Resource {
	return mesos.Resource{
		Name:   name,
		Type:   mesos.SCALAR.Enum(),
		Role:   proto.String("*"),
		Scalar: &mesos.Value_Scalar{Value: value},
	}
}

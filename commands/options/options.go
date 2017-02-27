package options

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go"
)

// Option mutates a mesos.TaskInfo in someway once
// it has been created.
type Option func(task *mesos.TaskInfo)

func AsDocker(docker bool) Option {
	return func(task *mesos.TaskInfo) {
		if docker {
			task.Container.Type = mesos.ContainerInfo_DOCKER.Enum()
			task.Container.Mesos = nil
		} else {
			task.Container.Type = mesos.ContainerInfo_MESOS.Enum()
			task.Container.Docker = nil
		}
	}
}

func WithPorts() Option {
	return func(task *mesos.TaskInfo) {
		if *task.Container.Type == mesos.ContainerInfo_DOCKER {
			for _, mapping := range task.Container.Docker.PortMappings {
				task.Resources = append(task.Resources, portOffer(mapping.HostPort))
			}
		}
	}
}

func portOffer(port uint32) mesos.Resource {
	return mesos.Resource{
		Name:   "ports",
		Type:   mesos.RANGES.Enum(),
		Role:   proto.String("*"),
		Ranges: &mesos.Value_Ranges{Range: []mesos.Value_Range{mesos.Value_Range{Begin: uint64(port), End: uint64(port)}}},
	}
}

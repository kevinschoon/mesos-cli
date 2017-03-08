package options

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go"
)

func Apply(task *mesos.TaskInfo, options ...Option) {
	for _, opt := range options {
		opt(task)
	}
}

// Option mutates a mesos.TaskInfo in someway once
// it has been created.
type Option func(task *mesos.TaskInfo)

// WithContainerizer configures a Mesos task to use a specific
// containerizer. There are only two options (mesos or docker)
func WithContainerizer(docker bool) Option {
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

// WithRole configures a Mesos task to run with a specific role
func WithRole(role string) Option {
	return func(task *mesos.TaskInfo) {
		for _, resource := range task.Resources {
			*resource.Role = role
		}
	}
}

// WithDefaultResources configures a Mesos task with base resources
// if existing resources are not specified.
func WithDefaultResources() Option {
	return func(task *mesos.TaskInfo) {
		resources := mesos.Resources(task.Resources)
		if cpus := resources.SumScalars(mesos.NamedResources("cpus")); cpus == nil {
			resource := mesos.Resource{
				Name:   "cpus",
				Type:   mesos.SCALAR.Enum(),
				Role:   proto.String("*"),
				Scalar: &mesos.Value_Scalar{Value: 0.1},
			}
			task.Resources = append(task.Resources, resource)
		}
		if memory := resources.SumScalars(mesos.NamedResources("mem")); memory == nil {
			resource := mesos.Resource{
				Name:   "mem",
				Type:   mesos.SCALAR.Enum(),
				Role:   proto.String("*"),
				Scalar: &mesos.Value_Scalar{Value: 64.0},
			}
			task.Resources = append(task.Resources, resource)
		}
		if disk := resources.SumScalars(mesos.NamedResources("disk")); disk == nil {
			resource := mesos.Resource{
				Name:   "disk",
				Type:   mesos.SCALAR.Enum(),
				Role:   proto.String("*"),
				Scalar: &mesos.Value_Scalar{Value: 64.0},
			}
			task.Resources = append(task.Resources, resource)
		}
	}
}

// WithPorts adds resources for each specified port mapping
// when configured with the Docker containerizer
func WithPorts() Option {
	return func(task *mesos.TaskInfo) {
		if *task.Container.Type == mesos.ContainerInfo_DOCKER {
			for _, mapping := range task.Container.Docker.PortMappings {
				resource := mesos.Resource{
					Name: "ports",
					Type: mesos.RANGES.Enum(),
					Role: proto.String("*"),
					Ranges: &mesos.Value_Ranges{
						Range: []mesos.Value_Range{
							mesos.Value_Range{
								Begin: uint64(mapping.HostPort),
								End:   uint64(mapping.HostPort),
							},
						},
					},
				}
				task.Resources = append(task.Resources, resource)
			}
		}
	}
}

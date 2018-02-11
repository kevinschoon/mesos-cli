package mesosfile

import (
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/api/v1/lib"
	"github.com/satori/go.uuid"
)

type Option func(*Group) *Group

func FrameworkID(id string) Option {
	return func(group *Group) *Group {
		group.executor.FrameworkID.Value = id
		return group
	}
}

func AgentID(id string) Option {
	return func(group *Group) *Group {
		for _, task := range group.Tasks {
			task.AgentID.Value = id
		}
		return group
	}
}

func Reset() Option {
	return func(group *Group) *Group {
		group.executor.ExecutorID.Value = uuid.NewV4().String()
		for _, task := range group.Tasks {
			task.TaskID.Value = uuid.NewV4().String()
		}
		return group
	}
}

func Init() Option {
	return func(group *Group) *Group {
		initTasks(group.Tasks)
		return group
	}
}

func Docker(docker bool) Option {
	return func(group *Group) *Group {
		for _, task := range group.Tasks {
			if docker {
				task.Container.Type = mesos.ContainerInfo_DOCKER.Enum()
			}
		}
		return group
	}
}

func Role(role string) Option {
	return func(group *Group) *Group {
		for _, task := range group.Tasks {
			for _, resource := range task.Resources {
				*resource.Role = role
			}
		}
		return group
	}
}

func Ports() Option {
	return func(group *Group) *Group {
		for _, task := range group.Tasks {
			switch *task.Container.Type {
			case mesos.ContainerInfo_DOCKER:
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
		return group
	}
}

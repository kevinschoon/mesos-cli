package calls

import (
	pb "github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go"
	"github.com/mesos/mesos-go/agent"
	"time"
)

func GetState() *agent.Call {
	return &agent.Call{
		Type: agent.Call_GET_STATE.Enum(),
	}
}

func GetContainers() *agent.Call {
	return &agent.Call{
		Type: agent.Call_GET_CONTAINERS.Enum(),
	}
}

func GetFrameworks() *agent.Call {
	return &agent.Call{
		Type: agent.Call_GET_FRAMEWORKS.Enum(),
	}
}

func GetExecutors() *agent.Call {
	return &agent.Call{
		Type: agent.Call_GET_EXECUTORS.Enum(),
	}
}

func GetTasks() *agent.Call {
	return &agent.Call{
		Type: agent.Call_GET_TASKS.Enum(),
	}
}

func GetFlags() *agent.Call {
	return &agent.Call{
		Type: agent.Call_GET_FLAGS.Enum(),
	}
}

func GetVersion() *agent.Call {
	return &agent.Call{
		Type: agent.Call_GET_VERSION.Enum(),
	}
}

func GetMetrics(timeout time.Duration) *agent.Call {
	return &agent.Call{
		Type: agent.Call_GET_METRICS.Enum(),
		GetMetrics: &agent.Call_GetMetrics{
			Timeout: &mesos.DurationInfo{
				Nanoseconds: timeout.Nanoseconds(),
			},
		},
	}
}

func SetLoggingLevel(level uint32, duration time.Duration) *agent.Call {
	return &agent.Call{
		Type: agent.Call_SET_LOGGING_LEVEL.Enum(),
		SetLoggingLevel: &agent.Call_SetLoggingLevel{
			Level: &level,
			Duration: &mesos.DurationInfo{
				Nanoseconds: duration.Nanoseconds(),
			},
		},
	}
}

func GetLoggingLevel() *agent.Call {
	return &agent.Call{
		Type: agent.Call_GET_LOGGING_LEVEL.Enum(),
	}
}

func ReadFile(path string, length, offset uint64) *agent.Call {
	return &agent.Call{
		Type: agent.Call_READ_FILE.Enum(),
		ReadFile: &agent.Call_ReadFile{
			Path:   path,
			Length: length,
			Offset: offset,
		},
	}
}

func ListFiles(path string) *agent.Call {
	return &agent.Call{
		Type: agent.Call_LIST_FILES.Enum(),
		ListFiles: &agent.Call_ListFiles{
			Path: pb.String(path),
		},
	}
}

func LaunchNestedContainer(command *mesos.CommandInfo, container *mesos.ContainerInfo, containerID *mesos.ContainerID) *agent.Call {
	return &agent.Call{
		Type: agent.Call_LAUNCH_NESTED_CONTAINER.Enum(),
		LaunchNestedContainer: &agent.Call_LaunchNestedContainer{
			Command:     command,
			Container:   container,
			ContainerId: containerID,
		},
	}
}

func WaitNestedContainer(containerID *mesos.ContainerID) *agent.Call {
	return &agent.Call{
		Type: agent.Call_WAIT_NESTED_CONTAINER.Enum(),
		WaitNestedContainer: &agent.Call_WaitNestedContainer{
			ContainerId: containerID,
		},
	}
}

func KillNestedContainer(containerID *mesos.ContainerID) *agent.Call {
	return &agent.Call{
		Type: agent.Call_KILL_NESTED_CONTAINER.Enum(),
		KillNestedContainer: &agent.Call_KillNestedContainer{
			ContainerId: containerID,
		},
	}
}

func LaunchNestedContainerSession(containerID *mesos.ContainerID, command *mesos.CommandInfo) *agent.Call {
	return &agent.Call{
		Type: agent.Call_LAUNCH_NESTED_CONTAINER_SESSION.Enum(),
		LaunchNestedContainerSession: &agent.Call_LaunchNestedContainerSession{
			ContainerId: containerID,
			Command:     command,
		},
	}
}

func AttachContainerInput(containerID *mesos.ContainerID, proc *agent.ProcessIO) *agent.Call {
	var atype *agent.Call_AttachContainerInput_Type
	switch {
	case containerID != nil && proc == nil:
		atype = agent.Call_AttachContainerInput_CONTAINER_ID.Enum()
	case containerID == nil && proc != nil:
		atype = agent.Call_AttachContainerInput_PROCESS_IO.Enum()
	default:
		panic("specify containerID OR processIO")
	}
	return &agent.Call{
		Type: agent.Call_ATTACH_CONTAINER_INPUT.Enum(),
		AttachContainerInput: &agent.Call_AttachContainerInput{
			Type:        atype,
			ContainerId: containerID,
			ProcessIo:   proc,
		},
	}
}

func AttachContainerOutput(containerID *mesos.ContainerID) *agent.Call {
	return &agent.Call{
		Type: agent.Call_ATTACH_CONTAINER_OUTPUT.Enum(),
		AttachContainerOutput: &agent.Call_AttachContainerOutput{
			ContainerId: containerID,
		},
	}
}

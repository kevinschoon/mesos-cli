package calls

import (
	pb "github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go"
	"github.com/mesos/mesos-go/master"
	"time"
)

func GetHealth() *master.Call {
	return &master.Call{
		Type: master.Call_GET_HEALTH.Enum(),
	}
}

func GetFlags() *master.Call {
	return &master.Call{
		Type: master.Call_GET_FLAGS.Enum(),
	}
}

func GetVersion() *master.Call {
	return &master.Call{
		Type: master.Call_GET_VERSION.Enum(),
	}
}

func GetMetrics() *master.Call {
	return &master.Call{
		Type: master.Call_GET_METRICS.Enum(),
	}
}

func GetLoggingLevel() *master.Call {
	return &master.Call{
		Type: master.Call_GET_LOGGING_LEVEL.Enum(),
	}
}

func SetLoggingLevel(level int, duration time.Duration) *master.Call {
	return &master.Call{
		Type: master.Call_SET_LOGGING_LEVEL.Enum(),
		SetLoggingLevel: &master.Call_SetLoggingLevel{
			Level: pb.Uint32(uint32(level)),
			Duration: &mesos.DurationInfo{
				Nanoseconds: duration.Nanoseconds(),
			},
		},
	}
}

func ListFiles(path string) *master.Call {
	return &master.Call{
		Type: master.Call_LIST_FILES.Enum(),
		ListFiles: &master.Call_ListFiles{
			Path: pb.String(path),
		},
	}
}

func ReadFile() *master.Call {
	return &master.Call{
		Type:     master.Call_READ_FILE.Enum(),
		ReadFile: &master.Call_ReadFile{},
	}
}

func GetState() *master.Call {
	return &master.Call{
		Type: master.Call_GET_STATE.Enum(),
	}
}

func GetAgents() *master.Call {
	return &master.Call{
		Type: master.Call_GET_AGENTS.Enum(),
	}
}

func GetFrameworks() *master.Call {
	return &master.Call{
		Type: master.Call_GET_FRAMEWORKS.Enum(),
	}
}

func GetExecutors() *master.Call {
	return &master.Call{
		Type: master.Call_GET_EXECUTORS.Enum(),
	}
}

func GetTasks() *master.Call {
	return &master.Call{
		Type: master.Call_GET_TASKS.Enum(),
	}
}

func GetRoles() *master.Call {
	return &master.Call{
		Type: master.Call_GET_ROLES.Enum(),
	}
}

func GetWeights() *master.Call {
	return &master.Call{
		Type: master.Call_GET_WEIGHTS.Enum(),
	}
}

func UpdateWeights() *master.Call {
	return &master.Call{
		Type:          master.Call_UPDATE_WEIGHTS.Enum(),
		UpdateWeights: &master.Call_UpdateWeights{},
	}
}

func GetMaster() *master.Call {
	return &master.Call{
		Type: master.Call_GET_MASTER.Enum(),
	}
}

func Subscribe() *master.Call {
	return &master.Call{
		Type: master.Call_SUBSCRIBE.Enum(),
	}
}

func ReserveResources() *master.Call {
	return &master.Call{
		Type: master.Call_RESERVE_RESOURCES.Enum(),
	}
}

func UnreserveResources() *master.Call {
	return &master.Call{
		Type: master.Call_UNRESERVE_RESOURCES.Enum(),
	}
}

func CreateVolumes() *master.Call {
	return &master.Call{
		Type:          master.Call_CREATE_VOLUMES.Enum(),
		CreateVolumes: &master.Call_CreateVolumes{},
	}
}

func DestoryVolumes() *master.Call {
	return &master.Call{
		Type:           master.Call_DESTROY_VOLUMES.Enum(),
		DestroyVolumes: &master.Call_DestroyVolumes{},
	}
}

func GetMaintenanceStatus() *master.Call {
	return &master.Call{
		Type: master.Call_GET_MAINTENANCE_STATUS.Enum(),
	}
}

func GetMaintenanceSchedule() *master.Call {
	return &master.Call{
		Type: master.Call_GET_MAINTENANCE_SCHEDULE.Enum(),
	}
}

func UpdateMaintenanceSchedule() *master.Call {
	return &master.Call{
		Type: master.Call_UPDATE_MAINTENANCE_SCHEDULE.Enum(),
	}
}

func StartMaintenance() *master.Call {
	return &master.Call{
		Type:             master.Call_START_MAINTENANCE.Enum(),
		StartMaintenance: &master.Call_StartMaintenance{},
	}
}

func StopMaintenance() *master.Call {
	return &master.Call{
		Type:            master.Call_STOP_MAINTENANCE.Enum(),
		StopMaintenance: &master.Call_StopMaintenance{},
	}
}

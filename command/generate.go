package command

import (
	"encoding/json"
	"fmt"
	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	"os"
	"strings"

	"github.com/jawher/mow.cli"
)

var masterCalls = map[master.Call_Type]func(*cli.Cmd){
	master.Call_GET_HEALTH: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_HEALTH,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_FLAGS: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_FLAGS,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_METRICS: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_METRICS,
			GetMetrics: &master.Call_GetMetrics{
				Timeout: &mesos.DurationInfo{},
			},
		}
		setFlags(call.GetMetrics, cmd)
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_LOGGING_LEVEL: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_LOGGING_LEVEL,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_SET_LOGGING_LEVEL: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:            master.Call_SET_LOGGING_LEVEL,
			SetLoggingLevel: &master.Call_SetLoggingLevel{},
		}
		// TODO
		setFlags(call.SetLoggingLevel, cmd)
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_LIST_FILES: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:      master.Call_LIST_FILES,
			ListFiles: &master.Call_ListFiles{},
		}
		setFlags(call.ListFiles, cmd)
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_READ_FILE: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:     master.Call_READ_FILE,
			ReadFile: &master.Call_ReadFile{},
		}
		// TODO
		setFlags(call.ReadFile, cmd)
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_STATE: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_STATE,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_AGENTS: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_AGENTS,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_FRAMEWORKS: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_FRAMEWORKS,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_EXECUTORS: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_EXECUTORS,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_TASKS: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_TASKS,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_ROLES: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_ROLES,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_WEIGHTS: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_WEIGHTS,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_UPDATE_WEIGHTS: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:          master.Call_UPDATE_WEIGHTS,
			UpdateWeights: &master.Call_UpdateWeights{},
		}
		setFlags(call.UpdateWeights, cmd)
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_MASTER: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_STATE,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_SUBSCRIBE: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_SUBSCRIBE,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_RESERVE_RESOURCES: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:             master.Call_RESERVE_RESOURCES,
			ReserveResources: &master.Call_ReserveResources{},
		}
		setFlags(call.ReserveResources, cmd)
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_UNRESERVE_RESOURCES: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:               master.Call_UNRESERVE_RESOURCES,
			UnreserveResources: &master.Call_UnreserveResources{},
		}
		setFlags(call.UnreserveResources, cmd)
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_CREATE_VOLUMES: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:          master.Call_CREATE_VOLUMES,
			CreateVolumes: &master.Call_CreateVolumes{},
		}
		setFlags(call.CreateVolumes, cmd)
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_DESTROY_VOLUMES: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:           master.Call_DESTROY_VOLUMES,
			DestroyVolumes: &master.Call_DestroyVolumes{},
		}
		setFlags(call.DestroyVolumes, cmd)
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_MAINTENANCE_STATUS: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_MAINTENANCE_STATUS,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_UPDATE_MAINTENANCE_SCHEDULE: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_MAINTENANCE_SCHEDULE,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_START_MAINTENANCE: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:             master.Call_START_MAINTENANCE,
			StartMaintenance: &master.Call_StartMaintenance{},
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_STOP_MAINTENANCE: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:            master.Call_STOP_MAINTENANCE,
			StopMaintenance: &master.Call_StopMaintenance{},
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_GET_QUOTA: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type: master.Call_GET_QUOTA,
		}
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_SET_QUOTA: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:     master.Call_SET_QUOTA,
			SetQuota: &master.Call_SetQuota{},
		}
		setFlags(call.SetQuota, cmd)
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
	master.Call_REMOVE_QUOTA: func(cmd *cli.Cmd) {
		call := &master.Call{
			Type:        master.Call_REMOVE_QUOTA,
			RemoveQuota: &master.Call_RemoveQuota{},
		}
		setFlags(call.RemoveQuota, cmd)
		cmd.Action = func() {
			maybe(json.NewEncoder(os.Stdout).Encode(call))
		}
	},
}

func MasterCMDs(cmd *cli.Cmd) {
	for i := 1; i < len(master.Call_Type_name); i++ {
		callType := master.Call_Type(int32(i))
		name := strings.ToLower(callType.String())
		if fn, ok := masterCalls[callType]; ok {
			cmd.Command(name, fmt.Sprintf("generate a %s call", name), fn)
		} else {
			cmd.Command(name, fmt.Sprintf("generate a %s call (NOT IMPLEMENTED)", name), fn)
		}
	}
}

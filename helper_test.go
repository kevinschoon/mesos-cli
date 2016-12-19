package main

import (
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetPorts(t *testing.T) {
	task := NewTask()
	assert.NoError(t, setPorts(task, []string{
		"8000",
		"8001:80",
		"8002:80/tcp",
		"53/udp",
	}))
	mappings := task.Container.Docker.PortMappings
	assert.Len(t, mappings, 4)
	assert.Equal(t, *mappings[0].HostPort, uint32(8000))
	assert.Equal(t, *mappings[0].ContainerPort, uint32(8000))
	assert.Equal(t, *mappings[0].Protocol, "tcp")
	assert.Equal(t, *mappings[1].HostPort, uint32(8001))
	assert.Equal(t, *mappings[1].ContainerPort, uint32(80))
	assert.Equal(t, *mappings[1].Protocol, "tcp")
	assert.Equal(t, *mappings[2].HostPort, uint32(8002))
	assert.Equal(t, *mappings[2].ContainerPort, uint32(80))
	assert.Equal(t, *mappings[2].Protocol, "tcp")
	assert.Equal(t, *mappings[3].HostPort, uint32(53))
	assert.Equal(t, *mappings[3].ContainerPort, uint32(53))
	assert.Equal(t, *mappings[3].Protocol, "udp")
}

func TestSetParameters(t *testing.T) {
	task := NewTask()
	assert.NoError(t, setParameters(task, []string{
		"dns=127.0.0.1",
	}))
	parameters := task.Container.Docker.Parameters
	assert.Len(t, parameters, 1)
	assert.Equal(t, *parameters[0].Key, "dns")
	assert.Equal(t, *parameters[0].Value, "127.0.0.1")
}

func TestSetVolumes(t *testing.T) {
	task := NewTask()
	assert.NoError(t, setVolumes(task, []string{
		"/var/run/docker.sock:/var/run/docker.sock",
		"/some/path:/path:ro",
	}))
	volumes := task.Container.Volumes
	assert.Len(t, volumes, 2)
	assert.Equal(t, *volumes[0].HostPath, "/var/run/docker.sock")
	assert.Equal(t, *volumes[0].ContainerPath, "/var/run/docker.sock")
	assert.Equal(t, volumes[0].Mode.Enum(), mesos.Volume_RW.Enum())
	assert.Equal(t, *volumes[1].HostPath, "/some/path")
	assert.Equal(t, *volumes[1].ContainerPath, "/path")
	assert.Equal(t, volumes[1].Mode.Enum(), mesos.Volume_RO.Enum())
}

func TestSetEnvironment(t *testing.T) {
	task := NewTask()
	assert.NoError(t, setEnvironment(task, []string{
		"SOME_VAR=SOME_VALUE",
	}))
	envs := task.Command.Environment.Variables
	assert.Len(t, envs, 1)
	assert.Equal(t, *envs[0].Name, "SOME_VAR")
	assert.Equal(t, *envs[0].Value, "SOME_VALUE")
}

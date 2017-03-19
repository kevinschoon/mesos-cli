package mesosfile

import (
	"encoding/json"
	"github.com/mesos/mesos-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	flat = []byte(`
	[
		{
			"name": "hello",
			"command": {
				"shell": true,
				"value": "echo hello"
			},
			"container": {
				"type": "DOCKER",
				"image": "alpine"
			}
		},
		{
				"name": "world",
				"command": {
					"shell": true,
					"value": "echo world"
				}
		}
	]
	`)

	nested = []byte(`
	[
			{
				"networks": [
					{
						"name": "hello",
						"port_mappings": [
							{
								"host_port": 31000,
								"container_port": 80,
								"protocol": "tcp"
							}
						]
					}
				],
				"tasks": [
				{
					"name": "hello",
					"command": {
						"shell": true,
						"value": "echo hello"
					},
					"container": {
						"volumes": [
							{
								"source": {
									"type": "SANDBOX_PATH",
									"sandbox_path": {
										"type": "PARENT",
										"path": "./"
									}
								}
							}
						]
					}
				},
				{
						"name": "world",
						"command": {
							"shell": true,
							"value": "echo world"
					}
				}
			]
		}
	]
	`)
)

func TestMesosfileFlat(t *testing.T) {
	mf := Mesosfile{}
	assert.NoError(t, json.Unmarshal(flat, &mf))
	assert.Len(t, mf, 2)
	assert.Equal(t, mf[0].Tasks[0].Name, "hello")
	assert.Equal(t, mf[1].Tasks[0].Name, "world")
	assert.Nil(t, mf[0].Tasks[0].Container.Mesos)
	assert.Nil(t, mf[1].Tasks[0].Container.Docker)
}

func TestMesosfileNested(t *testing.T) {
	mf := Mesosfile{}
	assert.NoError(t, json.Unmarshal(nested, &mf))
	assert.Len(t, mf, 1)
	assert.Equal(t, mf[0].Tasks[0].Name, "hello")
	assert.Equal(t, mf[0].Tasks[1].Name, "world")
	assert.Nil(t, mf[0].Tasks[0].Container.Docker)
	assert.Nil(t, mf[0].Tasks[1].Container.Docker)
	assert.Len(t, mf[0].Networks, 1)
	assert.Len(t, mf[0].executor.Container.NetworkInfos, 1)
	assert.Equal(t, *mf[0].Tasks[0].Container.Volumes[0].Source.Type, mesos.Volume_Source_SANDBOX_PATH)
	assert.Equal(t, *mf[0].Tasks[0].Container.Volumes[0].Source.SandboxPath.Type, mesos.Volume_Source_SandboxPath_PARENT)
	assert.Equal(t, mf[0].Tasks[0].Container.Volumes[0].Source.SandboxPath.Path, "./")
}

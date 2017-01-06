package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Agents returns a map of IDs to hostnames
func Agents(master string) (map[string]string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/master/slaves", master))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	agents := map[string]string{}
	for _, agent := range gjson.GetBytes(raw, "slaves").Array() {
		hostname := agent.Get("hostname").Str
		// Detect port agent is listening on
		split := strings.Split(agent.Get("pid").Str, ":")
		if len(split) != 2 {
			return nil, fmt.Errorf("Cannot detect port")
		}
		port, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			return nil, err
		}
		agents[agent.Get("id").Str] = fmt.Sprintf("%s:%d", hostname, port)
	}
	return agents, nil
}

// LogDir returns the directory path for following task output
func LogDir(hostname, executorId string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/slave(1)/state", hostname))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	for _, framework := range gjson.GetBytes(raw, "frameworks").Array() {
		for _, executor := range framework.Get("executors").Array() {
			if executor.Get("id").Str == executorId {
				return executor.Get("directory").Str, nil
			}
		}
	}
	return "", fmt.Errorf("Unable to find log directory")
}

// Resource returns the value of a resource
func Resource(name string, resources []*mesos.Resource) float64 {
	var value float64
	for _, resource := range resources {
		if resource.GetName() == name {
			value = resource.GetScalar().GetValue()
		}
	}
	return value
}

// Check if a Mesos resource offer can satisfy the Task
func Sufficent(task *mesos.TaskInfo, offer *mesos.Offer) bool {
	for _, resource := range offer.Resources {
		value := resource.GetScalar().GetValue()
		switch resource.GetName() {
		case "cpus":
			if value < Resource("cpus", task.Resources) {
				return false
			}
		case "mem":
			if value < Resource("mem", task.Resources) {
				return false
			}
		case "disk":
			if value < Resource("disk", task.Resources) {
				return false
			}
		}
	}
	return true
}

// NewTask returns a mesos.TaskInfo with sensibly
// populated default values.
func NewTask() *mesos.TaskInfo {
	task := &mesos.TaskInfo{
		// TODO: Generate unique taskid
		TaskId: &mesos.TaskID{Value: proto.String("mesos-cli")},
		Name:   proto.String("mesos-cli"),
		Command: &mesos.CommandInfo{
			Shell: proto.Bool(false),
			User:  proto.String("root"),
		},
		Container: &mesos.ContainerInfo{
			// Default to Mesos Containerizer
			Type:  mesos.ContainerInfo_MESOS.Enum(),
			Mesos: &mesos.ContainerInfo_MesosInfo{},
			Docker: &mesos.ContainerInfo_DockerInfo{
				Privileged:     proto.Bool(false),
				ForcePullImage: proto.Bool(false),
				Parameters:     []*mesos.Parameter{},
				PortMappings:   []*mesos.ContainerInfo_DockerInfo_PortMapping{},
				Network:        mesos.ContainerInfo_DockerInfo_BRIDGE.Enum(),
			},
		},
		Resources: []*mesos.Resource{
			&mesos.Resource{
				Name: proto.String("cpus"),
				Type: mesos.Value_SCALAR.Enum(),
				Scalar: &mesos.Value_Scalar{
					Value: proto.Float64(0.1),
				},
			},
			&mesos.Resource{
				Name: proto.String("mem"),
				Type: mesos.Value_SCALAR.Enum(),
				Scalar: &mesos.Value_Scalar{
					Value: proto.Float64(128.0),
				},
			},
			&mesos.Resource{
				Name: proto.String("disk"),
				Type: mesos.Value_SCALAR.Enum(),
				Scalar: &mesos.Value_Scalar{
					Value: proto.Float64(32.0),
				},
			},
		},
	}
	return task
}

func TaskFromJSON(task *mesos.TaskInfo, path string) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, task)
}

/* setPorts translates portMapping parameters into
Docker portmappings. Mappings can be accepted in several formats:
8000
8000:80
8000/tcp
8000:80/tcp
*/
func setPorts(task *mesos.TaskInfo, ports []string) error {
	mappings := []*mesos.ContainerInfo_DockerInfo_PortMapping{}
	for _, port := range ports {
		mapping := &mesos.ContainerInfo_DockerInfo_PortMapping{
			// Assume tcp
			Protocol: proto.String("tcp"),
		}
		// Check protocol
		split := strings.Split(port, "/")
		if len(split) > 2 {
			return fmt.Errorf("Bad port mapping %s", port)
		}
		if len(split) == 2 {
			fmt.Println(split)
			if !(split[1] == "tcp" || split[1] == "udp") {
				return fmt.Errorf("Bad protocol %s", port)
			}
			*mapping.Protocol = split[1]
			// Remove protocol
			port = strings.Replace(port, "/"+split[1], "", 1)
		}
		split = strings.Split(port, ":")
		if len(split) > 2 {
			return fmt.Errorf("Bad port mapping %s", port)
		}
		// 8000:80
		if len(split) == 2 {
			host, err := strconv.ParseUint(split[0], 10, 32)
			if err != nil {
				return err
			}
			mapping.HostPort = proto.Uint32(uint32(host))
			cont, err := strconv.ParseUint(split[1], 10, 32)
			if err != nil {
				return err
			}
			mapping.ContainerPort = proto.Uint32(uint32(cont))
		}
		// 8000
		if len(split) == 1 {
			host, err := strconv.ParseUint(split[0], 10, 32)
			if err != nil {
				return err
			}
			mapping.HostPort = proto.Uint32(uint32(host))
			mapping.ContainerPort = proto.Uint32(uint32(host))
		}
		mappings = append(mappings, mapping)
	}
	task.Container.Docker.PortMappings = mappings
	return nil
}

func setParameters(task *mesos.TaskInfo, params []string) error {
	parameters := []*mesos.Parameter{}
	for _, param := range params {
		split := strings.Split(param, "=")
		if len(split) != 2 {
			return fmt.Errorf("Invalid parameter: %s", param)
		}
		parameters = append(parameters, &mesos.Parameter{
			Key:   proto.String(split[0]),
			Value: proto.String(split[1]),
		})
	}
	task.Container.Docker.Parameters = parameters
	return nil
}

func setVolumes(task *mesos.TaskInfo, vols []string) error {
	volumes := []*mesos.Volume{}
	for _, vol := range vols {
		split := strings.Split(vol, ":")
		if len(split) < 2 || len(split) > 3 {
			return fmt.Errorf("Bad volume: %s", vol)
		}
		volume := &mesos.Volume{
			HostPath:      proto.String(split[0]),
			ContainerPath: proto.String(split[1]),
			Mode:          mesos.Volume_RW.Enum(),
		}
		if len(split) == 3 {
			switch split[2] {
			case "RO":
				volume.Mode = mesos.Volume_RO.Enum()
			case "ro":
				volume.Mode = mesos.Volume_RO.Enum()
			case "RW":
				volume.Mode = mesos.Volume_RW.Enum()
			case "rw":
				volume.Mode = mesos.Volume_RW.Enum()
			default:
				return fmt.Errorf("Bad volume: %s", vol)
			}
		}
		volumes = append(volumes, volume)
	}
	task.Container.Volumes = volumes
	return nil
}

func setEnvironment(task *mesos.TaskInfo, envs []string) error {
	variables := []*mesos.Environment_Variable{}
	for _, env := range envs {
		split := strings.Split(env, "=")
		if len(split) != 2 {
			return fmt.Errorf("Bad environment variable: %s", env)
		}
		variables = append(variables, &mesos.Environment_Variable{
			Name:  proto.String(split[0]),
			Value: proto.String(split[1]),
		})
	}
	task.Command.Environment = &mesos.Environment{
		Variables: variables,
	}
	return nil
}

func truncStr(s string, l int) string {
	runes := bytes.Runes([]byte(s))
	if len(runes) < l {
		return s
	}
	return string(runes[:l])
}

// Convenience types for cli so we may
// specify default values in one place
// as pass them to the cli parser.
type str struct {
	pt *string
}

func (s str) String() string {
	return *s.pt
}

func (s str) Set(other string) error {
	*s.pt = other
	return nil
}

type bl struct {
	pt *bool
}

func (b bl) String() string {
	if *b.pt {
		return "true"
	}
	return "false"
}

func (b bl) Set(other string) error {
	if other != "" {
		*b.pt = true
	}
	return nil
}

type flt struct {
	pt *float64
}

func (f flt) String() string {
	return fmt.Sprintf("%.1f", *f.pt)
}

func (f flt) Set(other string) error {
	v, err := strconv.ParseFloat(other, 64)
	if err != nil {
		return err
	}
	*f.pt = v
	return nil
}

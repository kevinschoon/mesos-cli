package flags

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/jawher/mow.cli"
	"github.com/mesos/mesos-go"
	"strconv"
	"strings"
)

var Flags = []Flag{
	Name,
	Command,
	Volumes,
	Image,
	Resources,
	Privileged,
	PortMapping,
	Parameters,
	NetworkMode,
}

type Flag func(*mesos.TaskInfo, *cli.Cmd)

// implements flag.Value
type flag struct {
	set func(string) error
	str func() string
}

func (f flag) Set(v string) error { return f.set(v) }
func (f flag) String() string     { return f.str() }

func Name(task *mesos.TaskInfo, cmd *cli.Cmd) {
	name := flag{
		set: func(v string) error {
			task.Name = v
			return nil
		},
		str: func() string {
			return task.Name
		},
	}
	cmd.VarOpt("name", name, "Friendly task name")
}

func Command(task *mesos.TaskInfo, cmd *cli.Cmd) {
	task.Command = &mesos.CommandInfo{
		User:  cmd.StringOpt("user", "root", "User to run as"),
		Shell: cmd.BoolOpt("shell", false, "Run as a shell command"),
		Environment: &mesos.Environment{
			Variables: []mesos.Environment_Variable{},
		},
	}
	value := flag{
		set: func(v string) error {
			if *task.Command.Shell {
				task.Command.Value = proto.String(v)
			} else {
				task.Command.Arguments = strings.Split(v, " ")
			}
			return nil
		},
		str: func() string { return "" },
	}
	envs := flag{
		set: func(v string) error {
			split := strings.Split(v, "=")
			if len(split) != 2 {
				return fmt.Errorf("Bad environment variable %s", v)
			}
			task.Command.Environment.Variables = append(
				task.Command.Environment.Variables,
				mesos.Environment_Variable{Name: split[0], Value: split[1]},
			)
			return nil
		},
		str: func() string {
			var s string
			for _, env := range task.Command.Environment.Variables {
				s += fmt.Sprintf("%s=%s", env.Name, env.Value)
			}
			return s
		},
	}
	cmd.VarOpt("e env", envs, "Environment variables")
	cmd.VarArg("CMD", value, "Command to run")
}

func Volumes(task *mesos.TaskInfo, cmd *cli.Cmd) {
	vols := flag{
		set: func(v string) error {
			// TODO Need to support image and other parameters
			split := strings.Split(v, ":")
			if len(split) < 2 {
				return fmt.Errorf("Bad volume: %s", v)
			}
			vol := mesos.Volume{
				HostPath:      split[0],
				ContainerPath: split[1],
			}
			if len(split) == 3 {
				split[2] = strings.ToUpper(split[2])
				if !(split[2] == "RW" || split[2] == "RO") {
					return fmt.Errorf("Bad volume: %s", v)
				}
				vol.Mode = mesos.Volume_Mode(mesos.Volume_Mode_value[split[2]]).Enum()
			} else {
				vol.Mode = mesos.RO.Enum()
			}
			task.Container.Volumes = append(task.Container.Volumes, vol)
			return nil
		},
		str: func() string {
			var s string
			for _, vol := range task.Container.Volumes {
				s += fmt.Sprintf("%s:%s", vol.HostPath, vol.ContainerPath)
			}
			return s
		},
	}
	cmd.VarOpt("v volume", vols, "Container volumes")
}

func Image(task *mesos.TaskInfo, cmd *cli.Cmd) {
	image := flag{
		set: func(v string) error {
			task.Container.Docker.Image = v
			return nil
		},
		str: func() string { return task.Container.Docker.Image },
	}
	cmd.VarOpt("i image", image, "Image to run")
}

func Resources(task *mesos.TaskInfo, cmd *cli.Cmd) {
	task.Resources = []mesos.Resource{
		mesos.Resource{
			Name:   "cpus",
			Type:   mesos.SCALAR.Enum(),
			Role:   proto.String("*"),
			Scalar: &mesos.Value_Scalar{Value: 0.1},
		},
		mesos.Resource{
			Name:   "mem",
			Type:   mesos.SCALAR.Enum(),
			Role:   proto.String("*"),
			Scalar: &mesos.Value_Scalar{Value: 64.0},
		},
		mesos.Resource{
			Name:   "disk",
			Type:   mesos.SCALAR.Enum(),
			Role:   proto.String("*"),
			Scalar: &mesos.Value_Scalar{Value: 64.0},
		},
	}
}

//
// Docker Only Options
//
func Privileged(task *mesos.TaskInfo, cmd *cli.Cmd) {
	task.Container.Docker.Privileged = cmd.BoolOpt(
		"privileged", false, "Run Docker in privileged mode",
	)
}

func PortMapping(task *mesos.TaskInfo, cmd *cli.Cmd) {
	ports := flag{
		set: func(v string) error {
			split := strings.Split(v, ":")
			if len(split) != 2 {
				return fmt.Errorf("Bad port mapping %s", v)
			}
			mapping := mesos.ContainerInfo_DockerInfo_PortMapping{
				Protocol: proto.String("tcp"),
			}
			host, err := strconv.ParseUint(split[0], 0, 32)
			if err != nil {
				return fmt.Errorf("Bad port mapping %s", v)
			}
			mapping.HostPort = uint32(host)
			split = strings.Split(split[1], "/")
			if len(split) == 2 {
				split[1] = strings.ToLower(split[1])
				if !(split[1] == "tcp" || split[1] == "udp") {
					return fmt.Errorf("Bad port mapping %s", v)
				}
				mapping.Protocol = proto.String(split[1])
			}
			cont, err := strconv.ParseUint(split[0], 0, 32)
			if err != nil {
				return fmt.Errorf("Bad port mapping %s", v)
			}
			mapping.ContainerPort = uint32(cont)
			task.Container.Docker.PortMappings = append(
				task.Container.Docker.PortMappings,
				mapping,
			)
			return nil
		},
		str: func() string {
			var s string
			for _, m := range task.Container.Docker.PortMappings {
				s += fmt.Sprintf("%d:%d/%s", m.HostPort, m.ContainerPort, *m.Protocol)
			}
			return s
		},
	}
	cmd.VarOpt("p port", ports, "Port mappings [Docker only]")
}

func Parameters(task *mesos.TaskInfo, cmd *cli.Cmd) {
	params := flag{
		set: func(v string) error {
			split := strings.Split(v, "=")
			if len(split) != 2 {
				return fmt.Errorf("Bad Docker parameter %s", v)
			}
			task.Container.Docker.Parameters = append(
				task.Container.Docker.Parameters,
				mesos.Parameter{Key: split[0], Value: split[1]},
			)
			return nil
		},
		str: func() string {
			var s string
			for _, param := range task.Container.Docker.Parameters {
				s += fmt.Sprintf("%s=%s", param.Key, param.Value)
			}
			return s
		},
	}
	cmd.VarOpt("param", params, "Docker parameters [Docker only]")
}

func NetworkMode(task *mesos.TaskInfo, cmd *cli.Cmd) {
	net := flag{
		set: func(v string) error {
			i, ok := mesos.ContainerInfo_DockerInfo_Network_value[strings.ToUpper(v)]
			if !ok {
				return fmt.Errorf("Bad network mode: %s", v)
			}
			task.Container.Docker.Network = mesos.ContainerInfo_DockerInfo_Network(i).Enum()
			return nil
		},
		str: func() string {
			return task.Container.Docker.Network.String()
		},
	}
	cmd.VarOpt("net", net, "Network Mode [Docker only]")
}

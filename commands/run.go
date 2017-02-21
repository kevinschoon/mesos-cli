package commands

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/jawher/mow.cli"
	"github.com/mesos/mesos-go"
	"github.com/vektorlab/mesos-cli/config"
	"github.com/vektorlab/mesos-cli/runner"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Run struct{}

func (_ Run) Name() string { return "run" }
func (_ Run) Desc() string { return "Run tasks on Mesos" }

func (r *Run) Init(profile Profile) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [CMD]"
		var (
			command  = cmd.StringArg("CMD", "", "Command to run")
			user     = cmd.StringOpt("user", "root", "User to run as")
			shell    = cmd.BoolOpt("shell", true, "Run as a shell command")
			hostname = cmd.StringOpt("master", "", "Mesos master")
			path     = cmd.StringOpt("path", "", "Path to a JSON file containing a Mesos TaskInfo")
			dump     = cmd.BoolOpt("json", false, "Dump the task to JSON instead of running it")
			docker   = cmd.BoolOpt("docker", false, "Run as a Docker container")
			image    = cmd.StringOpt("image", "", "Image to run")

			// Docker-only options
			privileged = cmd.BoolOpt("privileged", false, "Run in privileged mode [docker only]")
		)
		envs := &Envs{}
		cmd.VarOpt("e env", envs, "Environment variables")
		// Docker-only options
		net := &NetworkMode{}
		cmd.VarOpt("net", net, "Network Mode [Docker only]")
		params := &Parameters{}
		cmd.VarOpt("param", params, "Freeform Docker parameters [Docker only]")
		mappings := &PortMappings{}
		cmd.VarOpt("p port", mappings, "Port mappings [Docker only]")

		// TODO
		// resources
		// volumes
		// URIs
		// ...

		cmd.Action = func() {

			if *path != "" {
				raw, err := ioutil.ReadFile(*path)
				failOnErr(err)
				info := &mesos.TaskInfo{}
				failOnErr(json.Unmarshal(raw, info))
				profile().With(config.TaskInfo(info))
			}

			runner := runner.New(profile().With(
				config.Master(*hostname),
				config.Command(
					config.CommandOpts{
						Value: *command,
						User:  *user,
						Shell: *shell,
						Envs:  *envs,
					},
				),
				config.Container(
					config.ContainerOpts{
						Docker:       *docker,
						Privileged:   *privileged,
						Image:        *image,
						Parameters:   *params,
						NetworkMode:  net.mode,
						PortMappings: *mappings,
					},
				),
			))

			if *dump {
				raw, err := json.MarshalIndent(profile().TaskInfo, " ", " ")
				failOnErr(err)
				fmt.Println(string(raw))
				os.Exit(0)
			}

			failOnErr(runner.Run())
		}
	}
}

type NetworkMode struct {
	mode mesos.ContainerInfo_DockerInfo_Network
}

func (n *NetworkMode) Set(v string) error {
	mode, ok := mesos.ContainerInfo_DockerInfo_Network_value[strings.ToUpper(v)]
	if !ok {
		return fmt.Errorf("Bad network mode: %s", v)
	}
	n.mode = mesos.ContainerInfo_DockerInfo_Network(mode)
	return nil
}

func (n NetworkMode) String() string {
	mode, _ := mesos.ContainerInfo_DockerInfo_Network_name[int32(n.mode)]
	return mode
}

type PortMappings []mesos.ContainerInfo_DockerInfo_PortMapping

func (mappings *PortMappings) Set(v string) (err error) {
	split := strings.Split(v, ":")
	if len(split) != 2 {
		return fmt.Errorf("Bad port mapping %s", v)
	}
	mapping := mesos.ContainerInfo_DockerInfo_PortMapping{}
	host, err := strconv.ParseUint(split[0], 0, 32)
	if err != nil {
		return fmt.Errorf("Bad port mapping %s", v)
	}
	mapping.ContainerPort = uint32(host)
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
	*mappings = append(*mappings, mapping)
	return nil
}

func (mappings PortMappings) String() string {
	var s string
	for _, mapping := range mappings {
		s += fmt.Sprintf("%d:%d/%s", mapping.HostPort, mapping.ContainerPort, *mapping.Protocol)
	}
	return s
}

type Parameters []mesos.Parameter

func (params Parameters) String() string {
	var s string
	for _, param := range params {
		s += fmt.Sprintf("%s=%s", param.Key, param.Value)
	}
	return s
}

func (params *Parameters) Set(v string) error {
	split := strings.Split(v, "=")
	if len(split) != 2 {
		return fmt.Errorf("Bad Docker parameter %s", v)
	}
	*params = append(*params, mesos.Parameter{Key: split[0], Value: split[1]})
	return nil
}

type Envs []mesos.Environment_Variable

func (envs Envs) String() string {
	var s string
	for _, env := range envs {
		s += fmt.Sprintf("%s=%s", env.Name, env.Value)
	}
	return s
}

func (envs *Envs) Set(v string) error {
	split := strings.Split(v, "=")
	if len(split) != 2 {
		return fmt.Errorf("Bad environment variable %s", v)
	}
	*envs = append(*envs, mesos.Environment_Variable{Name: split[0], Value: split[1]})
	return nil
}

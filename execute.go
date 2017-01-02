package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/jawher/mow.cli"
	mesos "github.com/mesos/mesos-go/mesosproto"
)

func exec(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS] [ARG...]"

	var (
		arguments  = cmd.StringsArg("ARG", nil, "Command Arguments")
		taskPath   = cmd.StringOpt("task", "", "Path to a Mesos TaskInfo JSON file")
		parameters = cmd.StringsOpt("param", []string{}, "Docker parameters")
		image      = cmd.StringOpt("i image", "", "Docker image to run")
		volumes    = cmd.StringsOpt("v volume", []string{}, "Volume mappings")
		ports      = cmd.StringsOpt("p ports", []string{}, "Port mappings")
		envs       = cmd.StringsOpt("e env", []string{}, "Environment Variables")
		shell      = cmd.StringOpt("s shell", "", "Shell command to execute")
	)

	task := NewTask()
	cmd.VarOpt(
		"n name",
		str{pt: task.Name},
		"Task Name",
	)
	cmd.VarOpt(
		"u user",
		str{pt: task.Command.User},
		"User to run as",
	)
	cmd.VarOpt(
		"c cpus",
		flt{pt: task.Resources[0].Scalar.Value},
		"CPU Resources to allocate",
	)
	cmd.VarOpt(
		"m mem",
		flt{pt: task.Resources[1].Scalar.Value},
		"Memory Resources (mb) to allocate",
	)
	cmd.VarOpt(
		"d disk",
		flt{pt: task.Resources[2].Scalar.Value},
		"Disk Resources (mb) to allocate",
	)
	cmd.VarOpt(
		"privileged",
		bl{pt: task.Container.Docker.Privileged},
		"Give extended privileges to this container",
	)
	cmd.VarOpt(
		"f forcePullImage",
		bl{pt: task.Container.Docker.ForcePullImage},
		"Always pull the container image",
	)

	cmd.Before = func() {
		if *shell != "" {
			task.Command.Shell = proto.Bool(true)
			task.Command.Value = shell
		} else {
			for _, arg := range *arguments {
				*task.Command.Value += fmt.Sprintf(" %s", arg)
			}
		}
		if *taskPath != "" {
			failOnErr(TaskFromJSON(task, *taskPath))
		}
		failOnErr(setPorts(task, *ports))
		failOnErr(setVolumes(task, *volumes))
		failOnErr(setParameters(task, *parameters))
		failOnErr(setEnvironment(task, *envs))
		// Assuming that if image is specified the user wants
		// to run with the Docker containerizer. This is
		// not always the case as an image may be passed
		// to the Mesos containerizer as well.
		if *image != "" {
			task.Container.Mesos = nil
			task.Container.Type = mesos.ContainerInfo_DOCKER.Enum()
			task.Container.Docker.Image = image
		} else {
			task.Container.Docker = nil
		}
		// Nothing to do if not running a container
		// and no arguments are specified.
		if *image == "" && *taskPath == "" && len(*arguments) == 0 && *shell == "" {
			cmd.PrintHelp()
			cli.Exit(1)
		}
	}
	cmd.Action = func() { failOnErr(RunTask(config.Profile(), task)) }
}

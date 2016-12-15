package main

import (
	"flag"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/jawher/mow.cli"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"os"
)

const Version = "0.0.1"

func main() {
	app := cli.App("mesos-exec", "Execute Commands on Apache Mesos")
	app.Spec = "[OPTIONS] [ARG...]"
	app.Version("v version", fmt.Sprintf("mesos-exec version %s", Version))
	var (
		image     = app.StringOpt("i image", "", "Docker image to run")
		master    = app.StringOpt("master", "127.0.0.1:5050", "Master address <host:port>")
		arguments = app.StringsArg("ARG", nil, "Command Arguments")
		verbose   = app.IntOpt("V verbosity", 0, "Level of verbosity")
	)
	task := NewTask()
	app.VarOpt(
		"n name",
		str{pt: task.Name},
		"Task Name",
	)
	app.VarOpt(
		"s shell",
		bl{pt: task.Command.Shell},
		"Execute as shell command",
	)
	app.VarOpt(
		"u user",
		str{pt: task.Command.User},
		"User to run as",
	)
	app.VarOpt(
		"c cpus",
		flt{pt: task.Resources[0].Scalar.Value},
		"CPU Resources to allocate",
	)
	app.VarOpt(
		"m mem",
		flt{pt: task.Resources[1].Scalar.Value},
		"Memory Resources (mb) to allocate",
	)
	app.VarOpt(
		"d disk",
		flt{pt: task.Resources[2].Scalar.Value},
		"Disk Resources (mb) to allocate",
	)
	app.Before = func() {
		args := *arguments
		if *task.Command.Shell {
			cmd := ""
			if len(args) == 1 {
				task.Command.Value = proto.String(args[0])
			}
			if len(args) > 1 {
				for _, arg := range args {
					cmd += fmt.Sprintf(" %s", arg)
				}
				task.Command.Value = proto.String(cmd)
			}
		} else {
			task.Command.Arguments = args
		}
		// Assuming that if image is specified the user wants
		// to run with the Docker containerizer. This is
		// not always the case as an image may be passed
		// to the Mesos containerizer as well.
		if *image != "" {
			task.Container.Type = mesos.ContainerInfo_DOCKER.Enum()
			task.Container.Docker.Image = image
		}
		// Nothing to do if not running a container
		// and no arguments are specified.
		if *image == "" && len(args) == 0 {
			app.PrintHelp()
			cli.Exit(1)
		}
	}
	app.Action = func() {
		// This is done to satisfy the presumptuous golang/glog package
		// which assumes I am using flag and insists it be configured
		// with such. Since glog is used in go-mesos it is easiest to use
		// the same library for the moment. TODO
		flag.CommandLine.Set("v", string(*verbose))
		flag.CommandLine.Set("logtostderr", "1")
		flag.CommandLine.Parse([]string{})
		if err := RunTask(*master, task); err != nil {
			fmt.Errorf("Error: ", err.Error())
			os.Exit(1)
		}
	}
	app.Run(os.Args)
}

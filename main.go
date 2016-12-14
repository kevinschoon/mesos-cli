package main

import (
	"flag"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/jawher/mow.cli"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"os"
	"strconv"
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
		cpus      = app.StringOpt("c cpus", "0.1", "CPU Resources to allocate")
		mem       = app.StringOpt("m mem", "128.0", "Memory resources (mb) to allocate")
		disk      = app.StringOpt("d disk", "64.0", "Memory resources (mb) to allocate")
		level     = app.StringOpt("level", "0", "Logging level")
	)
	task := &mesos.TaskInfo{
		TaskId: &mesos.TaskID{Value: proto.String("mesos-exec")},
		Name:   app.StringOpt("n name", "mesos-exec", "Task Name"),
		Command: &mesos.CommandInfo{
			Shell: app.BoolOpt("s shell", false, "Execute as shell command"),
			User:  app.StringOpt("u user", "root", "User to run as"),
		},
		Container: &mesos.ContainerInfo{
			Type: mesos.ContainerInfo_MESOS.Enum(),
		},
	}
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
			task.Container.Docker = &mesos.ContainerInfo_DockerInfo{
				Image: image,
			}
		}
		// Nothing to do if not running a container
		// and no arguments are specified.
		if *image == "" && len(args) == 0 {
			app.PrintHelp()
			cli.Exit(1)
		}
		// Get resources
		task.Resources = resources(*cpus, *mem, *disk)
	}
	app.Action = func() {
		// This is done to satisfy the presumptuous golang/glog package
		// which assumes I am using flag and insists it be configured
		// with such. Since glog is used in go-mesos it is easiest to use
		// the same library for the moment. TODO
		flag.CommandLine.Set("v", *level)
		flag.CommandLine.Set("logtostderr", "1")
		flag.CommandLine.Parse([]string{})
		if err := RunTask(*master, task); err != nil {
			fmt.Errorf("Error: ", err.Error())
			os.Exit(1)
		}
	}
	app.Run(os.Args)
}

// Parse resource string args as float64 since mow.cli doesn't support it.
func resources(cpus, mem, disk string) []*mesos.Resource {
	toFloat := func(s string) float64 {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return v
	}
	return []*mesos.Resource{
		&mesos.Resource{
			Name: proto.String("cpus"),
			Type: mesos.Value_SCALAR.Enum(),
			Scalar: &mesos.Value_Scalar{
				Value: proto.Float64(toFloat(cpus)),
			},
		},
		&mesos.Resource{
			Name: proto.String("mem"),
			Type: mesos.Value_SCALAR.Enum(),
			Scalar: &mesos.Value_Scalar{
				Value: proto.Float64(toFloat(mem)),
			},
		},
		&mesos.Resource{
			Name: proto.String("disk"),
			Type: mesos.Value_SCALAR.Enum(),
			Scalar: &mesos.Value_Scalar{
				Value: proto.Float64(toFloat(disk)),
			},
		},
	}
}

package main

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gogo/protobuf/proto"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"os"
	"strings"
)

func ps(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS]"
	defaults := DefaultProfile()
	var (
		master   = cmd.StringOpt("master", defaults.Master, "Mesos Master")
		limit    = cmd.IntOpt("limit", 2000, "maximum number of tasks to return per request")
		max      = cmd.IntOpt("max", 250, "maximum number of tasks to list")
		order    = cmd.StringOpt("order", "desc", "accending or decending sort order [asc|desc]")
		truncate = cmd.BoolOpt("truncate", true, "truncate some values")

		all         = cmd.BoolOpt("all", false, "Show all tasks")
		frameworkID = cmd.StringOpt("framework", "", "Filter FrameworkID")
		fuzzy       = cmd.BoolOpt("fuzzy", true, "Fuzzy match Task name or Task ID prefix")
		name        = cmd.StringOpt("name", "", "Filter Task name")
		id          = cmd.StringOpt("id", "", "Filter Task ID")
		state       = cmd.StringsOpt("state", []string{"TASK_RUNNING"}, "Filter based on Task state")
	)

	cmd.Before = func() {
		*state = trimFlaged(*state, "--state")
	}

	cmd.Action = func() {
		tasks := make(chan *taskInfo)
		client := &Client{
			Hostname: config.Profile(WithMaster(*master)).Master,
		}
		filters, err := NewTaskFilters(&TaskFilterOptions{
			All:         *all,
			FrameworkID: *frameworkID,
			Fuzzy:       *fuzzy,
			ID:          *id,
			Name:        *name,
			States:      *state,
		})
		failOnErr(err)
		paginator := &TaskPaginator{
			limit: *limit,
			max:   *max,
			order: *order,
			tasks: tasks,
		}
		go func() {
			failOnErr(client.PaginateTasks(paginator, filters...))
		}()
		table := uitable.New()
		table.AddRow("ID", "FRAMEWORK", "STATE", "CPU", "MEM", "GPU", "DISK")
		for task := range tasks {
			frameworkID := task.FrameworkID
			if *truncate {
				frameworkID = truncStr(task.FrameworkID, 8)
			}
			table.AddRow(task.ID, frameworkID, task.State.String(), task.Resources.CPU, task.Resources.Mem, task.Resources.GPUs, task.Resources.Disk)
		}
		fmt.Println(table)
	}
}

func cat(cmd *cli.Cmd) {
	defaults := DefaultProfile()
	cmd.Spec = "[OPTIONS] FILE"
	var (
		master   = cmd.StringOpt("master", defaults.Master, "Mesos Master")
		lines    = cmd.IntOpt("n lines", 0, "Output the last N lines")
		tail     = cmd.BoolOpt("t tail", false, "Tail output")
		filename = cmd.StringArg("FILE", "", "Filename to retrieve")

		all         = cmd.BoolOpt("all", false, "Show all tasks")
		frameworkID = cmd.StringOpt("framework", "", "Filter FrameworkID")
		fuzzy       = cmd.BoolOpt("fuzzy", true, "Fuzzy match Task name or Task ID prefix")
		name        = cmd.StringOpt("name", "", "Filter Task name")
		id          = cmd.StringOpt("id", "", "Filter Task ID")
		state       = cmd.StringsOpt("state", []string{}, "Filter based on Task state")
	)

	cmd.Before = func() {
		*state = trimFlaged(*state, "--state")
	}

	cmd.Action = func() {
		client := &Client{
			Hostname: config.Profile(WithMaster(*master)).Master,
		}
		filters, err := NewTaskFilters(&TaskFilterOptions{
			All:         *all,
			FrameworkID: *frameworkID,
			Fuzzy:       *fuzzy,
			ID:          *id,
			Name:        *name,
			States:      *state,
		})
		failOnErr(err)
		task, err := client.FindTask(filters...)
		failOnErr(err)
		// Attempt to get the full agent state
		agent, err := client.Agent(task.AgentID)
		failOnErr(err)
		// Lookup executor information in agent state
		executor := findExecutor(agent, task.ID)
		if executor == nil {
			failOnErr(fmt.Errorf("could not resolve executor"))
		}
		client = &Client{Hostname: agent.FQDN()}
		var target *fileInfo
		files, err := client.Browse(executor.Directory)
		for _, file := range files {
			if file.Relative() == *filename {
				target = file
			}
		}
		if target == nil {
			failOnErr(fmt.Errorf("cannot find file %s", *filename))
		}
		fp := &FilePaginator{
			data:   make(chan *fileData),
			cancel: make(chan bool),
			path:   target.Path,
			tail:   *tail,
		}
		failOnErr(Monitor(client, os.Stdout, *lines, fp))
	}
}

func ls(cmd *cli.Cmd) {
	defaults := DefaultProfile()
	cmd.Spec = "[OPTIONS]"
	var (
		master   = cmd.StringOpt("master", defaults.Master, "Mesos Master")
		absolute = cmd.BoolOpt("a absolute", false, "Show absolute file paths")

		all         = cmd.BoolOpt("all", false, "Show all tasks")
		frameworkID = cmd.StringOpt("framework", "", "Filter FrameworkID")
		fuzzy       = cmd.BoolOpt("fuzzy", true, "Fuzzy match Task name or Task ID prefix")
		name        = cmd.StringOpt("name", "", "Filter Task name")
		id          = cmd.StringOpt("id", "", "Filter Task ID")
		state       = cmd.StringsOpt("state", []string{}, "Filter based on Task state")
	)

	cmd.Before = func() {
		*state = trimFlaged(*state, "--state")
	}

	cmd.Action = func() {
		client := &Client{
			Hostname: config.Profile(WithMaster(*master)).Master,
		}
		filters, err := NewTaskFilters(&TaskFilterOptions{
			All:         *all,
			FrameworkID: *frameworkID,
			Fuzzy:       *fuzzy,
			ID:          *id,
			Name:        *name,
			States:      *state,
		})
		failOnErr(err)
		// First attempt to resolve the task by ID
		task, err := client.FindTask(filters...)
		failOnErr(err)
		// Attempt to get the full agent state
		agent, err := client.Agent(task.AgentID)
		failOnErr(err)
		// Lookup executor information in agent state
		executor := findExecutor(agent, task.ID)
		if executor == nil {
			failOnErr(fmt.Errorf("could not resolve executor"))
		}
		client = &Client{Hostname: agent.FQDN()}
		files, err := client.Browse(executor.Directory)
		failOnErr(err)
		table := uitable.New()
		table.AddRow("UID", "GID", "MODE", "MODIFIED", "SIZE", "PATH")
		for _, file := range files {
			path := file.Relative()
			if *absolute {
				path = file.Path
			}
			table.AddRow(file.UID, file.GID, file.Mode, file.Modified().String(), fmt.Sprintf("%d", file.Size), path)
		}
		fmt.Println(table)
	}
}

func topCmd(cmd *cli.Cmd) {
	defaults := DefaultProfile()
	cmd.Spec = "[OPTIONS]"
	var master = cmd.StringOpt("master", defaults.Master, "Mesos Master")
	cmd.Action = func() {
		client := &Client{
			Hostname: config.Profile(WithMaster(*master)).Master,
		}
		failOnErr(RunTop(client))
	}
}

/*
func agents(cmd *cli.Cmd) {
	defaults := DefaultProfile()
	cmd.Spec = "[OPTIONS]"
	var master = cmd.StringOpt("master", defaults.Master, "Mesos Master")
	cmd.Action = func() {
		client := &Client{
			Hostname: config.Profile(WithMaster(*master)).Master,
		}
		agents, err := Agents(client)
		failOnErr(err)
		table := uitable.New()
		table.AddRow("ID", "FQDN", "VERSION", "UPTIME", "CPUS", "MEM", "GPUS", "DISK")
		for _, agent := range agents {
			table.AddRow(
				agent.ID,
				agent.FQDN(),
				agent.Version,
				agent.Uptime().String(),
				fmt.Sprintf("%.2f/%.2f", agent.UsedResources.CPU, agent.Resources.CPU),
				fmt.Sprintf("%.2f/%.2f", agent.UsedResources.Mem, agent.Resources.Mem),
				fmt.Sprintf("%.2f/%.2f", agent.UsedResources.GPUs, agent.Resources.GPUs),
				fmt.Sprintf("%.2f/%.2f", agent.UsedResources.Disk, agent.Resources.Disk),
			)
		}
		fmt.Println(table)
	}
}
*/

func run(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS] [ARG...]"
	defaults := DefaultProfile()
	var (
		master     = cmd.StringOpt("master", defaults.Master, "Mesos Master")
		arguments  = cmd.StringsArg("ARG", nil, "Command Arguments")
		taskPath   = cmd.StringOpt("task", "", "Path to a Mesos TaskInfo JSON file")
		parameters = cmd.StringsOpt("param", []string{}, "Docker parameters")
		image      = cmd.StringOpt("i image", "", "Docker image to run")
		volumes    = cmd.StringsOpt("v volume", []string{}, "Volume mappings")
		ports      = cmd.StringsOpt("p ports", []string{}, "Port mappings")
		envs       = cmd.StringsOpt("e env", []string{}, "Environment Variables")
		shell      = cmd.StringOpt("s shell", "", "Shell command to execute")
		tail       = cmd.BoolOpt("t tail", false, "Tail command output")
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
	cmd.Action = func() {
		failOnErr(RunTask(config.Profile(WithMaster(*master)), task, *tail))
	}
}

const (
	repository    string = "quay.io/vektorcloud/mesos:latest"
	containerName string = "mesos_cli"
)

// local attempts to launch a local Mesos cluster
// with github.com/vektorcloud/mesos.
func local(cmd *cli.Cmd) {
	var (
		container *docker.APIContainers
		image     *docker.APIImages
	)
	cmd.Spec = "[OPTIONS]"
	up := func(cmd *cli.Cmd) {
		var (
			remove = cmd.BoolOpt("rm remove", false, "Remove any existing local cluster")
			force  = cmd.BoolOpt("f force", false, "Force pull a new image from vektorcloud")
		)
		cmd.Action = func() {
			client, err := docker.NewClientFromEnv()
			failOnErr(err)
			image = getImage(repository, client)
			if image == nil || *force {
				failOnErr(client.PullImage(docker.PullImageOptions{Repository: repository}, docker.AuthConfiguration{}))
			}
			image = getImage(repository, client)
			if image == nil {
				failOnErr(fmt.Errorf("Cannot pull image %s", repository))
			}
			container = getContainer(containerName, client)
			if container != nil && *remove {
				failOnErr(client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true}))
				container = nil
			}
			if container == nil {
				_, err = client.CreateContainer(
					docker.CreateContainerOptions{
						Name: containerName,
						HostConfig: &docker.HostConfig{
							NetworkMode: "host",
							Binds: []string{
								"/var/run/docker.sock:/var/run/docker.sock:rw",
							},
						},
						Config: &docker.Config{
							Cmd:   []string{"mesos-local"},
							Image: repository,
						}})
				failOnErr(err)
				container = getContainer(containerName, client)
			}
			failOnErr(client.StartContainer(container.ID, &docker.HostConfig{}))
		}
	}
	down := func(cmd *cli.Cmd) {
		cmd.Action = func() {
			client, err := docker.NewClientFromEnv()
			failOnErr(err)
			if container = getContainer(containerName, client); container != nil {
				if container.State != "running" {
					fmt.Printf("container is in invalid state: %s\n", container.State)
					cli.Exit(1)
				}
			}
			fmt.Println("no countainer found")
			cli.Exit(1)
		}
	}
	status := func(cmd *cli.Cmd) {
		cmd.Action = func() {
			client, err := docker.NewClientFromEnv()
			failOnErr(err)
			if container = getContainer(containerName, client); container != nil {
				fmt.Printf("%s: %s\n", container.ID, container.State)
			} else {
				fmt.Println("no container found")
			}
			cli.Exit(0)
		}
	}
	rm := func(cmd *cli.Cmd) {
		cmd.Action = func() {
			client, err := docker.NewClientFromEnv()
			failOnErr(err)
			if container = getContainer(containerName, client); container != nil {
				fmt.Printf("removing container %s\n", container.ID)
				failOnErr(client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true}))
				cli.Exit(0)
			}
			fmt.Println("no container found")
			cli.Exit(1)
		}
	}
	cmd.Command("up", "Start the local cluster", up)
	cmd.Command("down", "Stop the local cluster", down)
	cmd.Command("status", "Display the status of the local cluster", status)
	cmd.Command("rm", "Remove the local cluster", rm)
}

func getImage(n string, client *docker.Client) *docker.APIImages {
	images, err := client.ListImages(docker.ListImagesOptions{All: true})
	failOnErr(err)
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == n {
				return &image
			}
		}
	}
	return nil
}

func getContainer(n string, client *docker.Client) *docker.APIContainers {
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	failOnErr(err)
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Replace(name, "/", "", 1) == n {
				return &container
			}
		}
	}
	return nil
}

package commands

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/config"
	"strings"
)

const (
	repository    string = "quay.io/vektorcloud/mesos:latest"
	containerName string = "mesos_cli"
)

type Local struct {
	*command
}

func NewLocal() Command {
	return Local{
		&command{"local", "Run a local Mesos cluster"},
	}
}

func (local Local) Init(cfg config.CfgFn) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
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
								Env:   []string{"MESOS_LOGGING_LEVEL=INFO"},
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

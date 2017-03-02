package local

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	docker "github.com/docker/docker/client"
	"go.uber.org/zap"
	"strings"
	"time"
)

// Client is a simple high level client for
// managing the quay.io/vektorcloud/mesos container.
type Client struct {
	docker *docker.Client
	log    *zap.Logger
}

func (c Client) FindImage(name string) (*types.ImageSummary, error) {
	var image *types.ImageSummary
	images, err := c.docker.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		return nil, err
	}
	for _, img := range images {
		image = &img
		c.log.Info("docker", zap.String("message", fmt.Sprintf("Found Docker image %s", image.ID)))
		break
	}
	return image, nil
}

func (c Client) FindContainer(name, id string) (*types.Container, error) {
	if !strings.HasPrefix(name, "/") {
		name = fmt.Sprintf("/%s", name)
	}
	containers, err := c.docker.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}
	var container *types.Container
loop:
	for _, cont := range containers {
		if cont.ID == id {
			container = &cont
			break
		}
		for _, n := range cont.Names {
			if name == n {
				container = &cont
				break loop
			}
		}
	}
	if container != nil {
		c.log.Info("docker",
			zap.String("message", fmt.Sprintf("Found container %s", container.ID)),
			zap.String("state", container.State),
		)
	}
	return container, nil
}

func (c Client) PullImage(ref string) error {
	c.log.Info("docker", zap.String("message", fmt.Sprintf("Pulling image %s", ref)))
	_, err := c.docker.ImagePull(context.Background(), ref, types.ImagePullOptions{All: true}) // TODO ?
	c.log.Info("docker", zap.String("message", "Pull complete"))
	return err
}

func (c Client) RemoveContainer(id string, force bool) error {
	c.log.Info("docker", zap.String("message", fmt.Sprintf("Removing container %s", id)))
	return c.docker.ContainerRemove(
		context.Background(),
		id,
		types.ContainerRemoveOptions{Force: force},
	)
}

func (c Client) CreateContainer(name string, image string, envs []string) (*types.Container, error) {
	c.log.Info("docker", zap.String("message", fmt.Sprintf("Creating new container %s", name)))
	resp, err := c.docker.ContainerCreate(
		context.Background(),
		&container.Config{
			Cmd:   strslice.StrSlice{"mesos-local"},
			Image: image,
			Env:   envs,
		},
		&container.HostConfig{
			NetworkMode: container.NetworkMode("host"),
			Binds:       []string{"/var/run/docker.sock:/var/run/docker.sock"},
		},
		&network.NetworkingConfig{}, name,
	)
	if err != nil {
		return nil, err
	}
	container, err := c.FindContainer("", resp.ID)
	if err != nil {
		return nil, err
	}
	if container == nil {
		return nil, fmt.Errorf("could not create container")
	}
	return container, nil
}

func (c Client) StartContainer(id string) error {
	c.log.Info("docker", zap.String("message", fmt.Sprintf("Starting container %s", id)))
	return c.docker.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
}

func (c Client) StopContainer(id string) error {
	c.log.Info("docker", zap.String("message", fmt.Sprintf("Stopping container %s", id)))
	d := 120 * time.Second
	return c.docker.ContainerStop(context.Background(), id, &d)
}

package local

import (
	"context"
	"fmt"
	docker "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	"github.com/docker/engine-api/types/strslice"
	"go.uber.org/zap"
	"strings"
)

const Timeout = 120

// Client is a simple high level client for
// managing the quay.io/vektorcloud/mesos container.
type Client struct {
	docker *docker.Client
	log    *zap.Logger
}

func (c Client) FindImage(name string) (*types.Image, error) {
	var image *types.Image
	images, err := c.docker.ImageList(context.Background(), types.ImageListOptions{MatchName: name})
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

func (c Client) PullImage(id, tag string) error {
	c.log.Info("docker", zap.String("message", fmt.Sprintf("Pulling image %s:%s", id, tag)))
	_, err := c.docker.ImagePull(context.Background(), types.ImagePullOptions{ImageID: id, Tag: tag}, nil)
	c.log.Info("docker", zap.String("message", "Pull complete"))
	return err
}

func (c Client) RemoveContainer(id string, force bool) error {
	c.log.Info("docker", zap.String("message", fmt.Sprintf("Removing container %s", id)))
	return c.docker.ContainerRemove(
		context.Background(),
		types.ContainerRemoveOptions{ContainerID: id, Force: force},
	)
}

func (c Client) CreateContainer(name string, image *types.Image, envs []string) (*types.Container, error) {
	c.log.Info("docker", zap.String("message", fmt.Sprintf("Creating new container %s", name)))
	resp, err := c.docker.ContainerCreate(
		context.Background(),
		&container.Config{
			Cmd:   strslice.StrSlice{"mesos-local"},
			Image: image.ID,
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
	return c.docker.ContainerStart(context.Background(), id)
}

func (c Client) StopContainer(id string) error {
	c.log.Info("docker", zap.String("message", fmt.Sprintf("Stopping container %s", id)))
	return c.docker.ContainerStop(context.Background(), id, Timeout)
}

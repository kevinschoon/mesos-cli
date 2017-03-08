package local

import (
	"fmt"
	docker "github.com/docker/docker/client"
	"github.com/jawher/mow.cli"
	"github.com/vektorlab/mesos-cli/config"
	"os"
)

const (
	Repository    string = "quay.io/vektorcloud/mesos:latest"
	ContainerName string = "mesos_cli"
)

// TODO: Add support for passing environment variables
// to the mesos container.
// Add support running in the "foreground" where stdout/stderr
// are streamed back to the caller.

type Local struct{}

func (_ Local) Name() string { return "local" }
func (_ Local) Desc() string { return "Run a local Mesos cluster" }

func (_ Local) Init(profile config.ProfileFn) func(*cli.Cmd) {
	var client *Client
	fn := func() *Client { return client }
	return func(cmd *cli.Cmd) {
		cmd.Before = func() {
			cli, err := docker.NewEnvClient()
			failOnErr(err)
			client = &Client{cli, profile().Log()}
		}
		cmd.Command("up", "Start a local Mesos cluster", Up(fn))
		cmd.Command("down", "Shutdown a running cluster", Down(fn))
		cmd.Command("status", "Return the status of the running cluster", Status(fn))
		cmd.Command("rm", "Remove the existing cluster", Remove(fn))
	}
}

func failOnErr(err error) {
	if err != nil {
		fmt.Printf("Encountered Error: %v\n", err)
		os.Exit(2)
	}
}

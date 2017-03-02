package local

import (
	"github.com/jawher/mow.cli"
)

func Up(fn func() *Client) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		var (
			force  = cmd.BoolOpt("f force", false, "Force pull the Mesos image")
			remove = cmd.BoolOpt("rm remove", false, "Remove any existing container")
			envs   = cmd.StringsOpt("e env", []string{"MESOS_LOGGING_LEVEL=INFO"}, "Environment variables")
		)
		cmd.Action = func() {
			client := fn()
			image, err := client.FindImage(Repository)
			failOnErr(err)
			if image == nil || *force {
				failOnErr(client.PullImage(Repository, "latest"))
				image, err = client.FindImage(Repository)
				failOnErr(err)
			}
			container, err := client.FindContainer(ContainerName, "")
			failOnErr(err)
			if *remove {
				if container != nil {
					failOnErr(client.RemoveContainer(container.ID, *force))
					container, err = client.CreateContainer(ContainerName, image, *envs)
					failOnErr(err)
				}
			}
			if container == nil {
				container, err = client.CreateContainer(ContainerName, image, *envs)
				failOnErr(err)
			}
			failOnErr(client.StartContainer(container.ID))
		}
	}
}

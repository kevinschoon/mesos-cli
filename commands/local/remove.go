package local

import (
	"fmt"
	"github.com/jawher/mow.cli"
)

func Remove(fn func() *Client) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		var (
			force = cmd.BoolOpt("force", false, "Force removal of the container")
		)
		cmd.Action = func() {
			client := fn()
			container, err := client.FindContainer(ContainerName, "")
			failOnErr(err)
			if container == nil {
				failOnErr(fmt.Errorf("Container not found"))
			}
			failOnErr(client.RemoveContainer(container.ID, *force))
		}
	}
}

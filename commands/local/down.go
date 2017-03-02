package local

import (
	"fmt"
	"github.com/jawher/mow.cli"
)

func Down(fn func() *Client) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			client := fn()
			container, err := client.FindContainer(ContainerName, "")
			failOnErr(err)
			if container == nil {
				failOnErr(fmt.Errorf("Container not found"))
			}
			if container.State == "stopped" {
				failOnErr(fmt.Errorf("Container already stopped"))
			}
			failOnErr(client.StopContainer(container.ID))
		}
	}
}

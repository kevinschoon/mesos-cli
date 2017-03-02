package local

import (
	"fmt"
	"github.com/jawher/mow.cli"
)

func Status(fn func() *Client) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			client := fn()
			cont, err := client.FindContainer(ContainerName, "")
			failOnErr(err)
			if cont == nil {
				failOnErr(fmt.Errorf("Container Not Found"))
			}
			fmt.Println(cont.State)
		}
	}
}

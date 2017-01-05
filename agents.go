package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
)

func agents(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS]"
	cmd.Action = func() {
		client := &Client{handler: DefaultHandler{hostname: config.Profile().Master}}
		agents, err := client.Agents()
		failOnErr(err)
		for _, agent := range agents {
			fmt.Println(*agent.Id.Value, *agent.Hostname)
		}
	}
}

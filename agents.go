package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"net/url"
)

func agents(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS]"
	defaults := DefaultProfile()
	var master = cmd.StringOpt("master", defaults.Master, "Mesos Master")
	cmd.Action = func() {
		client := &Client{
			Hostname: config.Profile(WithMaster(*master)).Master,
		}
		agents := struct {
			Agents []struct {
				ID       *string `json:"id"`
				Hostname *string `json:"hostname"`
			} `json:"slaves"`
		}{}
		failOnErr(client.Get(&url.URL{Path: "/master/slaves"}, &agents))
		for _, agent := range agents.Agents {
			fmt.Println(*agent.ID, *agent.Hostname)
		}
	}
}

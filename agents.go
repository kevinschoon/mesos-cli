package main

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	"net/url"
	"time"
)

type agentInfo struct {
	ID               string  `json:"id"`
	Hostname         string  `json:"hostname"`
	RegisteredTime   float64 `json:"registered_time"`
	ReRegisteredTime float64 `json:"reregistered_time"`
	Version          string  `json:"version"`
	Resources        struct {
		CPU  float64 `json:"cpus"`
		Mem  float64 `json:"mem"`
		Disk float64 `json:"disk"`
		GPUs float64 `json:"gpus"`
	} `json:"resources"`
	UsedResources struct {
		CPU  float64 `json:"cpus"`
		Mem  float64 `json:"mem"`
		Disk float64 `json:"disk"`
		GPUs float64 `json:"gpus"`
	} `json:"used_resources"`
}

func (a *agentInfo) Registered() time.Time {
	return time.Unix(int64(a.RegisteredTime), 0)
}

func (a *agentInfo) ReRegistered() time.Time {
	return time.Unix(int64(a.ReRegisteredTime), 0)
}

func (a *agentInfo) Uptime() time.Duration {
	return time.Since(a.ReRegistered())
}

func agents(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS]"
	defaults := DefaultProfile()
	var master = cmd.StringOpt("master", defaults.Master, "Mesos Master")
	cmd.Action = func() {
		client := &Client{
			Hostname: config.Profile(WithMaster(*master)).Master,
		}
		agents := struct {
			Agents []*agentInfo `json:"slaves"`
		}{}
		failOnErr(client.Get(&url.URL{Path: "/master/slaves"}, &agents))
		table := uitable.New()
		table.AddRow("ID", "HOSTNAME", "VERSION", "UPTIME", "CPUS", "MEM", "GPUS", "DISK")
		for _, agent := range agents.Agents {
			table.AddRow(
				agent.ID,
				agent.Hostname,
				agent.Version,
				agent.Uptime().String(),
				fmt.Sprintf("%.2f/%.2f", agent.UsedResources.CPU, agent.Resources.CPU),
				fmt.Sprintf("%.2f/%.2f", agent.UsedResources.Mem, agent.Resources.Mem),
				fmt.Sprintf("%.2f/%.2f", agent.UsedResources.GPUs, agent.Resources.GPUs),
				fmt.Sprintf("%.2f/%.2f", agent.UsedResources.Disk, agent.Resources.Disk),
			)
		}
		fmt.Println(table)
	}
}

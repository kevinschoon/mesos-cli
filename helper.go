package main

import (
	"fmt"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

// Agents returns a map of IDs to hostnames
func Agents(master string) (map[string]string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/master/slaves", master))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	agents := map[string]string{}
	for _, agent := range gjson.GetBytes(raw, "slaves").Array() {
		agents[agent.Get("id").Str] = agent.Get("hostname").Str
	}
	return agents, nil
}

// LogDir returns the directory path for following task output
func LogDir(hostname, executorId string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:5051/slave(1)/state", hostname))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	for _, framework := range gjson.GetBytes(raw, "frameworks").Array() {
		for _, executor := range framework.Get("executors").Array() {
			if executor.Get("id").Str == executorId {
				return executor.Get("directory").Str, nil
			}
		}
	}
	return "", fmt.Errorf("Unable to find log directory")
}

// Resource returns the value of a resource
func Resource(name string, resources []*mesos.Resource) float64 {
	var value float64
	for _, resource := range resources {
		if resource.GetName() == name {
			value = resource.GetScalar().GetValue()
		}
	}
	return value
}

// Check if a Mesos resource offer can satisfy the Task
func Sufficent(task *mesos.TaskInfo, offer *mesos.Offer) bool {
	for _, resource := range offer.Resources {
		value := resource.GetScalar().GetValue()
		switch resource.GetName() {
		case "cpus":
			if value < Resource("cpus", task.Resources) {
				return false
			}
		case "mem":
			if value < Resource("mem", task.Resources) {
				return false
			}
		case "disk":
			if value < Resource("disk", task.Resources) {
				return false
			}
		}
	}
	return true
}

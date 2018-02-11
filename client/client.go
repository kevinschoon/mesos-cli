package client

import (
	"context"
	"fmt"

	"github.com/mesanine/mesos-cli/config"
	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/httpcli"

	"github.com/mesos/mesos-go/api/v1/lib/httpcli/httpmaster"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	mcalls "github.com/mesos/mesos-go/api/v1/lib/master/calls"

	"github.com/mesos/mesos-go/api/v1/lib/agent"
	acalls "github.com/mesos/mesos-go/api/v1/lib/agent/calls"
	"github.com/mesos/mesos-go/api/v1/lib/httpcli/httpagent"
)

func NewAgent(profile *config.Profile, info *mesos.AgentInfo) acalls.Sender {
	endpoint := profile.Endpoint()
	endpoint.Host = fmt.Sprintf("%s:%d", info.GetHostname(), info.GetPort())
	endpoint.Path = config.OperatorAPIPath
	client := httpagent.NewSender(
		httpcli.New(httpcli.Endpoint(endpoint.String())).Send,
	)
	return client
}

type Operator struct {
	master mcalls.Sender
}

func NewOperator(profile *config.Profile) *Operator {
	return &Operator{
		master: httpmaster.NewSender(
			httpcli.New(httpcli.Endpoint(profile.Endpoint().String())).Send,
		),
	}
}

func (o Operator) CallMaster(call *master.Call) (*master.Response, error) {
	resp, err := o.master.Send(
		context.Background(),
		mcalls.NonStreaming(call),
	)
	if err != nil {
		return nil, err
	}
	mresp := &master.Response{}
	err = resp.Decode(mresp)
	if err != nil {
		return nil, err
	}
	return mresp, resp.Close()
}

func (o Operator) CallAgent(call *agent.Call) (*agent.Response, error) {
	return nil, nil
}

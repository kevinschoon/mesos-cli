package helper

import (
	"bytes"
	"fmt"
	"github.com/mesanine/mesos-cli/config"
	"github.com/mesos/mesos-go"
	"github.com/mesos/mesos-go/httpcli"
	"github.com/mesos/mesos-go/httpcli/operator"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
)

func newCaller(endpoint *url.URL, log *zap.Logger) operator.Caller {
	return operator.NewCaller(
		httpcli.New(
			httpcli.Endpoint(endpoint.String()),
			httpcli.RequestOptions(func(req *http.Request) {
				buf, _ := ioutil.ReadAll(req.Body)
				req.Body.Close()
				log.Debug(
					"http request",
					zap.String("url", req.URL.String()),
					zap.String("body", string(buf)),
				)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			},
			),
		))
}

func NewCaller(profile *config.Profile) operator.Caller {
	return newCaller(profile.Endpoint(), profile.Log())
}

func NewAgentCaller(profile *config.Profile, agent *mesos.AgentInfo) operator.Caller {
	endpoint := profile.Endpoint()
	endpoint.Host = fmt.Sprintf("%s:%d", agent.Hostname, agent.GetPort())
	endpoint.Path = config.OperatorAPIPath
	//endpoint.Path = fmt.Sprintf("slave(1)%s", config.OperatorAPIPath)
	return newCaller(endpoint, profile.Log())
}

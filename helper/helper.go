package helper

/*
import (
	"bytes"
	"fmt"
	"github.com/mesanine/mesos-cli/config"
	mesos "github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/httpcli"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
)

func newClient(endpoint *url.URL, log *zap.Logger) *httpcli.Client {
	return httpcli.New(
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
		}),
	)
}

func NewClient(profile *config.Profile) *httpcli.Client {
	return newClient(profile.Endpoint(), profile.Log())
}

func NewAgentClient(profile *config.Profile) *httpcli.Client {
	endpoint := profile.Endpoint()
	endpoint.Host = fmt.Sprintf("%s:%d", agent.Hostname, agent.GetPort())
	endpoint.Path = config.OperatorAPIPath
	//endpoint.Path = fmt.Sprintf("slave(1)%s", config.OperatorAPIPath)
	return newClient(endpoint, profile.Log())
}
*/

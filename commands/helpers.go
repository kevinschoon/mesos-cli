package commands

import (
	"bytes"
	"fmt"
	"github.com/mesos/mesos-go"
	"github.com/mesos/mesos-go/httpcli"
	"github.com/mesos/mesos-go/httpcli/operator"
	"github.com/vektorlab/mesos-cli/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func NewCaller(profile *config.Profile) operator.Caller {
	endpoint := url.URL{
		Scheme: profile.Scheme,
		Host:   profile.Master,
		Path:   config.OperatorAPIPath,
	}
	return operator.NewCaller(
		httpcli.New(
			httpcli.Endpoint(endpoint.String()),
			httpcli.RequestOptions(func(req *http.Request) {
				buf, _ := ioutil.ReadAll(req.Body)
				req.Body.Close()
				fmt.Println(string(buf))
				req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			},
			),
		))
}

func Scalar(name string, resources mesos.Resources) (v float64) {
	if scalar := resources.SumScalars(mesos.NamedResources(name)); scalar != nil {
		v = scalar.Value
	}
	return v
}

func truncStr(s string, l int) string {
	runes := bytes.Runes([]byte(s))
	if len(runes) < l {
		return s
	}
	return string(runes[:l])
}

func failOnErr(err error) {
	if err != nil {
		fmt.Printf("Encountered Error: %v\n", err)
		os.Exit(2)
	}
}

package commands

import (
	"bytes"
	"fmt"
	"github.com/mesos/mesos-go"
	"os"
)

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

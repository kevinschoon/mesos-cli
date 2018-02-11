package command

import (
	"fmt"

	"github.com/jawher/mow.cli"
	"github.com/mesanine/mesos-cli/command/flags"
	"os"
	"reflect"
)

/*
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
*/

func setFlags(call interface{}, cmd *cli.Cmd) {
	cv := reflect.ValueOf(call).Elem()
	for i := 0; i < cv.NumField(); i++ {
		name, flag, desc := flags.ToValue(cv.Field(i))
		cmd.VarOpt(name, flag, desc)
	}
}

func notImplemented(name string) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Action = func() {
			maybe(fmt.Errorf("%s is not yet implemented", name))
		}
	}
}

func maybe(err error) {
	if err != nil {
		fmt.Printf("Encountered Error: %v\n", err)
		os.Exit(2)
	}
}

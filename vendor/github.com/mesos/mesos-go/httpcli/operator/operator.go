package operator

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go"
	"github.com/mesos/mesos-go/agent"
	"github.com/mesos/mesos-go/encoding"
	"github.com/mesos/mesos-go/httpcli"
	"github.com/mesos/mesos-go/master"
)

// ErrInvalidCallType is caused by a miss-match of
// Agent or Master call types possibly indicating
// an API change between Mesos versions.
type ErrInvalidCallType struct {
	Name  string
	Value int32
}

func (e ErrInvalidCallType) Error() string {
	return fmt.Sprintf("ErrInvalidCallType %s %d", e.Name, e.Value)
}

type (
	client struct {
		*httpcli.Client
	}
	Caller interface {
		// Call issues an HTTP call to a Mesos Master or Agent operator API depending
		// on the type of message supplied. If the type is neither master.Call
		// or agent.Call we will panic.
		Call(proto.Message) (mesos.Response, error)
		// CallAgent issues a call to the Mesos Agent operator API
		CallAgent(call *agent.Call) (*agent.Response, error)
		// CallMaster issues a call to the Mesos Master operator API
		CallMaster(call *master.Call) (*master.Response, error)
	}
)

func (cli *client) Call(msg proto.Message) (mesos.Response, error) {
	switch call := msg.(type) {
	case *agent.Call:
		return cli.Do(call)
	case *master.Call:
		return cli.Do(call)
	}
	panic("Unknown caller message type")
}

func (cli *client) CallAgent(call *agent.Call) (*agent.Response, error) {
	callType, ok := agent.Call_Type_name[int32(*call.Type)]
	if !ok {
		return nil, ErrInvalidCallType{Value: int32(*call.Type)}
	}
	resp, err := cli.Call(call)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	message := &agent.Response{}
	err = resp.Decoder().Invoke(message)
	if err != nil {
		return nil, err
	}
	return message, checkType(callType, int32(*message.Type), agent.Response_Type_name)
}

func (cli *client) CallMaster(call *master.Call) (*master.Response, error) {
	callType, ok := master.Call_Type_name[int32(*call.Type)]
	if !ok {
		return nil, ErrInvalidCallType{Value: int32(*call.Type)}
	}
	resp, err := cli.Call(call)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	message := &master.Response{}
	err = resp.Decoder().Invoke(message)
	if err != nil {
		return nil, err
	}
	return message, checkType(callType, int32(*message.Type), master.Response_Type_name)
}

func NewCaller(cli *httpcli.Client) Caller {
	cli.With(
		httpcli.Codec(&encoding.JSONCodec),
	)
	return &client{cli}
}

// checkType validates the callType matches the respType
func checkType(callType string, respType int32, types map[int32]string) error {
	name, ok := types[respType]
	if !ok {
		return ErrInvalidCallType{Value: respType}
	}
	if name != callType {
		return ErrInvalidCallType{Name: name, Value: respType}
	}
	return nil
}

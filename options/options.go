package options

import (
	"fmt"
	"github.com/mesos/mesos-go"
	"github.com/satori/go.uuid"
	"strconv"
)

type TaskID struct {
	ID *mesos.TaskID
}

func NewTaskID() *TaskID {
	return &TaskID{
		ID: &mesos.TaskID{
			Value: uuid.NewV4().String(),
		},
	}
}

func (o *TaskID) String() string {
	return o.ID.Value
}

func (o *TaskID) Set(name string) error {
	o.ID.Value = name
	return nil
}

type States []*mesos.TaskState

func NewStates() States {
	return States{mesos.TASK_RUNNING.Enum()}
}

func (o *States) String() string {
	return fmt.Sprintf("%v", *o)
}

func (o *States) Set(name string) error {
	v, ok := mesos.TaskState_value[name]
	if !ok {
		return fmt.Errorf("Invalid state %s", name)
	}
	*o = append(*o, mesos.TaskState(v).Enum())
	return nil
}

func (o *States) Clear() {
	*o = []*mesos.TaskState{}
}

type ScalarResources struct {
	Name      string
	Resources mesos.Resources
}

func NewScalarResources(name string, resources mesos.Resources) *ScalarResources {
	return &ScalarResources{
		Name:      name,
		Resources: resources,
	}
}

func (o *ScalarResources) String() string {
	var value float64
	if scalar := o.Resources.SumScalars(mesos.NamedResources(o.Name)); scalar != nil {
		value = scalar.Value
	}
	return fmt.Sprintf("%.2f", value)
}

func (o *ScalarResources) Set(v string) error {
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	o.Resources = append(o.Resources, mesos.Resource{
		Name:   o.Name,
		Type:   mesos.SCALAR.Enum(),
		Scalar: &mesos.Value_Scalar{Value: value},
	})
	return nil
}

package main

import (
	"flag"
	"fmt"
	"github.com/mesos/mesos-go"
	"github.com/mesos/mesos-go/encoding"
	"github.com/mesos/mesos-go/httpcli"
	"github.com/mesos/mesos-go/httpcli/operator"
	"github.com/mesos/mesos-go/master"
	"github.com/mesos/mesos-go/master/calls"
	"github.com/mesos/mesos-go/master/events"
	"os"
)

func subscribe() {
	eventLoop := func(decoder encoding.Decoder, h events.Handler) (err error) {
		for err == nil {
			var e master.Event
			if err = decoder.Invoke(&e); err == nil {
				err = h.HandleEvent(&e)
			}
		}
		return err
	}

	handler := func() events.Handler {
		return events.NewMux(
			events.Handle(master.Event_AGENT_ADDED, events.HandlerFunc(func(e *master.Event) error {
				fmt.Printf("Agent Added: %s\n", e.GetAgentAdded().GetAgent().GetAgentInfo().GetID().Value)
				return nil
			})),
			events.Handle(master.Event_AGENT_REMOVED, events.HandlerFunc(func(e *master.Event) error {
				fmt.Printf("Agent Removed: %s\n", e.GetAgentRemoved().GetAgentId().Value)
				return nil
			})),
			events.Handle(master.Event_TASK_ADDED, events.HandlerFunc(func(e *master.Event) error {
				fmt.Printf("Task Added: %s", e.GetTaskAdded().GetTask().GetName())
				return nil
			})),
			events.Handle(master.Event_TASK_UPDATED, events.HandlerFunc(func(e *master.Event) error {
				taskID := e.GetTaskUpdated().GetStatus().GetTaskID().Value
				state := mesos.TaskState_name[int32(e.GetTaskUpdated().GetState())]
				fmt.Printf("Task Updated: %s [%s]\n", taskID, state)
				return nil
			})),
			events.Handle(master.Event_SUBSCRIBED, events.HandlerFunc(func(e *master.Event) error {
				fmt.Println("Subscribed")
				return nil
			})),
		)
	}
	var (
		client = httpcli.New(
			httpcli.Endpoint("http://localhost:5050/api/v1"),
			// httpcli.Codec(&encoding.FramingJSONCodec),
		)
		subscribe = calls.Subscribe()
	)
	resp, err := client.Do(subscribe, httpcli.Close(true))
	if resp != nil {
		defer resp.Close()
	}
	if err == nil {
		// Connected
		err = eventLoop(resp.Decoder(), handler())
	}
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func tasks() {
	op := operator.NewCaller(
		httpcli.New(
			httpcli.Endpoint("http://localhost:5050/api/v1"),
		),
	)
	resp, err := op.CallMaster(calls.GetTasks())
	if err != nil {
		fmt.Println("Error: ", err)
	}
	for _, task := range resp.GetTasks.Tasks {
		fmt.Printf("%s - %s: %s\n", task.Name, task.TaskID.Value, mesos.TaskState_name[int32(*task.State)])
	}
}

func main() {
	doSubscribe := flag.Bool("subscribe", false, "demo subscribe")
	flag.Parse()
	if *doSubscribe {
		subscribe()
	} else {
		tasks()
	}
}

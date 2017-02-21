package runner

import (
	"fmt"
	"github.com/mesos/mesos-go"
	"github.com/mesos/mesos-go/backoff"
	"github.com/mesos/mesos-go/encoding"
	"github.com/mesos/mesos-go/extras/scheduler/controller"
	"github.com/mesos/mesos-go/httpcli"
	"github.com/mesos/mesos-go/httpcli/httpsched"
	"github.com/mesos/mesos-go/scheduler"
	"github.com/mesos/mesos-go/scheduler/calls"
	"github.com/mesos/mesos-go/scheduler/events"
	"github.com/vektorlab/mesos-cli/config"
	"math/rand"
	"net/http"
	"time"
)

func LoggingCaller() calls.Decorator {
	return func(h calls.Caller) calls.Caller {
		return calls.CallerFunc(func(c *scheduler.Call) (mesos.Response, error) {
			fmt.Println("Calling: ", c)
			return h.Call(c)
		})
	}
}

type ErrRunnerFailed struct{ *scheduler.Event }

func (e ErrRunnerFailed) Error() string {
	et := scheduler.Event_Type_name[int32(*e.Type)]
	msg := fmt.Sprintf("Runner has failed (%s)", et)
	switch *e.Type {
	case scheduler.Event_ERROR:
		msg += fmt.Sprintf("\n %s", e.GetError().Message)
	case scheduler.Event_FAILURE:
		agentID := e.GetFailure().AgentID.Value
		executorID := e.GetFailure().ExecutorID.Value
		msg += fmt.Sprintf("\n agentID: %s, executorID: %s", agentID, executorID)
	}
	return msg
}

// Runner implements a high level client for running
// a task on Mesos. It accepts a mesos.TaskInfo, schedules
// it, blocks until it is complete, and returns nil if it
// completed successfully.
type Runner struct {
	shutdown  chan (struct{})
	caller    calls.Caller
	random    *rand.Rand
	framework *mesos.FrameworkInfo
	context   *controller.ContextAdapter
	done      bool
	task      *mesos.TaskInfo
	status    *mesos.TaskStatus
	scheduled bool
}

func (r Runner) getCaller() calls.Caller { return r.caller }

func (r *Runner) Mux() *events.Mux {
	return events.NewMux(
		events.DefaultHandler(events.HandlerFunc(controller.DefaultHandler)),
		events.MapFuncs(map[scheduler.Event_Type]events.HandlerFunc{
			scheduler.Event_ERROR: func(e *scheduler.Event) error {
				return ErrRunnerFailed{e}
			},
			scheduler.Event_FAILURE: func(e *scheduler.Event) error {
				return ErrRunnerFailed{e}
			},
			scheduler.Event_OFFERS: func(e *scheduler.Event) error {
				// Magic begins here
				offers := e.GetOffers().GetOffers()
				for _, offer := range offers {
					opts := calls.OfferOperations{}
					if mesos.Resources(offer.Resources).
						ContainsAll(mesos.Resources(r.task.Resources)) && !r.scheduled {
						// AgentID is assigned to TaskInfo
						r.task.AgentID = offer.GetAgentID()
						opts = append(opts, calls.OpLaunch(*r.task))
						r.scheduled = true
					}
					accept := calls.Accept(
						opts.WithOffers(offer.ID),
					).With(calls.RefuseSecondsWithJitter(r.random, 15*time.Second))
					_, err := r.caller.Call(accept)
					if err != nil {
						return err
					}
				}
				return nil
			},
			scheduler.Event_UPDATE: func(e *scheduler.Event) error {
				err := events.AcknowledgeUpdates(r.getCaller).HandleEvent(e)
				if err != nil {
					return err
				}
				status := e.GetUpdate().GetStatus()
				switch status.GetState() {
				case mesos.TASK_FINISHED:
					r.done = true
					fmt.Println("Task completed")
				}
				return nil
			},
			scheduler.Event_SUBSCRIBED: func(e *scheduler.Event) (err error) {
				frameworkID := e.GetSubscribed().GetFrameworkID().GetValue()
				r.caller = calls.FrameworkCaller(frameworkID).Apply(r.caller)
				fmt.Println("Subscribed!")
				return nil
			},
		}),
	)
}

func New(profile *config.Profile) *Runner {
	endpoint := profile.Endpoint()
	endpoint.Path = config.SchedulerAPIPath
	caller := httpsched.NewCaller(httpcli.New(
		httpcli.Endpoint(endpoint.String()),
		httpcli.Codec(&encoding.FramingJSONCodec),
		httpcli.Do(httpcli.With(
			httpcli.Transport(func(t *http.Transport) {
				t.ResponseHeaderTimeout = 15 * time.Second
				t.MaxIdleConnsPerHost = 2 // don't depend on go's default
			}),
		)),
	))
	caller = LoggingCaller().Apply(caller)
	return &Runner{
		task:     profile.TaskInfo,
		shutdown: make(chan struct{}),
		caller:   caller,
		random:   rand.New(rand.NewSource(time.Now().Unix())),
		framework: &mesos.FrameworkInfo{
			ID:   &mesos.FrameworkID{Value: ""},
			Name: "mesos-cli",
		},
	}
}

func (r *Runner) Run() error {
	return controller.New().Run(controller.Config{
		Context:            r,
		Framework:          r.framework,
		Caller:             r.caller,
		Handler:            r.Mux(),
		RegistrationTokens: backoff.Notifier(1*time.Second, 15*time.Second, r.shutdown),
	})
}

func (r *Runner) Done() bool { return r.done }

func (r *Runner) Error(err error) {
	r.shutdown <- struct{}{}
	r.done = true
}

func (r *Runner) FrameworkID() string {
	return r.framework.GetID().Value
}

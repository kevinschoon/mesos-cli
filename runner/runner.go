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
	"github.com/vektorlab/mesos-cli/state"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type ErrRunnerFailed struct{ *scheduler.Event }

func (e ErrRunnerFailed) Error() string {
	return fmt.Sprintf("Runner failed [%s]", e.Type.String())
}

func Mux(db *state.State, c *Context) events.Handler {
	return events.NewMux(
		events.DefaultHandler(
			events.HandlerFunc(
				func(e *scheduler.Event) (err error) {
					// Unknown scheduler event
					return ErrRunnerFailed{e}
				}),
		),
		events.MapFuncs(map[scheduler.Event_Type]events.HandlerFunc{
			scheduler.Event_SUBSCRIBED: func(e *scheduler.Event) (err error) {
				c.framework.ID.Value = e.Subscribed.FrameworkID.Value
				// Apply the mesos FrameworkID to our caller
				c.caller = calls.FrameworkCaller(c.framework.ID.Value).Apply(c.caller)
				return err
			},
			scheduler.Event_OFFERS: func(e *scheduler.Event) (err error) {
			loop:
				for _, offer := range e.GetOffers().GetOffers() {
					tasks := []mesos.TaskInfo{}
					flattened := mesos.Resources(offer.Resources).Flatten()
					for i := 0; i < db.Total(); i++ {
						task := db.Pop()
						if task == nil {
							break
						}
						if flattened.ContainsAll(mesos.Resources(task.Resources).Flatten()) {
							task.AgentID = offer.GetAgentID()
							tasks = append(tasks, *task)
							flattened.Subtract(task.Resources...)
						} else {
							db.Append(task)
						}
					}
					accept := calls.Accept(calls.OfferOperations{calls.OpLaunch(tasks...)}.WithOffers(offer.ID))
					_, err = c.caller.Call(accept)
					if err != nil {
						break loop
					}
				}
				return err
			},
			scheduler.Event_INVERSE_OFFERS: func(e *scheduler.Event) (err error) {
				return err
			},
			scheduler.Event_RESCIND: func(e *scheduler.Event) (err error) {
				return err
			},
			scheduler.Event_RESCIND_INVERSE_OFFER: func(e *scheduler.Event) (err error) {
				return err
			},
			scheduler.Event_UPDATE: func(e *scheduler.Event) (err error) {
				err = events.AcknowledgeUpdates(func() calls.Caller { return c.caller }).HandleEvent(e)
				db.Update(e.Update.Status)
				return err
			},
			scheduler.Event_MESSAGE: func(e *scheduler.Event) (err error) {
				return err
			},
			scheduler.Event_FAILURE: func(e *scheduler.Event) (err error) {
				return err
			},
			scheduler.Event_ERROR: func(e *scheduler.Event) (err error) {
				return err
			},
			scheduler.Event_HEARTBEAT: func(e *scheduler.Event) (err error) {
				return err
			},
		}),
	)
}

func caller(profile *config.Profile) calls.Caller {
	wrap := calls.Decorator(func(h calls.Caller) calls.Caller {
		return calls.CallerFunc(func(call *scheduler.Call) (mesos.Response, error) {
			LogCall(call, profile.Log())
			return h.Call(call)
		})
	})
	return wrap.Apply(
		httpsched.NewCaller(httpcli.New(
			httpcli.Endpoint(profile.Scheduler().String()),
			httpcli.Codec(&encoding.FramingJSONCodec),
			httpcli.Do(httpcli.With(
				httpcli.Transport(func(t *http.Transport) {
					t.ResponseHeaderTimeout = 15 * time.Second
					t.MaxIdleConnsPerHost = 2
				}),
			)),
		)),
	)
}

func handler(profile *config.Profile, db *state.State, ctx *Context) events.Handler {
	wrap := events.Decorator(
		func(h events.Handler) events.Handler {
			return events.HandlerFunc(
				func(e *scheduler.Event) error {
					LogEvent(e, profile.Log())
					return h.HandleEvent(e)
				},
			)
		},
	)
	return wrap.Apply(Mux(db, ctx))
}

func Run(profile *config.Profile, tasks []*mesos.TaskInfo) (err error) {
	var wg sync.WaitGroup
	// TODO: Expand to support multiple tasks and LaunchGroups/Nested Containers
	db := state.New(tasks, profile.Restart)
	sched := controller.New()
	ctx := &Context{
		caller:    caller(profile),
		framework: profile.Framework(),
		shutdown:  make(chan struct{}),
		random:    rand.New(rand.NewSource(time.Now().Unix())),
	}
	cfg := controller.Config{
		Context:            ctx,
		Framework:          ctx.framework,
		Caller:             ctx.caller,
		Handler:            handler(profile, db, ctx),
		RegistrationTokens: backoff.Notifier(1*time.Second, 15*time.Second, ctx.shutdown),
	}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err = sched.Run(cfg)
		if err != nil {
			profile.Log().Warn(
				"scheduler",
				zap.String("error", err.Error()),
			)
		} else {
			profile.Log().Info(
				"scheduler",
				zap.String("message", "scheduler has shut down"),
			)
		}
		db.Done()
	}()
	go func() {
		defer wg.Done()
		err = db.Monitor()
		if err != nil {
			profile.Log().Warn(
				"state",
				zap.String("error", err.Error()),
			)
		} else {
			profile.Log().Info("state",
				zap.String("message", "state db has shutdown"),
			)
		}
		//ctx.caller.Call(calls.TearDown()) TODO
		ctx.done = true
	}()
	wg.Wait()
	return err
}

type Context struct {
	caller    calls.Caller
	scheduled bool
	shutdown  chan struct{}
	done      bool
	err       error
	framework *mesos.FrameworkInfo
	random    *rand.Rand
}

func (c Context) Done() bool { return c.done }
func (c *Context) Error(err error) {
	c.err = err
	c.done = true
}
func (c Context) FrameworkID() string {
	return c.framework.ID.Value
}

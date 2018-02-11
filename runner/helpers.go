package runner

import (
	mesos "github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/scheduler"
	"go.uber.org/zap"
)

func LogEvent(e *scheduler.Event, log *zap.Logger) {
	eventType := zap.String("type", e.Type.String())
	if msg := log.Check(zap.DebugLevel, "EVENT"); msg != nil {
		msg.Write(
			eventType,
			zap.Any("message", e),
		)
		return
	}
	switch *e.Type {
	case scheduler.Event_SUBSCRIBED:
		master := e.Subscribed.MasterInfo
		framework := e.Subscribed.FrameworkID
		log.Info(
			"EVENT",
			eventType,
			zap.Any("master", master.Address),
			zap.String("FrameworkID", framework.Value),
		)
	case scheduler.Event_UPDATE:
		if terminal(*e.Update.Status.State) {
			log.Warn(
				"EVENT",
				eventType,
				zap.Any("event", e),
			)
		}
	}
}

func LogCall(c *scheduler.Call, log *zap.Logger) {
	callType := zap.String("type", c.Type.String())
	if msg := log.Check(zap.DebugLevel, "CALL"); msg != nil {
		msg.Write(
			callType,
			zap.Any("event", c),
		)
		return
	}
	switch *c.Type {
	case scheduler.Call_SUBSCRIBE:
		log.Info(
			"CALL",
			callType,
			zap.Any("framework", c.Subscribe.FrameworkInfo),
		)
	case scheduler.Call_ACCEPT:
		if len(c.Accept.Operations) > 0 {
			for _, op := range c.Accept.Operations {
				log.Info(
					"OPERATION",
					callType,
					zap.Any("operation", op),
					zap.Any("offers", c.Accept.OfferIDs),
				)
			}
		}
	}
}

func terminal(state mesos.TaskState) bool {
	switch state {
	case mesos.TASK_FAILED:
		return true
	case mesos.TASK_KILLED:
		return true
	case mesos.TASK_ERROR:
		return true
	case mesos.TASK_LOST:
		return true
	case mesos.TASK_DROPPED:
		return true
	//case mesos.TASK_UNREACHABLE:
	//	return true
	case mesos.TASK_GONE:
		return true
	case mesos.TASK_GONE_BY_OPERATOR:
		return true
	case mesos.TASK_UNKNOWN:
		return true
	}
	return false
}

package runner

import (
	"github.com/mesos/mesos-go/scheduler"
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
		log.Info(
			"EVENT",
			eventType,
			//zap.String("data", string(e.Update.Status.Data)),
			zap.String("state", e.Update.Status.State.String()),
			zap.String("message", e.Update.Status.Message),
		)
	case scheduler.Event_OFFERS:
		log.Info(
			"EVENT",
			eventType,
			//zap.Any("offers", e.Offers.Offers),
		)
	default:
		log.Info(
			"EVENT",
			eventType,
		)
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
		log.Info(
			"CALL",
			callType,
			zap.Any("operations", c.Accept.Operations),
			zap.Any("offers", c.Accept.OfferIDs),
		)
	default:
		log.Info(
			"CALL",
			callType,
		)
	}
}

package mailgunner

import (
	"github.com/dtynn/caesar/request"
)

const (
	EventTypeDelivered    = "delivered"
	EventTypeDropped      = "dropped"
	EventTypeBounced      = "bounced"
	EventTypeComplained   = "complained"
	EventTypeUnsubscribed = "unsubscribed"
	EventTypeClicked      = "clicked"
	EventTypeOpened       = "opened"
)

var (
	EventTypes = []string{
		EventTypeDelivered,
		EventTypeDropped,
		EventTypeBounced,
		EventTypeComplained,
		EventTypeUnsubscribed,
		EventTypeClicked,
		EventTypeOpened,
	}
)

func (this *HandlerMod) eventWebHooker(c *request.C) {
	if !CheckSignatureFromRequest(c.Req, this.cfg.APIKey) {
		c.Abort(0, errInvalidSignature)
		return
	}

	eventType := c.Req.PostFormValue("event")
	switch eventType {
	case EventTypeDelivered, EventTypeDropped, EventTypeBounced, EventTypeComplained, EventTypeUnsubscribed, EventTypeClicked, EventTypeOpened:

	default:
		this.log(c, errInvalidEventType.ErrorMsg, "type:", eventType)
		c.Abort(0, errInvalidEventType)
		return
	}

	go this.doEvent(eventType, c)

	return
}

func (this *HandlerMod) doEvent(eventType string, c *request.C) {
	hookers, ok := this.eventHookers[eventType]
	if !ok || hookers == nil || len(hookers) == 0 {
		return
	}

	for i, h := range hookers {
		if err := h(c.Req, this.cfg, eventType); err != nil {
			this.log(c, "Hooker:", i, "event hooker error:", err)
		}
	}
}

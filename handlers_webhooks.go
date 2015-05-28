package mailgunner

import (
	"net/http"

	"github.com/tukdesk/httputils/jsonutils"
	"github.com/zenazn/goji/web"
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

func (this *HandlerMod) eventWebHooker(c web.C, w http.ResponseWriter, r *http.Request) {
	if !CheckSignatureFromRequest(r, this.cfg.APIKey) {
		jsonutils.OutputJsonError(errInvalidSignature, w, r)
		return
	}

	eventType := r.PostFormValue("event")
	switch eventType {
	case EventTypeDelivered, EventTypeDropped, EventTypeBounced, EventTypeComplained, EventTypeUnsubscribed, EventTypeClicked, EventTypeOpened:

	default:
		this.log(&c, w, r, errInvalidEventType.ErrorMsg, "type:", eventType)
		jsonutils.OutputJsonError(errInvalidEventType, w, r)
		return
	}

	go this.doEvent(&c, w, r, eventType)

	return
}

func (this *HandlerMod) doEvent(c *web.C, w http.ResponseWriter, r *http.Request, eventType string) {
	hookers, ok := this.eventHookers[eventType]
	if !ok || hookers == nil || len(hookers) == 0 {
		return
	}

	for i, h := range hookers {
		if err := h(r, this.cfg, eventType); err != nil {
			this.log(c, w, r, "Hooker:", i, "event hooker error:", err)
		}
	}
}

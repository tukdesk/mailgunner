package mailgunner

import (
	"net/http"

	"github.com/dtynn/caesar/httputils/jsonutils"
)

const (
	// use status code 406 so that mailgun will not retry
	statusNotAcceptable = 406
)

var (
	// callbacks
	errInvalidSignature = jsonutils.NewAPIError(statusNotAcceptable, 990101, "invalid signature")
	errInvalidMessage   = jsonutils.NewAPIError(statusNotAcceptable, 990102, "invalid message")
	errInvalidEventType = jsonutils.NewAPIError(statusNotAcceptable, 990103, "invalid event type")

	// send
	errFromRequired    = jsonutils.NewAPIError(http.StatusBadRequest, 990104, "from required")
	errSubjectRequired = jsonutils.NewAPIError(http.StatusBadRequest, 990105, "subject required")
	errTextRequired    = jsonutils.NewAPIError(http.StatusBadRequest, 990106, "text required")
	errRcptsRequired   = jsonutils.NewAPIError(http.StatusBadRequest, 990107, "recipients required")
)

func newErrSendFailure(msg string) error {
	return jsonutils.NewAPIError(http.StatusInternalServerError, 990108, msg)
}

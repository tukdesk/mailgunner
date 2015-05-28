package mailgunner

import (
	"net/http"

	"github.com/mailgun/mailgun-go"
	"github.com/tukdesk/httputils/jsonutils"
	"github.com/zenazn/goji/web"
)

type SendArgs struct {
	From    string            `json:"from"`
	Subject string            `json:"subject"`
	Text    string            `json:"text"`
	Rcpts   []string          `json:"rcpts"`
	Headers map[string]string `json:"headers"`

	Timestamp string `json:"timestamp"`
	Token     string `json:"token"`
	Signature string `json:"signature"`
}

type SendResult struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

func (this *HandlerMod) send(c web.C, w http.ResponseWriter, r *http.Request) {
	args := &SendArgs{}
	if err := jsonutils.GetJsonArgsFromRequest(r, args); err != nil {
		jsonutils.OutputJsonError(newErrInvalidRequestBody(err.Error()), w, r)
		return
	}

	if !CheckSignature(this.cfg.APIKey, args.Token, args.Timestamp, args.Signature) {
		jsonutils.OutputJsonError(errInvalidSignature, w, r)
		return
	}

	if args.From == "" {
		jsonutils.OutputJsonError(errFromRequired, w, r)
		return
	}

	if args.Subject == "" {
		jsonutils.OutputJsonError(errSubjectRequired, w, r)
		return
	}

	if args.Text == "" {
		jsonutils.OutputJsonError(errTextRequired, w, r)
		return
	}

	if args.Rcpts == nil || len(args.Rcpts) == 0 {
		jsonutils.OutputJsonError(errRcptsRequired, w, r)
		return
	}

	mail := mailgun.NewMessage(args.From, args.Subject, args.Text, args.Rcpts...)
	if args.Headers != nil && len(args.Headers) > 0 {
		for field, value := range args.Headers {
			mail.AddHeader(field, value)
		}
	}

	msg, id, err := this.m.Send(mail)

	if err != nil {
		jsonutils.OutputJsonError(newErrSendFailure(err.Error()), w, r)
		return
	}

	res := &SendResult{
		Message: msg,
		Id:      id,
	}

	jsonutils.OutputJson(res, w, r)
	return
}

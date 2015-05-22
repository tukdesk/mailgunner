package mailgunner

import (
	"github.com/dtynn/caesar/httputils/jsonutils"
	"github.com/dtynn/caesar/request"
	"github.com/mailgun/mailgun-go"
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

func (this *handlers) send(c *request.C) {
	args := &SendArgs{}
	if err := jsonutils.GetJsonArgsFromContext(c, args); err != nil {
		c.Abort(0, err)
		return
	}

	if !CheckSignature(this.gunner.cfg.APIKey, args.Token, args.Timestamp, args.Signature) {
		c.Abort(0, errInvalidSignature)
		return
	}

	if args.From == "" {
		c.Abort(0, errFromRequired)
		return
	}

	if args.Subject == "" {
		c.Abort(0, errSubjectRequired)
		return
	}

	if args.Text == "" {
		c.Abort(0, errTextRequired)
		return
	}

	if args.Rcpts == nil || len(args.Rcpts) == 0 {
		c.Abort(0, errRcptsRequired)
		return
	}

	mail := mailgun.NewMessage(args.From, args.Subject, args.Text, args.Rcpts...)
	if args.Headers != nil && len(args.Headers) > 0 {
		for field, value := range args.Headers {
			mail.AddHeader(field, value)
		}
	}

	msg, id, err := this.gunner.m.Send(mail)
	if err != nil {
		c.Abort(0, newErrSendFailure(err.Error()))
		return
	}

	res := &SendResult{
		Message: msg,
		Id:      id,
	}

	jsonutils.OutputJsonWithC(res, c)
	return
}

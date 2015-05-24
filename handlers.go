package mailgunner

import (
	"fmt"

	"github.com/dtynn/caesar"
	"github.com/dtynn/caesar/httputils/jsonutils"
	"github.com/dtynn/caesar/request"
	"github.com/mailgun/mailgun-go"
)

type HandlerMod struct {
	cfg Config

	storers      []storer
	eventHookers map[string][]eventHooker

	m mailgun.Mailgun
}

func NewHandlerMod(cfg Config) (*HandlerMod, error) {
	if cfg.MailDomain == "" {
		return nil, fmt.Errorf("mail domain required")
	}

	if cfg.PublicAPIKey == "" {
		return nil, fmt.Errorf("public api key required")
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api key required")
	}

	h := &HandlerMod{
		cfg:          cfg,
		storers:      []storer{},
		eventHookers: map[string][]eventHooker{},
		m:            mailgun.NewMailgun(cfg.MailDomain, cfg.APIKey, cfg.PublicAPIKey),
	}

	return h, nil
}

func (this *HandlerMod) log(c *request.C, v ...interface{}) {
	if this.cfg.Debug {
		c.Logger.Info(v...)
	}
}

func (this *HandlerMod) AddStorers(s storer) {
	if s == nil {
		return
	}
	if this.storers == nil {
		this.storers = []storer{s}
		return
	}

	this.storers = append(this.storers, s)
	return
}

func (this *HandlerMod) AddEventHooker(eventType string, hooker eventHooker) {
	if hooker == nil {
		return
	}

	hookers, ok := this.eventHookers[eventType]
	if !ok {
		this.eventHookers[eventType] = []eventHooker{hooker}
		return
	}

	this.eventHookers[eventType] = append(hookers, hooker)
}

func (this *HandlerMod) Blueprint() *caesar.Blueprint {
	bp, _ := caesar.NewBlueprint("/mailgun")
	bp.Post("/send", this.send)
	bp.Post("/message/store", this.storeMessage)
	bp.Post("/webhook", this.eventWebHooker)
	bp.SetErrorHandler(jsonutils.OutputJsonError)
	bp.SetNotFoundHandler(jsonutils.RouteNotFound)
	return bp
}

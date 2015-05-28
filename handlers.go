package mailgunner

import (
	"fmt"
	"net/http"

	"github.com/mailgun/mailgun-go"
	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/tukdesk/httputils/jsonutils"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
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

func (this *HandlerMod) log(c *web.C, w http.ResponseWriter, r *http.Request, v ...interface{}) {
	logger := gojimiddleware.GetRequestLogger(c, w, r)
	logger.Info(v...)
	return
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

func (this *HandlerMod) RegisterMux(app *web.Mux) {
	m := web.New()
	m.Post("/send", this.send)
	m.Post("/message/store", this.storeMessage)
	m.Post("/webhook", this.eventWebHooker)
	m.NotFound(jsonutils.NotFoundHandler)
	m.Use(middleware.SubRouter)

	app.Handle("/mailgun/*", m)
	return
}

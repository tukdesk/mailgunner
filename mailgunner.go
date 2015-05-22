package mailgunner

import (
	"fmt"
	"net/http"

	"github.com/dtynn/caesar"
	"github.com/dtynn/caesar/httputils/jsonutils"
	"github.com/dtynn/caesar/request"
	"github.com/mailgun/mailgun-go"
)

type storer func(req *http.Request, cfg Config, msg *GunMessage) error
type eventHooker func(req *http.Request, cfg Config, eventType string) error

type Gunner struct {
	cfg      Config
	handlers *handlers

	storers      []storer
	eventHookers map[string][]eventHooker

	app *caesar.Caesar
	m   mailgun.Mailgun
}

func NewGunner(cfg Config) (*Gunner, error) {
	if cfg.Addr == "" {
		return nil, fmt.Errorf("listen addr required")
	}

	if cfg.MailDomain == "" {
		return nil, fmt.Errorf("mail domain required")
	}

	if cfg.PublicAPIKey == "" {
		return nil, fmt.Errorf("public api key required")
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api key required")
	}

	gunner := &Gunner{
		cfg:          cfg,
		storers:      []storer{},
		eventHookers: map[string][]eventHooker{},
	}

	// handlers
	hdls := &handlers{
		gunner: gunner,
	}
	gunner.handlers = hdls

	// app
	appCfg := &caesar.Config{
		Prefix: cfg.URLPrefix,
	}
	app := caesar.New()
	app.SetConfig(appCfg)

	app.Post("/mailgun/send", hdls.send)
	app.Post("/mailgun/message/store", hdls.storeMessage)
	app.Post("/mailgun/webhook", hdls.eventWebHooker)

	app.SetErrorHandler(jsonutils.OutputJsonError)
	app.SetNotFoundHandler(jsonutils.RouteNotFound)
	app.AddAfterRequest(request.TimerAfterHandler)

	gunner.app = app

	// mailgun api impl
	m := mailgun.NewMailgun(cfg.MailDomain, cfg.APIKey, cfg.PublicAPIKey)
	gunner.m = m

	return gunner, nil
}

func (this *Gunner) AddStorers(s storer) {
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

func (this *Gunner) AddEventHooker(eventType string, hooker eventHooker) {
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

func (this *Gunner) Run() error {
	return this.app.Run(this.cfg.Addr)
}

func (this *Gunner) Close() {
	this.app.Close()
}

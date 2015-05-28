package mailgunner

import (
	"net/http"

	"github.com/tukdesk/httputils/gojimiddleware"
)

type storer func(req *http.Request, cfg Config, msg *GunMessage) error
type eventHooker func(req *http.Request, cfg Config, eventType string) error

type Gunner struct {
	handlers *HandlerMod

	app *gojimiddleware.App
}

func NewGunner(cfg Config) (*Gunner, error) {
	// handlers
	handlers, err := NewHandlerMod(cfg)
	if err != nil {
		return nil, err
	}

	gunner := &Gunner{
		handlers: handlers,
	}

	return gunner, nil
}

func (this *Gunner) build() {
	if this.app != nil {
		return
	}

	app := gojimiddleware.NewApp()

	this.handlers.RegisterMux(app.Mux())

	app.Mux().Use(gojimiddleware.RequestLogger)
	app.Mux().Use(gojimiddleware.RequestTimer)

	this.app = app
	return
}

func (this *Gunner) AddStorers(s storer) {
	this.handlers.AddStorers(s)
	return
}

func (this *Gunner) AddEventHooker(eventType string, hooker eventHooker) {
	this.handlers.AddEventHooker(eventType, hooker)
	return
}

func (this *Gunner) Run(addr string) error {
	this.build()
	return this.app.Run(addr)
}

func (this *Gunner) Close() {
	this.app.Close()
}

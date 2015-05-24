package mailgunner

import (
	"net/http"

	"github.com/dtynn/caesar"
	"github.com/dtynn/caesar/request"
)

type storer func(req *http.Request, cfg Config, msg *GunMessage) error
type eventHooker func(req *http.Request, cfg Config, eventType string) error

type Gunner struct {
	handlers *HandlerMod

	app *caesar.Caesar
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

	app := caesar.New()

	bp := this.handlers.Blueprint()
	app.RegisterBlueprint(bp)
	app.AddAfterRequest(request.TimerAfterHandler)

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

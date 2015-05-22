package mailgunner

import (
	"github.com/dtynn/caesar/request"
)

type handlers struct {
	gunner *Gunner
}

func (this *handlers) log(c *request.C, v ...interface{}) {
	if this.gunner.cfg.Debug {
		c.Logger.Info(v...)
	}
}

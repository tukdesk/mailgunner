package mailgunner

import (
	"github.com/dtynn/caesar/request"
)

func (this *HandlerMod) storeMessage(c *request.C) {
	if !CheckSignatureFromRequest(c.Req, this.cfg.APIKey) {
		timestamp, token, signature := GetSignatureStuffsFromReq(c.Req)
		this.log(c, errInvalidSignature.ErrorMsg, "timestamp:", timestamp, "token:", token, "signature:", signature)
		c.Abort(0, errInvalidSignature)
		return
	}

	message, err := NewGunMessageFromRequest(c.Req)
	if err != nil {
		c.Abort(0, errInvalidMessage)
		return
	}

	go this.doStore(c, message)

	return
}

func (this *HandlerMod) doStore(c *request.C, message *GunMessage) {
	if this.storers == nil || len(this.storers) == 0 {
		return
	}

	for i, s := range this.storers {
		if err := s(c.Req, this.cfg, message); err != nil {
			this.log(c, "Storer:", i, "storing error:", err.Error())
		}
	}
}

package mailgunner

import (
	"net/http"

	"github.com/tukdesk/httputils/jsonutils"
	"github.com/zenazn/goji/web"
)

func (this *HandlerMod) storeMessage(c web.C, w http.ResponseWriter, r *http.Request) {
	if !CheckSignatureFromRequest(r, this.cfg.APIKey) {
		timestamp, token, signature := GetSignatureStuffsFromReq(r)
		this.log(&c, w, r, errInvalidSignature.ErrorMsg, "timestamp:", timestamp, "token:", token, "signature:", signature)
		jsonutils.OutputJsonError(errInvalidSignature, w, r)
		return
	}

	message, err := NewGunMessageFromRequest(r)
	if err != nil {
		jsonutils.OutputJsonError(errInvalidMessage, w, r)
		return
	}

	go this.doStore(&c, w, r, message)

	return
}

func (this *HandlerMod) doStore(c *web.C, w http.ResponseWriter, r *http.Request, message *GunMessage) {
	if this.storers == nil || len(this.storers) == 0 {
		return
	}

	for i, s := range this.storers {
		if err := s(r, this.cfg, message); err != nil {
			this.log(c, w, r, "Storer:", i, "storing error:", err.Error())
		}
	}
}

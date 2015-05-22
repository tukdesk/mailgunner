package mailgunner

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
)

func Sign(apiKey, data []byte) string {
	h := hmac.New(sha256.New, apiKey)
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GenerateSignature(apiKey, token, timestamp string) string {
	data := timestamp + token
	return Sign([]byte(apiKey), []byte(data))
}

func CheckSignature(apiKey, token, timestamp, signature string) bool {
	return GenerateSignature(apiKey, token, timestamp) == signature
}

func GetSignatureStuffsFromReq(req *http.Request) (string, string, string) {
	return req.PostFormValue("timestamp"), req.PostFormValue("token"), req.PostFormValue("signature")
}

func CheckSignatureFromRequest(req *http.Request, apiKey string) bool {
	timestamp, token, signature := GetSignatureStuffsFromReq(req)
	return CheckSignature(apiKey, token, timestamp, signature)
}

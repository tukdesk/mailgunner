package mailgunner

import (
	"encoding/json"
	"net/http"

	"github.com/mailgun/mailgun-go"
)

type GunMessage struct {
	message *mailgun.StoredMessage
}

func NewGunMessageFromRequest(req *http.Request) (*GunMessage, error) {

	message := &mailgun.StoredMessage{}
	message.Recipients = req.PostFormValue("recipient")
	message.Sender = req.PostFormValue("sender")
	message.From = req.PostFormValue("from")
	message.Subject = req.PostFormValue("subject")
	message.BodyPlain = req.PostFormValue("body-plain")
	message.StrippedText = req.PostFormValue("stripped-text")
	message.StrippedSignature = req.PostFormValue("stripped-signature")
	message.BodyHtml = req.PostFormValue("body-html")
	message.StrippedHtml = req.PostFormValue("stripped-html")
	message.MessageUrl = req.PostFormValue("message-url")

	// attachments
	attachments := []mailgun.StoredAttachment{}
	if attachmentsVal := req.PostFormValue("attachments"); len(attachmentsVal) > 0 {
		if err := json.Unmarshal([]byte(attachmentsVal), &attachments); err != nil {
			return nil, err
		}
	}
	message.Attachments = attachments

	// content id map
	contentIdMap := map[string]interface{}{}
	if contentIdMapVal := req.PostFormValue("content-id-map"); len(contentIdMapVal) > 0 {
		if err := json.Unmarshal([]byte(contentIdMapVal), &contentIdMap); err != nil {
			return nil, err
		}
	}
	message.ContentIDMap = contentIdMap

	// message header
	messageHeaders := [][]string{}
	if messageHeadersVal := req.PostFormValue("message-headers"); len(messageHeadersVal) > 0 {
		if err := json.Unmarshal([]byte(messageHeadersVal), &messageHeaders); err != nil {
			return nil, err
		}
	}
	message.MessageHeaders = messageHeaders

	return &GunMessage{
		message: message,
	}, nil
}

func (this *GunMessage) Message() *mailgun.StoredMessage {
	return this.message
}

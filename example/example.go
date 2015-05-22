package main

import (
	"log"
	"net/http"

	"github.com/tukdesk/mailgunner"
)

func main() {
	guunerCfg := mailgunner.Config{
		URLPrefix:    "",
		MailDomain:   "test.com",
		PublicAPIKey: "abcde",
		APIKey:       "key-fghijk",
		Addr:         "127.0.0.1:56666",
		Debug:        true,
	}
	gunner, err := mailgunner.NewGunner(guunerCfg)
	if err != nil {
		log.Fatalln(err)
	}

	gunner.AddStorers(echoStorer)
	for _, typ := range mailgunner.EventTypes {
		gunner.AddEventHooker(typ, echoHooker)
	}

	if err := gunner.Run(); err != nil {
		log.Fatalln(err)
	}
}

func echoStorer(req *http.Request, cfg mailgunner.Config, msg *mailgunner.GunMessage) error {
	log.Printf("echo storer %#v\n", msg.Message())
	return nil
}

func echoHooker(req *http.Request, cfg mailgunner.Config, eventType string) error {
	log.Println("echo hooker", eventType, req.PostForm)
	return nil
}

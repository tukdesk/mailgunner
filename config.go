package mailgunner

type Config struct {
	URLPrefix string
	Addr      string

	MailDomain   string
	PublicAPIKey string
	APIKey       string

	Debug bool
}

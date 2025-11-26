package types

type EmailConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Recipients []string
}

type Message struct {
	From    string
	To      []string
	Subject string
	Body    string
}

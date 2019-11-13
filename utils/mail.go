package utils

import (
	"os"
	"strconv"
	mail "github.com/go-mail/mail"
)

// Mail send mail in HTML format
func Mail(to []string, title string, content string) error {
	host := os.Getenv("MAIL_HOST")
	if host == "" {
		log.Panic().Str("field", "MAIL_HOST").Msg("env required")
	}
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		return err
	}
	if port == 0 {
		log.Panic().Str("field", "MAIL_PORT").Msg("env required")
	}
	user := os.Getenv("MAIL_USER")
	if user == "" {
		log.Panic().Str("field", "MAIL_USER").Msg("env required")
	}
	passwd := os.Getenv("MAIL_PASSWD")
	if passwd == "" {
		log.Panic().Str("field", "MAIL_PASSWD").Msg("env required")
	}

	// filter empty
	var list []string
	for _, t := range to {
		if t != "" {
			list = append(list, t)
		}
	}

	if len(list) == 0 {
		return nil
	}

	m := mail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", list...)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	d := mail.NewDialer(host, port, user, passwd)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	return d.DialAndSend(m)
}

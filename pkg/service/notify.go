package service

import (
	"crypto/tls"
	"github.com/go-gomail/gomail"
	"github.com/sirupsen/logrus"
)

type (
	NotifyConfig struct {
		Host          string
		HostUser      string
		HostPwd       string
		From          string
		Port          int
		TlsSkipVerify bool
	}

	Notify struct {
		logger *logrus.Logger
		config *NotifyConfig
	}
)

func NewNotify(logger *logrus.Logger, config *NotifyConfig) *Notify {
	return &Notify{logger, config}
}

func (s *Notify) SendMail(to, sub, body string) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.config.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", sub)
	msg.SetBody("text/html", body)
	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.HostUser, s.config.HostPwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: s.config.TlsSkipVerify}
	if err := d.DialAndSend(msg); err != nil {
		s.logger.Errorf("service: notifier send mail error: %v", err)
	}
}

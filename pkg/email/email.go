package email

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"github.com/hammer-code/lms-be/pkg/ngelog"
)

// SMTPSender interface for dependency injection
type SMTPSender interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

// DefaultSMTPSender implements SMTPSender using the real smtp.SendMail
type DefaultSMTPSender struct{}

func (s *DefaultSMTPSender) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return smtp.SendMail(addr, a, from, to, msg)
}

type SMTP struct {
	Email    string
	Password string
	Host     string
	Port     string
}

type Receiver struct {
	Email string
	Data  any
}

type PayloadEmail struct {
	HtmlTemplate *template.Template
	Subject      string
	Mime         string
	Sender       SMTP
	Receiver     []Receiver
	CC           []string
	smtpSender   SMTPSender // injected SMTP sender
}

type failedToSend struct {
	Subject  string
	Receiver Receiver
}

// Updated constructor with optional SMTPSender
func NewSendEmail(ctx context.Context, smtp SMTP, mime string, subject string, HtmlTemplate *template.Template) PayloadEmail {
	return PayloadEmail{
		HtmlTemplate: HtmlTemplate,
		Subject:      subject,
		Mime:         mime,
		Sender:       smtp,
		CC:           []string{},
		smtpSender:   &DefaultSMTPSender{}, // Use default implementation
	}
}

// New constructor for testing that allows injecting the SMTPSender
func NewSendEmailWithSender(ctx context.Context, smtp SMTP, mime string, subject string, HtmlTemplate *template.Template, sender SMTPSender) PayloadEmail {
	if sender == nil {
		sender = &DefaultSMTPSender{}
	}
	return PayloadEmail{
		HtmlTemplate: HtmlTemplate,
		Subject:      subject,
		Mime:         mime,
		Sender:       smtp,
		CC:           []string{},
		smtpSender:   sender,
	}
}

func (p *PayloadEmail) AddReceiver(ctx context.Context, receiver Receiver) error {
	if receiver.Email == "" {
		ngelog.Error(ctx, "email cannot be null", nil)
		return errors.New("email cannot be null")
	}

	p.Receiver = append(p.Receiver, receiver)
	return nil
}

func (p *PayloadEmail) AddCC(ctx context.Context, ccEmail string) error {
	if ccEmail == "" {
		ngelog.Error(ctx, "cc email cannot be null", nil)
		return errors.New("cc email cannot be null")
	}

	p.CC = append(p.CC, ccEmail)
	return nil
}

func (p *PayloadEmail) AddMultipleCC(ctx context.Context, ccEmails []string) error {
	for _, email := range ccEmails {
		if email == "" {
			ngelog.Error(ctx, "cc email cannot be null", nil)
			return errors.New("cc email cannot be null")
		}
	}

	p.CC = append(p.CC, ccEmails...)
	return nil
}

// Updated SendEmail method using the injected sender
func (p PayloadEmail) SendEmail(ctx context.Context) {
	var failedSendedEmails []failedToSend
	for _, receiver := range p.Receiver {

		// Create a temporary buffer to hold the parsed HTML
		var bodyContent string
		bodyBuffer := new(bytes.Buffer)
		// Execute the template with data
		if err := p.HtmlTemplate.Execute(bodyBuffer, receiver.Data); err != nil {
			ngelog.Error(ctx, fmt.Sprintf("Error executing template for %s", receiver.Email), err)
			failedSendedEmails = append(failedSendedEmails, failedToSend{
				Subject:  p.Subject,
				Receiver: receiver,
			})
			continue
		}

		bodyContent = bodyBuffer.String()
		
		// Build message headers
		headers := fmt.Sprintf("To: %s\r\n", receiver.Email)
		
		// Add CC header if there are CC recipients
		if len(p.CC) > 0 {
			headers += fmt.Sprintf("Cc: %s\r\n", strings.Join(p.CC, ", "))
		}
		
		headers += fmt.Sprintf("Subject: %s\r\n%s\r\n", p.Subject, p.Mime)
		
		message := []byte(headers + bodyContent + "\r\n")

		auth := smtp.PlainAuth("", p.Sender.Email, p.Sender.Password, p.Sender.Host)

		// Combine recipient and CC addresses for SMTP
		allRecipients := []string{receiver.Email}
		allRecipients = append(allRecipients, p.CC...)

		host := fmt.Sprintf("%s:%s", p.Sender.Host, p.Sender.Port)
		
		// Use the injected SMTP sender
		if err := p.smtpSender.SendMail(host, auth, p.Sender.Email, allRecipients, message); err != nil {
			ngelog.Error(ctx, fmt.Sprintf("failed to send to %s", receiver.Email), err)
			failedSendedEmails = append(failedSendedEmails, failedToSend{
				Subject:  p.Subject,
				Receiver: receiver,
			})
			continue
		}

		ngelog.Info(ctx, "success send email", ngelog.AddFields{
			"email":   receiver.Email,
			"subject": p.Subject,
			"cc":      p.CC,
		})
	}

	if len(failedSendedEmails) > 0 {
		ngelog.Error(ctx, "failed to send email", nil, ngelog.AddFields{
			"data": failedSendedEmails,
		})
	}
}
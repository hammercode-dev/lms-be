package usecase

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"os"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (us *usecase) ForgotPassword(ctx context.Context, emailForgot domain.ForgotPassword) (user domain.User, resetLink string, err error) {
	err = us.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		user, err = us.userRepo.FindByEmail(ctx, emailForgot.Email)
		if err != nil {
			logrus.Error("us.ForgotPassword: failed to get Email", err)
			return err
		}
		return nil
	})

	if err != nil {
		logrus.Error("us.ForgotPassword: failed to get Email", err)
		return
	}

	resetToken, err := us.jwt.GenerateAccessToken(ctx, &user)
	if err != nil {
		logrus.Error("us.ForgotPassword: failed to generate token", err)
		return
	}
	// link to reset password
	link := "http://localhost:8000/api/v1/auth/forgot_password?token=" + *resetToken

	htmlTmpl, err := os.ReadFile("./assets/kirim-email.html")
	if err != nil {
		logrus.Error("us.ForgotPassword: failed to read template file", err)
		return
	}

	tmpl, err := template.New("reset_email").Parse(string(htmlTmpl))
	if err != nil {
		logrus.Error("us.ForgotPassword: failed to parse template", err)
	}

	var bodyBuffer bytes.Buffer
	if err = tmpl.Execute(&bodyBuffer, link); err != nil {
		logrus.Error("us.ForgotPassword: failed to execute template", err)
		return
	}

	subject := "Reset Password"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"%s\r\n"+
		"%s\r\n", emailForgot.Email, subject, mime, bodyBuffer.String()))

	auth := smtp.PlainAuth("", us.cfg.SMTP_EMAIL, us.cfg.SMTP_PASSWORD, us.cfg.SMTP_HOST)

	host := fmt.Sprintf("%s:%s", us.cfg.SMTP_HOST, us.cfg.SMTP_PORT)
	if err := smtp.SendMail(host, auth, us.cfg.SMTP_EMAIL, []string{emailForgot.Email}, message); err != nil {
		return domain.User{}, "", err
	}

	return user, link, nil
}
